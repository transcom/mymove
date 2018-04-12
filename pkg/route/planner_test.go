package route

import (
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type PlannerSuite struct {
	suite.Suite
	logger *zap.Logger
}

type PlannerFullSuite struct {
	PlannerSuite
	planner Planner
}

func TestHandlerSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	var testSuite suite.TestingSuite
	if testing.Short() == false {
		testSuite = &PlannerFullSuite{PlannerSuite: PlannerSuite{logger: logger}}
	} else {
		testSuite = &PlannerSuite{logger: logger}
	}
	suite.Run(t, testSuite)
}
