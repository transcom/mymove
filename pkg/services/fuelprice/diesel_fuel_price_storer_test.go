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
	dateToTest := time.Date(2010, time.January, 10, 0, 0, 0, 0, time.UTC) // first Mon 1/2010 is 4th
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
	prePubDateTestClock := clock.NewMock()
	dateToTest = time.Date(2010, time.January, 2, 0, 0, 0, 0, time.UTC) // first Mon 1/2010 is 4th
	timeDiff = dateToTest.Sub(prePubDateTestClock.Now())
	prePubDateTestClock.Add(timeDiff)
	dieselFuelPriceStorer := NewDieselFuelPriceStorer(suite.DB(), prePubDateTestClock, mockedFetchFuelPriceData, "", "No data available yet this month")
	verrs, err := dieselFuelPriceStorer.StoreFuelPrices(numMonthsToVerify)
	suite.NoError(err)

	prePubDatePrices := []models.FuelEIADieselPrice{}
	err = queryForThisMonth.All(&prePubDatePrices)
	if err != nil {
		suite.logger.Error(err.Error())
	}
	suite.Empty(&prePubDatePrices)

	// Test case where there is no data for a given month (but should be)
	dieselFuelPriceStorer = NewDieselFuelPriceStorer(suite.DB(), testClock, mockedFetchFuelPriceData, "", "Data missing but expected")
	verrs, err = dieselFuelPriceStorer.StoreFuelPrices(numMonthsToVerify)
	suite.Error(err)

	// Test case where data is missing and is fetched and saved to db
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

	dieselFuelPriceStorer = NewDieselFuelPriceStorer(suite.DB(), testClock, mockedFetchFuelPriceData, "", "Stores all missing months")
	verrs, err = dieselFuelPriceStorer.StoreFuelPrices(numMonthsToVerify)
	suite.NoError(err)
	suite.Empty(verrs.Errors)

	// check that the records are added back in for previously missing ones
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

	// Test case where all desired data already exists in db
	dieselFuelPriceStorer = NewDieselFuelPriceStorer(suite.DB(), testClock, mockedFetchFuelPriceData, "", "No data needed")
	verrs, err = dieselFuelPriceStorer.StoreFuelPrices(numMonthsToVerify)
	suite.NoError(err)

	// Test case where an error message is returned from api
	dieselFuelPriceStorer = NewDieselFuelPriceStorer(suite.DB(), testClock, mockedFetchFuelPriceData, "", "Error")
	verrs, err = dieselFuelPriceStorer.StoreFuelPrices(numMonthsToVerify)
	suite.Error(err)

	// Test case where api returns unexpected JSON structure/value
	dieselFuelPriceStorer = NewDieselFuelPriceStorer(suite.DB(), testClock, mockedFetchFuelPriceData, "", "Unexpected response")
	verrs, err = dieselFuelPriceStorer.StoreFuelPrices(numMonthsToVerify)
	suite.Error(err)
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
	if re.MatchString("Stores all missing months") {
		//TODO: fix this one
		return EiaData{
			SeriesData: []EiaSeriesData{
				{
					Data: [][]interface{}{
						{
							[]string{"20100104"}, []float64{2.79},
						},
						{
							[]string{"20100111"}, []float64{2.81},
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
	if re.MatchString("Error") {
		return EiaData{
			OtherData: EiaOtherData{
				Error: "error message",
			},
		}, nil
	}
	if re.MatchString("Unexpected response") {
		return EiaData{
			OtherData: EiaOtherData{
				//UnknownInfo: []string{
				//	"some sort of response"
				//}
			},
		}, nil
	}
	return EiaData{}, nil

}
