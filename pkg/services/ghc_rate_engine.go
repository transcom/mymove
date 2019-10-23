package services

import (
	"time"

	"github.com/transcom/mymove/pkg/unit"
)

// DomesticLinehaulPricer prices domestic linehaul for a GHC move
//go:generate mockery -name DomesticLinehaulPricer
type DomesticLinehaulPricer interface {
	PriceDomesticLinehaul(moveDate time.Time, distance unit.Miles, weight unit.Pound, serviceArea string) (unit.Cents, error)
}

// DomesticPerWeightPricer domestic prices: origin and destination service area, SIT day 1, SIT Addt'l days
// domestic other prices: pack, unpack, and sit p/d costs for a GHC move
//go:generate mockery -name DomesticPerWeightPricer
type DomesticPerWeightPricer interface {
	PricePerWeightService(moveDate time.Time, distance unit.Miles, weight unit.Pound, serviceArea string, isDomesticOtherService bool) (unit.Cents, error)
}
