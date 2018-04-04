package rateengine

import (
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
	err = "whoops"
	fmt.Print(originZip)
	fmt.Print(destinationZip)
	// TODO (Rebecca): Lookup originZip to destinationZip mileage
	mileage = 1000
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
	err = "whoops"
	return baseLinehaulCharge, err
}

// Determine the Linehaul Factors (OLF and DLF)
func (re *RateEngine) determineLinehaulFactors(weight int, zip string) (linehaulFactor float64) {
	// TODO: Fetch origin service area code via originZip
	fmt.Print(zip)
	serviceArea := 101
	// TODO: Fetch linehaul factor for origin
	fmt.Print(serviceArea)
	linehaulFactor = 0.51
	// Calculate linehaulFactor for the trip distance
	return (weight / 100) * linehaulFactor
}

// Determine Shorthaul (SH) Charge (ONLY applies if shipment moves 800 miles and less)
func (re *RateEngine) determineShorthaulCharge(mileage int, cwt int) (shorthaulCharge float64, err error) {
	cwtMiles := mileage * cwt
	// TODO: shorthaulCharge will be a lookup
	shorthaulCharge = cwtMiles
	return shorthaulCharge
}

// Determine Linehaul Charge (LC) TOTAL
// Formula: LC= [BLH + OLF + DLF + SH] x InvdLH
func (re *RateEngine) determineLinehaulChargeTotal(originZip string, destinationZip string) (linehaulCharge float64, err error) {
	mileage := determineMileage(originZip, destinationZip)
	// TODO: Where is weight coming from?
	weight := 2000
	cwt := determineCWT(weight)
	baseLinehaulCharge := determineBaseLinehaul(mileage, weight)
	originLinehaulFactor := determineLinehaulFactors(weight, originZip)
	destinationLinehaulFactor := determineLinehaulFactors(weight, destinationZip)
	shorthaulCharge := determineShorthaulCharge(mileage, cwt)
	// TODO: Where is our discount coming from?
	discount := 0.41
	inverseDiscount := 1.0 - discount
	// TODO: Make real error
	err = "Whoops"
	return ((baseLinehaulCharge + originLinehaulFactor + destinationLinehaulFactor + shorthaulCharge) * inverseDiscount), err
}

// NewRateEngine creates a new RateEngine
func NewRateEngine(db *pop.Connection, logger *zap.Logger) *RateEngine {
	return &RateEngine{db: db, logger: logger}
}
