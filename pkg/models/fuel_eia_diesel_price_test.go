package models_test

import (
	"testing"
	"time"

	"github.com/facebookgo/clock"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestBasicFuelEIADieselPriceInstantiation() {
	testCases := map[string]struct {
		fuelEIADP    models.FuelEIADieselPrice
		expectedErrs map[string][]string
	}{
		"Successful Create": {
			fuelEIADP: models.FuelEIADieselPrice{
				ID:                          uuid.Must(uuid.NewV4()),
				PubDate:                     time.Now(),
				RateStartDate:               time.Now(),
				RateEndDate:                 time.Now().AddDate(0, 1, -1),
				EIAPricePerGallonMillicents: 320700,
				BaselineRate:                6,
			},
			expectedErrs: nil,
		},

		"Empty Fields": {
			fuelEIADP: models.FuelEIADieselPrice{},
			expectedErrs: map[string][]string{
				"pub_date":                          {"PubDate can not be blank."},
				"rate_start_date":                   {"RateStartDate can not be blank."},
				"rate_end_date":                     {"RateEndDate can not be blank."},
				"e_i_a_price_per_gallon_millicents": {"0 is not greater than 0."},
			},
		},

		"Bad baseline rate": {
			fuelEIADP: models.FuelEIADieselPrice{
				ID:                          uuid.Must(uuid.NewV4()),
				PubDate:                     time.Now(),
				RateStartDate:               time.Now(),
				RateEndDate:                 time.Now(),
				EIAPricePerGallonMillicents: 320700,
				BaselineRate:                102,
			},
			expectedErrs: map[string][]string{
				"baseline_rate": {"102 is not less than 101."},
			},
		},

		"Bad RateEndDate": {
			fuelEIADP: models.FuelEIADieselPrice{
				ID:                          uuid.Must(uuid.NewV4()),
				PubDate:                     time.Now(),
				RateStartDate:               time.Now(),
				RateEndDate:                 time.Now().AddDate(0, -1, 0),
				EIAPricePerGallonMillicents: 320700,
				BaselineRate:                6,
			},
			expectedErrs: map[string][]string{
				"rate_end_date": {"RateEndDate must be after RateStartDate."},
			},
		},
	}

	for name, test := range testCases {
		suite.T().Run(name, func(t *testing.T) {
			suite.verifyValidationErrors(&test.fuelEIADP, test.expectedErrs)
		})
	}

}

func (suite *ModelSuite) TestFuelEIADieselPriceOverlappingDatesConstraint() {
	now := time.Now()
	id := uuid.Must(uuid.NewV4())

	suite.T().Run("Overlapping Dates Constraint Test", func(t *testing.T) {

		// Test for overalapping start and end dates
		newFuelPrice := models.FuelEIADieselPrice{
			ID:                          id,
			PubDate:                     now,
			RateStartDate:               time.Date(2017, time.November, 15, 0, 0, 0, 0, time.UTC),
			RateEndDate:                 time.Date(2017, time.December, 14, 0, 0, 0, 0, time.UTC),
			EIAPricePerGallonMillicents: 320700,
			BaselineRate:                6,
		}
		verrs, err := suite.DB().ValidateAndCreate(&newFuelPrice)

		id = uuid.Must(uuid.NewV4())
		newFuelPrice = models.FuelEIADieselPrice{
			ID:                          id,
			PubDate:                     now,
			RateStartDate:               time.Date(2017, time.December, 15, 0, 0, 0, 0, time.UTC),
			RateEndDate:                 time.Date(2018, time.January, 14, 0, 0, 0, 0, time.UTC),
			EIAPricePerGallonMillicents: 320700,
			BaselineRate:                6,
		}
		verrs, err = suite.DB().ValidateAndCreate(&newFuelPrice)

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

		verrs, err = suite.DB().ValidateAndCreate(&newFuelPrice)
		suite.EqualError(err, "pq: conflicting key value violates exclusion constraint \"no_overlapping_rates\"")
		suite.Empty(verrs.Error())

	})

}

// Create multiple records covering a range of dates
// Can change the dates for start and end ranges and
// can create a default baseline and price to use via assertions
func (suite *ModelSuite) TestMakeFuelEIADieselPrices() {
	testdatagen.MakeFuelEIADieselPrices(suite.DB(), testdatagen.Assertions{})
	// or call testdatagen.MakeDefaultFuelEIADieselPrices(suite.DB())
	// to change the date range:
	//     assertions testdatagen.Assertions{}
	//     assertions.assertions.FuelEIADieselPrice.RateStartDate = time.Date(testdatagen.TestYear-1, time.July, 15, 0, 0, 0, 0, time.UTC)
	//     assertions.assertions.FuelEIADieselPrice.RateEndDate = time.Date(testdatagen.TestYear+1, time.July, 14, 0, 0, 0, 0, time.UTC)
	//     testdatagen.MakeFuelEIADieselPrices(suite.DB(), assertions)
}

// Create 1 record for the shipment date provided and use assertions
func (suite *ModelSuite) TestMakeFuelEIADieselPriceForDate() {
	rateStartDate := time.Date(2017, time.July, 15, 0, 0, 0, 0, time.UTC)
	assertions := testdatagen.Assertions{}
	assertions.FuelEIADieselPrice.RateStartDate = rateStartDate
	assertions.FuelEIADieselPrice.EIAPricePerGallonMillicents = unit.Millicents(695700)
	shipmentDate := assertions.FuelEIADieselPrice.RateStartDate.AddDate(0, 0, 10)

	testdatagen.MakeFuelEIADieselPriceForDate(suite.DB(), shipmentDate, assertions)
}

func (suite *ModelSuite) TestFetchMostRecentFuelPrices() {
	// Make fuel price records for the last twelve months
	clock := clock.NewMock()
	currentTime := clock.Now()
	for month := 0; month < 15; month++ {
		var shipmentDate time.Time

		shipmentDate = currentTime.AddDate(0, -(month - 1), 0)
		testdatagen.MakeDefaultFuelEIADieselPriceForDate(suite.DB(), shipmentDate)
	}

	fuelPrices, err := models.FetchMostRecentFuelPrices(suite.DB(), clock, 12)
	expectedNumFuelPrices := 12
	suite.NoError(err)
	suite.Equal(expectedNumFuelPrices, len(fuelPrices))

	// if the day is after the 15th
	currentTime = currentTime.Add(time.Hour * 24 * 20)
	fuelPrices, err = models.FetchMostRecentFuelPrices(suite.DB(), clock, 12)
	expectedNumFuelPrices = 12
	suite.NoError(err)
	suite.Equal(expectedNumFuelPrices, len(fuelPrices))
}

// Create 1 record for the shipment date provided
func (suite *ModelSuite) TestMakeDefaultFuelEIADieselPriceForDate() {
	shipmentDate := time.Now()
	testdatagen.MakeDefaultFuelEIADieselPriceForDate(suite.DB(), shipmentDate)
}
