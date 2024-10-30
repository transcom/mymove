package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
)

// ServiceItemParamKeyData contains service item parameter keys
type ServiceItemParamKeyData struct {
	planner          route.Planner
	lookups          map[models.ServiceItemParamName]ServiceItemParamKeyLookup
	MTOServiceItemID uuid.UUID
	MTOServiceItem   models.MTOServiceItem
	PaymentRequestID uuid.UUID
	MoveTaskOrderID  uuid.UUID
	ContractCode     string
	mtoShipmentID    *uuid.UUID
	paramCache       *ServiceParamsCache
}

func NewServiceItemParamKeyData(planner route.Planner, lookups map[models.ServiceItemParamName]ServiceItemParamKeyLookup, mtoServiceItem models.MTOServiceItem, mtoShipment models.MTOShipment, contractCode string) ServiceItemParamKeyData {
	return ServiceItemParamKeyData{
		planner:          planner,
		lookups:          lookups,
		MTOServiceItem:   mtoServiceItem,
		MTOServiceItemID: mtoServiceItem.ID,
		mtoShipmentID:    &mtoShipment.ID,
		MoveTaskOrderID:  mtoShipment.MoveTaskOrderID,
		ContractCode:     contractCode,
	}
}

// ServiceItemParamKeyLookup does lookup on service item parameter keys
type ServiceItemParamKeyLookup interface {
	lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error)
}

// We don't have comprehensive lookups for all SYSTEM and PRIME params so we need a list of those that do exist.
var ServiceItemParamsWithLookups = []models.ServiceItemParamName{
	models.ServiceItemParamNameActualPickupDate,
	models.ServiceItemParamNameRequestedPickupDate,
	models.ServiceItemParamNameReferenceDate,
	models.ServiceItemParamNameDistanceZip,
	models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
	models.ServiceItemParamNameWeightAdjusted,
	models.ServiceItemParamNameWeightBilled,
	models.ServiceItemParamNameWeightEstimated,
	models.ServiceItemParamNameWeightOriginal,
	models.ServiceItemParamNameWeightReweigh,
	models.ServiceItemParamNameZipPickupAddress,
	models.ServiceItemParamNameZipDestAddress,
	models.ServiceItemParamNameMTOAvailableToPrimeAt,
	models.ServiceItemParamNameServiceAreaOrigin,
	models.ServiceItemParamNameServiceAreaDest,
	models.ServiceItemParamNameContractCode,
	models.ServiceItemParamNameCubicFeetBilled,
	models.ServiceItemParamNamePSILinehaulDom,
	models.ServiceItemParamNamePSILinehaulDomPrice,
	models.ServiceItemParamNameEIAFuelPrice,
	models.ServiceItemParamNameServicesScheduleOrigin,
	models.ServiceItemParamNameServicesScheduleDest,
	models.ServiceItemParamNameSITScheduleOrigin,
	models.ServiceItemParamNameSITScheduleDest,
	models.ServiceItemParamNameSITServiceAreaDest,
	models.ServiceItemParamNameSITServiceAreaOrigin,
	models.ServiceItemParamNameNumberDaysSIT,
	models.ServiceItemParamNameZipSITDestHHGFinalAddress,
	models.ServiceItemParamNameZipSITDestHHGOriginalAddress,
	models.ServiceItemParamNameZipSITOriginHHGOriginalAddress,
	models.ServiceItemParamNameZipSITOriginHHGActualAddress,
	models.ServiceItemParamNameDistanceZipSITDest,
	models.ServiceItemParamNameDistanceZipSITOrigin,
	models.ServiceItemParamNameCubicFeetCrating,
	models.ServiceItemParamNameDimensionHeight,
	models.ServiceItemParamNameDimensionLength,
	models.ServiceItemParamNameDimensionWidth,
	models.ServiceItemParamNameStandaloneCrate,
	models.ServiceItemParamNameStandaloneCrateCap,
	models.ServiceItemParamNameLockedPriceCents,
}

