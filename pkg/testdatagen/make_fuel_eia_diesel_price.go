package testdatagen

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
	"time"
)

func getFuelDefaultPrices() []unit.Millicents {
	return []unit.Millicents{320700, 333800, 331300, 325200, 279200, 261700, 244200, 195700, 158100, 433800}
}

func getFuelDefaultBaselines() []int64 {
	return []int64{320700, 333800, 331300, 325200, 279200, 261700, 244200, 195700, 158100, 433800}
}

// MakeFuelEIADieselPrices creates a series of FuelEIADieselPrice records
func MakeFuelEIADieselPrices(db *pop.Connection) {
	now := time.Now()
	// Create rates for 4 months past today's date
	endYear, _, _ := now.AddDate(0, 4, 0).Date()
	// Create rates starting with oldestStartDate
	oldestYear := 2018
	oldestMonth := time.October
	oldestDayStart := 15
	oldestStartDate := time.Date(oldestYear, oldestMonth, oldestDayStart, 0, 0, 0, 0, time.UTC)

	// Pick a publish date that is the first Monday of the month of the oldestStartDate
	// Notes: This isn't important for testing
	oldestPubDate := time.Date(2018, 10, 01, 0, 0, 0, 0, time.UTC)

	nextStartDate := oldestStartDate
	nextEndDate := nextStartDate.AddDate(0, 1, -1)
	nextPubDate := oldestPubDate

	pricesMillicents := getFuelDefaultPrices()
	pricesLen := len(pricesMillicents)

	baselineRates := getFuelDefaultBaselines()
	baselineRateLen := len(baselineRates)

	recordCount := 0
	for y := oldestYear; y <= endYear; y++ {
		for i := 1; i <= 12; i++ {
			id := uuid.Must(uuid.NewV4())
			fuelPrice := models.FuelEIADieselPrice{
				ID:                          id,
				PubDate:                     nextPubDate,
				RateStartDate:               nextStartDate,
				RateEndDate:                 nextEndDate,
				EIAPricePerGallonMillicents: pricesMillicents[recordCount%pricesLen],
				BaselineRate:                baselineRates[recordCount%baselineRateLen],
			}

			nextPubDate = nextPubDate.AddDate(0, 0, 28)
			nextStartDate = nextEndDate.AddDate(0, 0, 1)
			nextEndDate = nextStartDate.AddDate(0, 1, -1)
			mustCreate(db, &fuelPrice)

			recordCount++
		}
	}
}

// MakeFuelEIADieselPriceForDate creates a single FuelEIADieselPrice record for the date provided
// Note: this function is only good for generating one (1) record and should not be called multiple times
// from the same function, unless care is taken to NOT use overlapping rate periods. If multiple rate periods
// (FuelEIADieselPrice) records are needed use function MakeFuelEIADieselPrices or some variation of it
// spanning a range of dates.
func MakeFuelEIADieselPriceForDate(db *pop.Connection, shipmentDate time.Time, assertions Assertions) models.FuelEIADieselPrice {

	// rate period: october 15 - november 14
	// shipment: october 12 --> september 15 - october 14
	// shipment: october 17 --> october 15 - november 14
	// shipment: october 20 --> october 15 - november 14

	// Setup default dates, prices, and baselines
	_, _, shipDay := shipmentDate.Date()
	var rateStartMonth time.Month
	var rateStartYear int
	rateStartDay := 15
	if shipDay < 15 {
		// get new start year and month
		rateStartYear, rateStartMonth, _ = shipmentDate.AddDate(0, -1, 0).Date()
	}

	// Check for assertions
	id := assertions.FuelEIADieselPrice.ID
	if isZeroUUID(id) {
		id = uuid.Must(uuid.NewV4())
	}

	rateStartDate := time.Date(rateStartYear, rateStartMonth, rateStartDay, 0, 0, 0, 0, time.UTC)
	if !assertions.FuelEIADieselPrice.RateStartDate.IsZero() {
		rateStartDate = assertions.FuelEIADieselPrice.RateStartDate
	}
	rateEndDate := rateStartDate.AddDate(0, 1, -1)

	if !assertions.FuelEIADieselPrice.RateEndDate.IsZero() {
		rateEndDate = assertions.FuelEIADieselPrice.RateEndDate
	}

	pubDate := assertions.FuelEIADieselPrice.PubDate
	if pubDate.IsZero() {
		pubDate = time.Date(rateStartYear, rateStartMonth, 3, 0, 0, 0, 0, time.UTC)
	}

	price := assertions.FuelEIADieselPrice.EIAPricePerGallonMillicents
	if price == unit.Millicents(0) {
		price = getFuelDefaultPrices()[0]
	}

	rate := assertions.FuelEIADieselPrice.BaselineRate
	if rate == 0 {
		rate = getFuelDefaultBaselines()[0]
	}

	// Create new FuelEIADieselPrice based on shipment date
	fuelPrice := models.FuelEIADieselPrice{
		ID:                          id,
		PubDate:                     pubDate,
		RateStartDate:               rateStartDate,
		RateEndDate:                 rateEndDate,
		EIAPricePerGallonMillicents: price,
		BaselineRate:                rate,
	}

	// Overwrite values with those from assertions
	mergeModels(&fuelPrice, assertions.FuelEIADieselPrice)

	mustCreate(db, &fuelPrice)

	return fuelPrice
}

// MakeDefaultFuelEIADieselPriceForDate creates a single FuelEIADieselPrice record with default values for a given shipmentDate
func MakeDefaultFuelEIADieselPriceForDate(db *pop.Connection, shipmentDate time.Time) models.FuelEIADieselPrice {
	return MakeFuelEIADieselPriceForDate(db, shipmentDate, Assertions{})
}
