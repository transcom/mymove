package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	FeatureFlagAccessCode       string = "feature-flag-access-code"
	FeatureFlagRoleBasedAuth    string = "feature-flag-role-based-auth"
	FeatureFlagConvertPPMsToGHC string = "feature-flag-convert-ppms-to-ghc"
)

// InitFeatureFlag initializes FeatureFlags command line flags
func InitFeatureFlags(flag *pflag.FlagSet) {
	flag.Bool(FeatureFlagAccessCode, false, "Flag (bool) to enable requires-access-code")
	flag.Bool(FeatureFlagRoleBasedAuth, false, "Flag (bool) to enable role-based-auth")
	flag.Bool(FeatureFlagConvertPPMsToGHC, false, "Flag (bool) to enable convert-ppms-to-ghc")
}

// CheckFeatureFlag validates Verbose command line flags
func CheckFeatureFlag(v *viper.Viper) error {
	return nil
}
