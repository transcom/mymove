package route

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// dtodPlanner holds configuration information to make calls to the GHC services (DTOD and RM).
type dtodPlanner struct {
	soapClient   SoapCaller
	dtodUsername string
	dtodPassword string
}

// TransitDistance calculates the distance between two valid addresses
func (p *dtodPlanner) TransitDistance(appCtx appcontext.AppContext, source *models.Address, destination *models.Address) (int, error) {
	// This might get retired after we transition over fully to GHC.
	panic("dtod does not need this method and this will be deprecated when the HERE planner is deprecated")
}

// LatLongTransitDistance calculates the distance between two sets of LatLong coordinates
func (p *dtodPlanner) LatLongTransitDistance(appCtx appcontext.AppContext, source LatLong, destination LatLong) (int, error) {
	// This might get retired after we transition over fully to GHC.
	panic("dtod does not need this method and this will be deprecated when the HERE planner is deprecated")
}

// Zip5TransitDistanceLineHaul calculates the distance between two valid Zip5s; it is used by the PPM flow
// and checks for minimum distance restriction as PPM doesn't allow short hauls.
func (p *dtodPlanner) Zip5TransitDistanceLineHaul(appCtx appcontext.AppContext, source string, destination string) (int, error) {
	// This might get retired after we transition over fully to GHC.
	panic("dtod does not need this method and this will be deprecated when the HERE planner is deprecated")
}

// Zip5TransitDistance calculates the distance between two valid Zip5s; it is used by the PPM flow
func (p *dtodPlanner) Zip5TransitDistance(appCtx appcontext.AppContext, source string, destination string) (int, error) {
	// This might get retired after we transition over fully to GHC.
	panic("dtod does not need this method and this will be deprecated when the HERE planner is deprecated")
}

// Zip3TransitDistance calculates the distance between two valid Zip5s; it is used by the PPM flow
func (p *dtodPlanner) Zip3TransitDistance(appCtx appcontext.AppContext, source string, destination string) (int, error) {
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

	dtod := NewDTODZip5Distance(p.dtodUsername, p.dtodPassword, p.soapClient)
	return dtod.DTODZip5Distance(appCtx, source, destination)
}

// NewDtodPlanner constructs and returns a Planner for GHC routing.
func NewDtodPlanner(soapClient SoapCaller, dtodUsername string, dtodPassword string) Planner {
	return &dtodPlanner{
		soapClient:   soapClient,
		dtodUsername: dtodUsername,
		dtodPassword: dtodPassword,
	}
}
