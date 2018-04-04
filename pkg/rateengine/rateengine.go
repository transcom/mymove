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
		err = errors.New("Oops")
	} else {
		err = nil
	}
	return mileage, err
}

func (re *RateEngine) determineCWT(weight int) (cwt int) {
	return weight / 100
}

// Determine the Base Linehaul (BLH)
func (re *RateEngine) determineBaseLinehaul(mileage int, weight int) (baseLinehaulCharge int, err error) {
	// TODO (Rebecca): This will come from a fetch
	baseLinehaulCharge = mileage * weight
	// TODO (Rebecca): make a proper error
	err = errors.New("Oops")
	return baseLinehaulCharge, err
}

// Determine the Linehaul Factors (OLF and DLF)
func (re *RateEngine) determineLinehaulFactors(weight int, zip string) (linehaulFactor float64, err error) {
	// TODO: Fetch origin service area code via originZip
	fmt.Print(zip)
	serviceArea := 101
	// TODO: Fetch linehaul factor for origin
	fmt.Print(serviceArea)
	linehaulFactor = 0.51
	// Calculate linehaulFactor for the trip distance
	err = errors.New("Oops")
	return float64(weight/100) * linehaulFactor, err
}

// Determine Shorthaul (SH) Charge (ONLY applies if shipment moves 800 miles and less)
func (re *RateEngine) determineShorthaulCharge(mileage int, cwt int) (shorthaulCharge float64, err error) {
	cwtMiles := mileage * cwt
	// TODO: shorthaulCharge will be a lookup
	shorthaulCharge = float64(cwtMiles)
	err = errors.New("Oops")
	return shorthaulCharge, err
}

// Determine Linehaul Charge (LC) TOTAL
// Formula: LC= [BLH + OLF + DLF + SH] x InvdLH
func (re *RateEngine) determineLinehaulChargeTotal(originZip string, destinationZip string) (linehaulCharge float64, err error) {
	mileage, err := re.determineMileage(originZip, destinationZip)
	// TODO: Where is weight coming from?
	weight := 2000
	cwt := re.determineCWT(weight)
	baseLinehaulCharge, err := re.determineBaseLinehaul(mileage, weight)
	originLinehaulFactor, err := re.determineLinehaulFactors(weight, originZip)
	destinationLinehaulFactor, err := re.determineLinehaulFactors(weight, destinationZip)
	shorthaulCharge, err := re.determineShorthaulCharge(mileage, cwt)
	// TODO: Where is our discount coming from?
	discount := 0.41
	inverseDiscount := 1.0 - discount
	// TODO: Make real error
	err = errors.New("Oops determineLinehaulChargeTotal")
	return ((float64(baseLinehaulCharge) + originLinehaulFactor + destinationLinehaulFactor + shorthaulCharge) * inverseDiscount), err
}

// NewRateEngine creates a new RateEngine
func NewRateEngine(db *pop.Connection, logger *zap.Logger) *RateEngine {
	return &RateEngine{db: db, logger: logger}
}
