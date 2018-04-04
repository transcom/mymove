package rateengine

import (
	"errors"
	"fmt"

	"github.com/gobuffalo/pop"
	"go.uber.org/zap"
)

// RateEngine encapsulates the TSP rate engine process
type RateEngine struct {
	db     *pop.Connection
	logger *zap.Logger
}

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

func (re *RateEngine) determineCWT(weight int) (cwt int) {
	return weight / 100
}

// Determine the Base Linehaul (BLH)
func (re *RateEngine) baseLinehaul(mileage int, weight int) (baseLinehaulCharge int, err error) {
	// TODO (Rebecca): This will come from a fetch
	baseLinehaulCharge = mileage * weight
	// TODO (Rebecca): make a proper error
	if baseLinehaulCharge == 0 {
		err = errors.New("Oops determineBaseLinehaul")
	} else {
		err = nil
	}
	return baseLinehaulCharge, err
}

// Determine the Linehaul Factors (OLF and DLF)
func (re *RateEngine) linehaulFactors(weight int, zip string) (linehaulFactor float64, err error) {
	// TODO: Fetch origin service area code via originZip
	fmt.Print(zip)
	serviceArea := 101
	// TODO: Fetch linehaul factor for origin
	fmt.Print(serviceArea)
	linehaulFactor = 0.51
	// Calculate linehaulFactor for the trip distance
	if linehaulFactor == 0 {
		err = errors.New("Oops determineLinehaulFactors")
	} else {
		err = nil
	}

	return float64(weight/100) * linehaulFactor, err
}

// Determine Shorthaul (SH) Charge (ONLY applies if shipment moves 800 miles and less)
func (re *RateEngine) shorthaulCharge(mileage int, cwt int) (shorthaulCharge float64, err error) {
	if mileage >= 800 {
		return 0, nil
	}

	cwtMiles := mileage * cwt
	// TODO: shorthaulCharge will be a lookup
	shorthaulCharge = float64(cwtMiles)
	if shorthaulCharge == 0 {
		err = errors.New("Oops determineShorthaulCharge")
	} else {
		err = nil
	}
	return shorthaulCharge, err
}

// Determine Linehaul Charge (LC) TOTAL
// Formula: LC= [BLH + OLF + DLF + SH] x InvdLH
func (re *RateEngine) linehaulChargeTotal(originZip string, destinationZip string) (linehaulCharge float64, err error) {
	mileage, err := re.determineMileage(originZip, destinationZip)
	// TODO: Where is weight coming from?
	weight := 2000
	cwt := re.determineCWT(weight)
	baseLinehaulCharge, err := re.baseLinehaul(mileage, weight)
	originLinehaulFactor, err := re.linehaulFactors(weight, originZip)
	destinationLinehaulFactor, err := re.linehaulFactors(weight, destinationZip)
	shorthaulCharge, err := re.shorthaulCharge(mileage, cwt)
	// TODO: Where is our discount coming from?
	discount := 0.41
	inverseDiscount := 1.0 - discount
	// TODO: Make real error
	if err != nil {
		err = errors.New("Oops determineLinehaulChargeTotal")
	}
	return ((float64(baseLinehaulCharge) + originLinehaulFactor + destinationLinehaulFactor + shorthaulCharge) * inverseDiscount), err
}

// NewRateEngine creates a new RateEngine
func NewRateEngine(db *pop.Connection, logger *zap.Logger) *RateEngine {
	return &RateEngine{db: db, logger: logger}
}
