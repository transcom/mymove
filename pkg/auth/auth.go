package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"

	"github.com/dgrijalva/jwt-go"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
)

const loginGovClientID = "urn:gov:gsa:openidconnect.profiles:sp:sso:dod:mymovemil"
const gothProviderType = "openid-connect"
const sessionExpiryInMinutes = 15

// RegisterProvider registers Login.gov with Goth, which uses
// auto-discovery to get the OpenID configuration
func RegisterProvider(jwtSecret, hostname, protocol, clientPort string) {
	if jwtSecret == "" {
		zap.L().Warn("Auth secret key environment variable not set")
	}

	// TODO: set the urls below as variables based on environment rather than hardcoding.
	provider, err := openidConnect.New(
		loginGovClientID,
		jwtSecret,
		fmt.Sprintf("%s://%s:%s/auth/login-gov/callback", protocol, hostname, clientPort),
		"https://idp.int.identitysandbox.gov/.well-known/openid-configuration",
	)

	if err != nil {
		zap.L().Error("Register Login.gov provider with Goth", zap.Error(err))
	}

	if provider != nil {
		goth.UseProviders(provider)
	}
}

// AuthorizationRedirectHandler constructs the Login.gov authentication URL and redirects to it
func AuthorizationRedirectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := getAuthorizationURL()
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// AuthorizationCallbackHandler handles the callback from the Login.gov authorization flow
func AuthorizationCallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authError := r.URL.Query().Get("error")

		// The user has either cancelled or declined to authorize the client
		if authError == "access_denied" {
			http.Redirect(w, r, "http://localhost:3000/landing", http.StatusTemporaryRedirect)
		}

		if authError == "invalid_request" {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		// TODO: validate the state is the same (pull from session)

		session, err := fetchToken(r)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		provider, err := goth.GetProvider(gothProviderType)
		if err != nil {
			zap.L().Error("Get Goth provider", zap.Error(err))
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		user, err := provider.FetchUser(session)
		if err != nil {
			zap.L().Error("Login.gov user info request", zap.Error(err))
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		landingURL := fmt.Sprintf("http://localhost:3000/landing?email=%s", user.RawData["email"])
		http.Redirect(w, r, landingURL, http.StatusTemporaryRedirect)
	}
}

func getAuthorizationURL() (string, error) {
	provider, err := goth.GetProvider(gothProviderType)
	if err != nil {
		zap.L().Error("Get Goth provider", zap.Error(err))
		return "", err
	}
	state := generateNonce()
	sess, err := provider.BeginAuth(state)
	if err != nil {
		zap.L().Error("Goth begin auth", zap.Error(err))
		return "", err
	}

	baseURL, err := sess.GetAuthURL()
	if err != nil {
		zap.L().Error("Goth get auth URL", zap.Error(err))
		return "", err
	}

	authURL, err := url.Parse(baseURL)
	if err != nil {
		zap.L().Error("Parse auth URL", zap.Error(err))
		return "", err
	}

	params := authURL.Query()
	params.Add("acr_values", "http://idmanagement.gov/ns/assurance/loa/1")
	params.Add("nonce", state)
	params.Set("scope", "openid email")

	authURL.RawQuery = params.Encode()
	return authURL.String(), err
}

func fetchToken(r *http.Request) (goth.Session, error) {
	// TODO: How can we get the token URL from Goth instead?
	tokenURL := "http://localhost:8000/api/openid_connect/token"
	clientAssertion, err := createClientAssertionJWT(tokenURL)
	if err != nil {
		return nil, err
	}

	params := url.Values{
		"client_assertion":      {clientAssertion},
		"client_assertion_type": {"urn:ietf:params:oauth:client-assertion-type:jwt-bearer"},
		"code":                  {r.URL.Query().Get("code")},
		"grant_type":            {"authorization_code"},
	}

	response, err := http.PostForm(tokenURL, params)
	if err != nil {
		// TODO
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		// TODO
	}

	type tokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		IDToken     string `json:"id_token"`
	}

	var parsedResponse tokenResponse
	json.Unmarshal(responseBody, &parsedResponse)

	// TODO: decode and validate ID Token

	// TODO: get goth session from storage instead of constructing a new one
	session := openidConnect.Session{
		AccessToken: parsedResponse.AccessToken,
		ExpiresAt:   time.Now().Add(time.Second * time.Duration(parsedResponse.ExpiresIn)),
		IDToken:     parsedResponse.IDToken,
	}

	return &session, err
}

func createClientAssertionJWT(tokenURL string) (string, error) {
	claims := &jwt.StandardClaims{
		Issuer:    loginGovClientID,
		Subject:   loginGovClientID,
		Audience:  tokenURL,
		Id:        generateNonce(),
		ExpiresAt: time.Now().Add(time.Minute * sessionExpiryInMinutes).Unix(),
	}

	key, isSet := os.LookupEnv(jwtPrivateKeyEnvVariable)
	if !isSet {
		// TODO
	}

	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(key))
	if err != nil {
		// TODO
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	jwt, err := token.SignedString(rsaKey)
	if err != nil {
		zap.L().Error("Signing JWT", zap.Error(err))
	}
	return jwt, err
}

func generateNonce() string {
	nonceBytes := make([]byte, 64)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 64; i++ {
		nonceBytes[i] = byte(random.Int63() % 256)
	}
	return base64.URLEncoding.EncodeToString(nonceBytes)
}
