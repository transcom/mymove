package rateengine

import "github.com/transcom/mymove/pkg/unit"

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

// Like min quantity pricer, but the provided rate is scaled by some factor and needs to be unscaled before application
type scaledRateMinimumQuantityPricer struct {
	scale int
	min   int
}

func newScaledRateMinimumQuantityPricer(scale, min int) scaledRateMinimumQuantityPricer {
	return scaledRateMinimumQuantityPricer{scale, min}
}

func (m scaledRateMinimumQuantityPricer) price(rate unit.Cents, q1 unit.BaseQuantity, discount *unit.DiscountRate) unit.Cents {
	// Need to divide the rate by the scale to get the actual rate
	rate = rate.Multiply(1 / m.scale)

	if qConv := q1.ToUnitFloat(); qConv < float64(m.min) {
		q1 = unit.BaseQuantityFromInt(m.min)
	}

	calculatedRate := rate.MultiplyFloat64(q1.ToUnitFloat())

	if discount != nil {
		calculatedRate = discount.Apply(calculatedRate)
	}

	return calculatedRate
}
