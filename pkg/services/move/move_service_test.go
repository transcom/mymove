package move

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type MoveServiceSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

// TestAppContext returns the AppContext for the test suite
func (suite *MoveServiceSuite) TestAppContext() appcontext.AppContext {
	return appcontext.NewAppContext(suite.DB(), suite.logger)
}

func TestMoveServiceSuite(t *testing.T) {

	hs := &MoveServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
