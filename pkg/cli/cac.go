package cli

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-input"

	"github.com/transcom/nom/pkg/pkcs11"
)

type ErrInvalidPath struct {
	Path string
}

func (e *ErrInvalidPath) Error() string {
	return fmt.Sprintf("invalid path %q", e.Path)
}

type ErrInvalidLabel struct {
	Cert string
	Key  string
}

func (e *ErrInvalidLabel) Error() string {
	return fmt.Sprintf("invalid cert label %q or key label %q", e.Cert, e.Key)
}

const (
	// CACFlag indicates that a CAC should be used
	CACFlag string = "cac"
	// PKCS11ModuleFlag is the location of the PCKS11 module to use with the smart card
	PKCS11ModuleFlag string = "pkcs11module"
	// TokenLabel is the Token Label to use with the smart card
	TokenLabelFlag string = "tokenlabel"
	// CertLabel is the Certificate Label to use with the smart card
	CertLabelFlag string = "certlabel"
	// KeyLabel is the Key Label to use with the smart card
	KeyLabelFlag string = "keylabel"
)

var pkcs11Modules = []string{
	"opensc-pkcs11.so",
	"cackey.dylib",
}

// InitCACFlags initializes the CAC Flags
func InitCACFlags(flag *pflag.FlagSet) {
	flag.Bool(CACFlag, false, "Use a CAC for authentication")
	flag.String(PKCS11ModuleFlag, "/usr/local/lib/pkcs11/opensc-pkcs11.so", "Smart card: Path to the PKCS11 module to use")
	flag.String(TokenLabelFlag, "", "Smart card: name of the token to use")
	flag.String(CertLabelFlag, "Certificate for PIV Authentication", "Smart card: label of the public cert")
	flag.String(KeyLabelFlag, "PIV AUTH key", "Smart card: label of the private key")
}

// CheckCAC validates CAC command line flags
func CheckCAC(v *viper.Viper) error {
	if v.GetBool(CACFlag) {
		pkcs11ModulePath := v.GetString(PKCS11ModuleFlag)
		if pkcs11ModulePath == "" {
			return fmt.Errorf("%q is invalid: %w", PKCS11ModuleFlag, &ErrInvalidPath{Path: pkcs11ModulePath})
		} else if _, err := os.Stat(pkcs11ModulePath); err != nil {
			return fmt.Errorf("%q is invalid: %w", PKCS11ModuleFlag, &ErrInvalidPath{Path: pkcs11ModulePath})
		}
		if pkcs11Base := path.Base(pkcs11ModulePath); !stringSliceContains(pkcs11Modules, pkcs11Base) {
			return fmt.Errorf("invalid PKCS11 module %s, expecting one of %q", pkcs11ModulePath, pkcs11Modules)
		}

		certLabel := v.GetString(CertLabelFlag)
		keyLabel := v.GetString(KeyLabelFlag)
		if certLabel == "" || keyLabel == "" {
			return fmt.Errorf("%q or %q is invalid: %w", CertLabelFlag, KeyLabelFlag, &ErrInvalidLabel{Cert: certLabel, Key: keyLabel})
		}
	}
	return nil
}

// GetCACStore retrieves the CAC store
// Call 'defer store.Close()' after retrieving the store
func GetCACStore(v *viper.Viper) (*pkcs11.Store, error) {
	pkcs11ModulePath := v.GetString(PKCS11ModuleFlag)
	tokenLabel := v.GetString(TokenLabelFlag)
	certLabel := v.GetString(CertLabelFlag)
	keyLabel := v.GetString(KeyLabelFlag)
	pkcsConfig := pkcs11.Config{
		Module:           pkcs11ModulePath,
		CertificateLabel: certLabel,
		PrivateKeyLabel:  keyLabel,
		TokenLabel:       tokenLabel,
	}

	store, errPKCS11New := pkcs11.New(pkcsConfig)
	if errPKCS11New != nil {
		return nil, errPKCS11New
	}

	inputUI := &input.UI{
		Writer: os.Stderr,
		Reader: os.Stdin,
	}

	pin, errUIAsk := inputUI.Ask("CAC PIN", &input.Options{
		Default:     "",
		HideOrder:   true,
		HideDefault: true,
		Required:    true,
		Loop:        true,
		Mask:        true,
		ValidateFunc: func(input string) error {
			matched, matchErr := regexp.Match("^\\d+$", []byte(input))
			if matchErr != nil {
				return matchErr
			}
			if !matched {
				return errors.New("Invalid PIN format")
			}
			return nil
		},
	})
	if errUIAsk != nil {
		return nil, errUIAsk
	}

	errLogin := store.Login(pin)
	if errLogin != nil {
		return nil, errLogin
	}
	return store, nil
}
