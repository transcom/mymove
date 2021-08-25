package mtoserviceitem

import (
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type MTOServiceItemServiceSuite struct {
	testingsuite.PopTestSuite
	testingsuite.AppContextTestHelper
	logger *zap.Logger
}

// TestAppContext returns the AppContext for the test suite
func (suite *MTOServiceItemServiceSuite) TestAppContext() appcontext.AppContext {
	return appcontext.NewAppContext(suite.AppContextTestHelper.CurrentTestContext(suite.T().Name()), suite.DB())
}

func (suite *MTOServiceItemServiceSuite) SetupTest() {
	err := suite.TruncateAll()
	suite.FatalNoError(err)
}

func TestMTOServiceItemServiceSuite(t *testing.T) {
	ts := &MTOServiceItemServiceSuite{
		PopTestSuite:         testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		AppContextTestHelper: testingsuite.NewAppContextTestHelper(),
		logger:               zap.NewNop(),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
