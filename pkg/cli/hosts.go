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
	// HTTPTSPServerNameFlag is the HTTP TSP Server Name Flag
	HTTPTSPServerNameFlag string = "http-tsp-server-name"
	// HTTPAdminServerNameFlag is the HTTP Admin Server Name Flag
	HTTPAdminServerNameFlag string = "http-admin-server-name"
	// HTTPOrdersServerNameFlag is the HTTP Orders Server Name Flag
	HTTPOrdersServerNameFlag string = "http-orders-server-name"
)

type errInvalidHost struct {
	Host string
}

func (e *errInvalidHost) Error() string {
	return fmt.Sprintf("invalid host '%s', must not contain whitespace, :, /, or \\", e.Host)
}

// InitHostFlags initializes the Hosts command line flags
func InitHostFlags(flag *pflag.FlagSet) {
	flag.String(HTTPMyServerNameFlag, "milmovelocal", "Hostname according to environment.")
	flag.String(HTTPOfficeServerNameFlag, "officelocal", "Hostname according to environment.")
	flag.String(HTTPTSPServerNameFlag, "tsplocal", "Hostname according to environment.")
	flag.String(HTTPAdminServerNameFlag, "adminlocal", "Hostname according to environment.")
	flag.String(HTTPOrdersServerNameFlag, "orderslocal", "Hostname according to environment.")
}

// CheckHosts validates the Hosts command line flags
func CheckHosts(v *viper.Viper) error {

	hostVars := []string{
		HTTPMyServerNameFlag,
		HTTPOfficeServerNameFlag,
		HTTPTSPServerNameFlag,
		HTTPAdminServerNameFlag,
		HTTPOrdersServerNameFlag,
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
