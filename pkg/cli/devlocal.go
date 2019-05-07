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
	if devlocalAuthEnabled := v.GetBool(DevlocalAuthFlag); devlocalAuthEnabled {
		// Check against the DB Environment
		allowedDBEnvironments := []string{
			EnvironmentExperimental,
			EnvironmentTest,
			EnvironmentDevlocal,
		}
		if !stringSliceContains(allowedDBEnvironments, environment) {
			return errors.Errorf("Devlocal Auth cannot run in the '%s' environment, only in %v", environment, allowedDBEnvironments)
		}

		// Check against My Server Names
		allowedMyServerNames := []string{
			"milmovelocal",
			"my.test.move.mil",
			"my.experimental.move.mil",
		}
		if serverName := v.GetString(HTTPMyServerNameFlag); !stringSliceContains(allowedMyServerNames, serverName) {
			return errors.Errorf("Devlocal Auth cannot run in the '%s' server name, only in %v", serverName, allowedMyServerNames)
		}

		// Check against Office Server Names
		allowedOfficeServerNames := []string{
			"officelocal",
			"office.test.move.mil",
			"office.experimental.move.mil",
		}
		if serverName := v.GetString(HTTPOfficeServerNameFlag); !stringSliceContains(allowedOfficeServerNames, serverName) {
			return errors.Errorf("Devlocal Auth cannot run in the '%s' server name, only in %v", serverName, allowedOfficeServerNames)
		}

		// Check against TSP Server Names
		allowedTSPServerNames := []string{
			"tsplocal",
			"tsp.test.move.mil",
			"tsp.experimental.move.mil",
		}
		if serverName := v.GetString(HTTPTSPServerNameFlag); !stringSliceContains(allowedTSPServerNames, serverName) {
			return errors.Errorf("Devlocal Auth cannot run in the '%s' server name, only in %v", serverName, allowedTSPServerNames)
		}

		// Check against Admin Server Names
		allowedAdminServerNames := []string{
			"adminlocal",
			"admin.test.move.mil",
			"admin.experimental.move.mil",
		}
		if serverName := v.GetString(HTTPAdminServerNameFlag); !stringSliceContains(allowedAdminServerNames, serverName) {
			return errors.Errorf("Devlocal Auth cannot run in the '%s' server name, only in %v", serverName, allowedAdminServerNames)
		}
	}
	return nil
}
