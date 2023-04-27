package rateengine

import (
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// RateEngine encapsulates the TSP rate engine process
// Deprecated: This is part of the old pre-GHC rate engine.
type RateEngine struct {
	move models.Move
}

// CostComputation represents the results of a computation.
// Deprecated: This is part of the old pre-GHC rate engine.
type CostComputation struct {
	LinehaulCostComputation
	NonLinehaulCostComputation

	SITFee      unit.Cents
	SITMax      unit.Cents
	GCC         unit.Cents
	LHDiscount  unit.DiscountRate
	SITDiscount unit.DiscountRate
	Weight      unit.Pound
}

// CostDetail holds the costComputation and a bool that signifies if the calculation is the winning (lowest cost) computation
// Deprecated: This is part of the old pre-GHC rate engine.
type CostDetail struct {
	Cost      CostComputation
	IsWinning bool
}

// CostDetails is a map of CostDetail
// Deprecated: This is part of the old pre-GHC rate engine.
type CostDetails map[string]*CostDetail

// ComputePPMMoveCosts uses zip codes to make two calculations for the price of a PPM move - once with the pickup zip and once with the current duty location zip - and returns both calcs.
// Deprecated: This is part of the old pre-GHC rate engine.
func (re *RateEngine) ComputePPMMoveCosts(appCtx appcontext.AppContext, weight unit.Pound, originPickupZip5 string, originDutyLocationZip5 string, destinationZip5 string, distanceMilesFromOriginPickupZip int, distanceMilesFromOriginDutyLocationZip int, date time.Time, daysInSit int) (costDetails CostDetails, err error) {
	errDeprecated := errors.New("ComputePPMMoveCosts function is deprecated")
	appCtx.Logger().Error("Invoking deprecated function", zap.Error(errDeprecated))
	return nil, errDeprecated
}

// GetWinningCostMove returns a costComputation of the winning calculation
// Deprecated: This is part of the old pre-GHC rate engine.
func GetWinningCostMove(costDetails CostDetails) CostComputation {
	if costDetails["pickupLocation"].IsWinning {
		return costDetails["pickupLocation"].Cost
	}
	return costDetails["originDutyLocation"].Cost
}

// GetNonWinningCostMove returns a costComputation of the non-winning calculation
// Deprecated: This is part of the old pre-GHC rate engine.
func GetNonWinningCostMove(costDetails CostDetails) CostComputation {
	if costDetails["pickupLocation"].IsWinning {
		return costDetails["originDutyLocation"].Cost
	}
	return costDetails["pickupLocation"].Cost
}

// NewRateEngine creates a new RateEngine
// Deprecated: This is part of the old pre-GHC rate engine.
func NewRateEngine(move models.Move) *RateEngine {
	return &RateEngine{move: move}
}
