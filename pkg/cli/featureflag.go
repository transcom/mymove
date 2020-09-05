package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// FeatureFlagAccessCode is the access-code feature flag
	FeatureFlagAccessCode string = "feature-flag-access-code"
	// FeatureFlagConvertProfileOrdersToGHC is the convert-ppms-to-ghc feature flag
	FeatureFlagConvertProfileOrdersToGHC string = "feature-flag-convert-profile-orders-to-ghc"
)

// InitFeatureFlags initializes FeatureFlags command line flags
func InitFeatureFlags(flag *pflag.FlagSet) {
	flag.Bool(FeatureFlagAccessCode, false, "Flag (bool) to enable requires-access-code")
	flag.Bool(FeatureFlagConvertProfileOrdersToGHC, false, "Flag (bool) to enable convert-profile-orders-to-ghc")
}

// CheckFeatureFlag validates Verbose command line flags
func CheckFeatureFlag(v *viper.Viper) error {
	return nil
}
