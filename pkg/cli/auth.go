package cli

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Set of flags used for authentication
const (
	//RA Summary: gosec - G101 - Password Management: Hardcoded Password
	//RA: This line was flagged because of use of the word "secret"
	//RA: This line is used to identify the name of the flag. ClientAuthSecretKeyFlag is the Client Auth Secret Key Flag.
	//RA: This value of this variable does not store an application secret.
	//RA Developer Status: Mitigated
	//RA Validator Status: Mitigated
	//RA Validator: jneuner@mitre.org
	//RA Modified Severity: CAT III
	// #nosec G101
	// ClientAuthSecretKeyFlag is the Client Auth Secret Key Flag
	ClientAuthSecretKeyFlag string = "client-auth-secret-key"

	// Okta API key flag
	//RA Summary: gosec - G101 - Password Management: Hardcoded Password
	//RA: This line was flagged because of use of the word "key"
	//RA: This line is used to identify the name of the flag. OktaApiKeyFlag is the Okta API Key Flag.
	//RA: This value of this variable does not store an application secret.
	//RA Developer Status: Mitigated
	//RA Validator Status: Mitigated
	//RA Validator: leodis.f.scott.civ@mail.mil
	//RA Modified Severity: CAT III
	// #nosec G101
	OktaAPIKeyFlag string = "okta-api-key"

	// Okta flags for local development environment that serves test-milmove.okta.mil
	// Okta tenant flags
	OktaTenantOrgURLFlag string = "okta-tenant-org-url"
	// OktaTenantCallbackPortFlag is the test-milmove Callback Port Flag
	OktaTenantCallbackPortFlag string = "okta-tenant-callback-port"
	// OktaTenantCallbackPortFlag is the test-milmove Callback Protocol Flag
	OktaTenantCallbackProtocolFlag string = "okta-tenant-callback-protocol"

	// Okta Customer client id and secret flags
	OktaCustomerClientIDFlag string = "okta-customer-client-id"
	OktaCustomerCallbackURL  string = "okta-customer-callback-url"

	// RA Summary: gosec - G101 - Password Management: Hardcoded Password
	// RA: This line was flagged because of use of the word "secret"
	// RA: This line is used to identify the name of the flag. OktaCustomerSecretKeyFlag is the Okta Customer Application Secret Key Flag.
	// RA: This value of this variable does not store an application secret.
	// RA Developer Status: RA Request
	// RA Validator Status: Mitigated
	// RA Validator: leodis.f.scott.civ@mail.mil
	// RA Modified Severity: CAT III
	// #nosec G101
	OktaCustomerSecretKeyFlag string = "okta-customer-secret-key"

	// Okta Office client id and secret flags
	OktaOfficeClientIDFlag string = "okta-office-client-id"
	OktaOfficeCallbackURL  string = "okta-office-callback-url"

	// RA Summary: gosec - G101 - Password Management: Hardcoded Password
	// RA: This line was flagged because of use of the word "secret"
	// RA: This line is used to identify the name of the flag. OktaOfficeSecretKeyFlag is the Okta Office Application Secret Key Flag.
	// RA: This value of this variable does not store an application secret.
	// RA Developer Status: RA Request
	// RA Validator Status: Mitigated
	// RA Validator: leodis.f.scott.civ@mail.mil
	// RA Modified Severity: CAT III
	// #nosec G101
	OktaOfficeSecretKeyFlag string = "okta-office-secret-key"

	// Okta Admin client id and secret flags
	OktaAdminClientIDFlag string = "okta-admin-client-id"
	OktaAdminCallbackURL  string = "okta-admin-callback-url"

	// RA Summary: gosec - G101 - Password Management: Hardcoded Password
	// RA: This line was flagged because of use of the word "secret"
	// RA: This line is used to identify the name of the flag. OktaAdminSecretKeyFlag is the Okta Admin Application Secret Key Flag.
	// RA: This value of this variable does not store an application secret.
	// RA Developer Status: RA Request
	// RA Validator Status: Mitigated
	// RA Validator: leodis.f.scott.civ@mail.mil
	// RA Modified Severity: CAT III
	// #nosec G101
	OktaAdminSecretKeyFlag  string = "okta-admin-secret-key"
	OktaOfficeGroupIDFlag   string = "okta-office-group-id"
	OktaCustomerGroupIDFlag string = "okta-customer-group-id"
)

