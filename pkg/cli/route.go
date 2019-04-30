package cli

import (
	"net/http"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/route"
)

const (
	// HEREMapsGeocodeEndpointFlag is the HERE Maps Geocode Endpoint Flag
	HEREMapsGeocodeEndpointFlag string = "here-maps-geocode-endpoint"
	// HEREMapsRoutingEndpointFlag is the HERE Maps Routing Endpoint Flag
	HEREMapsRoutingEndpointFlag string = "here-maps-routing-endpoint"
	// HEREMapsAppIDFlag is the HERE Maps App ID Flag
	HEREMapsAppIDFlag string = "here-maps-app-id"
	// HEREMapsAppCodeFlag is the HERE Maps App Code Flag
	HEREMapsAppCodeFlag string = "here-maps-app-code"

	// hereRequestTimeout is how long to wait on HERE request before timing out (15 seconds).
	hereRequestTimeout = time.Duration(15) * time.Second
)

// InitRouteFlags initializes Route command line flags
func InitRouteFlags(flag *pflag.FlagSet) {
	flag.String(HEREMapsGeocodeEndpointFlag, "", "URL for the HERE maps geocode endpoint")
	flag.String(HEREMapsRoutingEndpointFlag, "", "URL for the HERE maps routing endpoint")
	flag.String(HEREMapsAppIDFlag, "", "HERE maps App ID for this application")
	flag.String(HEREMapsAppCodeFlag, "", "HERE maps App API code")
}

// InitRoutePlanner validates Route Planner command line flags
func InitRoutePlanner(v *viper.Viper, logger Logger) route.Planner {
	hereClient := &http.Client{Timeout: hereRequestTimeout}
	return route.NewHEREPlanner(
		logger,
		hereClient,
		v.GetString(HEREMapsGeocodeEndpointFlag),
		v.GetString(HEREMapsRoutingEndpointFlag),
		v.GetString(HEREMapsAppIDFlag),
		v.GetString(HEREMapsAppCodeFlag))
}

// CheckRoute validates Route command line flags
func CheckRoute(v *viper.Viper) error {
	urlVars := []string{
		HEREMapsGeocodeEndpointFlag,
		HEREMapsRoutingEndpointFlag,
	}

	for _, c := range urlVars {
		err := ValidateURL(v, c)
		if err != nil {
			return err
		}
	}
	return nil
}
