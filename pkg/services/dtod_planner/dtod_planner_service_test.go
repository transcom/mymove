package dtod

import (
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

// DTODPlannerServiceSuite is a suite for testing DTOD planner
type DTODPlannerServiceSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *DTODPlannerServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestDTODPlannerServiceSuite(t *testing.T) {
	ts := &DTODPlannerServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
