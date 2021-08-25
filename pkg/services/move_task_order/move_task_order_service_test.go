package movetaskorder_test

import (
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type MoveTaskOrderServiceSuite struct {
	testingsuite.PopTestSuite
	testingsuite.AppContextTestHelper
	logger *zap.Logger
}

// TestAppContext returns the AppContext for the test suite
func (suite *MoveTaskOrderServiceSuite) TestAppContext() appcontext.AppContext {
	return appcontext.NewAppContext(suite.AppContextTestHelper.CurrentTestContext(suite.T().Name()), suite.DB())
}

func TestMoveTaskOrderServiceSuite(t *testing.T) {
	ts := &MoveTaskOrderServiceSuite{
		PopTestSuite:         testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		AppContextTestHelper: testingsuite.NewAppContextTestHelper(),
		logger:               zap.NewNop(), // Use a no-op logger during testing, // Use a no-op logger during testing
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
