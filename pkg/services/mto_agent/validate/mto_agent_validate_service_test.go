package mtoagentvalidate

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MTOAgentValidateServiceSuite struct {
	testingsuite.PopTestSuite
}

func (suite *MTOAgentValidateServiceSuite) SetupTest() {
	err := suite.TruncateAll()
	suite.FatalNoError(err)
}
func TestMTOAgentValidateServiceSuite(t *testing.T) {
	ts := &MTOAgentValidateServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
