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

	// DTODApiUsernameFlag is the DTOD API Username Flag
	DTODApiUsernameFlag string = "dtod-api-username"
	// DTODApiPasswordFlag is the DTOD API Password Flag
	DTODApiPasswordFlag string = "dtod-api-password"
	// DTODApiURLFlag is the DTOD API URL Flag
	DTODApiURLFlag string = "dtod-api-url"
	// DTODApiWSDLFlag is the DTOD API WSDL Flag
	DTODApiWSDLFlag string = "dtod-api-wsdl"

	// DTODUseMock is the DTOD Use Mock Flag
	DTODUseMock string = "dtod-use-mock"
)

// InitRouteFlags initializes Route command line flags
func InitRouteFlags(flag *pflag.FlagSet) {
	flag.String(HEREMapsGeocodeEndpointFlag, "", "URL for the HERE maps geocode endpoint")
	flag.String(HEREMapsRoutingEndpointFlag, "", "URL for the HERE maps routing endpoint")
	flag.String(HEREMapsAppIDFlag, "", "HERE maps App ID for this application")
	flag.String(HEREMapsAppCodeFlag, "", "HERE maps App API code")

	flag.String(DTODApiUsernameFlag, "", "DTOD api auth username")
	flag.String(DTODApiPasswordFlag, "", "DTOD api auth password")
	flag.String(DTODApiURLFlag, "", "URL for sending a SOAP request to DTOD")
	flag.String(DTODApiWSDLFlag, "", "WSDL for sending a SOAP request to DTOD")

	flag.Bool(DTODUseMock, false, "Whether to use a mocked version of DTOD")
}

// CheckRoute validates Route command line flags
func CheckRoute(v *viper.Viper) error {
	urlVars := []string{
		HEREMapsGeocodeEndpointFlag,
		HEREMapsRoutingEndpointFlag,
		DTODApiURLFlag,
		DTODApiWSDLFlag,
	}

	for _, c := range urlVars {
		err := ValidateURL(v, c)
		if err != nil {
			return err
		}
	}

	// TODO: Removing this check for now to see how Circle reacts.
	//if len(v.GetString(DTODApiUsernameFlag)) == 0 {
	//	return errors.Errorf("%s is missing", DTODApiUsernameFlag)
	//}
	//if len(v.GetString(DTODApiPasswordFlag)) == 0 {
	//	return errors.Errorf("%s is missing", DTODApiPasswordFlag)
	//}

	return nil
}