// ServiceParamLookupInitialize initializes service parameter lookup
func ServiceParamLookupInitialize(
	appCtx appcontext.AppContext,
	planner route.Planner,
	mtoServiceItem models.MTOServiceItem,
	paymentRequestID uuid.UUID,
	moveTaskOrderID uuid.UUID,
	paramCache *ServiceParamsCache,
) (*ServiceItemParamKeyData, error) {
	contract, err := fetchContractForMove(appCtx, moveTaskOrderID)
	if err != nil {
		return nil, err
	}
	s := ServiceItemParamKeyData{
		planner:          planner,
		lookups:          make(map[models.ServiceItemParamName]ServiceItemParamKeyLookup),
		MTOServiceItemID: mtoServiceItem.ID,
		MTOServiceItem:   mtoServiceItem,
		PaymentRequestID: paymentRequestID,
		MoveTaskOrderID:  moveTaskOrderID,
		paramCache:       paramCache,
		mtoShipmentID:    mtoServiceItem.MTOShipmentID,
		/*
			DefaultContractCode = TRUSS_TEST is temporarily being used here because the contract
			code is not currently accessible. This is caused by:
				- mtoServiceItem is not linked or associated with a contract record
				- MTO currently has a contractor_id but not a contract_id
			In order for this lookup's query to have accesss to a contract code there must be a contract_code field created on either the mtoServiceItem or the MTO models
			If it'll will be possible for a MTO to contain service items that are associated with different contracts
			then it would be ideal for the mtoServiceItem records to contain a contract code that can then be passed
			to this query. Otherwise the contract_code field could be added to the MTO.
		*/
		ContractCode: contract.Code,
	}

	//
	// Query and save PickupAddress & DestinationAddress upfront
	// s.serviceItemNeedsParamKey() could be used to check if the PickupAddress or DestinationAddress
	// can be used but it depends on the paramCache being set (not nil). It is possible to set the
	// paramCache to nil, especially during unit test, so not using that function for this part.
	//

	// Load data that is only used by a few service items
	var sitDestinationFinalAddress, sitDestinationOriginalAddress models.Address
	var serviceItemDimensions models.MTOServiceItemDimensions

	switch mtoServiceItem.ReService.Code {
	case models.ReServiceCodeDCRT, models.ReServiceCodeDUCRT, models.ReServiceCodeDCRTSA:
		err := appCtx.DB().Load(&mtoServiceItem, "Dimensions")
		if err != nil {
			return nil, err
		}
		serviceItemDimensions = mtoServiceItem.Dimensions
	case models.ReServiceCodeDDASIT, models.ReServiceCodeDDDSIT, models.ReServiceCodeDDFSIT, models.ReServiceCodeDDSFSC:
		// load destination address from final address on service item
		if mtoServiceItem.SITDestinationFinalAddressID != nil && *mtoServiceItem.SITDestinationFinalAddressID != uuid.Nil {
			err := appCtx.DB().Load(&mtoServiceItem, "SITDestinationFinalAddress")
			if err != nil {
				return nil, err
			}
			sitDestinationFinalAddress = *mtoServiceItem.SITDestinationFinalAddress
		}

		if mtoServiceItem.SITDestinationOriginalAddressID != nil && *mtoServiceItem.SITDestinationOriginalAddressID != uuid.Nil {
			err := appCtx.DB().Load(&mtoServiceItem, "SITDestinationOriginalAddress")
			if err != nil {
				return nil, err
			}
			sitDestinationOriginalAddress = *mtoServiceItem.SITDestinationOriginalAddress
		}
	}

	mtoServiceItem.SITDestinationFinalAddress = &sitDestinationFinalAddress
	mtoServiceItem.SITDestinationOriginalAddress = &sitDestinationOriginalAddress
	mtoServiceItem.Dimensions = serviceItemDimensions

	// Load shipment fields for service items that need them
	var mtoShipment models.MTOShipment
	var pickupAddress models.Address
	var destinationAddress models.Address

	if mtoServiceItem.ReService.Code != models.ReServiceCodeCS && mtoServiceItem.ReService.Code != models.ReServiceCodeMS {
		// Make sure there's an MTOShipment since that's nullable
		if mtoServiceItem.MTOShipmentID == nil {
			return nil, apperror.NewNotFoundError(uuid.Nil, "the shipment service item is missing a MTOShipmentID")
		}
		err := appCtx.DB().Eager("PickupAddress", "DestinationAddress", "StorageFacility", "PPMShipment").Find(&mtoShipment, mtoServiceItem.MTOShipmentID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return nil, apperror.NewNotFoundError(*mtoServiceItem.MTOShipmentID, "looking for MTOShipment")
			default:
				return nil, apperror.NewQueryError("MTOShipment", err, "")
			}
		}

		// Due to a bug in pop (https://github.com/gobuffalo/pop/issues/578), we cannot eager load the storage
		// facility's address as "StorageFacility.Address" because StorageFacility is a pointer.
		if mtoShipment.StorageFacility != nil {
			err = appCtx.DB().Load(mtoShipment.StorageFacility, "Address", "Address.Country")
			if err != nil {
				return nil, apperror.NewQueryError("Address", err, "")
			}
		}

		pickupAddress, err = getPickupAddressForService(mtoServiceItem.ReService.Code, mtoShipment)
		if err != nil {
			return nil, fmt.Errorf("not found looking for pickup address")
		}

		destinationAddress, err = getDestinationAddressForService(appCtx, mtoServiceItem.ReService.Code, mtoShipment)
		if err != nil {
			return nil, err
		}
	}

	mtoShipment.PickupAddress = &pickupAddress
	mtoShipment.DestinationAddress = &destinationAddress

	switch mtoServiceItem.ReService.Code {
	case models.ReServiceCodeDDASIT, models.ReServiceCodeDDDSIT, models.ReServiceCodeDDFSIT, models.ReServiceCodeDDSFSC, models.ReServiceCodeDOASIT, models.ReServiceCodeDOPSIT, models.ReServiceCodeDOFSIT, models.ReServiceCodeDOSFSC:
		err := appCtx.DB().Load(&mtoShipment, "SITDurationUpdates")
		if err != nil {
			return nil, err
		}
	}

	//
	// Set all lookup functions to "NOT IMPLEMENTED"
	//

	notImplementedLookup := NotImplementedLookup{}
	for _, key := range models.ValidServiceItemParamNames {
		s.lookups[key] = notImplementedLookup
	}

	// ReService code for current MTO Service Item
	serviceItemCode := mtoServiceItem.ReService.Code

	paramKeyLookups := InitializeLookups(appCtx, mtoShipment, mtoServiceItem)

	for _, paramKeyName := range ServiceItemParamsWithLookups {
		lookup, ok := paramKeyLookups[paramKeyName]
		if !ok {
			return nil, fmt.Errorf("no lookup was found for service item param key name %s", paramKeyName)
		}

		err := s.setLookup(appCtx, serviceItemCode, paramKeyName, lookup)
		if err != nil {
			return nil, err
		}
	}

	return &s, nil
}

