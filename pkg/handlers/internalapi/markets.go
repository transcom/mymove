package internalapi

import (
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

func payloadForMarkets(market *string) *internalmessages.ShipmentMarket {
	if market == nil {
		return nil
	}
	m := internalmessages.ShipmentMarket(*market)
	return &m
}
