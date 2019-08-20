package rateengine

import (
	"github.com/transcom/mymove/pkg/unit"
)

type pricer interface {
	price(rate unit.Cents, q1 unit.BaseQuantity, discount *unit.DiscountRate) unit.Cents
}

// Basic pricer, multiplies the rate against the provided quantity
type basicQuantityPricer struct{}

func newBasicQuantityPricer() basicQuantityPricer {
	return basicQuantityPricer{}
}

func (m basicQuantityPricer) price(rate unit.Cents, q1 unit.BaseQuantity, discount *unit.DiscountRate) unit.Cents {
	calculatedRate := rate.MultiplyFloat64(q1.ToUnitFloat())

	if discount != nil {
		calculatedRate = discount.Apply(calculatedRate)
	}

	return calculatedRate
}

// Like the basic pricer, but enforces a minimum value for the quantity
type minimumQuantityPricer struct {
	min int
}

func newMinimumQuantityPricer(min int) minimumQuantityPricer {
	return minimumQuantityPricer{min}
}

func (m minimumQuantityPricer) price(rate unit.Cents, q1 unit.BaseQuantity, discount *unit.DiscountRate) unit.Cents {
	if qConv := q1.ToUnitFloat(); qConv < float64(m.min) {
		q1 = unit.BaseQuantityFromInt(m.min)
	}

	calculatedRate := rate.MultiplyFloat64(q1.ToUnitFloat())

	if discount != nil {
		calculatedRate = discount.Apply(calculatedRate)
	}

	return calculatedRate
}

// Line the min quantity pricer, but multiplies rate by quantity / 100
type minimumQuantityHundredweightPricer struct {
	min int
}

func newMinimumQuantityHundredweightPricer(min int) minimumQuantityHundredweightPricer {
	return minimumQuantityHundredweightPricer{min}
}

func (m minimumQuantityHundredweightPricer) price(rate unit.Cents, q1 unit.BaseQuantity, discount *unit.DiscountRate) unit.Cents {
	if qConv := q1.ToUnitFloat(); qConv < float64(m.min) {
		q1 = unit.BaseQuantityFromInt(m.min)
	}

	calculatedRate := rate.MultiplyFloat64(q1.ToUnitFloat() / 100.0)

	if discount != nil {
		calculatedRate = discount.Apply(calculatedRate)
	}

	return calculatedRate
}

// Ignores quantity, just returns rate with discount applied
type flatRatePricer struct{}

func newFlatRatePricer() flatRatePricer {
	return flatRatePricer{}
}

func (m flatRatePricer) price(rate unit.Cents, q1 unit.BaseQuantity, discount *unit.DiscountRate) unit.Cents {
	calculatedRate := rate

	if discount != nil {
		calculatedRate = discount.Apply(calculatedRate)
	}

	return calculatedRate
}
