package route

import "github.com/transcom/mymove/pkg/models"

type testingPlanner struct {
	distance int
}

func (tp testingPlanner) TransitDistance(source *models.Address, destination *models.Address) (int, error) {
	return tp.distance, nil
}

func (tp testingPlanner) LatLongTransitDistance(source LatLong, dest LatLong) (int, error) {
	return tp.distance, nil
}

func (tp testingPlanner) Zip5TransitDistance(source string, destination string) (int, error) {
	return zip5TransitDistanceHelper(tp, source, destination)
}

// NewTestingPlanner constructs a route.Planner to be used when testing other code
func NewTestingPlanner(distance int) Planner {
	return testingPlanner{
		distance: distance,
	}
}
