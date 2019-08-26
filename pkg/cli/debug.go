package cli

import (
	"github.com/spf13/pflag"
)

const (
	//DebugPProfFlag enables the pprof debugging endpoints
	DebugPProfFlag string = "debug-pprof"
)

//  InitDebugFlags initializes the Debug command line flags
func InitDebugFlags(flag *pflag.FlagSet) {
	flag.Bool(DebugPProfFlag, false, "Enables the go pprof debugging endpoints")
}
