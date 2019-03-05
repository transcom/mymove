package fuelprice

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/facebookgo/clock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type FuelPriceServiceSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
	storer storage.FileStorer
}

func (suite *FuelPriceServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestFuelPriceSuite(t *testing.T) {
	// Use a no-op logger during testing
	logger := zap.NewNop()

	fakeS3 := storageTest.NewFakeS3Storage(true)

	hs := &FuelPriceServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       logger,
		storer:       fakeS3,
	}
	suite.Run(t, hs)
}

func (suite *FuelPriceServiceSuite) TestStoreFuelPrices() {
	testClock := clock.NewMock()
	dateToTest := time.Date(2010, time.January, 12, 0, 0, 0, 0, time.UTC) // first Mon 1/2010 is 4th
	timeDiff := dateToTest.Sub(testClock.Now())
	testClock.Add(timeDiff)
	currentDate := testClock.Now().UTC()
	// create fuel prices in db for last 15 months
	for month := 0; month < 15; month++ {
		var shipmentDate time.Time
		shipmentDate = currentDate.AddDate(0, -(month - 1), 0)
		testdatagen.MakeDefaultFuelEIADieselPriceForDate(suite.DB(), shipmentDate)
	}
	// remove this month's data
	thisMonthPrices := []models.FuelEIADieselPrice{}
	queryForThisMonth := suite.DB().RawQuery(
		"SELECT * FROM fuel_eia_diesel_prices WHERE (date_part('year', pub_date) = $1 "+
			"AND date_part('month', pub_date) = $2)", currentDate.Year(), int(currentDate.Month()))
	err := queryForThisMonth.All(&thisMonthPrices)
	if err != nil {
		suite.logger.Error(err.Error())
	}
	err = suite.DB().Destroy(&thisMonthPrices)
	if err != nil {
		suite.logger.Error("Error deleting eia diesel price", zap.Error(err))
	}

	numMonthsToVerify := 10

	// Test case where there is no data yet available for the current month (nor expected).
	suite.T().Run("no data yet available for current month", func(t *testing.T) {
		prePubDateTestClock := clock.NewMock()
		dateToTest = time.Date(2010, time.January, 2, 0, 0, 0, 0, time.UTC) // first Mon 1/2010 is 4th
		timeDiff = dateToTest.Sub(prePubDateTestClock.Now().UTC())
		prePubDateTestClock.Add(timeDiff)
		dieselFuelPriceStorer := NewDieselFuelPriceStorer(suite.DB(), suite.logger, prePubDateTestClock, mockedFetchFuelPriceData, "", "No data available yet this month")
		verrs, err := dieselFuelPriceStorer.StoreFuelPrices(numMonthsToVerify)
		suite.NoError(err)
		suite.Empty(verrs.Errors)

		prePubDatePrices := []models.FuelEIADieselPrice{}
		err = queryForThisMonth.All(&prePubDatePrices)
		if err != nil {
			suite.logger.Error(err.Error())
		}
		suite.Empty(&prePubDatePrices)
	})

	// Test case where there is no data for a given month (but should be)
	suite.T().Run("no data available for current month (though expected)", func(t *testing.T) {
		dieselFuelPriceStorer := NewDieselFuelPriceStorer(suite.DB(), suite.logger, testClock, mockedFetchFuelPriceData, "", "Data missing but expected")
		verrs, err := dieselFuelPriceStorer.StoreFuelPrices(numMonthsToVerify)
		suite.Error(err)
		suite.Empty(verrs.Errors)
	})

	// Test case where there is no data for a given month (but should be) and the first Monday is a holiday
	suite.T().Run("no data available for current month (though expected) and first Monday is a holiday", func(t *testing.T) {
		postMonHolidayTestClock := clock.NewMock()
		dateToTest = time.Date(2018, time.September, 5, 0, 0, 0, 0, time.UTC) // Labor Day 2018 Mon 3/3
		timeDiff = dateToTest.Sub(postMonHolidayTestClock.Now().UTC())
		postMonHolidayTestClock.Add(timeDiff)

		dieselFuelPriceStorer := NewDieselFuelPriceStorer(suite.DB(), suite.logger, postMonHolidayTestClock, mockedFetchFuelPriceData, "", "Data missing but expected")
		verrs, err := dieselFuelPriceStorer.StoreFuelPrices(numMonthsToVerify)
		suite.Error(err)
		suite.Empty(verrs.Errors)
	})

	// Test case where an error message is returned from api
	suite.T().Run("error message returned from api", func(t *testing.T) {
		dieselFuelPriceStorer := NewDieselFuelPriceStorer(suite.DB(), suite.logger, testClock, mockedFetchFuelPriceData, "", "Error")
		verrs, err := dieselFuelPriceStorer.StoreFuelPrices(numMonthsToVerify)
		suite.Error(err)
		suite.Empty(verrs.Errors)
	})
	// Test case where api returns unexpected JSON structure/value
	suite.T().Run("unexpected JSON structure returned from api", func(t *testing.T) {
		dieselFuelPriceStorer := NewDieselFuelPriceStorer(suite.DB(), suite.logger, testClock, mockedFetchFuelPriceData, "", "Unexpected response")
		verrs, err := dieselFuelPriceStorer.StoreFuelPrices(numMonthsToVerify)
		suite.Error(err)
		suite.Empty(verrs.Errors)
	})

	suite.T().Run("stores current month missing data", func(t *testing.T) {
		dieselFuelPriceStorer := NewDieselFuelPriceStorer(suite.DB(), suite.logger, testClock, mockedFetchFuelPriceData, "", "Stores current month missing data")
		verrs, err := dieselFuelPriceStorer.StoreFuelPrices(numMonthsToVerify)
		suite.NoError(err)
		suite.Empty(verrs.Errors)

		currentMonthPrices := []models.FuelEIADieselPrice{}
		err = queryForThisMonth.All(&currentMonthPrices)
		if err != nil {
			suite.logger.Error(err.Error())
		}
		suite.NotEmpty(currentMonthPrices)
		dbBaselineRate := currentMonthPrices[0].BaselineRate
		expectedBaselineRate := int64(3)
		suite.Equal(expectedBaselineRate, dbBaselineRate)
	})

	// Test case where data is missing from a prior month and saved to db
	suite.T().Run("stores data missing for a month that is not the current month", func(t *testing.T) {
		//remove a month other than current month
		priorMonthsToRemove := []models.FuelEIADieselPrice{}

		queryForPriorMonth := suite.DB().RawQuery(
			"SELECT * FROM fuel_eia_diesel_prices WHERE (date_part('year', pub_date) = $1 "+
				"AND date_part('month', pub_date) = $2)", currentDate.AddDate(0, -3, 0).Year(), int(currentDate.AddDate(0, -3, 0).Month()))
		err = queryForPriorMonth.All(&priorMonthsToRemove)
		if err != nil {
			suite.logger.Error(err.Error())
		}
		err = suite.DB().Destroy(&priorMonthsToRemove)
		if err != nil {
			suite.logger.Error("Error deleting eia diesel price", zap.Error(err))
		}

		dieselFuelPriceStorer := NewDieselFuelPriceStorer(suite.DB(), suite.logger, testClock, mockedFetchFuelPriceData, "", "Store previous month missing data")
		verrs, err := dieselFuelPriceStorer.StoreFuelPrices(numMonthsToVerify)
		suite.NoError(err)
		suite.Empty(verrs.Errors)

		// check that the records are added back in for previous month
		resultingFuelEIADeiselPrices := []models.FuelEIADieselPrice{}

		err = queryForThisMonth.All(&resultingFuelEIADeiselPrices)
		if err != nil {
			suite.logger.Error(err.Error())
		}
		suite.NotEmpty(&resultingFuelEIADeiselPrices)

		priorMonthPrices := []models.FuelEIADieselPrice{}
		err = queryForPriorMonth.All(&priorMonthPrices)
		if err != nil {
			suite.logger.Error(err.Error())
		}
		suite.NotEmpty(&priorMonthPrices)
		dbBaselineRate := priorMonthPrices[0].BaselineRate
		expectedBaselineRate := int64(1)
		suite.Equal(expectedBaselineRate, dbBaselineRate)
	})

	// Test case where all desired data already exists in db
	suite.T().Run("all desired data already exists in the db", func(t *testing.T) {
		dieselFuelPriceStorer := NewDieselFuelPriceStorer(suite.DB(), suite.logger, testClock, mockedFetchFuelPriceData, "", "No data needed")
		verrs, err := dieselFuelPriceStorer.StoreFuelPrices(numMonthsToVerify)
		suite.NoError(err)
		suite.Empty(verrs.Errors)
	})
}

