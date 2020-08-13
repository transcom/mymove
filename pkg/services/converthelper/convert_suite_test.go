package converthelper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type ConvertSuite struct {
	testingsuite.PopTestSuite
}

func (suite *ConvertSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestConvertHelperSuite(t *testing.T) {
	ts := &ConvertSuite{
		testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
