package models_test

import (
	"time"

	"github.com/benbjohnson/clock"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgerrcode"

	"github.com/transcom/mymove/pkg/db/dberr"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
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
				"pub_date":                        {"PubDate can not be blank."},
				"rate_start_date":                 {"RateStartDate can not be blank."},
				"rate_end_date":                   {"RateEndDate can not be blank."},
				"eia_price_per_gallon_millicents": {"0 is not greater than 0."},
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
		suite.Run(name, func() {
			suite.verifyValidationErrors(&test.fuelEIADP, test.expectedErrs)
		})
	}

}

func (suite *ModelSuite) TestFuelEIADieselPriceOverlappingDatesConstraint() {
	now := time.Now()
	id := uuid.Must(uuid.NewV4())

	suite.Run("Overlapping Dates Constraint Test", func() {

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
		suite.NoError(err)
		suite.NoVerrs(verrs)

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
		suite.NoError(err)
		suite.NoVerrs(verrs)

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

		suite.True(dberr.IsDBErrorForConstraint(err, pgerrcode.ExclusionViolation, "no_overlapping_rates"))
		suite.Empty(verrs.Error())
	})
}

func (suite *ModelSuite) TestFetchMostRecentFuelPrices() {
	// Make fuel price records for the last twelve months
	clock := clock.NewMock()
	currentTime := clock.Now()
	for month := 0; month < 15; month++ {
		shipmentDate := currentTime.AddDate(0, -(month - 1), 0)
		testdatagen.MakeDefaultFuelEIADieselPriceForDate(suite.DB(), shipmentDate)
	}

	fuelPrices, err := models.FetchMostRecentFuelPrices(suite.DB(), clock, 12)
	expectedNumFuelPrices := 12
	suite.NoError(err)
	suite.Equal(expectedNumFuelPrices, len(fuelPrices))

	// if the day is after the 15th
	clock.Add(time.Hour * 24 * 20)
	fuelPrices, err = models.FetchMostRecentFuelPrices(suite.DB(), clock, 12)
	expectedNumFuelPrices = 12
	suite.NoError(err)
	suite.Equal(expectedNumFuelPrices, len(fuelPrices))
}
