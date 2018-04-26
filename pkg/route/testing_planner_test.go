package route

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *PlannerSuite) TestNewTestPlanner() {
	t := suite.T()

	planner := NewTestingPlanner(1234)
	if planner == nil {
		t.Error("NewTestingPlanner returned nil")
	}

	addressOne := models.Address{
		StreetAddress1: "742 Evergreen Terrace",
		City:           "Springfield",
		State:          "Nevada",
		PostalCode:     "89011"}
	addressTwo := models.Address{
		StreetAddress1: "1 Transcom Towers",
		City:           "Scott Airforce Base",
		State:          "Illinois",
		PostalCode:     "62225-5357"}

	distance, err := planner.TransitDistance(&addressOne, &addressTwo)
	if err != nil {
		t.Error("Test route planner returned an error.")
	}
	if distance != 1234 {
		t.Errorf("Expected distance from test_planner should be 1234, got %d", distance)
	}
}
