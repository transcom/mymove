package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// IWSRBSHostFlag is the IWS RBS Host Flag
	IWSRBSHostFlag string = "iws-rbs-host"
	// IWSRBSEnabledFlag is the IWS RBS Enabled Flag
	IWSRBSEnabledFlag string = "iws-rbs-enabled"
)

// InitIWSFlags initializes CSRF command line flags
func InitIWSFlags(flag *pflag.FlagSet) {
	flag.String(IWSRBSHostFlag, "", "Hostname for the IWS RBS")
	flag.Bool(IWSRBSEnabledFlag, false, "enable the IWS RBS integration")
}

// CheckIWS validates IWS command line flags
func CheckIWS(v *viper.Viper) error {
	return ValidateHost(v, IWSRBSHostFlag)
}
