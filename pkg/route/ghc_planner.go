package route

import (
	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
)

// ghcPlanner holds configuration information to make calls to the GHC services (DTOD and RM).
type ghcPlanner struct {
	db     *pop.Connection
	logger Logger
}

// TransitDistance calculates the distance between two valid addresses
func (p *ghcPlanner) TransitDistance(source *models.Address, destination *models.Address) (int, error) {
	// This might get retired after we transition over fully to GHC.

	panic("implement me")
}

// LatLongTransitDistance calculates the distance between two sets of LatLong coordinates
func (p *ghcPlanner) LatLongTransitDistance(source LatLong, destination LatLong) (int, error) {
	// This might get retired after we transition over fully to GHC.

	panic("implement me")
}

// Zip5TransitDistanceLineHaul calculates the distance between two valid Zip5s; it is used by the PPM flow
// and checks for minimum distance restriction as PPM doesn't allow short hauls.
func (p *ghcPlanner) Zip5TransitDistanceLineHaul(source string, destination string) (int, error) {
	// This might get retired after we transition over fully to GHC.

	panic("implement me")
}

// Zip5TransitDistance calculates the distance between two valid Zip5s
func (p *ghcPlanner) Zip5TransitDistance(source string, destination string) (int, error) {
	// Placeholder for the DTOD-based zip5-to-zip5 distance. This will be determined by making
	// a SOAP call to DTOD using the provided source/destination zip5 and returning the
	// associated distance.
	//
	// It could be implemented as a service object if we expect reuse beyond the planner, or
	// unexported code in this package if we always expect to access it via the planner.

	panic("implement me")
}

// Zip3TransitDistance calculates the distance between two valid Zip3s
func (p *ghcPlanner) Zip3TransitDistance(source string, destination string) (int, error) {
	return randMcNallyZip3Distance(p.db, source, destination)
}

// NewGHCPlanner constructs and returns a Planner for GHC routing.
func NewGHCPlanner(db *pop.Connection, logger Logger) Planner {
	return &ghcPlanner{
		db:     db,
		logger: logger,
	}
}
