package dtod

import (
	"log"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

// DTODPlannerServiceSuite is a suite for testing DTOD planner
type DTODPlannerServiceSuite struct {
	testingsuite.BaseTestSuite
	logger Logger
}

func (suite *DTODPlannerServiceSuite) SetupTest() {
}

func TestDTODPlannerServiceSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	testSuite := &DTODPlannerServiceSuite{logger: logger}
	suite.Run(t, testSuite)
}
