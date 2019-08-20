package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
)

// InitRouteFlags initializes Route command line flags
func InitRouteFlags(flag *pflag.FlagSet) {
	flag.String(HEREMapsGeocodeEndpointFlag, "", "URL for the HERE maps geocode endpoint")
	flag.String(HEREMapsRoutingEndpointFlag, "", "URL for the HERE maps routing endpoint")
	flag.String(HEREMapsAppIDFlag, "", "HERE maps App ID for this application")
	flag.String(HEREMapsAppCodeFlag, "", "HERE maps App API code")
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