// InitAuthFlags initializes Auth command line flags
func InitAuthFlags(flag *pflag.FlagSet) {
	flag.String(ClientAuthSecretKeyFlag, "", "Client auth secret JWT key.")
	flag.String(OktaAPIKeyFlag, "", "The api key for updating okta values in MilMove.")

	flag.String(OktaTenantOrgURLFlag, "", "Okta tenant org URL.")
	flag.Int(OktaTenantCallbackPortFlag, 443, "The port for callback URLs.")
	flag.String(OktaTenantCallbackProtocolFlag, "https", "Protocol for non local environments.")
	flag.String(OktaCustomerClientIDFlag, "", "The client ID for the military customer app, aka 'my'.")
	flag.String(OktaCustomerCallbackURL, "", "The callback URL from logging in to the customer Okta app back to MilMove.")
	flag.String(OktaCustomerSecretKeyFlag, "", "The secret key for the miltiary customer app, aka 'my'.")
	flag.String(OktaOfficeClientIDFlag, "", "The client ID for the military Office app, aka 'my'.")
	flag.String(OktaOfficeCallbackURL, "", "The callback URL from logging in to the office Okta app back to MilMove.")
	flag.String(OktaOfficeSecretKeyFlag, "", "The secret key for the miltiary Office app, aka 'my'.")
	flag.String(OktaAdminClientIDFlag, "", "The client ID for the military Admin app, aka 'my'.")
	flag.String(OktaAdminCallbackURL, "", "The callback URL from logging in to the admin Okta app back to MilMove.")
	flag.String(OktaAdminSecretKeyFlag, "", "The secret key for the miltiary Admin app, aka 'my'.")
	flag.String(OktaOfficeGroupIDFlag, "", "The office group id for the Office app, aka 'office'.")
	flag.String(OktaCustomerGroupIDFlag, "", "The customer group id for the Customer app.")
}

// CheckAuth validates Auth command line flags
func CheckAuth(v *viper.Viper) error {

	if err := ValidateProtocol(v, OktaTenantCallbackProtocolFlag); err != nil {
		return err
	}

	if err := ValidatePort(v, OktaTenantCallbackPortFlag); err != nil {
		return err
	}

	clientIDVars := []string{
		OktaCustomerClientIDFlag,
		OktaOfficeClientIDFlag,
		OktaAdminClientIDFlag,
		OktaAPIKeyFlag,
	}

	secretKeyVars := []string{
		OktaCustomerSecretKeyFlag,
		OktaOfficeSecretKeyFlag,
		OktaAdminSecretKeyFlag,
	}

	groupIDVars := []string{
		OktaOfficeGroupIDFlag,
		OktaCustomerGroupIDFlag,
	}

	for _, c := range clientIDVars {
		clientID := v.GetString(c)
		{
			if len(clientID) == 0 {
				return errors.Errorf("%s is missing", c)
			}
		}
	}

	for _, s := range secretKeyVars {
		privateKey := v.GetString(s)
		if len(privateKey) == 0 {
			return errors.Errorf("%s is missing", s)
		}
	}

	for _, s := range groupIDVars {
		groupID := v.GetString(s)
		if len(groupID) == 0 {
			return errors.Errorf("%s is missing", s)
		}
	}

	return nil
}

// ParsePrivateKey takes a private key and parses it into an slice of individual keys
func ParsePrivateKey(str string) []string {

	privateKeyFormat := "-----BEGIN PRIVATE KEY-----\n%s\n-----END PRIVATE KEY-----"

	// https://tools.ietf.org/html/rfc7468#section-2
	//	- https://stackoverflow.com/questions/20173472/does-go-regexps-any-charcter-match-newline
	re := regexp.MustCompile(`(?s)([-]{5}BEGIN PRIVATE KEY[-]{5})(\s*)(.+?)(\s*)([-]{5}END PRIVATE KEY[-]{5})`)
	matches := re.FindAllStringSubmatch(str, -1)

	privateKeys := make([]string, 0, len(matches))
	for _, m := range matches {
		// each match will include a slice of strings starting with
		// (0) the full match, then
		// (1) -----BEGIN PRIVATE KEY-----,
		// (2) whitespace if any,
		// (3) base64-encoded privateKey data,
		// (4) whitespace if any, and then
		// (5) -----END PRIVATE KEY-----,
		privateKeys = append(privateKeys, fmt.Sprintf(privateKeyFormat, m[3]))
	}
	return privateKeys
}
