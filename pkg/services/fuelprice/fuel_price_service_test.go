package fuelprice

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/storage"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type FuelPriceServiceSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
	storer storage.FileStorer
}

func (suite *FuelPriceServiceSuite) SetupTest() {
	err := suite.TruncateAll()
	suite.FatalNoError(err)
}

func TestFuelPriceSuite(t *testing.T) {
	// Use a no-op logger during testing
	logger := zap.NewNop()

	fakeS3 := storageTest.NewFakeS3Storage(true)

	ts := &FuelPriceServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       logger,
		storer:       fakeS3,
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
