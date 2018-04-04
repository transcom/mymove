package rateengine

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// RateEngine encapsulates the TSP rate engine process
type RateEngine struct {
	db     *pop.Connection
	logger *zap.Logger
}


func (re *RateEngine) determineMileage(origin_zip string, destination_zip string) (mileage int, error) {
	// TODO (Rebecca): make a proper error
	err := "whoops"
	fmt.Print(origin_zip)
	fmt.Print(destination_zip)
	// TODO (Rebecca): Lookup origin_zip to destination_zip mileage
	mileage := 1000
	return mileage, err
}

func (re *RateEngine) determineCWT(weight int) (cwt int) {
	return weight/100
}

// Determine the Base Linehaul (BLH)
func (re *RateEngine) determineBaseLinehaul(mileage int, weight int) (base_linehaul_charge int, error) {
	// TODO (Rebecca): This will come from a fetch
	base_linehaul_charge := mileage * weight
	// TODO (Rebecca): make a proper error
	err := "whoops"
	return base_linehaul_charge, err
}

// Determine the Linehaul Factors (OLF and DLF)
func (re *RateEngine) determineLinehaulFactors(weight int, zip string) (linehaul_factor float64) {
	// TODO: Fetch origin service area code via origin_zip
	fmt.print(zip)
	service_area := 101
	// TODO: Fetch linehaul factor for origin
	fmt.print(service_area)
	linehaul_factor := 0.51
	// Calculate linehaul_factor for the trip distance
	return (weight/100) * linehaul_factor
}

// Determine Shorthaul (SH) Charge (ONLY applies if shipment moves 800 miles and less)
func (re *RateEngine) determineShorthaulCharge(mileage int, cwt int) (shorthaul_charge float64, error) {
	cwt_miles := mileage * cwt
	// TODO: shorthaul_charge will be a lookup
	shorthaul_charge := cwt_miles
	return shorthaul_charge
}

// Determine Linehaul Charge (LC) TOTAL
// Formula: LC= [BLH + OLF + DLF + SH] x InvdLH
func (re *RateEngine) determineLinehaulChargeTotal(origin_zip string, destination_zip string) (linehaul_charge float64, error) {
	mileage := determineMileage(origin_zip, destination_zip)
	// TODO: Where is weight coming from?
	weight := 2000
	cwt := determineCWT(weight)
	base_linehaul_charge := determineBaseLinehaul(mileage, weight)
	origin_linehaul_factor := determineLinehaulFactors(weight, origin_zip)
	destination_linehaul_factor := determineLinehaulFactors(weight, destination_zip)
	shorthaul_charge := determineShorthaulCharge(mileage, cwt)
	// TODO: Where is our discount coming from?
	discount := 0.41
	inverse_discount := 1.0-discount
	// TODO: Make real error
	err := 'Whoops'
	return ((base_linehaul_charge + origin_linehaul_factor + destination_linehaul_factor + shorthaul_charge) * inverse_discount, err)
}

// NewRateEngine creates a new RateEngine
func NewRateEngine(db *pop.Connection, logger *zap.Logger) *RateEngine {
	return &RateEngine{db: db, logger: logger}
}