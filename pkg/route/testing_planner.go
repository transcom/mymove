package route

import "github.com/transcom/mymove/pkg/models"

type testingPlanner struct {
	distance int
}

// TransitDistance calculates the distance between two valid addresses
func (tp testingPlanner) TransitDistance(source *models.Address, destination *models.Address) (int, error) {
	return tp.distance, nil
}

// LatLongTransitDistance calculates the distance between two sets of LatLong coordinates
func (tp testingPlanner) LatLongTransitDistance(source LatLong, dest LatLong) (int, error) {
	return tp.distance, nil
}

// Zip5TransitDistance calculates the distance between two valid Zip5s
func (tp testingPlanner) Zip5TransitDistance(source string, destination string) (int, error) {
	return zip5TransitDistanceHelper(tp, source, destination)
}

// Zip3TransitDistance calculates the distance between two valid Zip3s
func (tp testingPlanner) Zip3TransitDistance(source string, destination string) (int, error) {
	return zip3TransitDistanceHelper(tp, source, destination)
}

// NewTestingPlanner constructs a route.Planner to be used when testing other code
func NewTestingPlanner(distance int) Planner {
	return testingPlanner{
		distance: distance,
	}
}