func (s *ServiceItemParamKeyData) setLookup(appCtx appcontext.AppContext, serviceItemCode models.ReServiceCode, paramKey models.ServiceItemParamName, lookup ServiceItemParamKeyLookup) error {
	useKey, err := s.serviceItemNeedsParamKey(appCtx, serviceItemCode, paramKey)
	if useKey && err == nil {
		s.lookups[paramKey] = lookup
	} else if err != nil {
		return err
	}
	return nil
}

func InitializeLookups(appCtx appcontext.AppContext, shipment models.MTOShipment, serviceItem models.MTOServiceItem) map[models.ServiceItemParamName]ServiceItemParamKeyLookup {
	lookups := map[models.ServiceItemParamName]ServiceItemParamKeyLookup{}

	if serviceItem.SITDestinationOriginalAddress == nil {
		serviceItem.SITDestinationOriginalAddress = &models.Address{}
	}

	if serviceItem.SITDestinationFinalAddress == nil {
		serviceItem.SITDestinationFinalAddress = &models.Address{}
	}

	if serviceItem.SITOriginHHGOriginalAddress == nil {
		serviceItem.SITOriginHHGOriginalAddress = &models.Address{}
	}

	lookups[models.ServiceItemParamNameActualPickupDate] = ActualPickupDateLookup{
		MTOShipment: shipment,
	}

	lookups[models.ServiceItemParamNameRequestedPickupDate] = RequestedPickupDateLookup{
		MTOShipment: shipment,
	}

	lookups[models.ServiceItemParamNameReferenceDate] = ReferenceDateLookup{
		MTOShipment: shipment,
	}

	serviceDestinationAddress, err := GetDestinationForDistanceLookup(appCtx, shipment, &serviceItem)
	if err != nil {
		return nil
	}
	lookups[models.ServiceItemParamNameDistanceZip] = DistanceZipLookup{
		PickupAddress:      *shipment.PickupAddress,
		DestinationAddress: serviceDestinationAddress,
	}

	lookups[models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier] = FSCWeightBasedDistanceMultiplierLookup{
		MTOShipment: shipment,
	}

	lookups[models.ServiceItemParamNameWeightAdjusted] = WeightAdjustedLookup{
		MTOShipment: shipment,
	}

	lookups[models.ServiceItemParamNameWeightBilled] = WeightBilledLookup{
		MTOShipment: shipment,
	}

	lookups[models.ServiceItemParamNameWeightEstimated] = WeightEstimatedLookup{
		MTOShipment: shipment,
	}

	lookups[models.ServiceItemParamNameWeightOriginal] = WeightOriginalLookup{
		MTOShipment: shipment,
	}

	lookups[models.ServiceItemParamNameWeightReweigh] = WeightReweighLookup{
		MTOShipment: shipment,
	}

	lookups[models.ServiceItemParamNameZipPickupAddress] = ZipAddressLookup{
		Address: *shipment.PickupAddress,
	}

	lookups[models.ServiceItemParamNameZipDestAddress] = ZipAddressLookup{
		Address: *shipment.DestinationAddress,
	}

	lookups[models.ServiceItemParamNameMTOAvailableToPrimeAt] = MTOAvailableToPrimeAtLookup{}

	lookups[models.ServiceItemParamNameServiceAreaOrigin] = ServiceAreaLookup{
		Address: *shipment.PickupAddress,
	}

	lookups[models.ServiceItemParamNameServiceAreaDest] = ServiceAreaLookup{
		Address: *shipment.DestinationAddress,
	}

	lookups[models.ServiceItemParamNameContractCode] = ContractCodeLookup{}

	lookups[models.ServiceItemParamNameCubicFeetBilled] = CubicFeetBilledLookup{
		Dimensions: serviceItem.Dimensions,
	}

	lookups[models.ServiceItemParamNamePSILinehaulDom] = PSILinehaulDomLookup{
		MTOShipment: shipment,
	}

	lookups[models.ServiceItemParamNamePSILinehaulDomPrice] = PSILinehaulDomPriceLookup{
		MTOShipment: shipment,
	}

	lookups[models.ServiceItemParamNameEIAFuelPrice] = EIAFuelPriceLookup{
		MTOShipment: shipment,
	}

	lookups[models.ServiceItemParamNameServicesScheduleOrigin] = ServicesScheduleLookup{
		Address: *shipment.PickupAddress,
	}

	lookups[models.ServiceItemParamNameServicesScheduleDest] = ServicesScheduleLookup{
		Address: *shipment.DestinationAddress,
	}

	lookups[models.ServiceItemParamNameSITScheduleOrigin] = SITScheduleLookup{
		Address: *shipment.PickupAddress,
	}

	lookups[models.ServiceItemParamNameSITScheduleDest] = SITScheduleLookup{
		Address: *serviceItem.SITDestinationFinalAddress,
	}

	lookups[models.ServiceItemParamNameNumberDaysSIT] = NumberDaysSITLookup{
		MTOShipment: shipment,
	}

	lookups[models.ServiceItemParamNameZipSITDestHHGFinalAddress] = ZipAddressLookup{
		Address: *serviceItem.SITDestinationFinalAddress,
	}

	lookups[models.ServiceItemParamNameSITServiceAreaDest] = ServiceAreaLookup{
		Address: *serviceItem.SITDestinationFinalAddress,
	}

	lookups[models.ServiceItemParamNameSITServiceAreaOrigin] = ServiceAreaLookup{
		Address: *serviceItem.SITOriginHHGOriginalAddress,
	}

	lookups[models.ServiceItemParamNameZipSITDestHHGOriginalAddress] = ZipAddressLookup{
		Address: *serviceItem.SITDestinationOriginalAddress,
	}

	lookups[models.ServiceItemParamNameZipSITOriginHHGOriginalAddress] = ZipSITOriginHHGOriginalAddressLookup{
		ServiceItem: serviceItem,
	}

	lookups[models.ServiceItemParamNameZipSITOriginHHGActualAddress] = ZipSITOriginHHGActualAddressLookup{
		ServiceItem: serviceItem,
	}

	lookups[models.ServiceItemParamNameDistanceZipSITDest] = DistanceZipSITDestLookup{
		DestinationAddress:      *serviceItem.SITDestinationOriginalAddress,
		FinalDestinationAddress: *serviceItem.SITDestinationFinalAddress,
	}

	lookups[models.ServiceItemParamNameDistanceZipSITOrigin] = DistanceZipSITOriginLookup{
		ServiceItem: serviceItem,
	}

	lookups[models.ServiceItemParamNameCubicFeetCrating] = CubicFeetCratingLookup{
		Dimensions: serviceItem.Dimensions,
	}

	lookups[models.ServiceItemParamNameDimensionHeight] = DimensionHeightLookup{
		Dimensions: serviceItem.Dimensions,
	}

	lookups[models.ServiceItemParamNameDimensionLength] = DimensionLengthLookup{
		Dimensions: serviceItem.Dimensions,
	}

	lookups[models.ServiceItemParamNameDimensionWidth] = DimensionWidthLookup{
		Dimensions: serviceItem.Dimensions,
	}

	lookups[models.ServiceItemParamNameStandaloneCrate] = StandaloneCrateLookup{
		ServiceItem: serviceItem,
	}

	lookups[models.ServiceItemParamNameStandaloneCrateCap] = StandaloneCrateCapLookup{
		ServiceItem: serviceItem,
	}

	lookups[models.ServiceItemParamNameLockedPriceCents] = LockedPriceCentsLookup{
		ServiceItem: serviceItem,
	}

	return lookups
}

