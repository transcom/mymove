package cli

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// SessionIdleTimeoutInMinutesFlag sets the session's Idle Timeout in minutes
	SessionIdleTimeoutInMinutesFlag string = "session-idle-timeout-in-minutes"
	// SessionLifetimeInHoursFlag sets the session's absolute expiry in hours
	SessionLifetimeInHoursFlag string = "session-lifetime-in-hours"
)

// InitSessionFlags initializes SessionFlags command line flags
func InitSessionFlags(flag *pflag.FlagSet) {
	flag.Int(SessionIdleTimeoutInMinutesFlag, 15, "Session idle timeout in minutes")
	flag.Int(SessionLifetimeInHoursFlag, 24, "Session absolute expiry in hours")
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
	timeout := v.GetInt(flagname)

	if v.GetString(EnvironmentFlag) == EnvironmentProd {
		if timeout < 15 || timeout > 60 {
			return errors.Errorf("%s must be an integer between 15 and 60, got %d", SessionIdleTimeoutInMinutesFlag, timeout)
		}
	}

	return nil
}

// ValidateSessionLifetime validates session lifetime
func ValidateSessionLifetime(v *viper.Viper, flagname string) error {
	lifetime := v.GetInt(flagname)

	if lifetime < 12 {
		return errors.Errorf("%s must be at least 12 hours, got %d", SessionLifetimeInHoursFlag, lifetime)
	}

	return nil
}
