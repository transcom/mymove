package cli

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// SessionIdleTimeoutInMinutesFlag sets the session's Idle Timeout in minutes
	SessionIdleTimeoutInMinutesFlag string = "session-idle-timeout-in-minutes"
	// SessionLifetimeInHoursFlag sets the session's absolute expiry in hours
	SessionLifetimeInHoursFlag string = "session-lifetime-in-hours"

	// SessionIdleTimeoutInMinutes is the default idle timeout in minutes
	SessionIdleTimeoutInMinutes int = 15
	// SessionLifetimeInHours is the default session lifetime in hours
	SessionLifetimeInHours int = 24
)

// InitSessionFlags initializes SessionFlags command line flags
func InitSessionFlags(flag *pflag.FlagSet) {
	flag.Duration(SessionIdleTimeoutInMinutesFlag, (time.Duration(SessionIdleTimeoutInMinutes) * time.Minute), "Session idle timeout in minutes")
	flag.Duration(SessionLifetimeInHoursFlag, (time.Duration(SessionLifetimeInHours) * time.Hour), "Session absoluty expiry in hours")
}

// CheckSession validates session command line flags
func CheckSession(v *viper.Viper) error {
	return nil
}
