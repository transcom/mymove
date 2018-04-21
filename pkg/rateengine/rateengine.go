package rateengine

import (
	"time"

	"github.com/gobuffalo/pop"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/unit"
)

// RateEngine encapsulates the TSP rate engine process
type RateEngine struct {
	db      *pop.Connection
	logger  *zap.Logger
	planner route.Planner
}

// zip5ToZip3 takes two ZIP5 strings and returns the ZIP3 representation of them.
func (re *RateEngine) zip5ToZip3(originZip5 string, destinationZip5 string) (originZip3 string, destinationZip3 string) {
	originZip3 = originZip5[0:3]
	destinationZip3 = destinationZip5[0:3]
	return originZip3, destinationZip3
}

func (re *RateEngine) computePPM(weight unit.Pound, originZip5 string, destinationZip5 string,
	date time.Time, daysInSIT int, inverseDiscount float64) (unit.Cents, error) {

	originZip3, destinationZip3 := re.zip5ToZip3(originZip5, destinationZip5)

	// Linehaul charges
	mileage, err := re.determineMileage(originZip5, destinationZip5)
	if err != nil {
		re.logger.Error("Failed to determine mileage", zap.Error(err))
		return 0, err
	}
	baseLinehaulChargeCents, err := re.baseLinehaul(mileage, weight, date)
	if err != nil {
		re.logger.Error("Failed to determine base linehaul charge", zap.Error(err))
		return 0, err
	}
	originLinehaulFactorCents, err := re.linehaulFactors(weight.ToCWT(), originZip3, date)
	if err != nil {
		re.logger.Error("Failed to determine origin linehaul factor", zap.Error(err))
		return 0, err
	}
	destinationLinehaulFactorCents, err := re.linehaulFactors(weight.ToCWT(), destinationZip3, date)
	if err != nil {
		re.logger.Error("Failed to determine destination linehaul factor", zap.Error(err))
		return 0, err
	}
	shorthaulChargeCents, err := re.shorthaulCharge(mileage, weight.ToCWT(), date)
	if err != nil {
		re.logger.Error("Failed to determine shorthaul charge", zap.Error(err))
		return 0, err
	}

	// Non linehaul charges
	originServiceFee, err := re.serviceFeeCents(weight.ToCWT(), originZip3, date)
	if err != nil {
		re.logger.Error("Failed to determine origin service fee", zap.Error(err))
		return 0, err
	}
	destinationServiceFee, err := re.serviceFeeCents(weight.ToCWT(), destinationZip3, date)
	if err != nil {
		re.logger.Error("Failed to determine destination service fee", zap.Error(err))
		return 0, err
	}
	pack, err := re.fullPackCents(weight.ToCWT(), originZip3, date)
	if err != nil {
		re.logger.Error("Failed to determine full pack cost", zap.Error(err))
		return 0, err
	}
	unpack, err := re.fullUnpackCents(weight.ToCWT(), destinationZip3, date)
	if err != nil {
		re.logger.Error("Failed to determine full unpack cost", zap.Error(err))
		return 0, err
	}
	sit, err := re.sitCharge(weight.ToCWT(), daysInSIT, destinationZip3, date, true)
	if err != nil {
		return 0, err
	}

	ppmSubtotal := baseLinehaulChargeCents + originLinehaulFactorCents + destinationLinehaulFactorCents +
		shorthaulChargeCents + originServiceFee + destinationServiceFee + pack + unpack + sit

	gcc := ppmSubtotal.MultiplyFloat64(inverseDiscount)

	// PPMs only pay 95% of the best value
	ppmPayback := gcc.MultiplyFloat64(.95)

	re.logger.Info("PPM compensation total calculated",
		zap.Int("PPM compensation total", ppmPayback.Int()),
		zap.Int("PPM subtotal", ppmSubtotal.Int()),
		zap.Float64("inverse discount", inverseDiscount),
		zap.Int("base linehaul", baseLinehaulChargeCents.Int()),
		zap.Int("origin lh factor", originLinehaulFactorCents.Int()),
		zap.Int("destination lh factor", destinationLinehaulFactorCents.Int()),
		zap.Int("shorthaul", shorthaulChargeCents.Int()),
		zap.Int("origin service fee", originServiceFee.Int()),
		zap.Int("destination service fee", destinationServiceFee.Int()),
		zap.Int("pack fee", pack.Int()),
		zap.Int("unpack fee", unpack.Int()),
		zap.Int("sit fee", sit.Int()),
	)

	return ppmPayback, nil
}

// NewRateEngine creates a new RateEngine
func NewRateEngine(db *pop.Connection, logger *zap.Logger, planner route.Planner) *RateEngine {
	return &RateEngine{db: db, logger: logger, planner: planner}
}
