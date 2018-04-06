package rateengine

import (
	"errors"
	"fmt"
)

func (re *RateEngine) determineMileage(originZip string, destinationZip string) (mileage int, err error) {
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
func (re *RateEngine) linehaulFactors(cwt int, zip string) (linehaulFactorCents int, err error) {
	// TODO: Fetch origin service area code via originZip
	fmt.Print(zip)
	serviceArea := 101
	// TODO: Fetch linehaul factor for origin
	fmt.Print(serviceArea)
	// TODO: linehaul factors are in CENTS
	linehaulFactorCents = 51
	// Calculate linehaulFactorCents for the trip distance
	if linehaulFactorCents == 0 {
		err = errors.New("Oops determineLinehaulFactors")
	} else {
		err = nil
	}

	return cwt * linehaulFactorCents, err
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
func (re *RateEngine) linehaulChargeTotal(originZip string, destinationZip string) (linehaulChargeCents int, err error) {
	mileage, err := re.determineMileage(originZip, destinationZip)
	// TODO: Where is weight coming from?
	weight := 2000
	cwt := re.determineCWT(weight)
	baseLinehaulChargeCents, err := re.baseLinehaul(mileage, cwt)
	originLinehaulFactorCents, err := re.linehaulFactors(cwt, originZip)
	destinationLinehaulFactorCents, err := re.linehaulFactors(cwt, destinationZip)
	shorthaulChargeCents, err := re.shorthaulCharge(mileage, cwt)
	// TODO: Where is our discount coming from?
	discount := 0.41
	inverseDiscount := 1.0 - discount
	// TODO: Make real error
	if err != nil {
		err = errors.New("Oops determineLinehaulChargeTotal")
	}
	return int(float64(baseLinehaulChargeCents+originLinehaulFactorCents+destinationLinehaulFactorCents+shorthaulChargeCents) * inverseDiscount), err
}
