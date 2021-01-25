package accesscode

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type AccessCodeServiceSuite struct {
	testingsuite.PopTestSuite
}

func (suite *AccessCodeServiceSuite) SetupTest() {
	err := suite.TruncateAll()
	suite.FatalNoError(err)
}

func TestAccessCodeServiceSuite(t *testing.T) {
	ts := &AccessCodeServiceSuite{
		testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
