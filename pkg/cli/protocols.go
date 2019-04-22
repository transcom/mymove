package cli

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type errInvalidProtocol struct {
	Protocol string
}

func (e *errInvalidProtocol) Error() string {
	return fmt.Sprintf("invalid protocol %s, must be http or https", e.Protocol)
}

// CheckProtocols validates the Protocol command line flags
func CheckProtocols(v *viper.Viper) error {

	protocolVars := []string{
		LoginGovCallbackProtocolFlag,
		HTTPSDDCProtocolFlag,
	}

	for _, c := range protocolVars {
		if p := v.GetString(c); p != "http" && p != "https" {
			return errors.Wrap(&errInvalidProtocol{Protocol: p}, fmt.Sprintf("%s is invalid", c))
		}
	}

	return nil
}
