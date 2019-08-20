package cli

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// MaxBodySizeFlag is the maximum body size for requests
	MaxBodySizeFlag string = "max-body-size"

	// MaxBodySizeDefault is 20 mb
	MaxBodySizeDefault int64 = 200 * 1000 * 1000
)

// InitMiddlewareFlags initializes the Middleware command line flags
func InitMiddlewareFlags(flag *pflag.FlagSet) {
	flag.Int64(MaxBodySizeFlag, MaxBodySizeDefault, "The maximum request body size in bytes as int")
}

// CheckMiddleWare validates middleware command line flags
func CheckMiddleWare(v *viper.Viper) error {
	if maxBodySize := v.GetInt64(MaxBodySizeFlag); maxBodySize < int64(0) {
		return errors.Errorf("Max Body Size %d must be greater than zero", maxBodySize)
	}

	return nil
}
