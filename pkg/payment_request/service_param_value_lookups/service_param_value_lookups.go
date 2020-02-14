package serviceparamvaluelookups

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/primemessages"
)

type ServiceItemParamKeyData struct {
	lookups            map[string]ServiceItemParamKeyLookup
	PayloadServiceItem primemessages.ServiceItem
	MTOServiceItemID   uuid.UUID
	PaymentRequestID   uuid.UUID
	MoveTaskOrderID    uuid.UUID
}

type ServiceItemParamKeyLookup interface {
	lookup(keyData *ServiceItemParamKeyData) (string, error)
}

func ServiceParamLookupInitialize(
	mtoServiceItemID uuid.UUID,
	paymentRequestID uuid.UUID,
	moveTaskOrderID uuid.UUID,
) *ServiceItemParamKeyData {

	s := ServiceItemParamKeyData{
		lookups:          make(map[string]ServiceItemParamKeyLookup),
		MTOServiceItemID: mtoServiceItemID,
		PaymentRequestID: paymentRequestID,
		MoveTaskOrderID:  moveTaskOrderID,
	}

	s.lookups["RequestedPickupDate"] = RequestedPickupDateLookup{}
	s.lookups["WeightBilledActual"] = WeightBilledActualLookup{}
	s.lookups["WeightActual"] = WeightActualLookup{}
	s.lookups["WeightEstimated"] = WeightEstimatedLookup{}
	s.lookups["DistanceZip3"] = DistanceZip3Lookup{}
	s.lookups["ZipPickupAddress"] = ZipPickupAddressLookup{}
	s.lookups["ZipDestAddress"] = ZipDestAddressLookup{}
	s.lookups["ServiceAreaOrigin"] = ServiceAreaOriginLookup{}

	return &s
}

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
