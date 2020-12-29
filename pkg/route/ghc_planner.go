package route

import (
	"github.com/gobuffalo/pop/v5"
	"github.com/tiaguinho/gosoap"

	"github.com/transcom/mymove/pkg/models"
)

// SoapCaller provides an interface for the Call method of the gosoap Client so it can be mocked.
// NOTE: Placing this in a separate package/directory to avoid a circular dependency from an existing mock.
//go:generate mockery --name SoapCaller --outpkg ghcmocks --output ./ghcmocks
type SoapCaller interface {
	Call(m string, p gosoap.SoapParams) (res *gosoap.Response, err error)
}

// ghcPlanner holds configuration information to make calls to the GHC services (DTOD and RM).
type ghcPlanner struct {
	logger       Logger
	db           *pop.Connection
	soapClient   SoapCaller
	dtodUsername string
	dtodPassword string
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
	dtod := NewDTODZip5Distance(p.logger, p.dtodUsername, p.dtodPassword, p.soapClient)
	return dtod.DTODZip5Distance(source, destination)
}

// Zip3TransitDistance calculates the distance between two valid Zip3s
func (p *ghcPlanner) Zip3TransitDistance(source string, destination string) (int, error) {
	return randMcNallyZip3Distance(p.db, source, destination)
}

// NewGHCPlanner constructs and returns a Planner for GHC routing.
func NewGHCPlanner(logger Logger, db *pop.Connection, soapClient SoapCaller, dtodUsername string, dtodPassword string) Planner {
	return &ghcPlanner{
		logger:       logger,
		db:           db,
		soapClient:   soapClient,
		dtodUsername: dtodUsername,
		dtodPassword: dtodPassword,
	}
}
