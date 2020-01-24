package cli

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// DbRetryIntervalFlag is the DB retry interval flag
	DbRetryIntervalFlag string = "db-retry-interval"
	// DbRetryMaxFlag is the DB retry maximum flag
	DbRetryMaxFlag string = "db-retry-max"
)

// InitDatabaseRetryFlags initializes Database Retry command line flags
func InitDatabaseRetryFlags(flag *pflag.FlagSet) {
	flag.Duration(DbRetryIntervalFlag, time.Second*5, "Database retry interval duration")
	flag.Int(DbRetryMaxFlag, 5, "Database maximum retries before connection failure")
}

// CheckDatabaseRetry validates Database Retry command line flags
func CheckDatabaseRetry(v *viper.Viper) error {
	if retryInterval := v.GetDuration(DbRetryIntervalFlag); retryInterval < 1*time.Second {
		return fmt.Errorf("retry interval must be greater than 1 seconds, got %q", retryInterval)
	}

	if retryMax := v.GetInt(DbRetryMaxFlag); retryMax < 0 {
		return fmt.Errorf("retries must be greater than 0, got %q", retryMax)
	}

	return nil
}
