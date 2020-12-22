package route

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/models"
)

const (
	// hereRequestTimeout is how long to wait on HERE request before timing out (15 seconds).
	hereRequestTimeout = time.Duration(15) * time.Second
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

// zip5TransitDistanceHelper takes a source and destination zip5 and calculates the distance between them using a Zip5 to LatLong lookup, this is needed to support HHG short haul distance lookups
func zip5TransitDistanceHelper(planner Planner, source string, destination string) (int, error) {
	sLL, err := Zip5ToLatLong(source)
	if err != nil {
		return 0, err
	}
	dLL, err := Zip5ToLatLong(destination)
	if err != nil {
		return 0, err
	}
	distance, err := planner.LatLongTransitDistance(sLL, dLL)
	if err != nil {
		return 0, err
	}
	return distance, err
}

// zip5TransitDistanceHelper takes a source and destination zip5 and calculates the distance between them using a Zip5 to LatLong lookup and will throw an error if distance is less than 50, this is used by PPM code
// Ideally I don't think we should check for minimum distance here and should refactor code to use zip5TransitDistanceHelper over this helper over time.
func zip5TransitDistanceLineHaulHelper(planner Planner, source string, destination string) (int, error) {
	sLL, err := Zip5ToLatLong(source)
	if err != nil {
		return 0, err
	}
	dLL, err := Zip5ToLatLong(destination)
	if err != nil {
		return 0, err
	}
	distance, err := planner.LatLongTransitDistance(sLL, dLL)
	if err != nil {
		return 0, err
	}
	if distance < 50 {
		err = NewShortHaulError(sLL, dLL, distance)
	}
	return distance, err
}

// zip3TransitDistanceHelper takes a source and destination zip3 and calculates the distence between them using a Zip3 to LatLong lookup, this is intended for HHG long haul calculations with two differnet zip3s
func zip3TransitDistanceHelper(planner Planner, source string, destination string) (int, error) {
	sLL, err := Zip5ToZip3LatLong(source)
	if err != nil {
		return 0, err
	}
	dLL, err := Zip5ToZip3LatLong(destination)
	if err != nil {
		return 0, err
	}
	distance, err := planner.LatLongTransitDistance(sLL, dLL)
	if err != nil {
		return 0, err
	}
	return distance, err
}

// Planner is the interface needed by Handlers to be able to evaluate the distance to be used for move accounting
//go:generate mockery -name Planner
type Planner interface {
	TransitDistance(source *models.Address, destination *models.Address) (int, error)
	LatLongTransitDistance(source LatLong, destination LatLong) (int, error)
	// Zip5TransitDistanceLineHaul is used by PPM flow and checks for minimum distance restriction as PPM doesn't allow short hauls
	// New code should probably make the minimum checks after calling Zip5TransitDistance over using this method
	Zip5TransitDistanceLineHaul(source string, destination string) (int, error)
	Zip5TransitDistance(source string, destination string) (int, error)
	Zip3TransitDistance(source string, destination string) (int, error)
}

// InitRoutePlanner creates a new HERE route planner that adheres to the Planner interface
func InitRoutePlanner(v *viper.Viper, logger Logger) Planner {
	hereClient := &http.Client{Timeout: hereRequestTimeout}
	return NewHEREPlanner(
		logger,
		hereClient,
		v.GetString(cli.HEREMapsGeocodeEndpointFlag),
		v.GetString(cli.HEREMapsRoutingEndpointFlag),
		v.GetString(cli.HEREMapsAppIDFlag),
		v.GetString(cli.HEREMapsAppCodeFlag))
}

// InitGHCRoutePlanner creates a new GHC route planner that adheres to the Planner interface
func InitGHCRoutePlanner(v *viper.Viper, db *pop.Connection, logger Logger) Planner {
	return NewGHCPlanner(db, logger)
}
