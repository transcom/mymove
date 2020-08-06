package serviceparamvaluelookups

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

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
) *ServiceItemParamKeyData {

	s := ServiceItemParamKeyData{
		db:               db,
		planner:          planner,
		lookups:          make(map[string]ServiceItemParamKeyLookup),
		MTOServiceItemID: mtoServiceItemID,
		PaymentRequestID: paymentRequestID,
		MoveTaskOrderID:  moveTaskOrderID,
	}

	for _, key := range models.ValidServiceItemParamNames {
		s.lookups[key] = NotImplementedLookup{}
	}

	s.lookups[models.ServiceItemParamNameActualPickupDate.String()] = ActualPickupDateLookup{}
	s.lookups[models.ServiceItemParamNameContractCode.String()] = ContractCodeLookup{}
	s.lookups[models.ServiceItemParamNameDistanceZip3.String()] = DistanceZip3Lookup{}
	s.lookups[models.ServiceItemParamNameDistanceZip5.String()] = DistanceZip5Lookup{}
	s.lookups[models.ServiceItemParamNameEIAFuelPrice.String()] = EIAFuelPriceLookup{}
	s.lookups[models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier.String()] = FSCWeightBasedDistanceMultiplierLookup{}
	s.lookups[models.ServiceItemParamNameMTOAvailableToPrimeAt.String()] = MTOAvailableToPrimeAtLookup{}
	s.lookups[models.ServiceItemParamNameRequestedPickupDate.String()] = RequestedPickupDateLookup{}
	s.lookups[models.ServiceItemParamNameServiceAreaDest.String()] = ServiceAreaDestLookup{}
	s.lookups[models.ServiceItemParamNameServiceAreaOrigin.String()] = ServiceAreaOriginLookup{}
	s.lookups[models.ServiceItemParamNameServicesScheduleDest.String()] = ServicesScheduleDestLookup{}
	s.lookups[models.ServiceItemParamNameServicesScheduleOrigin.String()] = ServicesScheduleOriginLookup{}
	s.lookups[models.ServiceItemParamNameWeightActual.String()] = WeightActualLookup{}
	s.lookups[models.ServiceItemParamNameWeightBilledActual.String()] = WeightBilledActualLookup{}
	s.lookups[models.ServiceItemParamNameWeightEstimated.String()] = WeightEstimatedLookup{}
	s.lookups[models.ServiceItemParamNameZipDestAddress.String()] = ZipDestAddressLookup{}
	s.lookups[models.ServiceItemParamNameZipPickupAddress.String()] = ZipPickupAddressLookup{}

	return &s
}

// ServiceParamValue returns a service parameter value from a key
func (s *ServiceItemParamKeyData) ServiceParamValue(key string) (string, error) {
	if lookup, ok := s.lookups[key]; ok {
		value, err := lookup.lookup(s)
		if err != nil {
			return "", fmt.Errorf(" failed ServiceParamValue %sLookup with error %w", key, err)
		}
		return value, nil
	}
	return "", fmt.Errorf("  ServiceParamValue <%sLookup> does not exist for key: <%s>", key, key)
}
