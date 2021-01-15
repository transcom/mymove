package mtoshipment

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MTOShipmentServiceSuite struct {
	testingsuite.PopTestSuite
}

func (suite *MTOShipmentServiceSuite) SetupTest() {
	err := suite.TruncateAll()
	suite.FatalNoError(err)
}
func TestMTOShipmentServiceSuite(t *testing.T) {

	ts := &MTOShipmentServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
