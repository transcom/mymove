package fuelprice

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/storage"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type FuelPriceServiceSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
	storer storage.FileStorer
}

// AppContextForTest returns the AppContext for the test suite
func (suite *FuelPriceServiceSuite) AppContextForTest() appcontext.AppContext {
	return appcontext.NewAppContext(suite.DB(), suite.logger, nil)
}

func TestFuelPriceSuite(t *testing.T) {
	// Use a no-op logger during testing
	logger := zap.NewNop()

	fakeS3 := storageTest.NewFakeS3Storage(true)

	ts := &FuelPriceServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		logger:       logger,
		storer:       fakeS3,
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
