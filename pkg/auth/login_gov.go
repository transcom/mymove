package auth

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
	"go.uber.org/zap"
)

const gothProviderType = "openid-connect"

// LoginGovProvider facilitates generating URLs and parameters for interfacing with Login.gov
type LoginGovProvider struct {
	hostname  string
	secretKey string
	clientID  string
	logger    *zap.Logger
}

// NewLoginGovProvider returns a new LoginGovProvider
func NewLoginGovProvider(hostname string, secretKey string, clientID string, logger *zap.Logger) LoginGovProvider {
	return LoginGovProvider{
		hostname:  hostname,
		secretKey: secretKey,
		clientID:  clientID,
		logger:    logger,
	}
}

// RegisterProvider registers Login.gov with Goth, which uses
// auto-discovery to get the OpenID configuration
func (p LoginGovProvider) RegisterProvider(hostname string) {
	provider, err := openidConnect.New(
		p.clientID,
		p.secretKey,
		fmt.Sprintf("%s/auth/login-gov/callback", hostname),
		p.ConfigURL(),
	)

	if err != nil {
		p.logger.Error("Register Login.gov provider with Goth", zap.Error(err))
	}

	if provider != nil {
		goth.UseProviders(provider)
	}
}

func generateNonce() string {
	nonceBytes := make([]byte, 64)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 64; i++ {
		nonceBytes[i] = byte(random.Int63() % 256)
	}
	return base64.URLEncoding.EncodeToString(nonceBytes)
}

// AuthorizationURL returns a URL for login.gov authorization with required params
func (p LoginGovProvider) AuthorizationURL() (string, error) {
	provider, err := goth.GetProvider(gothProviderType)
	if err != nil {
		p.logger.Error("Get Goth provider", zap.Error(err))
		return "", err
	}
	state := generateNonce()
	sess, err := provider.BeginAuth(state)
	if err != nil {
		p.logger.Error("Goth begin auth", zap.Error(err))
		return "", err
	}

	baseURL, err := sess.GetAuthURL()
	if err != nil {
		p.logger.Error("Goth get auth URL", zap.Error(err))
		return "", err
	}

	authURL, err := url.Parse(baseURL)
	if err != nil {
		p.logger.Error("Parse auth URL", zap.Error(err))
		return "", err
	}

	params := authURL.Query()
	params.Add("acr_values", "http://idmanagement.gov/ns/assurance/loa/1")
	params.Add("nonce", state)
	params.Set("scope", "openid email")

	authURL.RawQuery = params.Encode()
	return authURL.String(), nil
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
	return logoutPath.String()
}

// TokenURL returns a full URL to retrieve a user token from login.gov
func (p LoginGovProvider) TokenURL() string {
	// TODO: Get the token endpoint URL from Goth instead when
	// https://github.com/markbates/goth/pull/207 is resolved
	return fmt.Sprintf("https://%s/api/openid_connect/token", p.hostname)
}

// TokenParams creates query params for use in the token endpoint
func (p LoginGovProvider) TokenParams(code string, expiry time.Time) (url.Values, error) {
	clientAssertion, err := p.createClientAssertionJWT(expiry)
	params := url.Values{
		"client_assertion":      {clientAssertion},
		"client_assertion_type": {"urn:ietf:params:oauth:client-assertion-type:jwt-bearer"},
		"code":                  {code},
		"grant_type":            {"authorization_code"},
	}

	return params, err
}

func (p LoginGovProvider) createClientAssertionJWT(expiry time.Time) (string, error) {
	claims := &jwt.StandardClaims{
		Issuer:    p.clientID,
		Subject:   p.clientID,
		Audience:  p.TokenURL(),
		Id:        generateNonce(),
		ExpiresAt: expiry.Unix(),
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
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	IDToken     string `json:"id_token"`
}

// ConfigURL returns the URL string of the token endpoint
func (p LoginGovProvider) ConfigURL() string {
	return fmt.Sprintf("https://%s/.well-known/openid-configuration", p.hostname)
}
