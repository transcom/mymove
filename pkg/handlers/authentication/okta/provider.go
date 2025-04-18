package okta

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/markbates/goth"
	gothOkta "github.com/markbates/goth/providers/okta"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/random"
)

const MilProviderName = "milProvider"
const OfficeProviderName = "officeProvider"
const AdminProviderName = "adminProvider"

// This provider struct will be unique to mil, office, and admin. When you run goth.getProvider() you'll now have the orgURL, clientID, and secret on hand
type Provider struct {
	gothOkta.Provider
	orgURL      string
	callbackURL string
	clientID    string
	secret      string
	logger      *zap.Logger
}

type Data struct {
	RedirectURL string
	Nonce       string
	GothSession goth.Session
}

// NewProvider creates a new instance of the Okta provider with the specified orgURL
func NewProvider(orgURL string, callbackURL string, clientID string, secret string) *Provider {
	return &Provider{
		orgURL:      orgURL,
		callbackURL: callbackURL,
		clientID:    clientID,
		secret:      secret,
	}
}

// This function will select the correct provider to use based on its set name.
func GetOktaProviderForRequest(r *http.Request) (*Provider, error) {
	session := auth.SessionFromRequestContext(r)

	// Default the provider name to the "MilProviderName" which is the customer application
	// It will update based on if office or admin app
	providerName := MilProviderName

	// Set the provider name based on of it is an office or admin app. Remember, the provider is selected by its name
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

	provider, err := ConvertGothProviderToOktaProvider(gothProvider)
	if err != nil {
		return nil, fmt.Errorf("okta provider was likely not wrapped properly when initially created and registered")
	}

	return provider, nil
}

func ConvertGothProviderToOktaProvider(gothProvider goth.Provider) (*Provider, error) {
	// Conduct type assertion to retrieve the Okta.Provider values from goth.Provider
	provider, ok := gothProvider.(*Provider)
	if !ok {
		// Received provider was not wrapped during its registration, it is of type goth.Provider but not of
		// the Okta Provider type that it should be
		return nil, fmt.Errorf("provided provider is not of the expected okta type")
	}
	return provider, nil
}

// This function will use the OktaProvider to return the correct authorization URL to use
func (op *Provider) AuthorizationURL(r *http.Request) (*Data, error) {

	// Retrieve the correct Okta Provider to use to get the correct authorization URL. This will choose from customer,
	// office, or admin domains and use their information to create the URL.
	provider, err := GetOktaProviderForRequest(r)
	if err != nil {
		op.logger.Error("Get Goth provider", zap.Error(err))
		return nil, err
	}

	// Generate a new state that will later be stored in a cookie for auth
	state := GenerateNonce()

	// Generate a session from the provider and state (nonce)
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

	return &Data{authURL.String(), state, sess}, nil
}

func NewOktaProvider(logger *zap.Logger) *Provider {
	return &Provider{
		logger: logger,
	}
}

// This function allows us to wrap new registered providers with the zap logger. The initial Okta provider is already wrapped
// This will wrap the gothOkta provider with our own version of OktaProvider (With added methods)
func WrapOktaProvider(provider *gothOkta.Provider, orgURL string, clientID string, secret string, callbackURL string, logger *zap.Logger) *Provider {
	return &Provider{
		Provider:    *provider,
		orgURL:      orgURL,
		clientID:    clientID,
		secret:      secret,
		callbackURL: callbackURL,
		logger:      logger,
	}
}

// Function to register all three providers at once.
func (op *Provider) RegisterProviders(v *viper.Viper) error {
	oktaTenantOrgURL := v.GetString(cli.OktaTenantOrgURLFlag)
	oktaCustomerCallbackURL := v.GetString(cli.OktaCustomerCallbackURL)
	oktaCustomerClientID := v.GetString(cli.OktaCustomerClientIDFlag)
	oktaCustomerSecretKey := v.GetString(cli.OktaCustomerSecretKeyFlag)
	oktaOfficeCallbackURL := v.GetString(cli.OktaOfficeCallbackURL)
	oktaOfficeClientID := v.GetString(cli.OktaOfficeClientIDFlag)
	oktaOfficeSecretKey := v.GetString(cli.OktaOfficeSecretKeyFlag)
	oktaAdminCallbackURL := v.GetString(cli.OktaAdminCallbackURL)
	oktaAdminClientID := v.GetString(cli.OktaAdminClientIDFlag)
	oktaAdminSecretKey := v.GetString(cli.OktaAdminSecretKeyFlag)

	// Declare OIDC scopes to be used within the providers
	scope := []string{"openid", "email", "profile"}

	// Register customer provider and pull values from env variables
	err := op.RegisterOktaProvider(MilProviderName, oktaTenantOrgURL, oktaCustomerCallbackURL, oktaCustomerClientID, oktaCustomerSecretKey, scope)
	if err != nil {
		op.logger.Error("Could not register customer okta provider", zap.Error(err))
		return err
	}

	// Register office provider
	err = op.RegisterOktaProvider(OfficeProviderName, oktaTenantOrgURL, oktaOfficeCallbackURL, oktaOfficeClientID, oktaOfficeSecretKey, scope)
	if err != nil {
		op.logger.Error("Could not register office okta provider", zap.Error(err))
		return err
	}

	// Register admin provider
	err = op.RegisterOktaProvider(AdminProviderName, oktaTenantOrgURL, oktaAdminCallbackURL, oktaAdminClientID, oktaAdminSecretKey, scope)
	if err != nil {
		op.logger.Error("Could not register admin okta provider", zap.Error(err))
		return err
	}

	return nil
}

