package rateengine

import (
	"github.com/transcom/mymove/pkg/unit"
)

// LinehaulCostComputation represents the results of a computation.
// Deprecated: This is part of the old pre-GHC rate engine.
type LinehaulCostComputation struct {
	BaseLinehaul              unit.Cents
	OriginLinehaulFactor      unit.Cents
	DestinationLinehaulFactor unit.Cents
	ShorthaulCharge           unit.Cents
	LinehaulChargeTotal       unit.Cents
	Mileage                   int
	FuelSurcharge             FeeAndRate
}
