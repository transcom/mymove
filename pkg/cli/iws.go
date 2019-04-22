package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/iws"
)

const (
	// IWSRBSHostFlag is the IWS RBS Host Flag
	IWSRBSHostFlag string = "iws-rbs-host"
)

// InitIWSFlags initializes CSRF command line flags
func InitIWSFlags(flag *pflag.FlagSet) {
	flag.String(IWSRBSHostFlag, "", "Hostname for the IWS RBS")
}

// InitRBSPersonLookup is the RBS Person Lookup service
func InitRBSPersonLookup(v *viper.Viper, logger logger) (*iws.RBSPersonLookup, error) {
	return iws.NewRBSPersonLookup(
		v.GetString(IWSRBSHostFlag),
		v.GetString("dod-ca-package"),
		v.GetString("move-mil-dod-tls-cert"),
		v.GetString("move-mil-dod-tls-key"))
}
