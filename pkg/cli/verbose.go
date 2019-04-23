package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// VerboseFlag is the Verbose Flag
	VerboseFlag string = "debug-logging"
)

// InitVerboseFlags initializes Verbose command line flags
func InitVerboseFlags(flag *pflag.FlagSet) {
	flag.BoolP(VerboseFlag, "v", false, "log messages at the debug level.")
}

// CheckVerbose validates Verbose command line flags
func CheckVerbose(v *viper.Viper) error {
	return nil
}
