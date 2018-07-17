package route

import (
	"net/url"
	"strings"

	"fmt"

	"github.com/transcom/mymove/pkg/models"
)

// LatLong is used to hold latitude and longitude as floats
type LatLong struct {
	Latitude  float32
	Longitude float32
}

// Coords returns a string with the comma separated co-ordinares, e.g "47.610,-122.107"
func (ll LatLong) Coords() string {
	return fmt.Sprintf("%f,%f", ll.Latitude, ll.Longitude)
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
	if address.Country != nil {
		s = append(s, *address.Country)
	}
	return url.QueryEscape(strings.Join(s, ","))
}

func zip5TransitDistanceHelper(planner Planner, source string, destination string) (int, error) {
	sLL, err := Zip5ToLatLong(source)
	if err != nil {
		return 0, err
	}
	dLL, err := Zip5ToLatLong(destination)
	if err != nil {
		return 0, err
	}
	return planner.LatLongTransitDistance(sLL, dLL)
}

// Planner is the interface needed by Handlers to be able to evaluate the distance to be used for move accounting
type Planner interface {
	TransitDistance(source *models.Address, destination *models.Address) (int, error)
	LatLongTransitDistance(source LatLong, destination LatLong) (int, error)
	Zip5TransitDistance(source string, destination string) (int, error)
}
