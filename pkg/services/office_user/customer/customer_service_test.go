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
	suite.DB().TruncateAll()
}

func TestAccessCodeServiceSuite(t *testing.T) {
	ts := &CustomerServiceSuite{
		testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
