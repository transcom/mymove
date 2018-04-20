package rateengine

import (
	"time"

	"github.com/gobuffalo/pop"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/unit"
)

// RateEngine encapsulates the TSP rate engine process
type RateEngine struct {
	db      *pop.Connection
	logger  *zap.Logger
	planner route.Planner
}

// CostComputation represents the results of a computation.
type CostComputation struct {
	BaseLinehaul              unit.Cents
	OriginLinehaulFactor      unit.Cents
	DestinationLinehaulFactor unit.Cents
	ShorthaulCharge           unit.Cents
	LinehaulChargeTotal       unit.Cents

	OriginServiceFee      unit.Cents
	DestinationServiceFee unit.Cents
	PackFee               unit.Cents
	UnpackFee             unit.Cents
	FullPackUnpackFee     unit.Cents

	GCC unit.Cents
}

// MarshalLogObject allows CostComputation to be logged by Zap.
func (c CostComputation) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt("BaseLinehaul", c.BaseLinehaul.Int())
	encoder.AddInt("OriginLinehaulFactor", c.OriginLinehaulFactor.Int())
	encoder.AddInt("DestinationLinehaulFactor", c.DestinationLinehaulFactor.Int())
	encoder.AddInt("ShorthaulCharge", c.ShorthaulCharge.Int())
	encoder.AddInt("LinehaulChargeTotal", c.LinehaulChargeTotal.Int())

	encoder.AddInt("OriginServiceFee", c.OriginServiceFee.Int())
	encoder.AddInt("DestinationServiceFee", c.DestinationServiceFee.Int())
	encoder.AddInt("PackFee", c.PackFee.Int())
	encoder.AddInt("UnpackFee", c.UnpackFee.Int())
	encoder.AddInt("FullPackUnpackFee", c.FullPackUnpackFee.Int())

	encoder.AddInt("GCC", c.GCC.Int())

	return nil
}

// zip5ToZip3 takes two ZIP5 strings and returns the ZIP3 representation of them.
func (re *RateEngine) zip5ToZip3(originZip5 string, destinationZip5 string) (originZip3 string, destinationZip3 string) {
	originZip3 = originZip5[0:3]
	destinationZip3 = destinationZip5[0:3]
	return originZip3, destinationZip3
}

// ComputePPM Calculates the cost of a PPM move.
func (re *RateEngine) ComputePPM(weight unit.Pound, originZip5 string, destinationZip5 string, date time.Time, inverseDiscount float64) (cost CostComputation, err error) {
	originZip3, destinationZip3 := re.zip5ToZip3(originZip5, destinationZip5)

	// Linehaul charges
	mileage, err := re.determineMileage(originZip5, destinationZip5)
	if err != nil {
		re.logger.Error("Failed to determine mileage", zap.Error(err))
		return
	}

	cost.BaseLinehaul, err = re.baseLinehaul(mileage, weight, date)
	if err != nil {
		re.logger.Error("Failed to determine base linehaul charge", zap.Error(err))
		return
	}

	cost.OriginLinehaulFactor, err = re.linehaulFactors(weight.ToCWT(), originZip3, date)
	if err != nil {
		re.logger.Error("Failed to determine origin linehaul factor", zap.Error(err))
		return
	}

	cost.DestinationLinehaulFactor, err = re.linehaulFactors(weight.ToCWT(), destinationZip3, date)
	if err != nil {
		re.logger.Error("Failed to determine destination linehaul factor", zap.Error(err))
		return
	}

	cost.ShorthaulCharge, err = re.shorthaulCharge(mileage, weight.ToCWT(), date)
	if err != nil {
		re.logger.Error("Failed to determine shorthaul charge", zap.Error(err))
		return
	}

	linehaulChargeSubtotal := cost.BaseLinehaul + cost.OriginLinehaulFactor +
		cost.DestinationLinehaulFactor + cost.ShorthaulCharge

	cost.LinehaulChargeTotal = linehaulChargeSubtotal.MultiplyFloat64(inverseDiscount)

	// Non linehaul charges
	originServiceFee, err := re.serviceFeeCents(weight.ToCWT(), originZip3, date)
	if err != nil {
		re.logger.Error("Failed to determine origin service fee", zap.Error(err))
		return
	}
	cost.OriginServiceFee = originServiceFee.MultiplyFloat64(inverseDiscount)

	destinationServiceFee, err := re.serviceFeeCents(weight.ToCWT(), destinationZip3, date)
	if err != nil {
		re.logger.Error("Failed to determine destination service fee", zap.Error(err))
		return
	}
	cost.DestinationServiceFee = destinationServiceFee.MultiplyFloat64(inverseDiscount)

	cost.PackFee, err = re.fullPackCents(weight.ToCWT(), originZip3, date)
	if err != nil {
		re.logger.Error("Failed to determine full pack cost", zap.Error(err))
		return
	}
	cost.UnpackFee, err = re.fullUnpackCents(weight.ToCWT(), destinationZip3, date)
	if err != nil {
		re.logger.Error("Failed to determine full unpack cost", zap.Error(err))
		return
	}

	cost.FullPackUnpackFee = (cost.PackFee + cost.UnpackFee).MultiplyFloat64(inverseDiscount)

	cost.GCC = cost.LinehaulChargeTotal + cost.OriginServiceFee + cost.DestinationServiceFee +
		cost.FullPackUnpackFee

	re.logger.Info("PPM cost computation", zap.Object("cost", cost))

	return cost, nil
}

// NewRateEngine creates a new RateEngine
func NewRateEngine(db *pop.Connection, logger *zap.Logger, planner route.Planner) *RateEngine {
	return &RateEngine{db: db, logger: logger, planner: planner}
}
