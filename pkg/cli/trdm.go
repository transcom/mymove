package cli

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// TRDMApiURLFlag is the TRDM API URL Flag
	TRDMApiURLFlag string = "trdm-api-url"
	// TRDMApiWSDLFlag is the TRDM API WSDL Flag
	TRDMApiWSDLFlag string = "trdm-api-wsdl"
	// TRDMUseMockFlag is the TRDM Use Mock Flag
	TRDMUseMockFlag    string = "trdm-use-mock"
	TRDMx509Cert       string = "trdm-x509-cert"
	TRDMx509PrivateKey string = "trdm-x509-privatekey"
	TRDMIsEnabled      string = "trdm-is-enabled"
)

// InitTRDMFlags initializes Route command line flags
func InitTRDMFlags(flag *pflag.FlagSet) {
	flag.String(TRDMApiURLFlag, "", "URL for sending a SOAP request to TRDM")
	flag.String(TRDMApiWSDLFlag, "", "WSDL for sending a SOAP request to TRDM")
	flag.String(TRDMx509Cert, "", "x509 certificate for TRDM web services")
	flag.String(TRDMx509PrivateKey, "", "x509 private key for TRDM web services")

	flag.Bool(TRDMUseMockFlag, false, "Whether to use a mocked version of TRDM")
	flag.Bool(TRDMIsEnabled, false, "Enable TRDM SOAP requests")
}

// CheckRoute validates Route command line flags
func CheckTRDM(v *viper.Viper) error {
	urlVars := []string{
		TRDMApiURLFlag,
		TRDMApiWSDLFlag,
		TRDMx509Cert,
		TRDMx509PrivateKey,
		TRDMIsEnabled,
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
