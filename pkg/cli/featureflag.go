package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// FeatureFlagAccessCode is the access-code feature flag
	FeatureFlagAccessCode string = "feature-flag-access-code"
	// FeatureFlagRoleBasedAuth is the role-based-auth feature flag
	FeatureFlagRoleBasedAuth string = "feature-flag-role-based-auth"
	// FeatureFlagConvertPPMsToGHC is the convert-ppms-to-ghc feature flag
	FeatureFlagConvertPPMsToGHC string = "feature-flag-convert-ppms-to-ghc"
	// FeatureFlagSupportEndpoints is the support-endpoints feature flag
	FeatureFlagSupportEndpoints = "feature-flag-support-endpoints"
)

// InitFeatureFlags initializes FeatureFlags command line flags
func InitFeatureFlags(flag *pflag.FlagSet) {
	flag.Bool(FeatureFlagAccessCode, false, "Flag (bool) to enable requires-access-code")
	flag.Bool(FeatureFlagRoleBasedAuth, false, "Flag (bool) to enable role-based-auth")
	flag.Bool(FeatureFlagConvertPPMsToGHC, false, "Flag (bool) to enable convert-ppms-to-ghc")
	flag.Bool(FeatureFlagSupportEndpoints, false, "Flag (bool) to enable support-endpoints")
}

// CheckFeatureFlag validates Verbose command line flags
func CheckFeatureFlag(v *viper.Viper) error {
	return nil
}
