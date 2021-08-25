package upload

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type UploadsServiceSuite struct {
	testingsuite.PopTestSuite
	testingsuite.AppContextTestHelper
	logger *zap.Logger
}

// TestAppContext returns the AppContext for the test suite
func (suite *UploadsServiceSuite) TestAppContext() appcontext.AppContext {
	return appcontext.NewAppContext(suite.AppContextTestHelper.CurrentTestContext(suite.T().Name()), suite.DB())
}

func TestUploadsServiceSuite(t *testing.T) {

	hs := &UploadsServiceSuite{
		PopTestSuite:         testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		AppContextTestHelper: testingsuite.NewAppContextTestHelper(),
		logger:               zap.NewNop(), // Use a no-op logger during testing, // Use a no-op logger during testing
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
