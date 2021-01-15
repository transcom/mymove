package customer

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type CustomerServiceSuite struct {
	testingsuite.PopTestSuite
}

func (suite *CustomerServiceSuite) SetupTest() {
	err := suite.TruncateAll()
	suite.FatalNoError(err)
}

func TestCustomerServiceSuite(t *testing.T) {
	ts := &CustomerServiceSuite{
		testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
