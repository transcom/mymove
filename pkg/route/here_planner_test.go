package route

import (
	"os"
)

type HereFullSuite struct {
	PlannerFullSuite
}

func (suite *HereFullSuite) SetupTest() {
	geocodeEndpoint := os.Getenv("HERE_MAPS_GEOCODE_ENDPOINT")
	routingEndpoint := os.Getenv("HERE_MAPS_ROUTING_ENDPOINT")
	testAppID := os.Getenv("HERE_MAPS_APP_ID")
	testAppCode := os.Getenv("HERE_MAPS_APP_CODE")
	if len(geocodeEndpoint) == 0 || len(routingEndpoint) == 0 ||
		len(testAppID) == 0 || len(testAppCode) == 0 {
		suite.T().Fatal("You must set HERE_... environment variables to run this test")
	}

	suite.planner = NewHEREPlanner(suite.logger, &geocodeEndpoint, &routingEndpoint, &testAppID, &testAppCode)
}
