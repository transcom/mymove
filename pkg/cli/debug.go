package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	//DebugPProfFlag enables the pprof debugging endpoints
	DebugPProfFlag string = "debug-pprof"
)

// InitDebugFlags initializes the Debug command line flags
func InitDebugFlags(flag *pflag.FlagSet) {
	flag.Bool(DebugPProfFlag, false, "Enables the go pprof debugging endpoints")
}

// CheckDebugFlags validates command line flags
func CheckDebugFlags(_ *viper.Viper) error {
	return nil
}
