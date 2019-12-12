package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	//RequiresAccessCode flag
	FeatureFlagAccessCode    string = "feature-flag-access-code"
	FeatureFlagRoleBasedAuth string = "feature-flag-role-based-auth"
)

// InitFeatureFlag initializes FeatureFlags command line flags
func InitFeatureFlags(flag *pflag.FlagSet) {
	flag.Bool(FeatureFlagAccessCode, false, "Flag (bool) to enable requires-access-code")
	flag.Bool(FeatureFlagRoleBasedAuth, false, "Flag (bool) to enable role-based-auth")
}

// CheckFeatureFlag validates Verbose command line flags
func CheckFeatureFlag(v *viper.Viper) error {
	return nil
}
