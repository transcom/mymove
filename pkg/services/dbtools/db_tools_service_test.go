package dbtools

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type DBToolsServiceSuite struct {
	testingsuite.PopTestSuite
	testingsuite.AppContextTestHelper
	logger Logger
}

func (suite *DBToolsServiceSuite) SetupTest() {
	err := suite.TruncateAll()
	suite.FatalNoError(err)
}

func TestDBToolsServiceSuite(t *testing.T) {
	ts := &DBToolsServiceSuite{
		PopTestSuite:         testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		AppContextTestHelper: testingsuite.NewAppContextTestHelper(),
		logger:               zap.NewNop(), // Use a no-op logger during testing, // Use a no-op logger during testing
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
