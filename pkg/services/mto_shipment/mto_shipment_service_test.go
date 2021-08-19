package mtoshipment

import (
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MTOShipmentServiceSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *MTOShipmentServiceSuite) SetupTest() {
	err := suite.TruncateAll()
	suite.FatalNoError(err)
}
func TestMTOShipmentServiceSuite(t *testing.T) {

	ts := &MTOShipmentServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
