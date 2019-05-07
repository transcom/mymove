package cli

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// EnvironmentFlag is the Environment Flag
	EnvironmentFlag string = "environment"

	// EnvironmentProd is the Production Environment name
	EnvironmentProd string = "prod"
	// EnvironmentStaging is the Staging Environment name
	EnvironmentStaging string = "staging"
	// EnvironmentExperimental is the Experimental Environment name
	EnvironmentExperimental string = "experimental"
	// EnvironmentTest is the Test Environment name
	EnvironmentTest string = "test"
	// EnvironmentDevelopment is the Development Environment name
	EnvironmentDevelopment string = "development"
)

var environments = []string{
	EnvironmentProd,
	EnvironmentStaging,
	EnvironmentExperimental,
	EnvironmentTest,
	EnvironmentDevelopment,
}

type errInvalidEnvironment struct {
	Environment string
}

func (e *errInvalidEnvironment) Error() string {
	return fmt.Sprintf("invalid environment %q, expecting one of %q", e.Environment, environments)
}

// InitEnvironmentFlags initializes the Environment command line flags
func InitEnvironmentFlags(flag *pflag.FlagSet) {
	flag.StringP(EnvironmentFlag, "e", EnvironmentDevelopment, fmt.Sprintf("The environment name, one of %v", environments))
}

// CheckEnvironment validates the Environment command line flags
func CheckEnvironment(v *viper.Viper) error {
	if environment := v.GetString(EnvironmentFlag); !stringSliceContains(environments, environment) {
		return fmt.Errorf("invalid environment %s, expecting one of %q", environment, environments)
	}
	return nil
}
