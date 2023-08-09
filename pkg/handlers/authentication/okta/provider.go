package okta

import (
	"encoding/base64"
	"math/rand"
	"net/http"
	"net/url"
	"os"

	"github.com/markbates/goth"
	gothOkta "github.com/markbates/goth/providers/okta"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/random"
)

const MilProviderName = "milProvider"
const OfficeProviderName = "officeProvider"
const AdminProviderName = "adminProvider"

type OktaProvider struct {
	gothOkta.Provider
	hostname string
	logger   *zap.Logger
}

type OktaData struct {
	RedirectURL string
	Nonce       string
	GothSession goth.Session
}

// This function will select the correct provider to use based on its set name.
func getOktaProviderForRequest(r *http.Request, oktaProvider OktaProvider) (goth.Provider, error) {
	session := auth.SessionFromRequestContext(r)

	// Default the provider name to the "MilProviderName" which is the customer application
	// It will update based on if office or admin app
	providerName := MilProviderName

	// Set the provider name based on of it is an office or admin app. Remember, the provider is slected by its name
	if session.IsOfficeApp() {
		providerName = OfficeProviderName
	} else if session.IsAdminApp() {
		providerName = AdminProviderName
	}

	// Retrieve the provider based on its name
	gothProvider, err := goth.GetProvider(providerName)
	if err != nil {
		return nil, err
	}

	return gothProvider, nil
}

func getProviderName(r *http.Request) string {
	session := auth.SessionFromRequestContext(r)

	// Set the provider name based on of it is an office or admin app. Remember, the provider is slected by its name
	if session.IsOfficeApp() {
		return OfficeProviderName
	} else if session.IsAdminApp() {
		return AdminProviderName
	}
	return MilProviderName
}

// ! This func will likely come back during continuation of the sessions story
// // This function will return the ClientID of the current provider
// func (op *OktaProvider) ClientID(r *http.Request) (string, error) {
// 	// Default the provider name to the "MilProviderName" which is the customer application
// 	providerName := getProviderName(r)

// 	// Retrieve the provider based on its name
// 	gothProvider, err := goth.GetProvider(providerName)
// 	if err != nil {
// 		return "", err
// 	}

// 	return gothProvider.ClientKey, nil
// }

// This function will use the OktaProvider to return the correct authorization URL to use
func (op *OktaProvider) AuthorizationURL(r *http.Request) (*OktaData, error) {

	// Retrieve the correct Okta Provider to use to get the correct authorization URL. This will choose from customer,
	// office, or admin domains and use their information to create the URL.
	provider, err := getOktaProviderForRequest(r, *op)
	if err != nil {
		op.logger.Error("Get Goth provider", zap.Error(err))
		return nil, err
	}

	// Generate a new state that will later be stored in a cookie for auth
	state := generateNonce()

	// Generate a session rom the provider and state (nonce)
	sess, err := provider.BeginAuth(state)
	if err != nil {
		op.logger.Error("Goth begin auth", zap.Error(err))
		return nil, err
	}

	// Use the goth.Session to generate the AuthURL. It knows this from the hostname. Currently we are not using a custom auth server
	// outside of "default" (Note for Okta, "default" doesn't mean the default server, it just means a server named default)
	baseURL, err := sess.GetAuthURL()
	if err != nil {
		op.logger.Error("Goth get auth URL", zap.Error(err))
		return nil, err
	}

	// Parse URL
	authURL, err := url.Parse(baseURL)
	if err != nil {
		op.logger.Error("Parse auth URL", zap.Error(err))
		return nil, err
	}

	params := authURL.Query()
	// Add the nonce and scope to the URL when getting ready to redirect to the login URL
	params.Add("nonce", state)
	params.Set("scope", "openid profile email")

	authURL.RawQuery = params.Encode()

	return &OktaData{authURL.String(), state, sess}, nil
}

func NewOktaProvider(logger *zap.Logger) *OktaProvider {
	return &OktaProvider{
		logger: logger,
	}
}

// This function allows us to wrap new registered providers with the zap logger. The initial Okta provider is already wrapped
// This will wrap the gothOkta provider with our own version of OktaProvider (With added methods)
func wrapOktaProvider(provider *gothOkta.Provider, logger *zap.Logger) *OktaProvider {
	return &OktaProvider{
		Provider: *provider,
		logger:   logger,
	}
}

