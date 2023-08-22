package cli

import (
	"fmt"
	"regexp"
	"strings"

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
	// LoginGovCallbackProtocolFlag is the Login.gov Callback Protocol Flag
	LoginGovCallbackProtocolFlag string = "login-gov-callback-protocol"
	// LoginGovCallbackPortFlag is the Login.gov Callback Port Flag
	LoginGovCallbackPortFlag string = "login-gov-callback-port"
	// LoginGovSecretKeyFlag is the Login.gov Secret Key Flag
	LoginGovSecretKeyFlag string = "login-gov-secret-key"
	// LoginGovMyClientIDFlag is the Login.gov My Client ID Flag
	LoginGovMyClientIDFlag string = "login-gov-my-client-id"
	// LoginGovOfficeClientIDFlag is the Login.gov Office Client ID Flag
	LoginGovOfficeClientIDFlag string = "login-gov-office-client-id"
	// LoginGovAdminClientIDFlag is the Login.gov Admin Client ID Flag
	LoginGovAdminClientIDFlag string = "login-gov-admin-client-id"
	// LoginGovHostnameFlag is the Login.gov Hostname Flag
	LoginGovHostnameFlag string = "login-gov-hostname"

	// NOTE: When constructing callback URLs, the
	// LoginGovCallbackProtocolFlag and LoginGovCallbackPortFlag are
	// *also* used for Okta

	// Okta flags for local development environment that serves test-milmove.okta.mil
	// Okta tenant flags
	OktaTenantOrgURLFlag string = "okta-tenant-org-url"
	// OktaTenantCallbackPortFlag is the test-milmove Callback Port Flag
	OktaTenantCallbackPortFlag string = "okta-tenant-callback-port"
	// OktaTenantCallbackPortFlag is the test-milmove Callback Protocol Flag
	OktaTenantCallbackProtocolFlag string = "okta-tenant-callback-protocol"

	// Okta Customer client id and secret flags
	OktaCustomerClientIDFlag string = "okta-customer-client-id"

	// RA Summary: gosec - G101 - Password Management: Hardcoded Password
	// RA: This line was flagged because of use of the word "secret"
	// RA: This line is used to identify the name of the flag. OktaCustomerSecretKeyFlag is the Okta Customer Application Secret Key Flag.
	// RA: This value of this variable does not store an application secret.
	// RA Developer Status: RA Request
	// RA Validator Status:
	// RA Validator:
	// RA Modified Severity:
	// #nosec G101
	OktaCustomerSecretKeyFlag string = "okta-customer-secret-key"

	// Okta Office client id and secret flags
	OktaOfficeClientIDFlag string = "okta-office-client-id"

	// RA Summary: gosec - G101 - Password Management: Hardcoded Password
	// RA: This line was flagged because of use of the word "secret"
	// RA: This line is used to identify the name of the flag. OktaOfficeSecretKeyFlag is the Okta Office Application Secret Key Flag.
	// RA: This value of this variable does not store an application secret.
	// RA Developer Status: RA Request
	// RA Validator Status:
	// RA Validator:
	// RA Modified Severity:
	// #nosec G101
	OktaOfficeSecretKeyFlag string = "okta-office-secret-key"

	// Okta Admin client id and secret flags
	OktaAdminClientIDFlag string = "okta-admin-client-id"

	// RA Summary: gosec - G101 - Password Management: Hardcoded Password
	// RA: This line was flagged because of use of the word "secret"
	// RA: This line is used to identify the name of the flag. OktaAdminSecretKeyFlag is the Okta Admin Application Secret Key Flag.
	// RA: This value of this variable does not store an application secret.
	// RA Developer Status: RA Request
	// RA Validator Status:
	// RA Validator:
	// RA Modified Severity:
	// #nosec G101
	OktaAdminSecretKeyFlag string = "okta-admin-secret-key"
)

type errInvalidClientID struct {
	ClientID string
}

func (e *errInvalidClientID) Error() string {
	return fmt.Sprintf("invalid client ID %s, must be of format 'urn:gov:gsa:openidconnect.profiles:sp:sso:dod:IDENTIFIER'", e.ClientID)
}

