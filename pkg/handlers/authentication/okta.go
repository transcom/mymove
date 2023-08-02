package authentication

import (
	"net/http"
	"net/url"
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/okta"
	"github.com/transcom/mymove/pkg/auth"
	"go.uber.org/zap"
)

const customerProviderName = "customerProvider"

// const officeProviderName = "officeProvider" //used in the login_gov.go
// const adminProviderName = "adminProvider" // used in login_gov.go

type OktaProvider struct {
	okta.Provider
	Logger *zap.Logger
}

type OktaData struct {
	RedirectURL string
	Nonce       string
}

func (op *OktaProvider) AuthorizationURL(r *http.Request) (*OktaData, error) {
	// if os.Getenv("OKTA_OAUTH2_ISSUER") == "" {
	// 	err := errors.New("Issuer not set")
	// 	op.logger.Error("Issuer not set", zap.Error(err))
	// 	return nil, err
	// }

	state := generateNonce()

	sess, err := op.BeginAuth(state)
	if err != nil {
		op.Logger.Error("Goth begin auth", zap.Error(err))
		return nil, err
	}

	baseURL, err := sess.GetAuthURL()
	if err != nil {
		op.Logger.Error("Goth get auth URL", zap.Error(err))
		return nil, err
	}

	authURL, err := url.Parse(baseURL)
	if err != nil {
		op.Logger.Error("Parse auth URL", zap.Error(err))
		return nil, err
	}

	// TODO: Verify CAC authenticator
	params := authURL.Query()
	session := auth.SessionFromRequestContext(r)
	// TODO: Switch away from idmanagement - This is login.gov
	if session.IsAdminApp() {
		// This specifies that a user has been authenticated with an HSPD12 credential, via their CAC. Both acr_values must be specified.
		params.Add("acr_values", "http://idmanagement.gov/ns/assurance/ial/1 http://idmanagement.gov/ns/assurance/aal/3?hspd12=true")
	} else {
		params.Add("acr_values", "http://idmanagement.gov/ns/assurance/loa/1")
	}
	params.Add("nonce", state)
	params.Set("scope", "openid email")

	authURL.RawQuery = params.Encode()

	return &OktaData{authURL.String(), state}, nil
}

func NewOktaProvider(logger *zap.Logger) *OktaProvider {
	return &OktaProvider{
		Provider: *okta.New(
			os.Getenv("OKTA_OAUTH2_CLIENT_ID"),
			os.Getenv("OKTA_OAUTH2_CLIENT_SECRET"),
			os.Getenv("OKTA_OAUTH2_ISSUER"),
			"http://milmovelocal:3000/",
			"openid", "profile", "email",
		),
		Logger: logger,
	}
}

// This function allows us to wrap new registered providers with the zap logger. The initial Okta provider is already wrapped
func wrapOktaProvider(provider *okta.Provider, logger *zap.Logger) *OktaProvider {
	return &OktaProvider{
		Provider: *provider,
		Logger:   logger,
	}
}

// This function allows us to register a new Okta provider with Goth. This is primarily used
// for the three different Okta applications we're supporting: Customer, Office, and Admin
func (op *OktaProvider) RegisterProvider(providerName string, clientID string, clientSecret string, issuerURL string, callbackURL string) error {

	oktaProvider := okta.NewCustomisedURL(clientID, clientSecret, callbackURL, issuerURL+"/v1/authorize", issuerURL+"/v1/token", issuerURL, issuerURL+"/v1/userinfo", "openid", "profile", "email")

	// set provider name for the Okta provider
	oktaProvider.SetName(providerName)

	goth.UseProviders(
		wrapOktaProvider(oktaProvider, op.Logger),
	)

	return nil
}

func (op *OktaProvider) RegisterProviders(customerHostname string, customerCallbackUrl string, customerClientID string, customerSecret string, officeHostname string, officeCallbackUrl string, officeClientID string, officeSecret string, adminHostname string, adminCallbackUrl string, adminClientID string, adminSecret string, callbackProtocol string, callbackPort int, oktaIssuer string) error {
	customerProvider := okta.New(customerClientID, customerSecret, oktaIssuer, customerCallbackUrl, "openid", "profile", "email")
	officeProvider := okta.New(officeClientID, officeSecret, oktaIssuer, officeCallbackUrl, "openid", "profile", "email")
	adminProvider := okta.New(adminClientID, adminSecret, oktaIssuer, adminCallbackUrl, "openid", "profile", "email")

	goth.UseProviders(
		wrapOktaProvider(customerProvider, op.Logger),
		wrapOktaProvider(officeProvider, op.Logger),
		wrapOktaProvider(adminProvider, op.Logger),
	)

	return nil
}
