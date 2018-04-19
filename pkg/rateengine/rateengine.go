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
	PPMPayback                unit.Cents
	PPMSubtotal               unit.Cents
	InverseDiscount           float64
	BaseLinehaul              unit.Cents
	OriginLinehaulFactor      unit.Cents
	DestinationLinehaulFactor unit.Cents
	ShorthaulCharge           unit.Cents
	OriginServiceFee          unit.Cents
	DestinationServiceFee     unit.Cents
	PackFee                   unit.Cents
	UnpackFee                 unit.Cents
}

// MarshalLogObject allows CostComputation to be logged by Zap.
func (c CostComputation) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt("PPMPayback", c.PPMPayback.Int())
	encoder.AddInt("PPMSubtotal", c.PPMSubtotal.Int())
	encoder.AddFloat64("InverseDiscount", c.InverseDiscount)
	encoder.AddInt("BaseLinehaul", c.BaseLinehaul.Int())
	encoder.AddInt("OriginLinehaulFactor", c.OriginLinehaulFactor.Int())
	encoder.AddInt("DestinationLinehaulFactor", c.DestinationLinehaulFactor.Int())
	encoder.AddInt("ShorthaulCharge", c.ShorthaulCharge.Int())
	encoder.AddInt("OriginServiceFee", c.OriginServiceFee.Int())
	encoder.AddInt("DestinationServiceFee", c.DestinationServiceFee.Int())
	encoder.AddInt("PackFee", c.PackFee.Int())
	encoder.AddInt("UnpackFee", c.UnpackFee.Int())

	return nil
}

// zip5ToZip3 takes two ZIP5 strings and returns the ZIP3 representation of them.
func (re *RateEngine) zip5ToZip3(originZip5 string, destinationZip5 string) (originZip3 string, destinationZip3 string) {
	originZip3 = originZip5[0:3]
	destinationZip3 = destinationZip5[0:3]
	return originZip3, destinationZip3
}

func (re *RateEngine) computePPM(weight unit.Pound, originZip5 string, destinationZip5 string, date time.Time, inverseDiscount float64) (cost CostComputation, err error) {
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
	// Non linehaul charges
	cost.OriginServiceFee, err = re.serviceFeeCents(weight.ToCWT(), originZip3)
	if err != nil {
		re.logger.Error("Failed to determine origin service fee", zap.Error(err))
		return
	}
	cost.DestinationServiceFee, err = re.serviceFeeCents(weight.ToCWT(), destinationZip3)
	if err != nil {
		re.logger.Error("Failed to determine destination service fee", zap.Error(err))
		return
	}
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

	cost.PPMSubtotal = cost.BaseLinehaul + cost.OriginLinehaulFactor + cost.DestinationLinehaulFactor +
		cost.ShorthaulCharge + cost.OriginServiceFee + cost.DestinationServiceFee + cost.PackFee + cost.UnpackFee

	ppmBestValue := cost.PPMSubtotal.MultiplyFloat64(inverseDiscount)

	// PPMs only pay 95% of the best value
	cost.PPMPayback = ppmBestValue.MultiplyFloat64(.95)

	re.logger.Info("PPM cost computation", zap.Object("computation", cost))

	return cost, nil
}

// NewRateEngine creates a new RateEngine
func NewRateEngine(db *pop.Connection, logger *zap.Logger, planner route.Planner) *RateEngine {
	return &RateEngine{db: db, logger: logger, planner: planner}
}
