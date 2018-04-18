package route

import "github.com/transcom/mymove/pkg/models"

type testingPlanner struct{}

func (tp testingPlanner) TransitDistance(source *models.Address, destination *models.Address) (int, error) {
	return 1234, nil
}

// NewTestingPlanner constructs a route.Planner to be used when testing other code
func NewTestingPlanner() Planner {
	return new(testingPlanner)
}
