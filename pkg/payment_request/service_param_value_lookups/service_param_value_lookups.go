package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"

	"github.com/gobuffalo/pop/v5"
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

	s := ServiceItemParamKeyData{
		db:               db,
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
		err = db.Eager("PickupAddress", "DestinationAddress").Find(&mtoShipment, mtoServiceItem.MTOShipmentID)
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
	err = s.setLookup(serviceItemCode, paramKey, ActualPickupDateLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameRequestedPickupDate
	err = s.setLookup(serviceItemCode, paramKey, RequestedPickupDateLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameDistanceZip5
	err = s.setLookup(serviceItemCode, paramKey, DistanceZip5Lookup{
		PickupAddress:      pickupAddress,
		DestinationAddress: destinationAddress,
	})
	if err != nil {
		return nil, err
	}
	paramKey = models.ServiceItemParamNameDistanceZip3
	err = s.setLookup(serviceItemCode, paramKey, DistanceZip3Lookup{
		PickupAddress:      pickupAddress,
		DestinationAddress: destinationAddress,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier
	err = s.setLookup(serviceItemCode, paramKey, FSCWeightBasedDistanceMultiplierLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameWeightBilledActual
	err = s.setLookup(serviceItemCode, paramKey, WeightBilledActualLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameWeightEstimated
	err = s.setLookup(serviceItemCode, paramKey, WeightEstimatedLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameWeightActual
	err = s.setLookup(serviceItemCode, paramKey, WeightActualLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameZipPickupAddress
	err = s.setLookup(serviceItemCode, paramKey, ZipAddressLookup{
		Address: pickupAddress,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameZipDestAddress
	err = s.setLookup(serviceItemCode, paramKey, ZipAddressLookup{
		Address: destinationAddress,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameMTOAvailableToPrimeAt
	err = s.setLookup(serviceItemCode, paramKey, MTOAvailableToPrimeAtLookup{})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameServiceAreaOrigin
	err = s.setLookup(serviceItemCode, paramKey, ServiceAreaLookup{
		Address: pickupAddress,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameServiceAreaDest
	err = s.setLookup(serviceItemCode, paramKey, ServiceAreaLookup{
		Address: destinationAddress,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameContractCode
	err = s.setLookup(serviceItemCode, paramKey, ContractCodeLookup{})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNamePSILinehaulDom
	err = s.setLookup(serviceItemCode, paramKey, PSILinehaulDomLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNamePSILinehaulDomPrice
	err = s.setLookup(serviceItemCode, paramKey, PSILinehaulDomPriceLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameEIAFuelPrice
	err = s.setLookup(serviceItemCode, paramKey, EIAFuelPriceLookup{
		MTOShipment: mtoShipment,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameServicesScheduleOrigin
	err = s.setLookup(serviceItemCode, paramKey, ServicesScheduleLookup{
		Address: pickupAddress,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameServicesScheduleDest
	err = s.setLookup(serviceItemCode, paramKey, ServicesScheduleLookup{
		Address: destinationAddress,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameSITScheduleOrigin
	err = s.setLookup(serviceItemCode, paramKey, SITScheduleLookup{
		Address: pickupAddress,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameSITScheduleDest
	err = s.setLookup(serviceItemCode, paramKey, SITScheduleLookup{
		Address: destinationAddress,
	})
	if err != nil {
		return nil, err
	}

	paramKey = models.ServiceItemParamNameDistanceZipSITDest
	err = s.setLookup(serviceItemCode, paramKey, DistanceZipSITDestLookup{
		DestinationAddress: destinationAddress,
	})
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *ServiceItemParamKeyData) setLookup(serviceItemCode models.ReServiceCode, paramKey models.ServiceItemParamName, lookup ServiceItemParamKeyLookup) error {
	useKey, err := s.serviceItemNeedsParamKey(serviceItemCode, paramKey)
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
func (s *ServiceItemParamKeyData) serviceItemNeedsParamKey(serviceItemCode models.ReServiceCode, paramKey models.ServiceItemParamName) (bool, error) {
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
			case models.ReServiceCodeDPK, models.ReServiceCodeDUPK:
				return false, nil
			}
		case models.ServiceItemParamNameZipPickupAddress:
			switch serviceItemCode {
			case models.ReServiceCodeDUPK:
				return false, nil
			}
		case models.ServiceItemParamNameZipDestAddress:
			switch serviceItemCode {
			case models.ReServiceCodeDPK:
				return false, nil
			}
		case models.ServiceItemParamNameServiceAreaOrigin:
			switch serviceItemCode {
			case models.ReServiceCodeDUPK:
				return false, nil
			}
		case models.ServiceItemParamNameServiceAreaDest:
			switch serviceItemCode {
			case models.ReServiceCodeDPK:
				return false, nil
			}
		case models.ServiceItemParamNameServicesScheduleOrigin:
			switch serviceItemCode {
			case models.ReServiceCodeDUPK:
				return false, nil
			}
		case models.ServiceItemParamNameServicesScheduleDest:
			switch serviceItemCode {
			case models.ReServiceCodeDPK:
				return false, nil
			}
		}
		return true, nil
	}

	useKey, err := s.paramCache.ServiceItemNeedsParamKey(serviceItemCode, paramKey)
	if err != nil {
		return false, fmt.Errorf("error with ParamKey: %s using ServiceItemNeedsParamKey() for ServiceItemCode %s: %w", paramKey, serviceItemCode, err)
	}
	return useKey, nil
}

// ServiceParamValue returns a service parameter value from a key
func (s *ServiceItemParamKeyData) ServiceParamValue(key models.ServiceItemParamName) (string, error) {

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
