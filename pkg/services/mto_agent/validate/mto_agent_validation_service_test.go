package mtoagentvalidate

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MTOAgentValidationServiceSuite struct {
	testingsuite.PopTestSuite
}

func (suite *MTOAgentValidationServiceSuite) SetupTest() {
	err := suite.TruncateAll()
	suite.FatalNoError(err)
}
func TestMTOAgentValidationServiceSuite(t *testing.T) {
	ts := &MTOAgentValidationServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
