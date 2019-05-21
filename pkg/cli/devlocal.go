package cli

import (
	"fmt"

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
	if v.GetBool(DevlocalAuthFlag) {
		// Check against the Environment
		allowedEnvironments := []string{
			EnvironmentDevelopment,
			EnvironmentTest,
			EnvironmentExperimental,
		}
		if environment := v.GetString(EnvironmentFlag); !stringSliceContains(allowedEnvironments, environment) {
			return errors.Errorf("Devlocal Auth cannot run in the '%s' environment, only in %v", environment, allowedEnvironments)
		}

		// Check against My Server Names
		allowedMyServerNames := []string{
			HTTPMyServerNameLocal,
			fmt.Sprintf("my.%s.move.mil", EnvironmentExperimental),
		}
		if serverName := v.GetString(HTTPMyServerNameFlag); !stringSliceContains(allowedMyServerNames, serverName) {
			return errors.Errorf("Devlocal Auth cannot run with the '%s' my server name, only in %v", serverName, allowedMyServerNames)
		}

		// Check against Office Server Names
		allowedOfficeServerNames := []string{
			HTTPOfficeServerNameLocal,
			fmt.Sprintf("office.%s.move.mil", EnvironmentExperimental),
		}
		if serverName := v.GetString(HTTPOfficeServerNameFlag); !stringSliceContains(allowedOfficeServerNames, serverName) {
			return errors.Errorf("Devlocal Auth cannot run with the '%s' office server name, only in %v", serverName, allowedOfficeServerNames)
		}

		// Check against TSP Server Names
		allowedTSPServerNames := []string{
			HTTPTSPServerNameLocal,
			fmt.Sprintf("tsp.%s.move.mil", EnvironmentExperimental),
		}
		if serverName := v.GetString(HTTPTSPServerNameFlag); !stringSliceContains(allowedTSPServerNames, serverName) {
			return errors.Errorf("Devlocal Auth cannot run with the '%s' tsp server name, only in %v", serverName, allowedTSPServerNames)
		}

		// Check against Admin Server Names
		allowedAdminServerNames := []string{
			HTTPAdminServerNameLocal,
			fmt.Sprintf("admin.%s.move.mil", EnvironmentExperimental),
		}
		if serverName := v.GetString(HTTPAdminServerNameFlag); !stringSliceContains(allowedAdminServerNames, serverName) {
			return errors.Errorf("Devlocal Auth cannot run with the '%s' admin server name, only in %v", serverName, allowedAdminServerNames)
		}
	}
	return nil
}
