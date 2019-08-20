package cli

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// MutualTLSPortFlag is the Mutual TLS Port Flag
	MutualTLSPortFlag string = "mutual-tls-port"
	// TLSPortFlag is the TLS Port Flag
	TLSPortFlag string = "tls-port"
	// NoTLSPortFlag is the No TLS Port Flag
	NoTLSPortFlag string = "no-tls-port"
)

type errInvalidPort struct {
	Port int
}

func (e *errInvalidPort) Error() string {
	return fmt.Sprintf("invalid port %d, must be > 0 and <= 65535", e.Port)
}

// InitPortFlags initializes Port command line flags
func InitPortFlags(flag *pflag.FlagSet) {
	flag.Int(MutualTLSPortFlag, 9443, "The `port` for the mutual TLS listener.")
	flag.Int(TLSPortFlag, 8443, "the `port` for the server side TLS listener.")
	flag.Int(NoTLSPortFlag, 8080, "the `port` for the listener not requiring any TLS.")
}

// CheckPorts validates the Port command line flags
func CheckPorts(v *viper.Viper) error {
	portVars := []string{
		MutualTLSPortFlag,
		TLSPortFlag,
		NoTLSPortFlag,
	}

	for _, c := range portVars {
		err := ValidatePort(v, c)
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidatePort validates a Port passed in from the command line
func ValidatePort(v *viper.Viper, flagname string) error {
	if p := v.GetInt(flagname); p <= 0 || p > 65535 {
		return errors.Wrap(&errInvalidPort{Port: p}, fmt.Sprintf("%s is invalid", flagname))
	}
	return nil
}