func mockedFetchFuelPriceData(url string) (data EiaData, err error) {
	// The url gets a querystring added to the url on the struct, so that needs to be removed to detect the string
	re := regexp.MustCompile(`%20`)
	url = re.ReplaceAllLiteralString(url, ` `)
	url = strings.Split(url, "?")[0]
	re = regexp.MustCompile(`^` + url + `.*`)

	if re.MatchString("No data available yet this month") {
		return EiaData{
			SeriesData: []EiaSeriesData{
				{
					Data: [][]interface{}{},
				},
			},
		}, nil
	}
	if re.MatchString("Data missing but expected") {
		return EiaData{
			SeriesData: []EiaSeriesData{
				{
					Data: [][]interface{}{},
				},
			},
		}, nil
	}
	if re.MatchString("Error") {
		return EiaData{
			OtherData: map[string]interface{}{
				"error": "error message",
			},
		}, nil
	}
	if re.MatchString("Unexpected response") {
		return EiaData{
			OtherData: map[string]interface{}{
				"rogueInfo": "Unexpected response from api",
			},
		}, nil
	}
	if re.MatchString("Stores current month missing data") {
		return EiaData{
			SeriesData: []EiaSeriesData{
				{
					Data: [][]interface{}{
						{
							"20100104", 2.797,
						},
						{
							"20100111", 2.81,
						},
					},
				},
			},
		}, nil
	}
	if re.MatchString("Store previous month missing data") {
		return EiaData{
			SeriesData: []EiaSeriesData{
				{
					Data: [][]interface{}{
						{
							"20091005", 2.582,
						},
						{
							"20091012", 3.41,
						},
						{
							"20091019", 3.32,
						},
						{
							"20091026", 3.07,
						},
					},
				},
			},
		}, nil
	}
	if re.MatchString("No data needed") {
		return EiaData{
			SeriesData: []EiaSeriesData{
				{
					Data: [][]interface{}{},
				},
			},
		}, nil
	}
	return EiaData{}, nil

}
