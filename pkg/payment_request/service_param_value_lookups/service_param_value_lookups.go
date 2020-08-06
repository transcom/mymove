package serviceparamvaluelookups

import (
	"fmt"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"

	"github.com/transcom/mymove/pkg/route"
)

// ServiceItemParamKeyData contains service item parameter keys
type ServiceItemParamKeyData struct {
	db               *pop.Connection
	planner          route.Planner
	lookups          map[string]ServiceItemParamKeyLookup
	MTOServiceItemID uuid.UUID
	PaymentRequestID uuid.UUID
	MoveTaskOrderID  uuid.UUID
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
) *ServiceItemParamKeyData {

	// Get the MTOServiceItem and associated MTOShipment
	var mtoServiceItem models.MTOServiceItem
	err := db.Eager("ReService").Find(&mtoServiceItem, mtoServiceItemID)
	if err != nil {
		// return fmt.Errorf("error in ServiceParamLookupInitialize() querying for MTOServiceItem using ID: %s with error: %w", mtoServiceItemID, err)
	}

	s := ServiceItemParamKeyData{
		db:               db,
		planner:          planner,
		lookups:          make(map[string]ServiceItemParamKeyLookup),
		MTOServiceItemID: mtoServiceItemID,
		PaymentRequestID: paymentRequestID,
		MoveTaskOrderID:  moveTaskOrderID,
		paramCache:       paramCache,
		mtoShipmentID:    mtoServiceItem.MTOShipmentID,
	}

	for _, key := range models.ValidServiceItemParamNames {
		s.lookups[key] = NotImplementedLookup{}
	}

	// ReService code for current MTO Service Item
	serviceItemCode := mtoServiceItem.ReService.Code

	paramKey := models.ServiceItemParamNameRequestedPickupDate
	useKey, err := s.ServiceItemNeedsParamKey(serviceItemCode, paramKey)
	if useKey && err == nil {
		s.lookups[paramKey.String()] = RequestedPickupDateLookup{}
	} else if err != nil {
		// TODO fmt and return error
		// return fmt.Errorf("error with ParamKey: %s using ServiceItemNeedsParamKey(): %w", paramKey, err)
	}

	paramKey = models.ServiceItemParamNameDistanceZip5
	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, paramKey)
	if useKey && err == nil {
		s.lookups[paramKey.String()] = DistanceZip5Lookup{}
	} else if err != nil {
		// TODO fmt and return error
		// return fmt.Errorf("error with ParamKey: %s using ServiceItemNeedsParamKey(): %w", paramKey, err)
	}

	paramKey = models.ServiceItemParamNameDistanceZip3
	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, paramKey)
	if useKey && err == nil {
		s.lookups[paramKey.String()] = DistanceZip3Lookup{}
	} else if err != nil {
		// TODO fmt and return error
		// return fmt.Errorf("error with ParamKey: %s using ServiceItemNeedsParamKey(): %w", paramKey, err)
	}

	paramKey = models.ServiceItemParamNameWeightBilledActual
	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, paramKey)
	if useKey && err == nil {
		s.lookups[paramKey.String()] = WeightBilledActualLookup{}
	}

	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNameWeightEstimated)
	if useKey && err == nil {
		s.lookups[models.ServiceItemParamNameWeightEstimated.String()] = WeightEstimatedLookup{}
	}

	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNameWeightActual)
	if useKey && err == nil {
		s.lookups[models.ServiceItemParamNameWeightActual.String()] = WeightActualLookup{}
	}

	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNameZipPickupAddress)
	if useKey && err == nil {
		s.lookups[models.ServiceItemParamNameZipPickupAddress.String()] = ZipPickupAddressLookup{}
	}

	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNameZipDestAddress)
	if useKey && err == nil {
		s.lookups[models.ServiceItemParamNameZipDestAddress.String()] = ZipDestAddressLookup{}
	}

	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNameMTOAvailableToPrimeAt)
	if useKey && err == nil {
		s.lookups[models.ServiceItemParamNameMTOAvailableToPrimeAt.String()] = MTOAvailableToPrimeAtLookup{}
	}

	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNameServiceAreaOrigin)
	if useKey && err == nil {
		s.lookups[models.ServiceItemParamNameServiceAreaOrigin.String()] = ServiceAreaOriginLookup{}
	}

	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNameServiceAreaDest)
	if useKey && err == nil {
		s.lookups[models.ServiceItemParamNameServiceAreaDest.String()] = ServiceAreaDestLookup{}
	}

	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNameContractCode)
	if useKey && err == nil {
		s.lookups[models.ServiceItemParamNameContractCode.String()] = ContractCodeLookup{}
	}

	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNamePSILinehaulDom)
	if useKey && err == nil {
		s.lookups[models.ServiceItemParamNamePSILinehaulDom.String()] = PSILinehaulDomLookup{}
	}

	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNamePSILinehaulDomPrice)
	if useKey && err == nil {
		s.lookups[models.ServiceItemParamNamePSILinehaulDomPrice.String()] = PSILinehaulDomPriceLookup{}
	}

	useKey, err = s.ServiceItemNeedsParamKey(serviceItemCode, models.ServiceItemParamNameEIAFuelPrice)
	if useKey && err == nil {
		s.lookups[models.ServiceItemParamNameEIAFuelPrice.String()] = EIAFuelPriceLookup{}
	}

	return &s
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
