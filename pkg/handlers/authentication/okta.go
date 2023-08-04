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

func NewOktaProvider(logger *zap.Logger) *OktaProvider {
	return &OktaProvider{
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
// TODO: Use viper instead of os environment variables
func (op *OktaProvider) RegisterProviders() error {
	// Declare OIDC scopes to be used within the providers
	scope := []string{"openid", "email"}
	// Register customer provider
	err := op.RegisterOktaProvider(customerProviderName, os.Getenv("OKTA_CUSTOMER_HOSTNAME"), os.Getenv("OKTA_CUSTOMER_CALLBACK_URL"), os.Getenv("OKTA_CUSTOMER_CLIENT_ID"), os.Getenv("OKTA_CUSTOMER_SECRET_KEY"), scope)
	if err != nil {
		return err
	}
	// Register office provider
	err = op.RegisterOktaProvider(officeProviderName, os.Getenv("OKTA_OFFICE_HOSTNAME"), os.Getenv("OKTA_OFFICE_CALLBACK_URL"), os.Getenv("OKTA_OFFICE_CLIENT_ID"), os.Getenv("OKTA_OFFICE_SECRET_KEY"), scope)
	if err != nil {
		return err
	}
	// Register admin provider
	err = op.RegisterOktaProvider(adminProviderName, os.Getenv("OKTA_ADMIN_HOSTNAME"), os.Getenv("OKTA_ADMIN_CALLBACK_URL"), os.Getenv("OKTA_ADMIN_CLIENT_ID"), os.Getenv("OKTA_ADMIN_SECRET_KEY"), scope)
	if err != nil {
		return err
	}

	return nil
}

// Create a new Okta provider and register it under the Goth providers
func (op *OktaProvider) RegisterOktaProvider(name string, hostname string, callbackUrl string, clientID string, secret string, scope []string) error {
	provider := okta.New(clientID, secret, hostname, callbackUrl, scope...)
	provider.SetName(name)
	goth.UseProviders(wrapOktaProvider(provider, op.Logger))

	// Check that the provider exists now
	err := verifyProvider(name)
	if err != nil {
		op.Logger.Error("Could not verify goth provider", zap.Error(err))
		return err
	}
	return nil
}

// Check if the provided provider name exists
func verifyProvider(name string) error {
	_, err := goth.GetProvider(name)
	if err != nil {
		return err
	}
	return nil
}
