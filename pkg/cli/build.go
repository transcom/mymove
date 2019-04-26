package cli

import (
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/spf13/pflag"
)

const (
	// BuildFlag is the Build Flag
	BuildFlag string = "build"
	// InterfaceFlag is the Interface Flag
	InterfaceFlag string = "interface"
	// GracefulShutdownTimeoutFlag is the Graceful Shutdown Timeout Flag
	GracefulShutdownTimeoutFlag string = "graceful-shutdown-timeout"
	// LogTaskMetadataFlag is the Log Task Metadata Flag
	LogTaskMetadataFlag string = "log-task-metadata"
)

// InitBuildFlags initializes the Build command line flags
func InitBuildFlags(flag *pflag.FlagSet) {
	flag.String(BuildFlag, "build", "the directory to serve static files from.")
	flag.String(InterfaceFlag, "", "The interface spec to listen for connections on. Default of empty string means all interfaces. Accepts 'localhost' or IPv4 addresses as well.")
	flag.Duration(GracefulShutdownTimeoutFlag, 25*time.Second, "The duration for which the server gracefully wait for existing connections to finish.  AWS ECS only gives you 30 seconds before sending SIGKILL.")
	flag.Bool(LogTaskMetadataFlag, false, "Fetch AWS Task Metadata and add to log.")
}

// CheckBuild validates the Build command line flags
func CheckBuild(v *viper.Viper) error {
	if buildDir := v.GetString(BuildFlag); len(buildDir) == 0 {
		return errors.Errorf("Build directory must not be empty")
	}

	iface := v.GetString(InterfaceFlag)
	if !(iface == "localhost" || iface == "") {
		addr, err := net.ResolveIPAddr("tcp", iface)
		if err != nil {
			return errors.Errorf("Unable to resolve IP address %s", iface)
		}
		if addr.IP.To4() == nil {
			return errors.Errorf("Expected IPv4 address, got %s", iface)
		}
	}

	if gracefulShutdownDuration := v.GetDuration(GracefulShutdownTimeoutFlag); gracefulShutdownDuration < 5*time.Second {
		return errors.Errorf("Graceful Shutdown Duration should not be less than 5 Seconds. Provided duration %q", gracefulShutdownDuration)
	}
	return nil
}
