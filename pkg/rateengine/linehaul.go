package rateengine

import (
	"errors"
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/models"
)

func (re *RateEngine) determineMileage(originZip int, destinationZip int) (mileage int, err error) {
	// TODO (Rebecca): make a proper error
	fmt.Print(originZip)
	fmt.Print(destinationZip)
	// TODO (Rebecca): Lookup originZip to destinationZip mileage
	mileage = 1000
	if mileage != 1000 {
		err = errors.New("Oops determineMileage")
	} else {
		err = nil
	}
	return mileage, err
}

// Determine the Base Linehaul (BLH)
func (re *RateEngine) baseLinehaul(mileage int, cwt int) (baseLinehaulChargeCents int, err error) {
	// TODO (Rebecca): This will come from a fetch
	baseLinehaulChargeCents = mileage * cwt
	// TODO (Rebecca): make a proper error
	if baseLinehaulChargeCents == 0 {
		err = errors.New("Oops determineBaseLinehaul")
	} else {
		err = nil
	}
	return baseLinehaulChargeCents, err
}

// Determine the Linehaul Factors (OLF and DLF)
func (re *RateEngine) linehaulFactors(cwt int, zip3 int, date time.Time) (linehaulFactorCents int, err error) {
	serviceArea, err := models.FetchTariff400ngServiceAreaForZip3(re.db, zip3)
	if err != nil {
		return 0.0, err
	}
	linehaulFactorCents, err = models.FetchTariff400ngLinehaulFactor(re.db, serviceArea.ServiceArea, date)
	if err != nil {
		return 0.0, err
	}
	return cwt * linehaulFactorCents, nil
}

// Determine Shorthaul (SH) Charge (ONLY applies if shipment moves 800 miles and less)
func (re *RateEngine) shorthaulCharge(mileage int, cwt int) (shorthaulChargeCents int, err error) {
	if mileage >= 800 {
		return 0, nil
	}

	cwtMiles := mileage * cwt
	// TODO: shorthaulChargeCents will be a lookup
	shorthaulChargeCents = cwtMiles
	if shorthaulChargeCents == 0 {
		err = errors.New("Oops determineShorthaulCharge")
	} else {
		err = nil
	}
	return shorthaulChargeCents, err
}

// Determine Linehaul Charge (LC) TOTAL
// Formula: LC= [BLH + OLF + DLF + {SH}] x InvdLH
func (re *RateEngine) linehaulChargeTotal(originZip int, destinationZip int, date time.Time) (linehaulChargeCents int, err error) {
	mileage, err := re.determineMileage(originZip, destinationZip)
	// TODO: Where is weight coming from?
	weight := 2000
	cwt := re.determineCWT(weight)
	baseLinehaulChargeCents, err := re.baseLinehaul(mileage, cwt)
	originLinehaulFactorCents, err := re.linehaulFactors(cwt, originZip, date)
	destinationLinehaulFactorCents, err := re.linehaulFactors(cwt, destinationZip, date)
	shorthaulChargeCents, err := re.shorthaulCharge(mileage, cwt)
	if err != nil {
		return 0, err
	}
	return int(baseLinehaulChargeCents + originLinehaulFactorCents + destinationLinehaulFactorCents + shorthaulChargeCents), err
}
