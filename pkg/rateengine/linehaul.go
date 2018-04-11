package rateengine

import (
	"errors"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
)

func (re *RateEngine) determineMileage(originZip int, destinationZip int) (mileage int, err error) {
	// TODO (Rebecca): Lookup originZip to destinationZip mileage using API of choice
	mileage = 1000
	if mileage != 1000 {
		err = errors.New("Oops determineMileage")
	} else {
		err = nil
	}
	return mileage, err
}

// Determine the Base Linehaul (BLH)
func (re *RateEngine) baseLinehaul(mileage int, cwt int, date time.Time) (baseLinehaulChargeCents int, err error) {
	baseLinehaulChargeCents, err = models.FetchBaseLinehaulRate(re.db, mileage, cwt, date)
	if err != nil {
		re.logger.Error("Base Linehaul query didn't complete: ", zap.Error(err))
	}

	return baseLinehaulChargeCents, err
}

// Determine the Linehaul Factors (OLF and DLF)
func (re *RateEngine) linehaulFactors(cwt int, zip3 int, date time.Time) (linehaulFactorCents int, err error) {
	serviceArea, err := models.FetchTariff400ngServiceAreaForZip3(re.db, zip3)
	if err != nil {
		return 0, err
	}
	linehaulFactorCents, err = models.FetchTariff400ngLinehaulFactor(re.db, serviceArea.ServiceArea, date)
	if err != nil {
		return 0, err
	}
	return cwt * linehaulFactorCents, nil
}

// Determine Shorthaul (SH) Charge (ONLY applies if shipment moves 800 miles and less)
func (re *RateEngine) shorthaulCharge(mileage int, cwt int, date time.Time) (shorthaulChargeCents int, err error) {
	if mileage >= 800 {
		return 0, nil
	}
	re.logger.Debug("Shipment qualifies for shorthaul fee",
		zap.Int("miles", mileage))

	cwtMiles := mileage * cwt
	shorthaulChargeCents, err = models.FetchShorthaulRateCents(re.db, cwtMiles, date)

	return shorthaulChargeCents, err
}

// Determine Linehaul Charge (LC) TOTAL
// Formula: LC= [BLH + OLF + DLF + [SH]
func (re *RateEngine) linehaulChargeTotal(weight int, originZip int, destinationZip int, date time.Time) (linehaulChargeCents int, err error) {
	mileage, err := re.determineMileage(originZip, destinationZip)
	cwt := re.determineCWT(weight)
	baseLinehaulChargeCents, err := re.baseLinehaul(mileage, cwt, date)
	if err != nil {
		return 0, err
	}
	originLinehaulFactorCents, err := re.linehaulFactors(cwt, originZip, date)
	if err != nil {
		return 0, err
	}
	destinationLinehaulFactorCents, err := re.linehaulFactors(cwt, destinationZip, date)
	if err != nil {
		return 0, err
	}
	shorthaulChargeCents, err := re.shorthaulCharge(mileage, cwt, date)
	if err != nil {
		return 0, err
	}

	linehaulChargeCents = baseLinehaulChargeCents + originLinehaulFactorCents + destinationLinehaulFactorCents + shorthaulChargeCents
	re.logger.Info("Linehaul charge total calculated",
		zap.Int("linehaul total", linehaulChargeCents),
		zap.Int("linehaul", baseLinehaulChargeCents),
		zap.Int("origin lh factor", originLinehaulFactorCents),
		zap.Int("destination lh factor", destinationLinehaulFactorCents),
		zap.Int("shorthaul", shorthaulChargeCents))

	return linehaulChargeCents, err
}
