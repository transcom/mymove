package route

import (
	"github.com/transcom/mymove/pkg/models"
)

// ghcPlanner holds configuration information to make calls to the GHC services (DTOD and RM).
type ghcPlanner struct {
	logger Logger
}

// TransitDistance calculates the distance between two valid addresses
func (p *ghcPlanner) TransitDistance(source *models.Address, destination *models.Address) (int, error) {
	panic("implement me")
}

// LatLongTransitDistance calculates the distance between two sets of LatLong coordinates
func (p *ghcPlanner) LatLongTransitDistance(source LatLong, destination LatLong) (int, error) {
	panic("implement me")
}

// Zip5TransitDistanceLineHaul calculates the distance between two valid Zip5s; it is used by the PPM flow
// and checks for minimum distance restriction as PPM doesn't allow short hauls.
func (p *ghcPlanner) Zip5TransitDistanceLineHaul(source string, destination string) (int, error) {
	panic("implement me")
}

// Zip3TransitDistance calculates the distance between two valid Zip3s
func (p *ghcPlanner) Zip5TransitDistance(source string, destination string) (int, error) {
	panic("implement me")
}

// Zip5TransitDistance calculates the distance between two valid Zip5s
func (p *ghcPlanner) Zip3TransitDistance(source string, destination string) (int, error) {
	panic("implement me")
}

// NewGHCPlanner constructs and returns a Planner for GHC routing.
func NewGHCPlanner(logger Logger) Planner {
	return &ghcPlanner{
		logger: logger,
	}
}
