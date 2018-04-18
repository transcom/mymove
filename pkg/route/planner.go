package route

import (
	"net/url"
	"strings"

	"github.com/transcom/mymove/pkg/models"
)

// LatLong is used to hold latitude and longitude as floats
type LatLong struct {
	Latitude  float32
	Longitude float32
}

// urlencodeAddress converts an address into a comma separated string which is safely encoded to include it in a URL
func urlencodeAddress(address *models.Address) string {
	s := []string{address.StreetAddress1}
	if address.StreetAddress2 != nil {
		s = append(s, *address.StreetAddress2)
	}
	if address.StreetAddress3 != nil {
		s = append(s, *address.StreetAddress3)
	}
	s = append(s, address.City, address.State, address.PostalCode)
	return url.QueryEscape(strings.Join(s, ","))
}

// Planner is the interface needed by Handlers to be able to evaluate the distance to be used for move accounting
type Planner interface {
	TransitDistance(source *models.Address, destination *models.Address) (int, error)
}
