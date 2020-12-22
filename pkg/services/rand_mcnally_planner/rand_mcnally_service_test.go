package randmcnally

import (
	"log"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

// RandMcNallyPlannerServiceSuite is a suite for testing Rand McNally planner
type RandMcNallyPlannerServiceSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *RandMcNallyPlannerServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestRandMcNallyPlannerServiceSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	testSuite := &RandMcNallyPlannerServiceSuite{
		testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger,
	}
	suite.Run(t, testSuite)
	testSuite.PopTestSuite.TearDown()
}