func GetDestinationForDistanceLookup(appCtx appcontext.AppContext, mtoShipment models.MTOShipment, mtoServiceItem *models.MTOServiceItem) (models.Address, error) {
	if mtoServiceItem == nil || mtoShipment.ShipmentType != models.MTOShipmentTypeHHG || (mtoServiceItem.ReService.Code != models.ReServiceCodeDLH && mtoServiceItem.ReService.Code != models.ReServiceCodeDSH && mtoServiceItem.ReService.Code != models.ReServiceCodeFSC) {
		return *mtoShipment.DestinationAddress, nil
	}
	shipmentCopy := mtoShipment
	err := appCtx.DB().Eager("DeliveryAddressUpdate.Status", "DeliveryAddressUpdate.UpdatedAt", "DeliveryAddressUpdate.OriginalAddress", "DeliveryAddressUpdate.NewAddress", "MTOServiceItems", "DestinationAddress").Find(&shipmentCopy, mtoShipment.ID)
	if err != nil {
		return models.Address{}, apperror.NewNotFoundError(shipmentCopy.ID, "MTOShipment not found in Destination For Distance Lookup")
	}
	if len(shipmentCopy.MTOServiceItems) == 0 {
		return *mtoShipment.DestinationAddress, nil
	}

	var result models.Address
	for _, si := range shipmentCopy.MTOServiceItems {
		siCopy := si
		err := appCtx.DB().Eager("ReService").Find(&siCopy, siCopy.ID)
		if err != nil {
			return models.Address{}, apperror.NewNotFoundError(siCopy.ID, "MTOServiceItem not found in Destination For Distance Lookup")
		}

		switch siCopy.ReService.Code {
		case models.ReServiceCodeDDASIT, models.ReServiceCodeDDDSIT, models.ReServiceCodeDDFSIT, models.ReServiceCodeDDSFSC:
			if shipmentCopy.DeliveryAddressUpdate != nil && shipmentCopy.DeliveryAddressUpdate.Status == models.ShipmentAddressUpdateStatusApproved {
				if shipmentCopy.DeliveryAddressUpdate.UpdatedAt.After(*siCopy.ApprovedAt) {
					return shipmentCopy.DeliveryAddressUpdate.OriginalAddress, nil
				}
				return shipmentCopy.DeliveryAddressUpdate.NewAddress, nil
			}
		}
	}
	if shipmentCopy.DeliveryAddressUpdate.Status == models.ShipmentAddressUpdateStatusApproved {
		result = shipmentCopy.DeliveryAddressUpdate.NewAddress
	} else {
		result = *shipmentCopy.DestinationAddress
	}
	return result, nil
}

