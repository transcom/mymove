package services

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"

	"github.com/transcom/mymove/pkg/unit"
)

// ServiceItemPricer prices a generic service item for a GHC move
//go:generate mockery -name ServiceItemPricer
type ServiceItemPricer interface {
	PriceServiceItem(item models.PaymentServiceItem) (unit.Cents, error)
	// PriceServiceItemByID(paymentServiceItemID uuid.UUID) (unit.Cents, error)
}

// DomesticLinehaulPricer prices domestic linehaul for a GHC move
//go:generate mockery -name DomesticLinehaulPricer
type DomesticLinehaulPricer interface {
	PriceDomesticLinehaul(moveDate time.Time, distance unit.Miles, weight unit.Pound, serviceArea string) (unit.Cents, error)
}

// DomesticShorthaulPricer prices the domestic shorthaul for a GHC Move
//go:generate mockery -name DomesticShorthaulPricer
type DomesticShorthaulPricer interface {
	PriceDomesticShorthaul(moveDate time.Time, distance unit.Miles, weight unit.Pound, serviceArea string) (unit.Cents, error)
}

// DomesticServiceAreaPricer domestic prices: origin and destination service area, SIT day 1, SIT Addt'l days
//go:generate mockery -name DomesticServiceAreaPricer
type DomesticServiceAreaPricer interface {
	PriceDomesticServiceArea(moveDate time.Time, weight unit.Pound, serviceArea string, servicesCode string) (unit.Cents, error)
}

//DomesticFuelSurchargePricer prices fuel surcharge for domestic GHC moves
//go:generate mockery -name DomesticFuelSurchargePricer
type DomesticFuelSurchargePricer interface {
	PriceDomesticFuelSurcharge(moveDate time.Time, planner route.Planner, weight unit.Pound, source string, destination string) (unit.Cents, error)
}
