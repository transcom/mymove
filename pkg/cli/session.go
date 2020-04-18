package cli

import (
	"time"

	"github.com/pkg/errors"
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
	flag.Duration(SessionLifetimeInHoursFlag, (time.Duration(SessionLifetimeInHours) * time.Hour), "Session absolute expiry in hours")
}

// CheckSession validates session command line flags
func CheckSession(v *viper.Viper) error {
	if err := ValidateSessionTimeout(v, SessionIdleTimeoutInMinutesFlag); err != nil {
		return err
	}

	if err := ValidateSessionLifetime(v, SessionLifetimeInHoursFlag); err != nil {
		return err
	}

	return nil
}

// ValidateSessionTimeout validates session idle timeout
func ValidateSessionTimeout(v *viper.Viper, flagname string) error {
	environment := v.GetString(EnvironmentFlag)
	timeout := v.GetDuration(flagname)

	if environment == EnvironmentProd && (timeout < 15 || timeout > 60) {
		return errors.Errorf("%s must be an integer between 15 and 60", SessionIdleTimeoutInMinutesFlag)
	}

	return nil
}

// ValidateSessionLifetime validates session lifetime
func ValidateSessionLifetime(v *viper.Viper, flagname string) error {
	environment := v.GetString(EnvironmentFlag)
	lifetime := v.GetDuration(flagname)

	if environment == EnvironmentProd && lifetime < 12 {
		return errors.Errorf("%s must be at least 12 hours in production", SessionLifetimeInHoursFlag)
	}

	if lifetime < 1 {
		return errors.Errorf("%s must be at least 1", SessionLifetimeInHoursFlag)
	}

	return nil
}
