package cli

import (
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/spf13/pflag"
)

const (
	// InterfaceFlag is the Interface Flag
	InterfaceFlag string = "interface"
	// GracefulShutdownTimeoutFlag is the Graceful Shutdown Timeout Flag
	GracefulShutdownTimeoutFlag string = "graceful-shutdown-timeout"

	// The default graceful shutdown duration
	DefaultGracefulShutdownDuration = time.Second * 25

	// The minimum graceful shutdown duration
	MinimumGracefulShutdownDuration = time.Second * 5
)

// InitWebserverFlags initializes the webserver command line flags
func InitWebserverFlags(flag *pflag.FlagSet) {
	flag.String(InterfaceFlag, "", "The interface spec to listen for connections on. Default of empty string means all interfaces. Accepts 'localhost' or IPv4 addresses as well.")
	flag.Duration(GracefulShutdownTimeoutFlag, DefaultGracefulShutdownDuration, "The duration for which the server gracefully wait for existing connections to finish.  AWS ECS only gives you 30 seconds before sending SIGKILL.")
}

// CheckWebserver validates the webserver command line flags
func CheckWebserver(v *viper.Viper) error {
	if str := v.GetString(InterfaceFlag); len(str) > 0 && str != "localhost" {
		addr, err := net.ResolveIPAddr("tcp", str)
		if err != nil {
			return errors.Errorf("Unable to resolve IP address %q", str)
		}
		if addr.IP.To4() == nil {
			return errors.Errorf("Expected IPv4 address, got %q", str)
		}
	}
	if d := v.GetDuration(GracefulShutdownTimeoutFlag); d < MinimumGracefulShutdownDuration {
		return errors.Errorf("Graceful Shutdown Duration should not be less than 5 Seconds. Provided duration %q", d)
	}
	return nil
}
