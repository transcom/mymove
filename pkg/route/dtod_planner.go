package route

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// dtodPlanner holds configuration information to make calls to the GHC services (DTOD and RM).
type dtodPlanner struct {
	dtodPlannerMileage DTODPlannerMileage
}

// TransitDistance calculates the distance between two valid addresses
func (p *dtodPlanner) TransitDistance(_ appcontext.AppContext, _ *models.Address, _ *models.Address) (int, error) {
	// This might get retired after we transition over fully to GHC.
	panic("dtod does not need this method and this will be deprecated when the HERE planner is deprecated")
}

// LatLongTransitDistance calculates the distance between two sets of LatLong coordinates
func (p *dtodPlanner) LatLongTransitDistance(_ appcontext.AppContext, _ LatLong, _ LatLong) (int, error) {
	// This might get retired after we transition over fully to GHC.
	panic("dtod does not need this method and this will be deprecated when the HERE planner is deprecated")
}

// Zip5TransitDistanceLineHaul calculates the distance between two valid Zip5s; it is used by the PPM flow
// and checks for minimum distance restriction as PPM doesn't allow short hauls.
func (p *dtodPlanner) Zip5TransitDistanceLineHaul(_ appcontext.AppContext, _ string, _ string) (int, error) {
	// This might get retired after we transition over fully to GHC.
	panic("dtod does not need this method and this will be deprecated when the HERE planner is deprecated")
}

// Zip5TransitDistance calculates the distance between two valid Zip5s; it is used by the PPM flow
func (p *dtodPlanner) Zip5TransitDistance(_ appcontext.AppContext, _ string, _ string) (int, error) {
	// This might get retired after we transition over fully to GHC.
	panic("dtod does not need this method and this will be deprecated when the HERE planner is deprecated")
}

// Zip3TransitDistance calculates the distance between two valid Zip5s; it is used by the PPM flow
func (p *dtodPlanner) Zip3TransitDistance(_ appcontext.AppContext, _ string, _ string) (int, error) {
	// This might get retired after we transition over fully to GHC.
	panic("dtod does not need this method and this will be deprecated when the HERE planner is deprecated")
}

// ZipTransitDistance calculates the distance between two valid Zips
func (p *dtodPlanner) ZipTransitDistance(appCtx appcontext.AppContext, source string, destination string) (int, error) {
	if len(source) < 5 {
		source = fmt.Sprintf("%05s", source)
	}
	if len(destination) < 5 {
		destination = fmt.Sprintf("%05s", destination)
	}

	return p.dtodPlannerMileage.DTODZip5Distance(appCtx, source, destination)
}

// NewDTODPlanner constructs and returns a Planner for GHC routing.
func NewDTODPlanner(dtodPlannerMileage DTODPlannerMileage) Planner {
	return &dtodPlanner{
		dtodPlannerMileage: dtodPlannerMileage,
	}
}
