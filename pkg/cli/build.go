package cli

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// BuildRootFlag is the build root flag
	BuildRootFlag string = "build-root"

	// DefaultBuildRoot Path to the build directory
	DefaultBuildRoot string = "build"
)

var (
	errMissingBuildRoot = errors.New("missing build root")
)

// InitBuildFlags initializes the build command line flags
func InitBuildFlags(flag *pflag.FlagSet) {
	flag.StringP(BuildRootFlag, "b", DefaultBuildRoot, "the root directory containing files to serve")
}

// CheckBuild validates the build command line flags
func CheckBuild(v *viper.Viper) error {
	if str := v.GetString(BuildRootFlag); len(str) == 0 {
		return errMissingBuildRoot
	}
	return nil
}
