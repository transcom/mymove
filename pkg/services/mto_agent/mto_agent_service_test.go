package mtoagent

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MTOAgentServiceSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func (suite *MTOAgentServiceSuite) SetupTest() {
	err := suite.TruncateAll()
	suite.FatalNoError(err)
}
func TestMTOAgentServiceSuite(t *testing.T) {
	ts := &MTOAgentServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
