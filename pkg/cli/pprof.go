package cli

import (
	"github.com/spf13/pflag"
)

const (
	//PProfFlag enables the pprof debugging endpoints
	EnablePProfFlag string = "enable-pprof"
)

// InitSwaggerFlags initializes the Swagger command line flags
func InitPProfFlags(flag *pflag.FlagSet) {
	flag.Bool(EnablePProfFlag, false, "Enables the go pprof debugging endpoints")
}
