package cli

import (
	"encoding/hex"

	"github.com/pkg/errors"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// CSRFAuthKeyFlag is the CSRF Auth Key Flag
	CSRFAuthKeyFlag string = "csrf-auth-key"
)

// InitCSRFFlags initializes CSRF command line flags
func InitCSRFFlags(flag *pflag.FlagSet) {
	flag.String(CSRFAuthKeyFlag, "", "CSRF Auth Key, 32 byte long")
}

// CheckCSRF validates CSRF command line flags
func CheckCSRF(v *viper.Viper) error {

	csrfAuthKey, err := hex.DecodeString(v.GetString(CSRFAuthKeyFlag))
	if err != nil {
		return errors.Wrap(err, "Error decoding CSRF Auth Key")
	}
	if len(csrfAuthKey) != 32 {
		return errors.Errorf("CSRF Auth Key is not 32 bytes. Auth Key length: %d", len(csrfAuthKey))
	}

	return nil
}
