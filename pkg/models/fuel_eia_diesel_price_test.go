package models_test

import (
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
	"time"
)

func (suite *ModelSuite) TestBasicFuelEIADieselPriceInstantiation() {
	now := time.Now()
	id := uuid.Must(uuid.NewV4())
	newFuelPrice := models.FuelEIADieselPrice{
		ID:                          id,
		PubDate:                     now,
		RateStartDate:               now,
		RateEndDate:                 now.AddDate(0, 1, -1),
		EIAPricePerGallonMillicents: 320700,
		BaselineRate:                6,
	}

	verrs, err := suite.db.ValidateAndCreate(&newFuelPrice)

	suite.Nil(err, "Error writing to the db.")
	suite.False(verrs.HasAny(), "Error validating model")
}

func (suite *ModelSuite) TestEmptyFuelEIADieselPriceInstantiation() {
	newFuelPrice := models.FuelEIADieselPrice{}

	expErrors := map[string][]string{
		"pub_date":                          {"PubDate can not be blank."},
		"rate_start_date":                   {"RateStartDate can not be blank."},
		"rate_end_date":                     {"RateEndDate can not be blank."},
		"e_i_a_price_per_gallon_millicents": {"0 is not greater than 0."},
		"baseline_rate":                     {"0 is not greater than 0."},
	}
	suite.verifyValidationErrors(&newFuelPrice, expErrors)
}

func (suite *ModelSuite) TestBadDatesFuelEIADieselPriceInstantiation() {
	now := time.Now()
	id := uuid.Must(uuid.NewV4())

	// Test for bad start and end dates within the same record
	newFuelPrice := models.FuelEIADieselPrice{
		ID:                          id,
		PubDate:                     now,
		RateStartDate:               now,
		RateEndDate:                 now.AddDate(0, -1, 0),
		EIAPricePerGallonMillicents: 320700,
		BaselineRate:                6,
	}

	expErrors := map[string][]string{
		"rate_end_date": {"RateEndDate must be after RateStartDate."},
	}

	suite.verifyValidationErrors(&newFuelPrice, expErrors)

	// Clear expected errors
	expErrors = map[string][]string{}

	id = uuid.Must(uuid.NewV4())
	// Test for overalapping start and end dates
	newFuelPrice = models.FuelEIADieselPrice{
		ID:                          id,
		PubDate:                     now,
		RateStartDate:               time.Date(2017, time.November, 15, 0, 0, 0, 0, time.UTC),
		RateEndDate:                 time.Date(2017, time.December, 14, 0, 0, 0, 0, time.UTC),
		EIAPricePerGallonMillicents: 320700,
		BaselineRate:                6,
	}
	verrs, err := suite.db.ValidateAndCreate(&newFuelPrice)

	id = uuid.Must(uuid.NewV4())
	newFuelPrice = models.FuelEIADieselPrice{
		ID:                          id,
		PubDate:                     now,
		RateStartDate:               time.Date(2017, time.December, 15, 0, 0, 0, 0, time.UTC),
		RateEndDate:                 time.Date(2018, time.January, 14, 0, 0, 0, 0, time.UTC),
		EIAPricePerGallonMillicents: 320700,
		BaselineRate:                6,
	}
	verrs, err = suite.db.ValidateAndCreate(&newFuelPrice)

	// Overlapping record should cause en error
	id = uuid.Must(uuid.NewV4())
	newFuelPrice = models.FuelEIADieselPrice{
		ID:                          id,
		PubDate:                     now,
		RateStartDate:               time.Date(2017, time.December, 20, 0, 0, 0, 0, time.UTC),
		RateEndDate:                 time.Date(2018, time.January, 14, 0, 0, 0, 0, time.UTC),
		EIAPricePerGallonMillicents: 320700,
		BaselineRate:                6,
	}

	verrs, err = suite.db.ValidateAndCreate(&newFuelPrice)
	suite.EqualError(err, "pq: conflicting key value violates exclusion constraint \"fuel_eia_diesel_prices_daterange_excl\"")
	suite.Empty(verrs.Error())

}

// Create multiple records covering a range of dates
// Can change the dates for start and end ranges and
// can create a default baseline and price to use via assertions
func (suite *ModelSuite) TestMakeFuelEIADieselPrices() {
	testdatagen.MakeFuelEIADieselPrices(suite.db, testdatagen.Assertions{})
	// or call testdatagen.MakeDefaultFuelEIADieselPrices(suite.db)
	// to change the date range:
	//     assertions testdatagen.Assertions{}
	//     assertions.assertions.FuelEIADieselPrice.RateStartDate = time.Date(testdatagen.TestYear-1, time.July, 15, 0, 0, 0, 0, time.UTC)
	//     assertions.assertions.FuelEIADieselPrice.RateEndDate = time.Date(testdatagen.TestYear+1, time.July, 14, 0, 0, 0, 0, time.UTC)
	//     testdatagen.MakeFuelEIADieselPrices(suite.db, assertions)
}

// Create 1 record for the shipment date provided and use assertions
func (suite *ModelSuite) TestMakeFuelEIADieselPriceForDate() {
	rateStartDate := time.Date(2017, time.July, 15, 0, 0, 0, 0, time.UTC)
	assertions := testdatagen.Assertions{}
	assertions.FuelEIADieselPrice.RateStartDate = rateStartDate
	assertions.FuelEIADieselPrice.EIAPricePerGallonMillicents = unit.Millicents(695700)
	shipmentDate := assertions.FuelEIADieselPrice.RateStartDate.AddDate(0, 0, 10)

	testdatagen.MakeFuelEIADieselPriceForDate(suite.db, shipmentDate, assertions)
}

// Create 1 record for the shipment date provided
func (suite *ModelSuite) TestMakeDefaultFuelEIADieselPriceForDate() {
	shipmentDate := time.Now()
	testdatagen.MakeDefaultFuelEIADieselPriceForDate(suite.db, shipmentDate)
}
