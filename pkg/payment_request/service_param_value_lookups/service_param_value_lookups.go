package service_param_value_lookups

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/gen/primemessages"
)

type serviceItemParamKeyData struct {
	lookups map[string]ServiceItemParamKeyLookup
	PayloadServiceItem primemessages.ServiceItem
	MTOServiceItemID uuid.UUID
	PaymentRequestID uuid.UUID
	MoveTaskOrderID uuid.UUID
}

type ServiceItemParamKeyLookup interface {
	lookup(keyData *serviceItemParamKeyData) (string, error)
}

func ServiceParamLookupInitialize (
	payloadServiceItem primemessages.ServiceItem,
	mtoServiceItemID uuid.UUID,
	paymentRequestID uuid.UUID,
	moveTaskOrderID uuid.UUID,
	) *serviceItemParamKeyData {

	s := serviceItemParamKeyData{
		lookups: make(map[string]ServiceItemParamKeyLookup),
		PayloadServiceItem: payloadServiceItem,
		MTOServiceItemID: mtoServiceItemID,
		PaymentRequestID: paymentRequestID,
		MoveTaskOrderID: moveTaskOrderID,
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

func (s *serviceItemParamKeyData) ServiceParamValue (key string) (string, error) {
	lookup := s.lookups[key]
	value, err := lookup.lookup(s)
	if err != nil {
		return "", fmt.Errorf(" failed ServiceParamValue %sLookup with error %w", key, err)
	}
	return value, nil
}