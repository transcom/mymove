package fuelprice

import (
	"github.com/facebookgo/clock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
	"time"

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
	// run the function
	numMonthsToVerify := 10
	verrs, err := DieselFuelPriceStorer{DB: suite.DB(), Clock: testClock, FetchFuelData: FetchFuelPriceData, EiaKey: "", URL: "all months"}.StoreFuelPrices(numMonthsToVerify)
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
}

func mockedFetchFuelPriceData(url string) (data *EiaData, err error) {
	// TODO: build out structs to match test scenarios
	switch url {
	case "all months":
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
	case "None To Fetch":
		return &EiaData{}, nil
	case "No data available":
		return &EiaData{}, nil
	case "Error":
		return &EiaData{}, nil
	default:
		return nil, nil
	}
}
