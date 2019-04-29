package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// FeatureFlag is the Feature Flag
	FeatureFlag string = "feature-flag"
)

// InitFeatureFlag initializes FeatureFlags command line flags
func InitFeatureFlag(flag *pflag.FlagSet) {
	flag.BoolP(FeatureFlag, "f", false, "Flag (bool) to enable the feature-flag")
}

// CheckFeatureFlag validates Verbose command line flags
func CheckFeatureFlag(v *viper.Viper) error {
	return nil
}