// serviceItemNeedsParamKey wrapper for using paramCache.ServiceItemNeedsParamKey, if s.paramCache is nil
// we are not using the ParamCache and all lookups will be initialized and all param lookups will run their own
// database queries
func (s *ServiceItemParamKeyData) serviceItemNeedsParamKey(appCtx appcontext.AppContext, serviceItemCode models.ReServiceCode, paramKey models.ServiceItemParamName) (bool, error) {
	if s.paramCache == nil {
		// We used to turn some (but not nearly all) lookups on and off with a big switch here if the cache was not
		// on.  But that had a few issues.  First, it wasn't keeping up with the latest service to param mappings
		// (which are stored in the database and challenging to keep in sync here).  Second, it didn't appear to be
		// helping us a lot as it's only controlling whether the lookup goes in a map of lookups (and the map already
		// has as many entries as we have lookups due to the NotImplementedLookup we set for all params by default).
		// Only the appropriate lookups are called (elsewhere) regardless of what happens here.  So, at least until
		// we rethink the cache, just allow all lookups to be set we don't have a cache.
		return true, nil
	}

	useKey, err := s.paramCache.ServiceItemNeedsParamKey(appCtx, serviceItemCode, paramKey)
	if err != nil {
		return false, fmt.Errorf("error with ParamKey: %s using ServiceItemNeedsParamKey() for ServiceItemCode %s: %w", paramKey, serviceItemCode, err)
	}
	return useKey, nil
}

