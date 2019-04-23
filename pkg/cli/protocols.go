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
		err := ValidateProtocol(v, c)
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateProtocol validates a Protocol passed in from the command line
func ValidateProtocol(v *viper.Viper, flagname string) error {
	if p := v.GetString(flagname); p != "http" && p != "https" {
		return errors.Wrap(&errInvalidProtocol{Protocol: p}, fmt.Sprintf("%s is invalid", flagname))
	}
	return nil
}
