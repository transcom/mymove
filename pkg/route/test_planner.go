package route

import "github.com/transcom/mymove/pkg/models"

type testPlanner struct{}

func (tp testPlanner) TransitDistance(source *models.Address, destination *models.Address) (int, error) {
	return 1234, nil
}

// NewTestPlanner constructs a route.Planner to be used when testing other code
func NewTestPlanner() Planner {
	return new(testPlanner)
}
