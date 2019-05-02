package cli

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// DevlocalAuthFlag is the Devlocal Auth Flag
	DevlocalAuthFlag string = "devlocal-auth"
)

// InitDevlocalFlags initializes the Devlocal command line flags
func InitDevlocalFlags(flag *pflag.FlagSet) {
	flag.Bool(DevlocalAuthFlag, false, "Enable the devlocal auth system for logging in without Login.gov.")
}

// CheckDevlocal validates the Devlocal command line flags
func CheckDevlocal(v *viper.Viper) error {
	environment := v.GetString(EnvironmentFlag)
	allowedEnvironments := []string{EnvironmentExperimental, EnvironmentTest, EnvironmentDevlocal}
	if devlocalAuthEnabled := v.GetBool(DevlocalAuthFlag); devlocalAuthEnabled && !stringSliceContains(allowedEnvironments, environment) {
		return errors.Errorf("Devlocal Auth cannot run in the '%s' environment, only in %v", environment, allowedEnvironments)
	}
	return nil
}