// Create a new Okta provider and register it under the Goth providers
func (op *Provider) RegisterOktaProvider(name string, orgURL string, callbackURL string, clientID string, secret string, scope []string) error {
	// Use goth to create a new provider
	provider := gothOkta.New(clientID, secret, orgURL, callbackURL, scope...)
	// Set the name manualy
	provider.SetName(name)
	// Assign to the active goth providers in a type asserted format based on our Provider struct
	goth.UseProviders(WrapOktaProvider(provider, orgURL, clientID, secret, callbackURL, op.logger))

	// Check that the provider exists now. The previous functions do not have error handling
	err := VerifyProvider(name)
	if err != nil {
		op.logger.Error("Could not verify goth provider", zap.Error(err))
		return err
	}
	return nil
}

// Check if the provided provider name exists
func VerifyProvider(name string) error {
	provider, err := goth.GetProvider(name)
	fmt.Println(provider)
	if err != nil {
		return err
	}
	return nil
}

func (op *Provider) SetSecret(secret string) {
	op.secret = secret
}

func (op *Provider) GetSecret() string {
	return op.secret
}

func (op *Provider) SetClientID(ID string) {
	op.clientID = ID
}

func (op *Provider) GetClientID() string {
	return op.clientID
}

func (op *Provider) SetOrgURL(orgURL string) {
	op.orgURL = orgURL
}

func (op *Provider) GetOrgURL() string {
	return op.orgURL
}

func (op *Provider) GetTokenURL() string {
	return op.orgURL + "/oauth2/default/v1/token"
}
func (op *Provider) GetAuthURL() string {
	return op.orgURL + "/oauth2/default/v1/authorize"
}
func (op *Provider) GetUserInfoURL() string {
	return op.orgURL + "/oauth2/default/v1/userinfo"
}
func (op *Provider) SetCallbackURL(URL string) {
	op.callbackURL = URL
}
func (op *Provider) GetCallbackURL() string {
	return op.callbackURL
}
func (op *Provider) GetIssuerURL() string {
	return op.orgURL + "/oauth2/default"
}
func (op *Provider) GetLogoutURL() string {
	return op.orgURL + "/oauth2/default/v1/logout"
}
func (op *Provider) GetRevokeURL() string {
	return op.orgURL + "/oauth2/default/v1/revoke"
}
func (op *Provider) GetSessionsURL() string {
	return op.orgURL + "/oauth2/default/v1/sessions"
}
func (op *Provider) GetJWKSURL() string {
	return op.orgURL + "/oauth2/default/.well-known/jwks.json"
}
func (op *Provider) GetOpenIDConfigURL() string {
	return op.orgURL + "/oauth2/default/.well-known/openid-configuration"
}
func (op *Provider) GetUsersURL() string {
	return op.orgURL + "/api/v1/users/"
}
func (op *Provider) GetUserURL(oktaUserID string) string {
	return op.orgURL + "/api/v1/users/" + oktaUserID
}
func (op *Provider) GetCreateUserURL(activate string) string {
	return op.orgURL + "/api/v1/users/?activate=" + url.QueryEscape(activate)
}
func (op *Provider) GetCreateAccountURL(activate string) string {
	return op.orgURL + "/api/v1/users/?activate=" + url.QueryEscape(activate)
}
func (op *Provider) GetUserGroupsURL(userID string) string {
	return op.orgURL + "/api/v1/users/" + userID + "/groups"
}
func (op *Provider) AddUserToGroupURL(groupID string, userID string) string {
	return op.orgURL + "/api/v1/groups/" + groupID + "/users/" + userID
}

// TokenURL returns a full URL to retrieve a user token from okta.mil
func (op Provider) TokenURL(r *http.Request) string {
	session := auth.SessionFromRequestContext(r)

	tokenURL := session.Hostname + "/oauth2/default/v1/token"
	op.logger.Info("Session", zap.String("tokenUrl", tokenURL))

	return tokenURL
}

func GenerateNonce() string {
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
