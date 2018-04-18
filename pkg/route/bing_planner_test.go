package route

import (
	"os"
)

type BingFullSuite struct {
	PlannerFullSuite
}

// It is only run as part of the server_test_all target and requires
// BING_MAPS_ENDPOINT & BING_MAPS_KEY environment variables to be set
func (suite *BingFullSuite) SetupTest() {
	testEndpoint := os.Getenv("BING_MAPS_ENDPOINT")
	testKey := os.Getenv("BING_MAPS_KEY")
	if len(testEndpoint) == 0 || len(testKey) == 0 {
		suite.T().Fatal("You must set BING_MAPS_ENDPOINT and BING_MAPS_KEY to run this test")
	}
	suite.planner = NewBingPlanner(suite.logger, &testEndpoint, &testKey)
}
