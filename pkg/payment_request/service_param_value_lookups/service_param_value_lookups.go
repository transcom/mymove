package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"

	"github.com/transcom/mymove/pkg/models"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"

	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
)

// ServiceItemParamKeyData contains service item parameter keys
type ServiceItemParamKeyData struct {
	db               *pop.Connection
	planner          route.Planner
	lookups          map[string]ServiceItemParamKeyLookup
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
	lookup(keyData *ServiceItemParamKeyData) (string, error)
}

// ServiceParamLookupInitialize initializes service parameter lookup
func ServiceParamLookupInitialize(
	db *pop.Connection,
	planner route.Planner,
	mtoServiceItemID uuid.UUID,
	paymentRequestID uuid.UUID,
	moveTaskOrderID uuid.UUID,
	paramCache *ServiceParamsCache,
) (*ServiceItemParamKeyData, error) {

	// Get the MTOServiceItem
	var mtoServiceItem models.MTOServiceItem
	err := db.Eager("ReService").Find(&mtoServiceItem, mtoServiceItemID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, services.NewNotFoundError(mtoServiceItemID, "looking for MTOServiceItemID")
		default:
			return nil, err
		}
	}

	var mtoShipment models.MTOShipment
	var pickupAddress models.Address
	var destinationAddress models.Address
	switch mtoServiceItem.ReService.Code {
	case models.ReServiceCodeCS, models.ReServiceCodeMS:
		// Do nothing, these service items don't use the MTOShipment
	default:
		// Make sure there's an MTOShipment since that's nullable
		if mtoServiceItem.MTOShipmentID == nil {
			return nil, services.NewNotFoundError(uuid.Nil, "looking for MTOShipmentID")
		}
		err := db.Eager("PickupAddress", "DestinationAddress").Find(&mtoShipment, mtoServiceItem.MTOShipmentID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return nil, services.NewNotFoundError(mtoServiceItemID, "looking for MTOServiceItemID")
			default:
				return nil, err
			}
		}

		if mtoServiceItem.ReService.Code != models.ReServiceCodeDUPK {
			if mtoShipment.PickupAddressID == nil {
				return nil, services.NewNotFoundError(uuid.Nil, "looking for PickupAddressID")
			}
			pickupAddress = *mtoShipment.PickupAddress
		}

		if mtoServiceItem.ReService.Code != models.ReServiceCodeDPK {
			if mtoShipment.DestinationAddressID == nil {
				return nil, services.NewNotFoundError(uuid.Nil, "looking for DestinationAddressID")
			}
			destinationAddress = *mtoShipment.DestinationAddress
		}
	}

	s := ServiceItemParamKeyData{
		db:               db,
		planner:          planner,
		lookups:          make(map[string]ServiceItemParamKeyLookup),
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

	for _, key := range models.ValidServiceItemParamNames {
		s.lookups[key] = NotImplementedLookup{}
	}

	// ReService code for current MTO Service Item
	serviceItemCode := mtoServiceItem.ReService.Code

	s.lookups[models.ServiceItemParamNameActualPickupDate.String()] = ActualPickupDateLookup{
		MTOShipment: mtoShipment,
	}

	paramKey := models.ServiceItemParamNameRequestedPickupDate
	useKey, err := s.ServiceItemNeedsParamKey(serviceItemCode, paramKey)
	if useKey && err == nil {
		s.lookups[paramKey.String()] = RequestedPickupDateLookup{
			MTOShipment: mtoShipment,
		}
	} else if err != nil {
		// TODO fmt and return error
		// return fmt.Errorf("error with ParamKey: %s using ServiceItemNeedsParamKey(): %w", paramKey, err)
	}

	paramKey = models.ServiceItemParamNameDistanceZip5
	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, paramKey)
	if useKey && err == nil {

	} else if err != nil {
		// TODO fmt and return error
		// return fmt.Errorf("error with ParamKey: %s using ServiceItemNeedsParamKey(): %w", paramKey, err)
	}

	paramKey = models.ServiceItemParamNameDistanceZip3
	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, paramKey)
	if useKey && err == nil {

	} else if err != nil {
		// TODO fmt and return error
		// return fmt.Errorf("error with ParamKey: %s using ServiceItemNeedsParamKey(): %w", paramKey, err)
	}

	if mtoServiceItem.ReService.Code != models.ReServiceCodeDPK && mtoServiceItem.ReService.Code != models.ReServiceCodeDUPK {
		s.lookups[models.ServiceItemParamNameDistanceZip5.String()] = DistanceZip5Lookup{
			PickupAddress:      pickupAddress,
			DestinationAddress: destinationAddress,
		}
		s.lookups[models.ServiceItemParamNameDistanceZip3.String()] = DistanceZip3Lookup{
			PickupAddress:      pickupAddress,
			DestinationAddress: destinationAddress,
		}
	}

	s.lookups[models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier.String()] = FSCWeightBasedDistanceMultiplierLookup{

		MTOShipment: mtoShipment,
	}

	paramKey = models.ServiceItemParamNameWeightBilledActual
	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, paramKey)
	if useKey && err == nil {
		s.lookups[paramKey.String()] = WeightBilledActualLookup{
			MTOShipment: mtoShipment,
		}
	}

	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNameWeightEstimated)
	if useKey && err == nil {
		s.lookups[models.ServiceItemParamNameWeightEstimated.String()] = WeightEstimatedLookup{
			MTOShipment: mtoShipment,
		}
	}

	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNameWeightActual)
	if useKey && err == nil {
		s.lookups[models.ServiceItemParamNameWeightActual.String()] = WeightActualLookup{
			MTOShipment: mtoShipment,
		}
	}

	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNameZipPickupAddress)
	if useKey && err == nil {
		if mtoServiceItem.ReService.Code != models.ReServiceCodeDUPK {
			s.lookups[models.ServiceItemParamNameZipPickupAddress.String()] = ZipAddressLookup{
				Address: pickupAddress,
			}
		}
	}

	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNameZipDestAddress)
	if useKey && err == nil {
		if mtoServiceItem.ReService.Code != models.ReServiceCodeDPK {
			s.lookups[models.ServiceItemParamNameZipDestAddress.String()] = ZipAddressLookup{
				Address: destinationAddress,
			}
		}
	}

	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNameMTOAvailableToPrimeAt)
	if useKey && err == nil {
		s.lookups[models.ServiceItemParamNameMTOAvailableToPrimeAt.String()] = MTOAvailableToPrimeAtLookup{}
	}

	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNameServiceAreaOrigin)
	if useKey && err == nil {
		if mtoServiceItem.ReService.Code != models.ReServiceCodeDUPK {
			s.lookups[models.ServiceItemParamNameServiceAreaOrigin.String()] = ServiceAreaLookup{
				Address: pickupAddress,
			}
		}
	}

	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNameServiceAreaDest)
	if useKey && err == nil {
		if mtoServiceItem.ReService.Code != models.ReServiceCodeDPK {
			s.lookups[models.ServiceItemParamNameServiceAreaDest.String()] = ServiceAreaLookup{
				Address: destinationAddress,
			}
		}
	}

	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNameContractCode)
	if useKey && err == nil {
		s.lookups[models.ServiceItemParamNameContractCode.String()] = ContractCodeLookup{}
	}

	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNamePSILinehaulDom)
	if useKey && err == nil {
		s.lookups[models.ServiceItemParamNamePSILinehaulDom.String()] = PSILinehaulDomLookup{
			MTOShipment: mtoShipment,
		}
	}

	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNamePSILinehaulDomPrice)
	if useKey && err == nil {
		s.lookups[models.ServiceItemParamNamePSILinehaulDomPrice.String()] = PSILinehaulDomPriceLookup{
			MTOShipment: mtoShipment,
		}
	}
	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNameEIAFuelPrice)
	if useKey && err == nil {
		s.lookups[models.ServiceItemParamNameEIAFuelPrice.String()] = EIAFuelPriceLookup{
			MTOShipment: mtoShipment,
		}
	}

	if mtoServiceItem.ReService.Code != models.ReServiceCodeDUPK {
		s.lookups[models.ServiceItemParamNameServicesScheduleOrigin.String()] = ServicesScheduleLookup{
			Address: pickupAddress,
		}
	}


	if mtoServiceItem.ReService.Code != models.ReServiceCodeDPK {
		s.lookups[models.ServiceItemParamNameServicesScheduleDest.String()] = ServicesScheduleLookup{
			Address: destinationAddress,
		}
	}


	return &s, nil
}

// ServiceItemNeedsParamKey wrapper for using paramCache.ServiceItemNeedsParamKey, if s.paramCache is nil
// we are not using the ParamCache and all lookups be initialized and all will run their own
// database queries
func (s *ServiceItemParamKeyData) ServiceItemNeedsParamKey(serviceItemCode models.ReServiceCode, paramKey models.ServiceItemParamName) (bool, error) {
	if s.paramCache == nil {
		return true, nil
	}

	return s.paramCache.ServiceItemNeedsParamKey(serviceItemCode, paramKey)
}

// ServiceParamValue returns a service parameter value from a key
func (s *ServiceItemParamKeyData) ServiceParamValue(key string) (string, error) {

	// Check cache for lookup value
	if s.paramCache != nil && s.mtoShipmentID != nil {
		paramCacheValue := s.paramCache.ParamValue(*s.mtoShipmentID, key)
		if paramCacheValue != nil {
			return *paramCacheValue, nil
		}
	}

	if lookup, ok := s.lookups[key]; ok {
		value, err := lookup.lookup(s)
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
