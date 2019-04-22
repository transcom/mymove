package cli

import (
	"time"

	"github.com/spf13/pflag"
)

const (
	// BuildFlag is the Build Flag
	BuildFlag string = "build"
	// ConfigDirFlag is the Config Dir Flag
	ConfigDirFlag string = "config-dir"
	// InterfaceFlag is the Interface Flag
	InterfaceFlag string = "interface"
	// GracefulShutdownTimeoutFlag is the Graceful Shutdown Timeout Flag
	GracefulShutdownTimeoutFlag string = "graceful-shutdown-timeout"
)

// InitBuildFlags initializes the Build command line flags
func InitBuildFlags(flag *pflag.FlagSet) {
	flag.String(BuildFlag, "build", "the directory to serve static files from.")
	flag.String(ConfigDirFlag, "config", "The location of server config files")
	flag.String(InterfaceFlag, "", "The interface spec to listen for connections on. Default is all.")
	flag.Duration(GracefulShutdownTimeoutFlag, 25*time.Second, "The duration for which the server gracefully wait for existing connections to finish.  AWS ECS only gives you 30 seconds before sending SIGKILL.")
}
