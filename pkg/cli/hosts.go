package cli

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// HTTPMyServerNameFlag is the HTTP My Server Name Flag
	HTTPMyServerNameFlag string = "http-my-server-name"
	// HTTPOfficeServerNameFlag is the HTTP Office Server Name Flag
	HTTPOfficeServerNameFlag string = "http-office-server-name"
	// HTTPAdminServerNameFlag is the HTTP Admin Server Name Flag
	HTTPAdminServerNameFlag string = "http-admin-server-name"
	// HTTPOrdersServerNameFlag is the HTTP Orders Server Name Flag
	HTTPOrdersServerNameFlag string = "http-orders-server-name"
	// HTTPPrimeServerNameFlag is the HTTP GHC Prime Server Name Flag
	HTTPPrimeServerNameFlag string = "http-prime-server-name"
	// HTTPPPTASServerNameFlag is the HTTP PPTAS Prime Server Name Flag
	HTTPPPTASServerNameFlag string = "http-pptas-server-name"
	// HTTPMyServerNameLocal is the HTTP My Server Name for Local
	HTTPMyServerNameLocal string = "milmovelocal"
	// HTTPOfficeServerNameLocal is the HTTP Office Server Name for Local
	HTTPOfficeServerNameLocal string = "officelocal"
	// HTTPAdminServerNameLocal is the HTTP Admin Server Name for Local
	HTTPAdminServerNameLocal string = "adminlocal"
	// HTTPOrdersServerNameLocal is the HTTP Orders Server Name for Local
	HTTPOrdersServerNameLocal string = "orderslocal"
	// HTTPPrimeServerNameLocal is the HTTP GHC Prime Server Name for Local
	HTTPPrimeServerNameLocal string = "primelocal"
)

type errInvalidHost struct {
	Host string
}

func (e *errInvalidHost) Error() string {
	return fmt.Sprintf("invalid host '%s', must not contain whitespace, :, /, or \\", e.Host)
}

// InitHostFlags initializes the Hosts command line flags
func InitHostFlags(flag *pflag.FlagSet) {
	flag.String(HTTPMyServerNameFlag, HTTPMyServerNameLocal, "Hostname according to environment.")
	flag.String(HTTPOfficeServerNameFlag, HTTPOfficeServerNameLocal, "Hostname according to environment.")
	flag.String(HTTPAdminServerNameFlag, HTTPAdminServerNameLocal, "Hostname according to environment.")
	flag.String(HTTPOrdersServerNameFlag, HTTPOrdersServerNameLocal, "Hostname according to environment.")
	flag.String(HTTPPrimeServerNameFlag, HTTPPrimeServerNameLocal, "Hostname according to environment.")
	flag.String(HTTPPPTASServerNameFlag, HTTPPrimeServerNameLocal, "Hostname according to environment.")
}

// CheckHosts validates the Hosts command line flags
func CheckHosts(v *viper.Viper) error {

	hostVars := []string{
		HTTPMyServerNameFlag,
		HTTPOfficeServerNameFlag,
		HTTPAdminServerNameFlag,
		HTTPOrdersServerNameFlag,
		HTTPPrimeServerNameFlag,
		HTTPPPTASServerNameFlag,
	}

	for _, c := range hostVars {
		err := ValidateHost(v, c)
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateHost validates a Hostname passed in from the command line
func ValidateHost(v *viper.Viper, flagname string) error {
	invalidChars := ":/\\ \t\n\v\f\r"
	if h := v.GetString(flagname); len(h) == 0 || strings.ContainsAny(h, invalidChars) {
		return errors.Wrap(&errInvalidHost{Host: h}, fmt.Sprintf("%s is invalid", flagname))
	}
	return nil
}
