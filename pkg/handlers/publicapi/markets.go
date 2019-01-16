package publicapi

import (
	"github.com/transcom/mymove/pkg/gen/apimessages"
)

func payloadForMarkets(market *string) *apimessages.ShipmentMarket {
	if market == nil {
		return nil
	}
	m := apimessages.ShipmentMarket(*market)
	return &m
}
