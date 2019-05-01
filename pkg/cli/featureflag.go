package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	//NewPPMFlow flag
	NewPPMFlow string = "new-ppm-flow"
	//RequireAccessCode flag
	RequireAccessCode string = "require-access-code"
)

// InitFeatureFlag initializes FeatureFlags command line flags
func InitFeatureFlag(flag *pflag.FlagSet) {
	flag.Bool(NewPPMFlow, false, "Flag (bool) to enable the new-ppm-flow")
	flag.Bool(RequireAccessCode, false, "Flag (bool) to enable the require-access-code")
}

// CheckFeatureFlag validates Verbose command line flags
func CheckFeatureFlag(v *viper.Viper) error {
	return nil
}
