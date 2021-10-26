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
func (suite *SitExtensionServiceSuite) TestAppContext() appcontext.AppContext {
	return appcontext.NewAppContext(suite.DB(), suite.logger)
}

func TestSitExtensionServiceSuite(t *testing.T) {
	testService := &SitExtensionServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		logger:       zap.NewNop(), // no-op logger used during testing
	}
	suite.Run(t, testService)
	testService.PopTestSuite.TearDown()
}
