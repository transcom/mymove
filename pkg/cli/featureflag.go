package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// FeatureFlagAccessCode is the access-code feature flag
	FeatureFlagAccessCode string = "feature-flag-access-code"
)

// InitFeatureFlags initializes FeatureFlags command line flags
func InitFeatureFlags(flag *pflag.FlagSet) {
	flag.Bool(FeatureFlagAccessCode, false, "Flag (bool) to enable requires-access-code")
}

// CheckFeatureFlag validates Verbose command line flags
func CheckFeatureFlag(v *viper.Viper) error {
	return nil
}
