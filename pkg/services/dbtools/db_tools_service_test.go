package dbtools

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type DBToolsServiceSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *DBToolsServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestDBToolsServiceSuiteServiceSuite(t *testing.T) {
	ts := &DBToolsServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
