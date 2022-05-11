package authentication

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/random"
)

const milProviderName = "milProvider"
const officeProviderName = "officeProvider"
const adminProviderName = "adminProvider"

func getLoginGovProviderForRequest(r *http.Request) (*openidConnect.Provider, error) {
	session := auth.SessionFromRequestContext(r)
	providerName := milProviderName
	if session.IsOfficeApp() {
		providerName = officeProviderName
	} else if session.IsAdminApp() {
		providerName = adminProviderName
	}
	gothProvider, err := goth.GetProvider(providerName)
	if err != nil {
		return nil, err
	}
	return gothProvider.(*openidConnect.Provider), nil
}

// LoginGovProvider facilitates generating URLs and parameters for interfacing with Login.gov
type LoginGovProvider struct {
	hostname  string
	secretKey string
	logger    *zap.Logger
}

// NewLoginGovProvider returns a new LoginGovProvider
func NewLoginGovProvider(hostname string, secretKey string, logger *zap.Logger) LoginGovProvider {
	return LoginGovProvider{
		hostname:  hostname,
		secretKey: secretKey,
		logger:    logger,
	}
}

func (p LoginGovProvider) getOpenIDProvider(hostname string, clientID string, callbackProtocol string, callbackPort int) (goth.Provider, error) {
	return openidConnect.New(
		clientID,
		p.secretKey,
		fmt.Sprintf("%s://%s:%d/auth/login-gov/callback", callbackProtocol, hostname, callbackPort),
		fmt.Sprintf("https://%s/.well-known/openid-configuration", p.hostname),
	)
}

// RegisterProvider registers Login.gov with Goth, which uses
// auto-discovery to get the OpenID configuration
func (p LoginGovProvider) RegisterProvider(milHostname string, milClientID string, officeHostname string, officeClientID string, adminHostname string, adminClientID string, callbackProtocol string, callbackPort int) error {

	milProvider, err := p.getOpenIDProvider(milHostname, milClientID, callbackProtocol, callbackPort)
	if err != nil {
		p.logger.Error("getting open_id provider", zap.String("host", milHostname), zap.Error(err))
		return err
	}
	milProvider.SetName(milProviderName)
	officeProvider, err := p.getOpenIDProvider(officeHostname, officeClientID, callbackProtocol, callbackPort)
	if err != nil {
		p.logger.Error("getting open_id provider", zap.String("host", officeHostname), zap.Error(err))
		return err
	}
	officeProvider.SetName(officeProviderName)
	adminProvider, err := p.getOpenIDProvider(adminHostname, adminClientID, callbackProtocol, callbackPort)
	if err != nil {
		p.logger.Error("getting open_id provider", zap.String("host", adminHostname), zap.Error(err))
		return err
	}
	adminProvider.SetName(adminProviderName)
	goth.UseProviders(milProvider, officeProvider, adminProvider)
	return nil
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

// LoginGovData contains the URL and State nonce used to redirect a user
// login.gov for authentication
type LoginGovData struct {
	RedirectURL string
	Nonce       string
}

// AuthorizationURL returns a URL for login.gov authorization with required params
func (p LoginGovProvider) AuthorizationURL(r *http.Request) (*LoginGovData, error) {
	provider, err := getLoginGovProviderForRequest(r)
	if err != nil {
		p.logger.Error("Get Goth provider", zap.Error(err))
		return nil, err
	}
	state := generateNonce()
	sess, err := provider.BeginAuth(state)
	if err != nil {
		p.logger.Error("Goth begin auth", zap.Error(err))
		return nil, err
	}

	baseURL, err := sess.GetAuthURL()
	if err != nil {
		p.logger.Error("Goth get auth URL", zap.Error(err))
		return nil, err
	}

	authURL, err := url.Parse(baseURL)
	if err != nil {
		p.logger.Error("Parse auth URL", zap.Error(err))
		return nil, err
	}

	params := authURL.Query()
	params.Add("acr_values", "http://idmanagement.gov/ns/assurance/loa/1")
	params.Add("nonce", state)
	params.Set("scope", "openid email")

	authURL.RawQuery = params.Encode()
	return &LoginGovData{authURL.String(), state}, nil
}

// LogoutURL returns a full URL to log out of login.gov with required params
func (p LoginGovProvider) LogoutURL(redirectURL string, idToken string) string {
	logoutPath, _ := url.Parse(fmt.Sprintf("https://%s/openid_connect/logout", p.hostname))
	// Parameters taken from https://developers.login.gov/oidc/#logout
	params := url.Values{
		"id_token_hint":            {idToken},
		"post_logout_redirect_uri": {redirectURL},
		"state":                    {generateNonce()},
	}

	logoutPath.RawQuery = params.Encode()
	strLogoutPath := logoutPath.String()
	p.logger.Info("Logout path", zap.String("strLogoutPath", strLogoutPath))

	return strLogoutPath
}

// TokenURL returns a full URL to retrieve a user token from login.gov
func (p LoginGovProvider) TokenURL() string {
	// TODO: Get the token endpoint URL from Goth instead when
	// https://github.com/markbates/goth/pull/207 is resolved
	tokenURL := fmt.Sprintf("https://%s/api/openid_connect/token", p.hostname)
	p.logger.Info("LoginGovProvider", zap.String("tokenUrl", tokenURL))

	return tokenURL
}

// TokenParams creates query params for use in the token endpoint
func (p LoginGovProvider) TokenParams(code string, clientID string, expiry time.Time) (url.Values, error) {
	clientAssertion, err := p.createClientAssertionJWT(clientID, expiry)
	params := url.Values{
		"client_assertion":      {clientAssertion},
		"client_assertion_type": {"urn:ietf:params:oauth:client-assertion-type:jwt-bearer"},
		"code":                  {code},
		"grant_type":            {"authorization_code"},
	}

	return params, err
}

func (p LoginGovProvider) createClientAssertionJWT(clientID string, expiry time.Time) (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:    clientID,
		Subject:   clientID,
		Audience:  jwt.ClaimStrings([]string{p.TokenURL()}),
		ID:        generateNonce(),
		ExpiresAt: jwt.NewNumericDate(expiry),
	}

	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(p.secretKey))
	if err != nil {
		p.logger.Error("JWT parse private key from PEM", zap.Error(err))
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	jwt, err := token.SignedString(rsaKey)
	if err != nil {
		p.logger.Error("Signing JWT", zap.Error(err))
	}
	return jwt, err
}

// LoginGovTokenResponse is a struct for parsing responses from the token endpoint
type LoginGovTokenResponse struct {
	Error       string `json:"error"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	IDToken     string `json:"id_token"`
}
