package rateengine

import (
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// LinehaulCostComputation represents the results of a computation.
type LinehaulCostComputation struct {
	BaseLinehaul              unit.Cents
	OriginLinehaulFactor      unit.Cents
	DestinationLinehaulFactor unit.Cents
	ShorthaulCharge           unit.Cents
	LinehaulChargeTotal       unit.Cents
	Mileage                   int
	FuelSurcharge             FeeAndRate
}

// Scale scales a cost computation by a multiplicative factor
func (c *LinehaulCostComputation) Scale(factor float64) {
	c.BaseLinehaul = c.BaseLinehaul.MultiplyFloat64(factor)
	c.OriginLinehaulFactor = c.OriginLinehaulFactor.MultiplyFloat64(factor)
	c.DestinationLinehaulFactor = c.DestinationLinehaulFactor.MultiplyFloat64(factor)
	c.ShorthaulCharge = c.ShorthaulCharge.MultiplyFloat64(factor)
	c.LinehaulChargeTotal = c.LinehaulChargeTotal.MultiplyFloat64(factor)
}

func (re *RateEngine) determineMileage(originZip5 string, destinationZip5 string) (mileage int, err error) {
	mileage, err = re.planner.Zip5TransitDistance(originZip5, destinationZip5)
	if err != nil {
		re.logger.Error("Failed to get distance from planner - %v", zap.Error(err),
			zap.String("origin_zip5", originZip5), zap.String("destination_zip5", destinationZip5))
	}
	return mileage, err
}

// Determine the Base Linehaul (BLH)
func (re *RateEngine) baseLinehaul(mileage int, weight unit.Pound, date time.Time) (baseLinehaulChargeCents unit.Cents, err error) {
	baseLinehaulChargeCents, err = models.FetchBaseLinehaulRate(re.db, mileage, weight, date)
	if err != nil {
		re.logger.Error("Base Linehaul query didn't complete: ", zap.Error(err))
	}

	return baseLinehaulChargeCents, err
}

// Determine the Linehaul Factors (OLF and DLF)
func (re *RateEngine) linehaulFactors(cwt unit.CWT, zip3 string, date time.Time) (linehaulFactorCents unit.Cents, err error) {
	serviceArea, err := models.FetchTariff400ngServiceAreaForZip3(re.db, zip3, date)
	if err != nil {
		return 0, err
	}
	return serviceArea.LinehaulFactor.Multiply(cwt.Int()), nil
}

// Determine Shorthaul (SH) Charge (ONLY applies if shipment moves 800 miles and less)
func (re *RateEngine) shorthaulCharge(mileage int, cwt unit.CWT, date time.Time) (shorthaulChargeCents unit.Cents, err error) {
	if mileage >= 800 {
		return 0, nil
	}
	re.logger.Debug("Shipment qualifies for shorthaul fee",
		zap.Int("miles", mileage))

	cwtMiles := mileage * cwt.Int()
	shorthaulChargeCents, err = models.FetchShorthaulRateCents(re.db, cwtMiles, date)

	return shorthaulChargeCents, err
}

// Determine Linehaul Charge (LC) TOTAL
// Formula: LC= [BLH + OLF + DLF + [SH]
func (re *RateEngine) linehaulChargeComputation(weight unit.Pound, originZip5 string, destinationZip5 string, date time.Time) (cost LinehaulCostComputation, err error) {
	cwt := weight.ToCWT()
	originZip3 := Zip5ToZip3(originZip5)
	destinationZip3 := Zip5ToZip3(destinationZip5)
	mileage, err := re.determineMileage(originZip5, destinationZip5)
	if err != nil {
		return cost, errors.Wrap(err, "Failed to determine mileage")
	}
	cost.Mileage = mileage

	cost.BaseLinehaul, err = re.baseLinehaul(mileage, weight, date)
	if err != nil {
		return cost, errors.Wrap(err, "Failed to determine base linehaul charge")
	}
	cost.OriginLinehaulFactor, err = re.linehaulFactors(cwt, originZip3, date)
	if err != nil {
		return cost, errors.Wrap(err, "Failed to determine origin linehaul factor")
	}
	cost.DestinationLinehaulFactor, err = re.linehaulFactors(cwt, destinationZip3, date)
	if err != nil {
		return cost, errors.Wrap(err, "Failed to determine destination linehaul factor")
	}
	cost.ShorthaulCharge, err = re.shorthaulCharge(mileage, cwt, date)
	if err != nil {
		return cost, errors.Wrap(err, "Failed to determine shorthaul charge")
	}
	cost.FuelSurcharge, err = re.fuelSurchargeComputation()
	if err != nil {
		return cost, errors.Wrap(err, "Failed to calculate fuel surcharge")
	}

	cost.LinehaulChargeTotal = cost.BaseLinehaul +
		cost.OriginLinehaulFactor +
		cost.DestinationLinehaulFactor +
		cost.ShorthaulCharge

	re.logger.Info("Linehaul charge total calculated",
		zap.Int("linehaul total", cost.LinehaulChargeTotal.Int()),
		zap.Int("linehaul", cost.BaseLinehaul.Int()),
		zap.Int("origin lh factor", cost.OriginLinehaulFactor.Int()),
		zap.Int("destination lh factor", cost.DestinationLinehaulFactor.Int()),
		zap.Int("shorthaul", cost.ShorthaulCharge.Int()),
	)

	return cost, err
}

// Calculate the fuel surcharge and return the result
// TODO: Fill this in with the actual formula for the fuel surcharge. Returning 0 as a stub/placeholder right now.
func (re *RateEngine) fuelSurchargeComputation() (fuelSurcharge FeeAndRate, err error) {
	return FeeAndRate{Fee: unit.Cents(0), Rate: unit.Millicents(0)}, err
}
