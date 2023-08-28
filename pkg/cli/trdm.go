package cli

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// TRDMApiURLFlag is the TRDM API URL Flag
	TRDMApiURLFlag string = "TRDM-api-url"
	// TRDMApiWSDLFlag is the TRDM API WSDL Flag
	TRDMApiWSDLFlag string = "TRDM-api-wsdl"

	// TRDMUseMockFlag is the TRDM Use Mock Flag
	TRDMUseMockFlag string = "TRDM-use-mock"
)

// InitTRDMFlags initializes Route command line flags
func InitTRDMFlags(flag *pflag.FlagSet) {
	flag.String(TRDMApiURLFlag, "", "URL for sending a SOAP request to TRDM")
	flag.String(TRDMApiWSDLFlag, "", "WSDL for sending a SOAP request to TRDM")

	flag.Bool(TRDMUseMockFlag, false, "Whether to use a mocked version of TRDM")
}

// CheckRoute validates Route command line flags
func CheckTRDM(v *viper.Viper) error {
	urlVars := []string{
		TRDMApiURLFlag,
		TRDMApiWSDLFlag,
	}

	for _, c := range urlVars {
		err := ValidateURL(v, c)
		if err != nil {
			return err
		}
	}

	if v.GetBool(TRDMUseMockFlag) {
		// Check against the Environment
		allowedEnvironments := []string{
			EnvironmentDevelopment,
			EnvironmentTest,
			EnvironmentExp,
			EnvironmentExperimental,
			EnvironmentDemo,
			EnvironmentLoadtest,
			EnvironmentReview,
			EnvironmentStg, // WARN: This is enabled only while the TRDM service is down.
		}
		if environment := v.GetString(EnvironmentFlag); !stringSliceContains(allowedEnvironments, environment) {
			return errors.Errorf("cannot mock TRDM with the '%s' environment, only in %v", environment, allowedEnvironments)
		}
	}

	return nil
}
