package rateengine

import (
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

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

// MarshalLogObject allows LinehaulCostComputation to be logged by Zap.
func (c LinehaulCostComputation) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt("BaseLinehaul", c.BaseLinehaul.Int())
	encoder.AddInt("OriginLinehaulFactor", c.OriginLinehaulFactor.Int())
	encoder.AddInt("DestinationLinehaulFactor", c.DestinationLinehaulFactor.Int())
	encoder.AddInt("ShorthaulCharge", c.ShorthaulCharge.Int())

	return nil
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
func (re *RateEngine) linehaulChargeComputation(weight unit.Pound, originZip5 string, destinationZip5 string, pickupDate time.Time) (cost LinehaulCostComputation, err error) {
	cwt := weight.ToCWT()
	originZip3 := Zip5ToZip3(originZip5)
	destinationZip3 := Zip5ToZip3(destinationZip5)
	mileage, err := re.determineMileage(originZip5, destinationZip5)
	if err != nil {
		return cost, errors.Wrap(err, "Failed to determine mileage")
	}
	cost.Mileage = mileage

	cost.BaseLinehaul, err = re.baseLinehaul(mileage, weight, pickupDate)
	if err != nil {
		return cost, errors.Wrap(err, "Failed to determine base linehaul charge")
	}
	cost.OriginLinehaulFactor, err = re.linehaulFactors(cwt, originZip3, pickupDate)
	if err != nil {
		return cost, errors.Wrap(err, "Failed to determine origin linehaul factor")
	}
	cost.DestinationLinehaulFactor, err = re.linehaulFactors(cwt, destinationZip3, pickupDate)
	if err != nil {
		return cost, errors.Wrap(err, "Failed to determine destination linehaul factor")
	}
	cost.ShorthaulCharge, err = re.shorthaulCharge(mileage, cwt, pickupDate)
	if err != nil {
		return cost, errors.Wrap(err, "Failed to determine shorthaul charge")
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
func (re *RateEngine) fuelSurchargeComputation(totalLinehaulCost unit.Cents, bookDate time.Time) (fuelSurcharge FeeAndRate, err error) {
	fuelEIADieselPriceSlice := []models.FuelEIADieselPrice{}

	// Changing the format of the date to remove the time portion so it plays nicely with db
	bookDateString := bookDate.Format("2006-01-02")

	query := re.db.Where("rate_start_date <= ?", bookDateString).Where("rate_end_date >= ?", bookDateString)
	err1 := query.All(&fuelEIADieselPriceSlice)
	if err1 != nil {
		re.logger.Error(err1.Error())
	}

	// We expect to only retrieve one value from the FuelEIADieselPrice table. There should be only one valid date range for a given bookDate.
	// If we get more than one, something is wrong.
	if len(fuelEIADieselPriceSlice) > 1 {
		re.logger.Error("Got back multiple values from FuelEIADieselPrice when we should have only gotten one.")
	}

	if len(fuelEIADieselPriceSlice) == 0 {
		re.logger.Error("Query failed to find an applicable FuelEIADieselPrice")
	}

	fuelEIADieselPrice := fuelEIADieselPriceSlice[0]
	fuelSurchargePercentage := float64(fuelEIADieselPrice.BaselineRate) / 100
	fee := totalLinehaulCost.MultiplyFloat64(fuelSurchargePercentage)

	return FeeAndRate{Fee: unit.Cents(fee), Rate: fuelEIADieselPrice.EIAPricePerGallonMillicents}, err
}
