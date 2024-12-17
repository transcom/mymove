package route

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// hhgPlanner holds configuration information to make calls to the GHC services (DTOD and RM).
type hhgPlanner struct {
	dtodPlannerMileage DTODPlannerMileage
}

// TransitDistance calculates the distance between two valid addresses
func (p *hhgPlanner) TransitDistance(_ appcontext.AppContext, _ *models.Address, _ *models.Address) (int, error) {
	// This might get retired after we transition over fully to GHC.
	panic("the HHG planner does not need this method and this will be deprecated when the HERE planner is deprecated")
}

// LatLongTransitDistance calculates the distance between two sets of LatLong coordinates
func (p *hhgPlanner) LatLongTransitDistance(_ appcontext.AppContext, _ LatLong, _ LatLong) (int, error) {
	// This might get retired after we transition over fully to GHC.
	panic("the HHG planner does not need this method and this will be deprecated when the HERE planner is deprecated")
}

// Zip5TransitDistanceLineHaul calculates the distance between two valid Zip5s; it is used by the PPM flow
// and checks for minimum distance restriction as PPM doesn't allow short hauls.
func (p *hhgPlanner) Zip5TransitDistanceLineHaul(_ appcontext.AppContext, _ string, _ string) (int, error) {
	// This might get retired after we transition over fully to GHC.
	panic("the HHG planner does not need this method and this will be deprecated when the HERE planner is deprecated")
}

// Zip5TransitDistance calculates the distance between two valid Zip5s; it is used by the PPM flow
func (p *hhgPlanner) Zip5TransitDistance(_ appcontext.AppContext, _ string, _ string) (int, error) {
	// This might get retired after we transition over fully to GHC.
	panic("the HHG planner does not need this method and this will be deprecated when the HERE planner is deprecated")
}

// Zip3TransitDistance calculates the distance between two valid Zip5s; it is used by the PPM flow
func (p *hhgPlanner) Zip3TransitDistance(_ appcontext.AppContext, _ string, _ string) (int, error) {
	// This might get retired after we transition over fully to GHC.
	panic("the HHG planner does not need this method and this will be deprecated when the HERE planner is deprecated")
}

// ZipTransitDistance calculates the distance between two valid Zips
func (p *hhgPlanner) ZipTransitDistance(appCtx appcontext.AppContext, source string, destination string, useDTOD bool) (int, error) {
	sourceZip5 := source
	if len(source) < 5 {
		sourceZip5 = fmt.Sprintf("%05s", source)
	}
	destZip5 := destination
	if len(destination) < 5 {
		destZip5 = fmt.Sprintf("%05s", destination)
	}
	sourceZip3 := sourceZip5[0:3]
	destZip3 := destZip5[0:3]

	if sourceZip3 == destZip3 || useDTOD {
		if sourceZip5 == destZip5 {
			return 1, nil
		}
		return p.dtodPlannerMileage.DTODZip5Distance(appCtx, source, destination)
	}

	// Get reZip3s for origin and destination to compare base point cities.
	// Dont throw/return errors from this. If we dont find them, we'll just use randMcNallyZip3Distance
	sourceReZip3, sErr := models.FetchReZip3Item(appCtx.DB(), sourceZip3)
	if sErr != nil {
		appCtx.Logger().Error("Failed to fetch the reZip3 item for sourceZip3", zap.Error(sErr))
	}
	destinationReZip3, dErr := models.FetchReZip3Item(appCtx.DB(), destZip3)
	if dErr != nil {
		appCtx.Logger().Error("Failed to fetch the reZip3 item for destinationZip3", zap.Error(dErr))
	}

	// Different zip3, same base point city, use DTOD
	if sourceReZip3 != nil && destinationReZip3 != nil && sourceReZip3.BasePointCity == destinationReZip3.BasePointCity {
		return p.dtodPlannerMileage.DTODZip5Distance(appCtx, source, destination)
	}

	return randMcNallyZip3Distance(appCtx, sourceZip3, destZip3)
}

// NewHHGPlanner constructs and returns a Planner for GHC routing.
func NewHHGPlanner(dtodPlannerMileage DTODPlannerMileage) Planner {
	return &hhgPlanner{
		dtodPlannerMileage: dtodPlannerMileage,
	}
}
