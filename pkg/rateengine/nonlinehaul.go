package rateengine

import (
	"github.com/transcom/mymove/pkg/unit"
)

// FeeAndRate holds the rate lookup and calculated fee (non-discounted)
// Deprecated: This is part of the old pre-GHC rate engine.
type FeeAndRate struct {
	Fee  unit.Cents
	Rate unit.Millicents
}

// NonLinehaulCostComputation represents the results of a computation.
// Deprecated: This is part of the old pre-GHC rate engine.
type NonLinehaulCostComputation struct {
	OriginService      FeeAndRate
	DestinationService FeeAndRate
	Pack               FeeAndRate
	Unpack             FeeAndRate
}