// ServiceParamValue returns a service parameter value from a key
func (s *ServiceItemParamKeyData) ServiceParamValue(appCtx appcontext.AppContext, key models.ServiceItemParamName) (string, error) {
	// NOTE: turning off param cache for now since we have a bug (MB-9497) that will likely require rethinking
	// how we cache.  Also, the cache does not seem to be having the impact we first thought it might.

	// Check cache for lookup value
	// if s.paramCache != nil && s.mtoShipmentID != nil {
	// 	paramCacheValue := s.paramCache.ParamValue(*s.mtoShipmentID, key)
	// 	if paramCacheValue != nil {
	// 		return *paramCacheValue, nil
	// 	}
	// }

	if lookup, ok := s.lookups[key]; ok {
		value, err := lookup.lookup(appCtx, s)
		if err != nil {
			switch err.(type) {
			case apperror.EventError:
				return "", err
			}
			return "", fmt.Errorf(" failed ServiceParamValue %sLookup with error %w", key, err)
		}
		// Save param value to cache
		// NOTE: although cache is not being checked above, continuing to cache values so existing tests don't break.
		if s.paramCache != nil && s.mtoShipmentID != nil {
			s.paramCache.addParamValue(*s.mtoShipmentID, key, value)
		}
		return value, nil
	}
	return "", fmt.Errorf("  ServiceParamValue <%sLookup> does not exist for key: <%s>", key, key)
}

func getPickupAddressForService(serviceCode models.ReServiceCode, mtoShipment models.MTOShipment) (models.Address, error) {
	// Determine which address field we should be using for pickup based on the shipment type.
	var ptrPickupAddress *models.Address
	var addressType string
	switch mtoShipment.ShipmentType {
	case models.MTOShipmentTypeHHGOutOfNTSDom:
		addressType = "storage facility"
		if mtoShipment.StorageFacility != nil {
			ptrPickupAddress = &mtoShipment.StorageFacility.Address
		}
	default:
		addressType = "pickup"
		ptrPickupAddress = mtoShipment.PickupAddress
	}

	// Determine if that address is valid based on which service we're pricing.
	switch serviceCode {
	case models.ReServiceCodeDUPK:
		// Pickup address isn't needed
		return models.Address{}, nil
	default:
		if ptrPickupAddress == nil || ptrPickupAddress.ID == uuid.Nil {
			return models.Address{}, apperror.NewNotFoundError(uuid.Nil, fmt.Sprintf("looking for %s address", addressType))
		}
		return *ptrPickupAddress, nil
	}
}

