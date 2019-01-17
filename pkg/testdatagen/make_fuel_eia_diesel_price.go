package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func getFuelDefaultPrices() []unit.Millicents {
	return []unit.Millicents{320700, 333800, 331300, 325200, 279200, 261700, 244200, 195700, 158100, 433800}
}

func getFuelDefaultBaselines() []int64 {
	return []int64{1, 4, 10, 7, 9, 20, 12, 23, 8, 3}
}

// getFuelDefaultDateRange returns the default start date and end year to use for
// creating fuel prices
func getFuelDefaultDateRange() (start time.Time, end time.Time) {
	now := time.Now()
	// Set the end date as 1 year out
	endDate := now.AddDate(1, 0, 0)
	// Create rates starting with oldestStartDate
	oldestYear := 2018
	oldestMonth := time.January
	oldestDayStart := 15
	oldestStartDate := time.Date(oldestYear, oldestMonth, oldestDayStart, 0, 0, 0, 0, time.UTC)

	return oldestStartDate, endDate
}

// MakeFuelEIADieselPrices creates a series of FuelEIADieselPrice records
// using a default range or the start and end dates from assertions
// also can set the fuel prices and baselines from the assertions
func MakeFuelEIADieselPrices(db *pop.Connection, assertions Assertions) {

	// Get the default range. The default will start from Jan 2018 until the
	// the present Year
	oldestStartDate, rateEndDate := getFuelDefaultDateRange()

	// Override the default date range with assertions
	if !assertions.FuelEIADieselPrice.RateStartDate.IsZero() {
		oldestStartDate = assertions.FuelEIADieselPrice.RateStartDate
		rateEndDate = oldestStartDate.AddDate(1, 0, 0)
	}

	if !assertions.FuelEIADieselPrice.RateEndDate.IsZero() {
		rateEndDate = assertions.FuelEIADieselPrice.RateEndDate
	}

	// If the start and end dates are not valid (i.e., end date is less than start date)
	// then correct the end date
	if oldestStartDate.After(rateEndDate) || oldestStartDate.Equal(rateEndDate) {
		rateEndDate = oldestStartDate.AddDate(1, 0, 0)
	}

	startYear := oldestStartDate.Year()
	endYear := rateEndDate.Year()

	// Pick a publish date that is early in the month of the oldestStartDate
	oldestPubDate := time.Date(oldestStartDate.Year(), oldestStartDate.Month(), 03, 0, 0, 0, 0, time.UTC)

	// Set up variables for creating records in for loop
	nextStartDate := oldestStartDate
	nextEndDate := nextStartDate.AddDate(0, 1, -1)
	nextPubDate := oldestPubDate

	// Get prices and baseline numbers
	pricesMillicents := getFuelDefaultPrices()
	baselineRates := getFuelDefaultBaselines()

	// Use assertions for prices and baseline numbers if assertions are set
	price := assertions.FuelEIADieselPrice.EIAPricePerGallonMillicents
	if price != unit.Millicents(0) {
		pricesMillicents = []unit.Millicents{price}
	}

	// TODO: rate of 0 is valid, but we're assuming that
	// TODO: we wouldn't pick 0 to put into an assertion
	rate := assertions.FuelEIADieselPrice.BaselineRate
	if rate != 0 {
		baselineRates = []int64{rate}
	}

	pricesLen := len(pricesMillicents)
	baselineRateLen := len(baselineRates)

	// Create records 12 months at a time starting from nextStartDate
	// The minimum amount of records generated will be 48 months past
	// nextStartDate
	recordCount := 0
	for y := startYear; y <= endYear; y++ {
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

// MakeDefaultFuelEIADieselPrices creates a single FuelEIADieselPrice record with default values for a given shipmentDate
func MakeDefaultFuelEIADieselPrices(db *pop.Connection) {
	MakeFuelEIADieselPrices(db, Assertions{})
}
