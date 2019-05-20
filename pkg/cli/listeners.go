package cli

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// MutualTLSListenerFlag is the Mutual TLS Listener Flag
	MutualTLSListenerFlag string = "mutual-tls-enabled"
	// TLSListenerFlag is the TLS Listener Flag
	TLSListenerFlag string = "tls-enabled"
	// NoTLSListenerFlag is the No TLS Listener Flag
	NoTLSListenerFlag string = "no-tls-enabled"
)

// InitListenerFlags initializes Listener command line flags
func InitListenerFlags(flag *pflag.FlagSet) {
	flag.Bool(MutualTLSListenerFlag, false, "enable the mutual TLS listener.")
	flag.Bool(TLSListenerFlag, false, "enable the server side TLS listener.")
	flag.Bool(NoTLSListenerFlag, false, "enable the listener not requiring any TLS.")
}

// CheckListeners validates the Listener command line flags
func CheckListeners(v *viper.Viper) error {
	mutualTLSEnabled := v.GetBool(MutualTLSListenerFlag)
	tlsEnabled := v.GetBool(TLSListenerFlag)
	noTLSEnabled := v.GetBool(NoTLSListenerFlag)

	if (!mutualTLSEnabled) && (!tlsEnabled) && (!noTLSEnabled) {
		return errors.New("no listener enabled")
	}

	return nil
}
