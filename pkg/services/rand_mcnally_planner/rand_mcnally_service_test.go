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
	testingsuite.BaseTestSuite
	logger Logger
}

func (suite *RandMcNallyPlannerServiceSuite) SetupTest() {
}

func TestRandMcNallyPlannerServiceSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	testSuite := &RandMcNallyPlannerServiceSuite{logger: logger}
	suite.Run(t, testSuite)
}
