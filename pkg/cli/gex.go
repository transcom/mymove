package cli

import (
	"fmt"
	"net/url"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// GEXBasicAuthUsernameFlag is the GEX Basic Auth Username Flag
	GEXBasicAuthUsernameFlag string = "gex-basic-auth-username"
	//RA Summary: gosec - G101 - Password Management: Hardcoded Password
	//RA: This line was flagged because of use of the word "password"
	//RA: This line is used to identify the name of the flag. GEXBasicAuthPasswordFlag is the GEX Basic Auth Password Flag.
	//RA: This value of this variable does not store an application password.
	//RA Developer Status: Mitigated
	//RA Validator Status: Mitigated
	//RA Validator: jneuner@mitre.org
	//RA Modified Severity: CAT III

	// GEXBasicAuthPasswordFlag is the GEX Basic Auth Password Flag #nosec G101
	GEXBasicAuthPasswordFlag string = "gex-basic-auth-password"
	// GEXSendProdInvoiceFlag is the GEX Send Prod Invoice Flag
	GEXSendProdInvoiceFlag string = "gex-send-prod-invoice"
	// GEXURLFlag is the GEX URL FLag
	GEXURLFlag string = "gex-url"
	// SendToSyncada is the flag to control if we try sending files to syncada or not
	SendToSyncada string = "send-to-syncada"
)

var gexHostnames = []string{
	"gexweba.daas.dla.mil",
	"gexwebb.daas.dla.mil",
}

var gexPaths = []string{
	"/msg_data/submit",
	"/msg_data/submit/",
}

var gexChannels = []string{
	"",
	"TRANSCOM-DPS-MILMOVE-GHG-IN-IGC-RCOM",
	"TRANSCOM-DPS-MILMOVE-CPS-IN-USBANK-RCOM",
}

// InitGEXFlags initializes GEX command line flags
func InitGEXFlags(flag *pflag.FlagSet) {
	flag.String(GEXBasicAuthUsernameFlag, "", "GEX api auth username")
	flag.String(GEXBasicAuthPasswordFlag, "", "GEX api auth password")
	flag.Bool(GEXSendProdInvoiceFlag, false, "Flag (bool) for EDI Invoices to signify if they should be sent with Production or Test indicator")
	flag.Bool(SendToSyncada, false, "Flag (bool) for turning on or off sending EDI 858s to syncada, default false")

	flag.String(GEXURLFlag, "", "URL for sending an HTTP POST request to GEX")
}

// CheckGEX validates GEX command line flags
func CheckGEX(v *viper.Viper) error {
	gexURL := v.GetString(GEXURLFlag)

	if len(gexURL) == 0 {
		return nil
	}

	// Parse the URL and check it
	u, parseErr := url.Parse(gexURL)
	if parseErr != nil {
		return parseErr
	}

	if u.Scheme != "https" {
		return fmt.Errorf("invalid gexURL Scheme %s, expecting https", u.Scheme)
	}

	if !stringSliceContains(gexHostnames, u.Hostname()) {
		return fmt.Errorf("invalid gexUrl Hostname %s, expecting one of %q", u.Hostname(), gexHostnames)
	}

	if !stringSliceContains(gexPaths, u.Path) {
		return fmt.Errorf("invalid gexUrl Path %s, expecting one of %q", u.Path, gexPaths)
	}

	channel := u.Query().Get("channel")
	if !stringSliceContains(gexChannels, channel) {
		return fmt.Errorf("invalid gexUrl channel query parameter %s, expecting one of %q", channel, gexChannels)
	}

	if len(v.GetString(GEXBasicAuthUsernameFlag)) == 0 {
		return fmt.Errorf("GEX_BASIC_AUTH_USERNAME is missing")
	}
	if len(v.GetString(GEXBasicAuthPasswordFlag)) == 0 {
		return fmt.Errorf("GEX_BASIC_AUTH_PASSWORD is missing")
	}

	return nil
}
