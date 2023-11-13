package cli

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// TRDMApiGatewayURLFlag is the TRDM API Gateway URL Flag
	TRDMApiGatewayURLFlag string = "trdm-api-gateway-url"
	// TRDMUseMockFlag is the TRDM Use Mock Flag
	TRDMUseMockFlag string = "trdm-use-mock"
	// FF to enable or disable TRDM soap requests
	TRDMIsEnabledFlag string = "trdm-is-enabled"
)

// InitTRDMFlags initializes Route command line flags
func InitTRDMFlags(flag *pflag.FlagSet) {
	flag.String(TRDMApiGatewayURLFlag, "", "URL for sending a REST request to the TRDM gateway")

	flag.Bool(TRDMUseMockFlag, false, "Whether to use a mocked version of TRDM")
	flag.Bool(TRDMIsEnabledFlag, false, "Enable TRDM data requests")
}

// CheckRoute validates Route command line flags
func CheckTRDM(v *viper.Viper) error {
	urlVars := []string{
		TRDMApiGatewayURLFlag,
		TRDMIsEnabledFlag,
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
