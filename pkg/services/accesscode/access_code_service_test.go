package accesscode

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type AccessCodeServiceSuite struct {
	testingsuite.PopTestSuite
	testingsuite.AppContextTestHelper
	logger *zap.Logger
}

// TestAppContext returns the AppContext for the test suite
func (suite *AccessCodeServiceSuite) TestAppContext() appcontext.AppContext {
	return appcontext.NewAppContext(suite.AppContextTestHelper.CurrentTestContext(suite.T().Name()), suite.DB())
}

func (suite *AccessCodeServiceSuite) SetupTest() {
	err := suite.TruncateAll()
	suite.FatalNoError(err)
}

func TestAccessCodeServiceSuite(t *testing.T) {
	ts := &AccessCodeServiceSuite{
		PopTestSuite:         testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		AppContextTestHelper: testingsuite.NewAppContextTestHelper(),
		logger:               zap.NewNop(), // Use a no-op logger during testing,
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
