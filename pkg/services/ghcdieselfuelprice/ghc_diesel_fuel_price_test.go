package ghcdieselfuelprice

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type GhcDieselFuelPriceServiceSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func (suite *GhcDieselFuelPriceServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestGhcDieselFuelPriceServiceSuite(t *testing.T) {
	// Use a no-op logger during testing
	logger := zap.NewNop()

	ts := &GhcDieselFuelPriceServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       logger,
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
