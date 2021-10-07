package sitextension

import (
	"testing"

	"github.com/transcom/mymove/pkg/appcontext"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type SitExtensionServiceSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

// TestAppContext returns the AppContext for the test suite
func (suite *SitExtensionServiceSuite) AppContextForTest() appcontext.AppContext {
	return appcontext.NewAppContext(suite.DB(), suite.logger, nil)
}

func (suite *SitExtensionServiceSuite) SetupTest() {
	err := suite.TruncateAll()
	suite.FatalNoError(err)
}

func TestSitExtensionServiceSuite(t *testing.T) {
	testService := &SitExtensionServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(), // no-op logger used during testing
	}
	suite.Run(t, testService)
	testService.PopTestSuite.TearDown()
}
