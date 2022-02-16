package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
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

// ServiceItemParamKeyLookup does lookup on service item parameter keys
type ServiceItemParamKeyLookup interface {
	lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error)
}

// ServiceParamLookupInitialize initializes service parameter lookup
func ServiceParamLookupInitialize(
	appCtx appcontext.AppContext,
	planner route.Planner,
	mtoServiceItemID uuid.UUID,
	paymentRequestID uuid.UUID,
	moveTaskOrderID uuid.UUID,
	paramCache *ServiceParamsCache,
) (*ServiceItemParamKeyData, error) {

	// Get the MTOServiceItem
	var mtoServiceItem models.MTOServiceItem
	err := appCtx.DB().Eager("ReService").Find(&mtoServiceItem, mtoServiceItemID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(mtoServiceItemID, "looking for MTOServiceItem")
		default:
			return nil, apperror.NewQueryError("MTOServiceItem", err, "")
		}
	}

	s := ServiceItemParamKeyData{
		planner:          planner,
		lookups:          make(map[models.ServiceItemParamName]ServiceItemParamKeyLookup),
		MTOServiceItemID: mtoServiceItemID,
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
		ContractCode: ghcrateengine.DefaultContractCode,
	}

	//
	// Query and save PickupAddress & DestinationAddress upfront
	// s.serviceItemNeedsParamKey() could be used to check if the PickupAddress or DestinationAddress
	// can be used but it depends on the paramCache being set (not nil). It is possible to set the
	// paramCache to nil, especially during unit test, so not using that function for this part.
	//

	// Load data that is only used by a few service items
	var sitDestinationFinalAddress models.Address
	var serviceItemDimensions models.MTOServiceItemDimensions

	switch mtoServiceItem.ReService.Code {
	case models.ReServiceCodeDCRT, models.ReServiceCodeDUCRT, models.ReServiceCodeDCRTSA:
		err = appCtx.DB().Load(&mtoServiceItem, "Dimensions")
		if err != nil {
			return nil, err
		}
		serviceItemDimensions = mtoServiceItem.Dimensions
	case models.ReServiceCodeDDASIT, models.ReServiceCodeDDDSIT, models.ReServiceCodeDDFSIT:
		// load destination address from final address on service item
		if mtoServiceItem.SITDestinationFinalAddressID != nil && *mtoServiceItem.SITDestinationFinalAddressID != uuid.Nil {
			err = appCtx.DB().Load(&mtoServiceItem, "SITDestinationFinalAddress")
			if err != nil {
				return nil, err
			}
			sitDestinationFinalAddress = *mtoServiceItem.SITDestinationFinalAddress
		}
	}

	// Load shipment fields for service items that need them
	var mtoShipment models.MTOShipment
	var pickupAddress models.Address
	var destinationAddress models.Address

	if mtoServiceItem.ReService.Code != models.ReServiceCodeCS && mtoServiceItem.ReService.Code != models.ReServiceCodeMS {
		// Make sure there's an MTOShipment since that's nullable
		if mtoServiceItem.MTOShipmentID == nil {
			return nil, apperror.NewNotFoundError(uuid.Nil, "looking for MTOShipment")
		}
		err = appCtx.DB().Eager("PickupAddress", "DestinationAddress", "StorageFacility").Find(&mtoShipment, mtoServiceItem.MTOShipmentID)
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
			err = appCtx.DB().Load(mtoShipment.StorageFacility, "Address")
			if err != nil {
				return nil, apperror.NewQueryError("Address", err, "")
			}
		}

		pickupAddress, err = getPickupAddressForService(mtoServiceItem.ReService.Code, mtoShipment)
		if err != nil {
			return nil, err
		}

		destinationAddress, err = getDestinationAddressForService(mtoServiceItem.ReService.Code, mtoShipment)
		if err != nil {
			return nil, err
		}
	}

	switch mtoServiceItem.ReService.Code {
	case models.ReServiceCodeDDASIT, models.ReServiceCodeDDDSIT, models.ReServiceCodeDDFSIT, models.ReServiceCodeDOASIT, models.ReServiceCodeDOPSIT, models.ReServiceCodeDOFSIT:
		err = appCtx.DB().Load(&mtoShipment, "SITExtensions")
		if err != nil {
			return nil, err
		}
	}

	//
	// Set all lookup functions to "NOT IMPLEMENTED"
	//

	for _, key := range models.ValidServiceItemParamNames {
		s.lookups[key] = NotImplementedLookup{}
	}

	//
	// Begin setting lookup functions if they are needed for the given ReServiceCode
	//

	var paramKey models.ServiceItemParamName

	// ReService code for current MTO Service Item
	serviceItemCode := mtoServiceItem.ReService.Code

	paramKey = models.ServiceItemParamNameActualPickupDate
	err = s.setLookup(appCtx, serviceItemCode, paramKey, ActualPickupDateLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameRequestedPickupDate
	err = s.setLookup(appCtx, serviceItemCode, paramKey, RequestedPickupDateLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameReferenceDate
	err = s.setLookup(appCtx, serviceItemCode, paramKey, ReferenceDateLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameDistanceZip5
	err = s.setLookup(appCtx, serviceItemCode, paramKey, DistanceZip5Lookup{
		PickupAddress:      pickupAddress,
		DestinationAddress: destinationAddress,
	})
	if err != nil {
		return nil, err
	}
	paramKey = models.ServiceItemParamNameDistanceZip3
	err = s.setLookup(appCtx, serviceItemCode, paramKey, DistanceZip3Lookup{
		PickupAddress:      pickupAddress,
		DestinationAddress: destinationAddress,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier
	err = s.setLookup(appCtx, serviceItemCode, paramKey, FSCWeightBasedDistanceMultiplierLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameWeightAdjusted
	err = s.setLookup(appCtx, serviceItemCode, paramKey, WeightAdjustedLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameWeightBilled
	err = s.setLookup(appCtx, serviceItemCode, paramKey, WeightBilledLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameWeightEstimated
	err = s.setLookup(appCtx, serviceItemCode, paramKey, WeightEstimatedLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameWeightOriginal
	err = s.setLookup(appCtx, serviceItemCode, paramKey, WeightOriginalLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameWeightReweigh
	err = s.setLookup(appCtx, serviceItemCode, paramKey, WeightReweighLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameZipPickupAddress
	err = s.setLookup(appCtx, serviceItemCode, paramKey, ZipAddressLookup{
		Address: pickupAddress,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameZipDestAddress
	err = s.setLookup(appCtx, serviceItemCode, paramKey, ZipAddressLookup{
		Address: destinationAddress,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameMTOAvailableToPrimeAt
	err = s.setLookup(appCtx, serviceItemCode, paramKey, MTOAvailableToPrimeAtLookup{})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameServiceAreaOrigin
	err = s.setLookup(appCtx, serviceItemCode, paramKey, ServiceAreaLookup{
		Address: pickupAddress,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameServiceAreaDest
	err = s.setLookup(appCtx, serviceItemCode, paramKey, ServiceAreaLookup{
		Address: destinationAddress,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameContractCode
	err = s.setLookup(appCtx, serviceItemCode, paramKey, ContractCodeLookup{})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameCubicFeetBilled
	err = s.setLookup(appCtx, serviceItemCode, paramKey, CubicFeetBilledLookup{
		Dimensions: serviceItemDimensions,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNamePSILinehaulDom
	err = s.setLookup(appCtx, serviceItemCode, paramKey, PSILinehaulDomLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNamePSILinehaulDomPrice
	err = s.setLookup(appCtx, serviceItemCode, paramKey, PSILinehaulDomPriceLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameEIAFuelPrice
	err = s.setLookup(appCtx, serviceItemCode, paramKey, EIAFuelPriceLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameServicesScheduleOrigin
	err = s.setLookup(appCtx, serviceItemCode, paramKey, ServicesScheduleLookup{
		Address: pickupAddress,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameServicesScheduleDest
	err = s.setLookup(appCtx, serviceItemCode, paramKey, ServicesScheduleLookup{
		Address: destinationAddress,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameSITScheduleOrigin
	err = s.setLookup(appCtx, serviceItemCode, paramKey, SITScheduleLookup{
		Address: pickupAddress,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameSITScheduleDest
	err = s.setLookup(appCtx, serviceItemCode, paramKey, SITScheduleLookup{
		Address: destinationAddress,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameNumberDaysSIT
	err = s.setLookup(appCtx, serviceItemCode, paramKey, NumberDaysSITLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameZipSITDestHHGFinalAddress
	err = s.setLookup(appCtx, serviceItemCode, paramKey, ZipAddressLookup{
		Address: sitDestinationFinalAddress,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameZipSITOriginHHGOriginalAddress
	err = s.setLookup(appCtx, serviceItemCode, paramKey, ZipSITOriginHHGOriginalAddressLookup{
		ServiceItem: mtoServiceItem,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameZipSITOriginHHGActualAddress
	err = s.setLookup(appCtx, serviceItemCode, paramKey, ZipSITOriginHHGActualAddressLookup{
		ServiceItem: mtoServiceItem,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameDistanceZipSITDest
	err = s.setLookup(appCtx, serviceItemCode, paramKey, DistanceZipSITDestLookup{
		DestinationAddress:      destinationAddress,
		FinalDestinationAddress: sitDestinationFinalAddress,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameDistanceZipSITOrigin
	err = s.setLookup(appCtx, serviceItemCode, paramKey, DistanceZipSITOriginLookup{
		ServiceItem: mtoServiceItem,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameCubicFeetCrating
	err = s.setLookup(appCtx, serviceItemCode, paramKey, CubicFeetCratingLookup{
		Dimensions: serviceItemDimensions,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameDimensionHeight
	err = s.setLookup(appCtx, serviceItemCode, paramKey, DimensionHeightLookup{
		Dimensions: serviceItemDimensions,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameDimensionLength
	err = s.setLookup(appCtx, serviceItemCode, paramKey, DimensionLengthLookup{
		Dimensions: serviceItemDimensions,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameDimensionWidth
	err = s.setLookup(appCtx, serviceItemCode, paramKey, DimensionWidthLookup{
		Dimensions: serviceItemDimensions,
	})
	if err != nil {
		return nil, err
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

// serviceItemNeedsParamKey wrapper for using paramCache.ServiceItemNeedsParamKey, if s.paramCache is nil
// we are not using the ParamCache and all lookups will be initialized and all param lookups will run their own
// database queries
func (s *ServiceItemParamKeyData) serviceItemNeedsParamKey(appCtx appcontext.AppContext, serviceItemCode models.ReServiceCode, paramKey models.ServiceItemParamName) (bool, error) {
	if s.paramCache == nil {

		/*
				If we are presetting any lookups to maximize use vs having many queries or to make the lookup functions
			    more DRY. Then the values that have been identified as needing presets need to be checked here
			   	if the paramCache is nil.

			   	These are the fields which are preset and should be called out if it is needed by each service item. The default is
				to return true if the paramCache is nil. These checks will return false if the field is not used by service item:
					- Address
					- PickupAddress
					- DestinationAddress
		*/
		switch paramKey {
		case models.ServiceItemParamNameDistanceZip5, models.ServiceItemParamNameDistanceZip3:
			switch serviceItemCode {
			case models.ReServiceCodeDPK, models.ReServiceCodeDNPK, models.ReServiceCodeDUPK:
				return false, nil
			}
		case models.ServiceItemParamNameZipPickupAddress:
			switch serviceItemCode {
			case models.ReServiceCodeDUPK:
				return false, nil
			}
		case models.ServiceItemParamNameZipDestAddress:
			switch serviceItemCode {
			case models.ReServiceCodeDPK, models.ReServiceCodeDNPK:
				return false, nil
			}
		case models.ServiceItemParamNameServiceAreaOrigin:
			switch serviceItemCode {
			case models.ReServiceCodeDUPK:
				return false, nil
			}
		case models.ServiceItemParamNameServiceAreaDest:
			switch serviceItemCode {
			case models.ReServiceCodeDPK, models.ReServiceCodeDNPK:
				return false, nil
			}
		case models.ServiceItemParamNameServicesScheduleOrigin:
			switch serviceItemCode {
			case models.ReServiceCodeDUPK:
				return false, nil
			}
		case models.ServiceItemParamNameServicesScheduleDest:
			switch serviceItemCode {
			case models.ReServiceCodeDPK, models.ReServiceCodeDNPK:
				return false, nil
			}
		}
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

	// Check cache for lookup value
	if s.paramCache != nil && s.mtoShipmentID != nil {
		paramCacheValue := s.paramCache.ParamValue(*s.mtoShipmentID, key)
		if paramCacheValue != nil {
			return *paramCacheValue, nil
		}
	}

	if lookup, ok := s.lookups[key]; ok {
		value, err := lookup.lookup(appCtx, s)
		if err != nil {
			return "", fmt.Errorf(" failed ServiceParamValue %sLookup with error %w", key, err)
		}
		// Save param value to cache
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

func getDestinationAddressForService(serviceCode models.ReServiceCode, mtoShipment models.MTOShipment) (models.Address, error) {
	// Determine which address field we should be using for destination based on the shipment type.
	var ptrDestinationAddress *models.Address
	var addressType string
	switch mtoShipment.ShipmentType {
	case models.MTOShipmentTypeHHGIntoNTSDom:
		addressType = "storage facility"
		if mtoShipment.StorageFacility != nil {
			ptrDestinationAddress = &mtoShipment.StorageFacility.Address
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
