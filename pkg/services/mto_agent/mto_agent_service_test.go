package mtoagent

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MTOAgentServiceSuite struct {
	testingsuite.PopTestSuite
}

func (suite *MTOAgentServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}
func TestMTOAgentServiceSuite(t *testing.T) {
	ts := &MTOAgentServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
