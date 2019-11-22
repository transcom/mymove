package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	//RequiresAccessCode flag
	FeatureFlagAccessCode    string = "feature-flag-access-code"
	FeatureFlagRoleBasedAuth string = "feature-flag-role-based-auth"
	// Enable or Disable CloudFront distribution
	CFEnableDistribution string = "enable-cf-distribution"
)

// InitFeatureFlag initializes FeatureFlags command line flags
func InitFeatureFlags(flag *pflag.FlagSet) {
	flag.Bool(FeatureFlagAccessCode, false, "Flag (bool) to enable requires-access-code")
	flag.Bool(FeatureFlagRoleBasedAuth, false, "Flag (bool) to enable role-based-auth")
	flag.Bool(CFEnableDistribution, false, "Flag (bool) to enable enable-cf-distribution")
}

// CheckFeatureFlag validates Verbose command line flags
func CheckFeatureFlag(v *viper.Viper) error {
	return nil
}
