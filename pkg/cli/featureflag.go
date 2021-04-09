package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// FeatureFlagAccessCode determines whether or not service members are prompted for an access code before they start onboarding
	FeatureFlagAccessCode string = "feature-flag-access-code"
	// FeatureFlagServiceCounseling controls whether moves get routed to service counseling
	FeatureFlagServiceCounseling string = "feature-flag-service-counseling"
)

// InitFeatureFlags initializes FeatureFlags command line flags
func InitFeatureFlags(flag *pflag.FlagSet) {
	flag.Bool(FeatureFlagAccessCode, false, "Flag (bool) to enable requires-access-code")
	flag.Bool(FeatureFlagServiceCounseling, false, "Flag (bool) to enable service counseling")
}

// CheckFeatureFlag validates Verbose command line flags
func CheckFeatureFlag(v *viper.Viper) error {
	return nil
}
