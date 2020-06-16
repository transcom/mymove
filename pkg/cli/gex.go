package cli

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// GEXBasicAuthUsernameFlag is the GEX Basic Auth Username Flag
	GEXBasicAuthUsernameFlag string = "gex-basic-auth-username"
	// GEXBasicAuthPasswordFlag is the GEX Basic Auth Password Flag #nosec G101
	GEXBasicAuthPasswordFlag string = "gex-basic-auth-password"
	// GEXSendProdInvoiceFlag is the GEX Send Prod Invoice Flag
	GEXSendProdInvoiceFlag string = "gex-send-prod-invoice"
	// GEXURLFlag is the GEX URL FLag
	GEXURLFlag string = "gex-url"
)

var allGEXURLs = []string{
	"", // empty string allowed for backwards compatibility
	"https://gexweba.daas.dla.mil/msg_data/submit/",
	"https://gexweba.daas.dla.mil:443/msg_data/submit?channel=TRANSCOM-DPS-MILMOVE-GHG-IN-IGC-RCOM",
}

// InitGEXFlags initializes GEX command line flags
func InitGEXFlags(flag *pflag.FlagSet) {
	flag.String(GEXBasicAuthUsernameFlag, "", "GEX api auth username")
	flag.String(GEXBasicAuthPasswordFlag, "", "GEX api auth password")
	flag.Bool(GEXSendProdInvoiceFlag, false, "Flag (bool) for EDI Invoices to signify if they should be sent with Production or Test indicator")
	flag.String(GEXURLFlag, "", "URL for sending an HTTP POST request to GEX")
}

// CheckGEX validates GEX command line flags
func CheckGEX(v *viper.Viper) error {
	gexURL := v.GetString(GEXURLFlag)

	if !stringSliceContains(allGEXURLs, gexURL) {
		return fmt.Errorf("invalid gexUrl %s, expecting one of %q", gexURL, allGEXURLs)
	}

	if len(gexURL) > 0 {
		if len(v.GetString(GEXBasicAuthUsernameFlag)) == 0 {
			return fmt.Errorf("GEX_BASIC_AUTH_USERNAME is missing")
		}
		if len(v.GetString(GEXBasicAuthPasswordFlag)) == 0 {
			return fmt.Errorf("GEX_BASIC_AUTH_PASSWORD is missing")
		}
	}

	return nil
}
