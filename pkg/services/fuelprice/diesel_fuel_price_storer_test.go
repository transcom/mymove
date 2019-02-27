package fuelprice

import (
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

	//var fuelPriceTestCases := []struct{
	//	url string
	//	numMonthsToCheck int
	//	expectedNumRecords int
	//}{
	//	{
	//		url: "gets all missing months",
	//		numMonthsToCheck: 12,
	//		expectedNumRecords: 12,
	//	},
	//}

	testClock := clock.NewMock()
	dateToTest := time.Date(2010, time.January, 10, 0, 0, 0, 0, time.UTC)
	timeDiff := dateToTest.Sub(testClock.Now())
	testClock.Add(timeDiff)
	currentDate := testClock.Now()
	// create fuel prices in db for last 15 months
	for month := 0; month < 15; month++ {
		var shipmentDate time.Time
		shipmentDate = currentDate.AddDate(0, -(month - 1), 0)
		testdatagen.MakeDefaultFuelEIADieselPriceForDate(suite.DB(), shipmentDate)
	}
	// remove this month's data
	fuelEIADeiselPrices := []models.FuelEIADieselPrice{}
	queryForThisMonth := suite.DB().RawQuery(
		"SELECT * FROM fuel_eia_diesel_prices WHERE (date_part('year', pub_date) = $1 "+
			"AND date_part('month', pub_date) = $2)", currentDate.Year(), int(currentDate.Month()))
	err := queryForThisMonth.All(&fuelEIADeiselPrices)
	if err != nil {
		suite.logger.Error(err.Error())
	}

	//remove a different month's data (not next to this month)
	queryForPriorMonth := suite.DB().RawQuery(
		"SELECT * FROM fuel_eia_diesel_prices WHERE (date_part('year', pub_date) = $1 "+
			"AND date_part('month', pub_date) = $2)", currentDate.AddDate(0, -3, 0).Year(), int(currentDate.AddDate(0, -3, 0).Month()))
	err = queryForPriorMonth.All(&fuelEIADeiselPrices)
	if err != nil {
		suite.logger.Error(err.Error())
	}
	err = suite.DB().Destroy(&fuelEIADeiselPrices)
	if err != nil {
		suite.logger.Error("Error deleting eia diesel price", zap.Error(err))
	}
	nowInDB := []models.FuelEIADieselPrice{}
	suite.DB().All(&nowInDB)

	// Test case where data is missing from the most recent months
	// run the function
	numMonthsToVerify := 10
	dieselFuelPriceStorer := NewDieselFuelPriceStorer(suite.DB(), testClock, FetchFuelPriceData, "", "gets all missing months")
	verrs, err := dieselFuelPriceStorer.StoreFuelPrices(numMonthsToVerify)
	suite.NoError(err, "error when creating diesel prices")
	suite.Empty(verrs.Errors, "validation error when creating diesel prices")

	// check that the given number of prior months have fuel data
	resultingFuelEIADeiselPrices := []models.FuelEIADieselPrice{}

	// check that the records are added back in for each of months previously missing
	err = queryForThisMonth.All(&resultingFuelEIADeiselPrices)
	if err != nil {
		suite.logger.Error(err.Error())
	}
	suite.NotEmpty(&resultingFuelEIADeiselPrices)

	err = queryForPriorMonth.All(&fuelEIADeiselPrices)
	if err != nil {
		suite.logger.Error(err.Error())
	}
	suite.NotEmpty(&fuelEIADeiselPrices)

	// Test case where there is no data yet available for the current month (nor expected)
	dieselFuelPriceStorer = NewDieselFuelPriceStorer(suite.DB(), testClock, FetchFuelPriceData, "", "No data available")
	verrs, err = dieselFuelPriceStorer.StoreFuelPrices(numMonthsToVerify)
	suite.NoError(err)

	// Test case where there is no data for a given month (but should be)
	dieselFuelPriceStorer = NewDieselFuelPriceStorer(suite.DB(), testClock, FetchFuelPriceData, "", "Data missing")
	verrs, err = dieselFuelPriceStorer.StoreFuelPrices(numMonthsToVerify)

	// Test case where all desired data already exists in db
	dieselFuelPriceStorer = NewDieselFuelPriceStorer(suite.DB(), testClock, FetchFuelPriceData, "", "No data needed")
	verrs, err = dieselFuelPriceStorer.StoreFuelPrices(numMonthsToVerify)
	suite.NoError(err)

	// Test case where an error message is returned from api
	dieselFuelPriceStorer = NewDieselFuelPriceStorer(suite.DB(), testClock, FetchFuelPriceData, "", "Error")
	verrs, err = dieselFuelPriceStorer.StoreFuelPrices(numMonthsToVerify)
	suite.Error(err)

	// Test case where api returns unexpected JSON structure/value
	dieselFuelPriceStorer = NewDieselFuelPriceStorer(suite.DB(), testClock, FetchFuelPriceData, "", "Unexpected response")
	verrs, err = dieselFuelPriceStorer.StoreFuelPrices(numMonthsToVerify)
	suite.Error(err)
}

func mockedFetchFuelPriceData(url string) (data *EiaData, err error) {
	// TODO: build out structs to match test scenarios
	switch url {
	case "gets all missing months":
		return &EiaData{
			//SeriesData: []EiaSeriesData{
			//	{
			//		Data: [][]interface{
			//				[string{"20100104"}, float64{2.79}],
			//
			//				//[]{
			//				//"20100111", 2.81,
			//				//},
			//		},
			//	},
			//},
		}, nil
	case "No data needed":
		return &EiaData{}, nil
	case "No data available":
		return &EiaData{}, nil
	case "Unexpected response":
		return &EiaData{}, nil
	case "Error":
		return &EiaData{}, nil
	default:
		return nil, nil
	}
}
