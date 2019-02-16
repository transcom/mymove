package fuelprice

import (
	"fmt"
	"testing"
	"time"

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
	// create fuel prices in db for last 15 months
	for month := 0; month < 15; month++ {
		var shipmentDate time.Time
		shipmentDate = time.Now().AddDate(0, -(month - 1), 0)
		testdatagen.MakeDefaultFuelEIADieselPriceForDate(suite.DB(), shipmentDate)
	}
	// remove this month's data
	fuelEIADeiselPrices := []models.FuelEIADieselPrice{}

	queryForThisMonth := suite.DB().RawQuery(
		"SELECT * FROM fuel_eia_diesel_prices WHERE (DATEPART(year, pub_date) = $1"+
			"AND DATEPART(month, pub_date = $2))", time.Now().Year(), time.Now().Month())
	err := queryForThisMonth.All(&fuelEIADeiselPrices)
	if err != nil {
		suite.logger.Error(err.Error())
	}

	err = suite.DB().Destroy(&fuelEIADeiselPrices)
	if err != nil {
		suite.logger.Error("Error deleting eia diesel price", zap.Error(err))
	}

	//remove a different month's data (not next to this month)
	queryForPriorMonth := suite.DB().RawQuery(
		"SELECT * FROM fuel_eia_diesel_prices WHERE (DATEPART(yy, pub_date) = $1"+
			"AND DATEPART(mm, pub_date = $2))", time.Now().AddDate(0, -5, 0).Year(), time.Now().AddDate(0, -5, 0).Month())
	err = queryForPriorMonth.All(&fuelEIADeiselPrices)
	if err != nil {
		suite.logger.Error(err.Error())
	}

	err = suite.DB().Destroy(&fuelEIADeiselPrices)
	if err != nil {
		suite.logger.Error("Error deleting eia diesel price", zap.Error(err))
	}

	// run the function
	verrs, err := FuelPriceStorer{DB: suite.DB()}.StoreFuelPrices()
	suite.NoError(err, "error when creating invoice")
	suite.Empty(verrs.Errors, "validation error when creating diesel prices")

	// check that the last twelve months have fuel data

	// check that the records are added back in for the months removed
	err = queryForThisMonth.All(&fuelEIADeiselPrices)
	if err != nil {
		suite.logger.Error(err.Error())
	}
	suite.NotEmpty(&fuelEIADeiselPrices)

	err = queryForPriorMonth.All(&fuelEIADeiselPrices)
	if err != nil {
		suite.logger.Error(err.Error())
	}
	suite.NotEmpty(&fuelEIADeiselPrices)
}
