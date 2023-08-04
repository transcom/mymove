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

// This function will select the correct provider to use based on its set name.
func getOktaProviderForRequest(r *http.Request, oktaProvider OktaProvider) (goth.Provider, error) {
	session := auth.SessionFromRequestContext(r)
	providerName := customerProviderName
	if session.IsOfficeApp() {
		providerName = officeProviderName
	} else if session.IsAdminApp() {
		providerName = adminProviderName
	}
	gothProvider, err := goth.GetProvider(providerName)
	if err != nil {
		return nil, err
	}

	return gothProvider, nil
}

// This function will use the OktaProvider to return the correct authorization URL to use
func (op *OktaProvider) AuthorizationURL(r *http.Request) (*OktaData, error) {

	// Retrieve the correct Okta Provider to use to get the correct authorization URL. This will choose from customer,
	// office, or admin domains.
	provider, err := getOktaProviderForRequest(r, *op)
	if err != nil {
		op.Logger.Error("Get Goth provider", zap.Error(err))
		return nil, err
	}

	state := generateNonce()

	sess, err := provider.BeginAuth(state)
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

	params := authURL.Query()
	params.Add("nonce", state)
	params.Set("scope", "openid email")

	authURL.RawQuery = params.Encode()

	return &OktaData{authURL.String(), state}, nil
}

// TODO: Clean up when working on callback
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

// Function to register all three providers at once.
// TODO: Split this function up
func (op *OktaProvider) RegisterProviders(customerHostname string, customerCallbackUrl string, customerClientID string, customerSecret string, officeHostname string, officeCallbackUrl string, officeClientID string, officeSecret string, adminHostname string, adminCallbackUrl string, adminClientID string, adminSecret string, callbackProtocol string, callbackPort int, oktaIssuer string) error {
	customerProvider := okta.New(customerClientID, customerSecret, oktaIssuer, customerCallbackUrl, "openid", "profile", "email")
	officeProvider := okta.New(officeClientID, officeSecret, oktaIssuer, officeCallbackUrl, "openid", "profile", "email")
	adminProvider := okta.New(adminClientID, adminSecret, oktaIssuer, adminCallbackUrl, "openid", "profile", "email")
	customerProvider.SetName(customerProviderName)
	officeProvider.SetName(officeProviderName)
	adminProvider.SetName(adminProviderName)
	goth.UseProviders(
		wrapOktaProvider(customerProvider, op.Logger),
		wrapOktaProvider(officeProvider, op.Logger),
		wrapOktaProvider(adminProvider, op.Logger),
	)

	return nil
}
