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
	clock := clock.NewMock()
	currentDate := clock.Now().AddDate(40, 0, 0)
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
			"AND date_part('month', pub_date) = $2)", currentDate.AddDate(0, -5, 0).Year(), int(currentDate.AddDate(0, -5, 0).Month()))
	err = queryForPriorMonth.All(&fuelEIADeiselPrices)
	if err != nil {
		suite.logger.Error(err.Error())
	}
	err = suite.DB().Destroy(&fuelEIADeiselPrices)
	if err != nil {
		suite.logger.Error("Error deleting eia diesel price", zap.Error(err))
	}

	// run the function
	numMonthsToVerify := 12
	verrs, err := DieselFuelPriceStorer{DB: suite.DB(), Clock: clock}.StoreFuelPrices(numMonthsToVerify)
	suite.NoError(err, "error when creating invoice")
	suite.Empty(verrs.Errors, "validation error when creating diesel prices")

	// check that the last twelve months have fuel data
	resultingFuelEIADeiselPrices := []models.FuelEIADieselPrice{}
	// check that the records are added back in for the months removed
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
