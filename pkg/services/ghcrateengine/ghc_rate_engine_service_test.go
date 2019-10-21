package ghcrateengine

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type GHCRateEngineServiceSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *GHCRateEngineServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestAccessCodeServiceSuite(t *testing.T) {
	ts := &GHCRateEngineServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
