package route

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/tiaguinho/gosoap"

	"github.com/transcom/mymove/pkg/appcontext"
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
func zip5TransitDistanceHelper(appCtx appcontext.AppContext, planner Planner, source string, destination string) (int, error) {
	sLL, err := Zip5ToLatLong(source)
	if err != nil {
		return 0, err
	}
	dLL, err := Zip5ToLatLong(destination)
	if err != nil {
		return 0, err
	}
	distance, err := planner.LatLongTransitDistance(appCtx, sLL, dLL)
	if err != nil {
		return 0, err
	}
	return distance, err
}

// zip5TransitDistanceHelper takes a source and destination zip5 and calculates the distance between them using a Zip5 to LatLong lookup and will throw an error if distance is less than 50, this is used by PPM code
// Ideally I don't think we should check for minimum distance here and should refactor code to use zip5TransitDistanceHelper over this helper over time.
func zip5TransitDistanceLineHaulHelper(appCtx appcontext.AppContext, planner Planner, source string, destination string) (int, error) {
	sLL, err := Zip5ToLatLong(source)
	if err != nil {
		return 0, err
	}
	dLL, err := Zip5ToLatLong(destination)
	if err != nil {
		return 0, err
	}
	distance, err := planner.LatLongTransitDistance(appCtx, sLL, dLL)
	if err != nil {
		return 0, err
	}
	if distance < 50 {
		err = NewShortHaulError(sLL, dLL, distance)
	}
	return distance, err
}

// zip3TransitDistanceHelper takes a source and destination zip3 and calculates the distence between them using a Zip3 to LatLong lookup, this is intended for HHG long haul calculations with two differnet zip3s
func zip3TransitDistanceHelper(appCtx appcontext.AppContext, planner Planner, source string, destination string) (int, error) {
	sLL, err := Zip5ToZip3LatLong(source)
	if err != nil {
		return 0, err
	}
	dLL, err := Zip5ToZip3LatLong(destination)
	if err != nil {
		return 0, err
	}
	distance, err := planner.LatLongTransitDistance(appCtx, sLL, dLL)
	if err != nil {
		return 0, err
	}
	return distance, err
}

// SoapCaller provides an interface for the Call method of the gosoap Client so it can be mocked.
// NOTE: Placing this in a separate package/directory to avoid a circular dependency from an existing mock.
//go:generate mockery --name SoapCaller --outpkg ghcmocks --output ./ghcmocks --disable-version-string
type SoapCaller interface {
	Call(m string, p gosoap.SoapParams) (res *gosoap.Response, err error)
}

// Planner is the interface needed by Handlers to be able to evaluate the distance to be used for move accounting
//go:generate mockery --name Planner --disable-version-string
type Planner interface {
	TransitDistance(appCtx appcontext.AppContext, source *models.Address, destination *models.Address) (int, error)
	LatLongTransitDistance(appCtx appcontext.AppContext, source LatLong, destination LatLong) (int, error)
	// Zip5TransitDistanceLineHaul is used by PPM flow and checks for minimum distance restriction as PPM doesn't allow short hauls
	// New code should probably make the minimum checks after calling Zip5TransitDistance over using this method
	Zip5TransitDistanceLineHaul(appCtx appcontext.AppContext, source string, destination string) (int, error)
	ZipTransitDistance(appCtx appcontext.AppContext, source string, destination string) (int, error)
	Zip3TransitDistance(appCtx appcontext.AppContext, source string, destination string) (int, error)
	Zip5TransitDistance(appCtx appcontext.AppContext, source string, destination string) (int, error)
}

// InitRoutePlanner creates a new HERE route planner that adheres to the Planner interface
func InitRoutePlanner(v *viper.Viper) Planner {
	hereClient := &http.Client{Timeout: hereRequestTimeout}
	return NewHEREPlanner(
		hereClient,
		v.GetString(cli.HEREMapsGeocodeEndpointFlag),
		v.GetString(cli.HEREMapsRoutingEndpointFlag),
		v.GetString(cli.HEREMapsAppIDFlag),
		v.GetString(cli.HEREMapsAppCodeFlag))
}

// InitHHGRoutePlanner creates a new HHG route planner that adheres to the Planner interface
func InitHHGRoutePlanner(appCtx appcontext.AppContext, v *viper.Viper, tlsConfig *tls.Config) (Planner, error) {
	dtodPlannerMileage, err := initDTODPlannerMileage(appCtx, v, tlsConfig, "HHG")
	if err != nil {
		return nil, err
	}

	return NewHHGPlanner(dtodPlannerMileage), nil
}

// InitDTODRoutePlanner creates a new DTOD route planner that adheres to the Planner interface
func InitDTODRoutePlanner(appCtx appcontext.AppContext, v *viper.Viper, tlsConfig *tls.Config) (Planner, error) {
	dtodPlannerMileage, err := initDTODPlannerMileage(appCtx, v, tlsConfig, "DTOD")
	if err != nil {
		return nil, err
	}

	return NewDTODPlanner(dtodPlannerMileage), nil
}

func initDTODPlannerMileage(appCtx appcontext.AppContext, v *viper.Viper, tlsConfig *tls.Config, plannerType string) (DTODPlannerMileage, error) {
	dtodUseMock := v.GetBool(cli.DTODUseMockFlag)

	var dtodPlannerMileage DTODPlannerMileage
	if dtodUseMock {
		appCtx.Logger().Info(fmt.Sprintf("Using mocked DTOD for %s route planner", plannerType))
		dtodPlannerMileage = NewMockDTODZip5Distance()
	} else {
		appCtx.Logger().Info(fmt.Sprintf("Using real DTOD for %s route planner", plannerType))
		tr := &http.Transport{TLSClientConfig: tlsConfig}
		httpClient := &http.Client{Transport: tr, Timeout: time.Duration(30) * time.Second}

		dtodWSDL := v.GetString(cli.DTODApiWSDLFlag)
		dtodURL := v.GetString(cli.DTODApiURLFlag)
		dtodAPIUsername := v.GetString(cli.DTODApiUsernameFlag)
		dtodAPIPassword := v.GetString(cli.DTODApiPasswordFlag)

		soapClient, err := gosoap.SoapClient(dtodWSDL, httpClient)
		if err != nil {
			return nil, fmt.Errorf("unable to create SOAP client: %w", err)
		}
		soapClient.URL = dtodURL

		dtodPlannerMileage = NewDTODZip5Distance(dtodAPIUsername, dtodAPIPassword, soapClient)
	}

	return dtodPlannerMileage, nil
}