// Function to register all three providers at once.
// TODO: Use viper instead of os environment variables
func (op *OktaProvider) RegisterProviders() error {

	// Declare OIDC scopes to be used within the providers
	scope := []string{"openid", "email", "profile"}

	// Register customer provider
	err := op.RegisterOktaProvider(MilProviderName, os.Getenv("OKTA_CUSTOMER_HOSTNAME"), os.Getenv("OKTA_CUSTOMER_CALLBACK_URL"), os.Getenv("OKTA_CUSTOMER_CLIENT_ID"), os.Getenv("OKTA_CUSTOMER_SECRET_KEY"), scope)
	if err != nil {
		op.logger.Error("Could not register customer okta provider", zap.Error(err))
		return err
	}
	// Register office provider
	err = op.RegisterOktaProvider(OfficeProviderName, os.Getenv("OKTA_OFFICE_HOSTNAME"), os.Getenv("OKTA_OFFICE_CALLBACK_URL"), os.Getenv("OKTA_OFFICE_CLIENT_ID"), os.Getenv("OKTA_OFFICE_SECRET_KEY"), scope)
	if err != nil {
		op.logger.Error("Could not register office okta provider", zap.Error(err))
		return err
	}
	// Register admin provider
	err = op.RegisterOktaProvider(AdminProviderName, os.Getenv("OKTA_ADMIN_HOSTNAME"), os.Getenv("OKTA_ADMIN_CALLBACK_URL"), os.Getenv("OKTA_ADMIN_CLIENT_ID"), os.Getenv("OKTA_ADMIN_SECRET_KEY"), scope)
	if err != nil {
		op.logger.Error("Could not register admin okta provider", zap.Error(err))
		return err
	}

	return nil
}

// Create a new Okta provider and register it under the Goth providers
func (op *OktaProvider) RegisterOktaProvider(name string, hostname string, callbackUrl string, clientID string, secret string, scope []string) error {
	// Use goth to create a new provider
	provider := gothOkta.New(clientID, secret, hostname, callbackUrl, scope...)
	// Set the name manualy
	provider.SetName(name)
	// Wrap
	wrap := wrapOktaProvider(provider, op.logger)
	// Set hostname
	wrap.SetHostname(hostname)
	// Assign to the active goth providers
	goth.UseProviders(wrap)

	// Check that the provider exists now. The previous functions do not have error handling
	err := verifyProvider(name)
	if err != nil {
		op.logger.Error("Could not verify goth provider", zap.Error(err))
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

func (op OktaProvider) SetHostname(hostname string) {
	op.hostname = hostname
}

func (op OktaProvider) GetHostname() string {
	return op.hostname
}

// TokenURL returns a full URL to retrieve a user token from okta.mil
func (op OktaProvider) TokenURL(r *http.Request) string {
	session := auth.SessionFromRequestContext(r)

	tokenURL := session.Hostname + "/oauth2/default/v1/token"
	op.logger.Info("Session", zap.String("tokenUrl", tokenURL))

	return tokenURL
}

// LogoutURL returns a full URL to log out of login.gov with required params
// !Ensure proper testing after sessions have been handled
// TODO: Ensure works as intended
func (op OktaProvider) LogoutURL(hostname string, redirectURL string, clientId string) (string, error) {
	logoutPath, _ := url.Parse(hostname + "/oauth2/v1/logout")
	// Parameters taken from https://developers.login.gov/oidc/#logout
	params := url.Values{
		"client_id":                {clientId},
		"post_logout_redirect_uri": {redirectURL},
		"state":                    {generateNonce()},
	}

	logoutPath.RawQuery = params.Encode()
	strLogoutPath := logoutPath.String()
	op.logger.Info("Logout path", zap.String("strLogoutPath", strLogoutPath))

	return strLogoutPath, nil
}

func generateNonce() string {
	nonceBytes := make([]byte, 64)
	//RA Summary: gosec - G404 - Insecure random number source (rand)
	//RA: gosec detected use of the insecure package math/rand rather than the more secure cryptographically secure pseudo-random number generator crypto/rand.
	//RA: This particular usage is mitigated by sourcing the seed from crypto/rand in order to create the new random number using math/rand.
	//RA Developer Status: Mitigated
	//RA Validator: jneuner@mitre.org
	//RA Validator Status: Mitigated
	//RA Modified Severity: CAT III
	// #nosec G404
	randomInt := rand.New(random.NewCryptoSeededSource())
	for i := 0; i < 64; i++ {
		nonceBytes[i] = byte(randomInt.Int63() % 256)
	}
	return base64.URLEncoding.EncodeToString(nonceBytes)
}