// InitAuthFlags initializes Auth command line flags
func InitAuthFlags(flag *pflag.FlagSet) {
	flag.String(ClientAuthSecretKeyFlag, "", "Client auth secret JWT key.")

	flag.String(LoginGovCallbackProtocolFlag, "https", "Protocol for non local environments.")
	flag.Int(LoginGovCallbackPortFlag, 443, "The port for callback urls.")
	flag.String(LoginGovSecretKeyFlag, "", "Login.gov auth secret JWT key.")
	flag.String(LoginGovMyClientIDFlag, "", "Client ID registered with login gov.")
	flag.String(LoginGovOfficeClientIDFlag, "", "Client ID registered with login gov.")
	flag.String(LoginGovAdminClientIDFlag, "", "Client ID registered with login gov.")
	flag.String(LoginGovHostnameFlag, "secure.login.gov", "Hostname for communicating with login gov.")

	flag.String(OktaTenantOrgURLFlag, "", "Okta tenant org URL.")
	flag.Int(OktaTenantCallbackPortFlag, 443, "The port for callback URLs.")
	flag.String(OktaTenantCallbackProtocolFlag, "https", "Protocol for non local environments.")
	flag.String(OktaCustomerClientIDFlag, "", "The client ID for the military customer app, aka 'my'.")
	flag.String(OktaCustomerSecretKeyFlag, "", "The secret key for the miltiary customer app, aka 'my'.")
	flag.String(OktaOfficeClientIDFlag, "", "The client ID for the military Office app, aka 'my'.")
	flag.String(OktaOfficeSecretKeyFlag, "", "The secret key for the miltiary Office app, aka 'my'.")
	flag.String(OktaAdminClientIDFlag, "", "The client ID for the military Admin app, aka 'my'.")
	flag.String(OktaAdminSecretKeyFlag, "", "The secret key for the miltiary Admin app, aka 'my'.")
}

// CheckAuth validates Auth command line flags
func CheckAuth(v *viper.Viper) error {

	if err := ValidateProtocol(v, LoginGovCallbackProtocolFlag); err != nil {
		return err
	}

	if err := ValidatePort(v, LoginGovCallbackPortFlag); err != nil {
		return err
	}

	secureLoginGov := "secure.login.gov"
	sandboxLoginGov := "idp.int.identitysandbox.gov"
	if loginGovHostname := v.GetString(LoginGovHostnameFlag); loginGovHostname != secureLoginGov && loginGovHostname != sandboxLoginGov {
		return errors.Wrap(&errInvalidHost{Host: loginGovHostname}, fmt.Sprintf("%s is invalid, expected %s or %s", LoginGovHostnameFlag, secureLoginGov, sandboxLoginGov))
	}

	loginGovClientIDVars := []string{
		LoginGovMyClientIDFlag,
		LoginGovOfficeClientIDFlag,
		LoginGovAdminClientIDFlag,
	}

	for _, c := range loginGovClientIDVars {
		err := ValidateClientID(v, c)
		if err != nil {
			return err
		}
	}

	privateKey := v.GetString(LoginGovSecretKeyFlag)
	if len(privateKey) == 0 {
		return errors.Errorf("%s is missing", LoginGovSecretKeyFlag)
	}

	keys := ParsePrivateKey(privateKey)
	if len(keys) == 0 {
		return errors.Errorf("%s is missing key block", LoginGovSecretKeyFlag)
	}

	if err := ValidateProtocol(v, OktaTenantCallbackProtocolFlag); err != nil {
		return err
	}

	if err := ValidatePort(v, OktaTenantCallbackPortFlag); err != nil {
		return err
	}

	oktaClientIDVars := []string{
		OktaCustomerClientIDFlag,
		OktaOfficeClientIDFlag,
		OktaAdminClientIDFlag,
	}

	oktaSecretKeyVars := []string{
		OktaCustomerSecretKeyFlag,
		OktaOfficeSecretKeyFlag,
		OktaAdminSecretKeyFlag,
	}

	for _, c := range oktaClientIDVars {
		clientID := v.GetString(c)
		{
			if len(clientID) == 0 {
				return errors.Errorf("%s is missing", c)
			}
		}
	}

	for _, s := range oktaSecretKeyVars {
		privateKey := v.GetString(s)
		if len(privateKey) == 0 {
			return errors.Errorf("%s is missing", s)
		}
	}

	return nil
}

// ValidateClientID validates a proper Login.gov ClientID was passed
func ValidateClientID(v *viper.Viper, flagname string) error {
	clientID := v.GetString(flagname)
	clientIDParts := strings.Split(clientID, ":")
	clientIDLen := 8
	if len(clientIDParts) != clientIDLen {
		return errors.Wrap(&errInvalidClientID{ClientID: clientID}, fmt.Sprintf("%s is invalid due to length, found %d parts, expected %d. ClientID was %s.", flagname, len(clientIDParts), clientIDLen, clientID))
	}
	openIDFormat := []string{"urn", "gov", "gsa", "openidconnect.profiles", "sp", "sso", "dod"}
	for i, v := range clientIDParts {
		if i == 7 {
			break
		}
		if v != openIDFormat[i] {
			return errors.Wrap(&errInvalidClientID{ClientID: clientID}, fmt.Sprintf("%s is not using OpenID connect", flagname))
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
