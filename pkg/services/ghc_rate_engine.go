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

// DomesticShorthaulPricer prices the domestic shorthaul for a GHC Move
type DomesticShorthaulPricer interface {
	PriceDomesticShorthaul(moveDate time.Time, distance unit.Miles, weight unit.Pound, serviceArea string) (unit.Cents, error)
}

// DomesticServiceAreaPricer domestic prices: origin and destination service area, SIT day 1, SIT Addt'l days

//go:generate mockery -name DomesticPerWeightPricer
type DomesticServiceAreaPricer interface {
	PriceDomesticServiceArea(moveDate time.Time, weight unit.Pound, serviceArea string, servicesCode string) (unit.Cents, error)
}
