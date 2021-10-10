package dbtools

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type DBToolsServiceSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

// AppContextForTest returns the AppContext for the test suite
func (suite *DBToolsServiceSuite) AppContextForTest() appcontext.AppContext {
	return appcontext.NewAppContext(suite.DB(), suite.logger, nil)
}

func TestDBToolsServiceSuite(t *testing.T) {
	ts := &DBToolsServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