func getDestinationAddressForService(appCtx appcontext.AppContext, serviceCode models.ReServiceCode, mtoShipment models.MTOShipment) (models.Address, error) {
	// Determine which address field we should be using for destination based on the shipment type.
	var ptrDestinationAddress *models.Address
	var addressType string
	switch mtoShipment.ShipmentType {
	case models.MTOShipmentTypeHHGIntoNTSDom:
		addressType = "storage facility"
		if mtoShipment.StorageFacility != nil {
			ptrDestinationAddress = &mtoShipment.StorageFacility.Address
		}
	case models.MTOShipmentTypeHHG:
		shipmentCopy := mtoShipment
		err := appCtx.DB().Eager("DeliveryAddressUpdate.OriginalAddress", "DeliveryAddressUpdate.NewAddress", "MTOServiceItems", "DestinationAddress").Find(&shipmentCopy, mtoShipment.ID)
		if err != nil {
			return models.Address{}, apperror.NewNotFoundError(shipmentCopy.ID, "MTOShipment not found in Destination Address For service")
		}
		for _, si := range shipmentCopy.MTOServiceItems {
			siCopy := si
			err := appCtx.DB().Eager("ReService").Find(&siCopy, siCopy.ID)
			if err != nil {
				return models.Address{}, apperror.NewNotFoundError(siCopy.ID, "MTOServiceItem not found in Destination Address For service")
			}

			switch siCopy.ReService.Code {
			case models.ReServiceCodeDDASIT, models.ReServiceCodeDDDSIT, models.ReServiceCodeDDFSIT, models.ReServiceCodeDDSFSC:
				if shipmentCopy.DeliveryAddressUpdate != nil && shipmentCopy.DeliveryAddressUpdate.Status == models.ShipmentAddressUpdateStatusApproved {
					if siCopy.ApprovedAt != nil && shipmentCopy.DeliveryAddressUpdate.UpdatedAt.After(*siCopy.ApprovedAt) {
						return shipmentCopy.DeliveryAddressUpdate.OriginalAddress, nil
					} else {
						return shipmentCopy.DeliveryAddressUpdate.NewAddress, nil
					}
				}
			}
		}
		if shipmentCopy.DeliveryAddressUpdate.Status == models.ShipmentAddressUpdateStatusApproved {
			return shipmentCopy.DeliveryAddressUpdate.NewAddress, nil
		} else if serviceCode == models.ReServiceCodeDPK || serviceCode == models.ReServiceCodeDNPK {
			return models.Address{}, nil
		} else {
			return *shipmentCopy.DestinationAddress, nil
		}
	default:
		addressType = "destination"
		ptrDestinationAddress = mtoShipment.DestinationAddress
	}

	// Determine if that address is valid based on which service we're pricing.
	switch serviceCode {
	case models.ReServiceCodeDPK, models.ReServiceCodeDNPK:
		// Destination address isn't needed
		return models.Address{}, nil
	default:
		if ptrDestinationAddress == nil || ptrDestinationAddress.ID == uuid.Nil {
			return models.Address{}, apperror.NewNotFoundError(uuid.Nil, fmt.Sprintf("looking for %s address", addressType))
		}
		return *ptrDestinationAddress, nil
	}
}

func fetchContractForMove(appCtx appcontext.AppContext, moveID uuid.UUID) (models.ReContract, error) {
	var move models.Move
	err := appCtx.DB().Find(&move, moveID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.ReContract{}, apperror.NewNotFoundError(moveID, "looking for Move")
		}
		return models.ReContract{}, err
	}

	if move.AvailableToPrimeAt == nil {
		return models.ReContract{}, apperror.NewConflictError(moveID, "unable to pick contract because move is not available to prime")
	}

	return FetchContract(appCtx, *move.AvailableToPrimeAt)
}

func FetchContract(appCtx appcontext.AppContext, date time.Time) (models.ReContract, error) {
	var contractYear models.ReContractYear
	err := appCtx.DB().EagerPreload("Contract").Where("? between start_date and end_date", date).
		First(&contractYear)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.ReContract{}, apperror.NewNotFoundError(uuid.Nil, fmt.Sprintf("no contract year found for %s", date.String()))
		}
		return models.ReContract{}, err
	}

	return contractYear.Contract, nil
}
