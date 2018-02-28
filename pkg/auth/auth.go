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
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/pkg/errors"
)

const gothProviderType = "openid-connect"
const sessionExpiryInMinutes = 15
const sessionRenewalTimeInMinutes = sessionExpiryInMinutes - 1

// JwtCookieName is the key we're storing cookies under in our cookie store
const JwtCookieName = "JWT_COOKIE"

// TODO: Use a stronger auth key for this
var cookieStore = sessions.NewCookieStore([]byte("supersecretstring"))

// UserClaims does a thing
type UserClaims struct {
	Email   string `json:"email"`
	IDToken string `json:"id_token"`
	jwt.StandardClaims
}

func getExpiryTimeFromMinutes(min int64) int64 {
	return time.Now().Add(time.Second * time.Duration(60*min)).Unix()
}

func shouldRenewForClaims(claims UserClaims) bool {
	exp := claims.StandardClaims.ExpiresAt
	renewal := getExpiryTimeFromMinutes(sessionRenewalTimeInMinutes)
	fmt.Println(exp)
	fmt.Println(renewal)
	return exp < renewal
}

func parseClaimsFromTokenString(ss string, secret string) (claims *UserClaims, err error) {
	token, err := jwt.ParseWithClaims(ss, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(secret))
		return &rsaKey.PublicKey, err
	})
	if err != nil {
		return
	}
	if !token.Valid {
		err = errors.New("Token failed validation")
		return
	}

	claims, ok := token.Claims.(*UserClaims)
	if claims == nil || !ok {
		err = errors.New("Failed extracting claims from token")
	}

	return
}

func signedTokenStringWithUserInfo(email string, idToken string, expiry int64, secret string) (ss string, err error) {
	claims := UserClaims{
		email,
		idToken,
		jwt.StandardClaims{
			ExpiresAt: expiry,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(secret))
	if err != nil {
		err = errors.Wrap(err, "Parsing RSA key from PEM")
		return
	}

	ss, err = token.SignedString(rsaKey)
	if err != nil {
		err = errors.Wrap(err, "Signing string with token")
		return
	}

	return
}

func saveTokenStringToStore(ss string, w http.ResponseWriter, r *http.Request) (err error) {
	// This will return a new session if non exists
	s, err := cookieStore.Get(r, JwtCookieName)
	if err != nil {
		err = errors.Wrap(err, "Fetching session from cookie store")
		return
	}

	// Flashes are one-time use data, so retrieving them clears the Flash
	// store so we start with a clean slate
	s.Flashes()
	s.AddFlash(ss)

	err = s.Save(r, w)
	if err != nil {
		err = errors.Wrap(err, "Saving session back to store")
		return
	}

	return
}

// RegisterProvider registers Login.gov with Goth, which uses
// auto-discovery to get the OpenID configuration
func RegisterProvider(logger *zap.Logger, loginGovSecretKey, hostname, loginGovClientID string) {
	if loginGovSecretKey == "" {
		logger.Warn("Login.gov secret key must be set.")
	}
	provider, err := openidConnect.New(
		loginGovClientID,
		loginGovSecretKey,
		fmt.Sprintf("%s/auth/login-gov/callback", hostname),
		"https://idp.int.identitysandbox.gov/.well-known/openid-configuration",
	)

	if err != nil {
		logger.Error("Register Login.gov provider with Goth", zap.Error(err))
	}

	if provider != nil {
		goth.UseProviders(provider)
	}
}

// AuthorizationRedirectHandler handles redirection
type AuthorizationRedirectHandler struct {
	logger *zap.Logger
}

// UserAuthMiddleware attempts to decrypt the provided token or redirects to landing page
func UserAuthMiddleware(clientAuthSecretKey string, hostname string) func(next http.Handler) http.Handler {
	redirectURL := fmt.Sprintf("%s/landing", hostname)
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			s, err := cookieStore.Get(r, JwtCookieName)
			if err != nil {
				zap.L().Error("Getting session from store", zap.Error(err))
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}

			// Our token is stored as a flash message under the default key ("_flash")
			flashes := s.Flashes()
			if s.IsNew || len(flashes) == 0 {
				fmt.Println("No prior token found")
				http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
				return
			}

			// Exract user info from the JWT
			ss := flashes[0].(string)
			claims, err := parseClaimsFromTokenString(ss, clientAuthSecretKey)
			if err != nil {
				fmt.Println("Parsing claims from token", zap.Error(err))
				http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
				return
			}

			if shouldRenewForClaims(*claims) {
				// Renew the token
				expiry := getExpiryTimeFromMinutes(sessionExpiryInMinutes)
				ss, err := signedTokenStringWithUserInfo(claims.Email, claims.IDToken, expiry, clientAuthSecretKey)
				if err != nil {
					zap.L().Error("Generating signed token string", zap.Error(err))
				}
				if err := saveTokenStringToStore(ss, w, r); err != nil {
					zap.L().Error("Saving token to store", zap.Error(err))
				}
			}

			// And put the user info on the request context
			context.Set(r, "email", claims.Email)
			context.Set(r, "id_token", claims.IDToken)

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(mw)
	}
}

// NewAuthorizationRedirectHandler creates a new AuthorizationRedirectHandler
func NewAuthorizationRedirectHandler(logger *zap.Logger) *AuthorizationRedirectHandler {
	handler := AuthorizationRedirectHandler{
		logger: logger,
	}
	return &handler
}

// AuthorizationLogoutHandler handles logging the user out of login.gov
func AuthorizationLogoutHandler(hostname string) http.HandlerFunc {
	logoutURL := "https://idp.int.identitysandbox.gov/openid_connect/logout"

	return func(w http.ResponseWriter, r *http.Request) {
		parsedURL, err := url.Parse(logoutURL)
		if err != nil {
			zap.L().Error("Parse logout URL", zap.Error(err))
		}

		s, err := cookieStore.Get(r, JwtCookieName)
		if err != nil || s.IsNew {
			zap.L().Error("Getting session from store", zap.Error(err))
		}
		// This kills the cookie
		s.Options.MaxAge = -1
		s.Save(r, w)

		redirectURL := fmt.Sprintf("%s/landing", hostname)
		idToken := context.Get(r, "id_token").(string)

		params := parsedURL.Query()
		params.Add("id_token_hint", idToken)
		params.Add("post_logout_redirect_uri", redirectURL)
		params.Set("state", generateNonce())

		parsedURL.RawQuery = params.Encode()

		http.Redirect(w, r, parsedURL.String(), http.StatusTemporaryRedirect)
	}
}

// AuthorizationRedirectHandler constructs the Login.gov authentication URL and redirects to it
func (h *AuthorizationRedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url, err := getAuthorizationURL(h.logger)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// AuthorizationCallbackHandler processes a callback from login.gov
type AuthorizationCallbackHandler struct {
	loginGovSecretKey string
	loginGovClientID  string
	hostname          string
	logger            *zap.Logger
}

// NewAuthorizationCallbackHandler creates a new AuthorizationCallbackHandler
func NewAuthorizationCallbackHandler(loginGovSecretKey string, loginGovClientID string, hostname string, logger *zap.Logger) *AuthorizationCallbackHandler {
	handler := AuthorizationCallbackHandler{
		loginGovSecretKey: loginGovSecretKey,
		loginGovClientID:  loginGovClientID,
		hostname:          hostname,
		logger:            logger,
	}
	return &handler
}

// AuthorizationCallbackHandler handles the callback from the Login.gov authorization flow
func (h *AuthorizationCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	authError := r.URL.Query().Get("error")

	// The user has either cancelled or declined to authorize the client
	if authError == "access_denied" {
		http.Redirect(w, r, fmt.Sprintf("%s/landing", h.hostname), http.StatusTemporaryRedirect)
		return
	}

	if authError == "invalid_request" {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	provider, err := goth.GetProvider(gothProviderType)
	if err != nil {
		h.logger.Error("Get Goth provider", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// TODO: validate the state is the same (pull from session)
	code := r.URL.Query().Get("code")
	session, err := fetchToken(h.logger, code, h.loginGovSecretKey, h.loginGovClientID)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	user, err := provider.FetchUser(session)
	if err != nil {
		h.logger.Error("Login.gov user info request", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	expiry := getExpiryTimeFromMinutes(sessionExpiryInMinutes)
	ss, err := signedTokenStringWithUserInfo(user.Email, session.IDToken, expiry, clientAuthSecretKey)
	if err != nil {
		zap.L().Error("Generating signed token string", zap.Error(err))
	}

	if err := saveTokenStringToStore(ss, w, r); err != nil {
		zap.L().Error("Saving token to store", zap.Error(err))
	}

	landingURL := fmt.Sprintf("%s/landing?email=%s", h.hostname, user.RawData["email"])
	http.Redirect(w, r, landingURL, http.StatusTemporaryRedirect)
}

func getAuthorizationURL(logger *zap.Logger) (string, error) {
	provider, err := goth.GetProvider(gothProviderType)
	if err != nil {
		logger.Error("Get Goth provider", zap.Error(err))
		return "", err
	}
	state := generateNonce()
	sess, err := provider.BeginAuth(state)
	if err != nil {
		logger.Error("Goth begin auth", zap.Error(err))
		return "", err
	}

	baseURL, err := sess.GetAuthURL()
	if err != nil {
		logger.Error("Goth get auth URL", zap.Error(err))
		return "", err
	}

	authURL, err := url.Parse(baseURL)
	if err != nil {
		logger.Error("Parse auth URL", zap.Error(err))
		return "", err
	}

	params := authURL.Query()
	params.Add("acr_values", "http://idmanagement.gov/ns/assurance/loa/1")
	params.Add("nonce", state)
	params.Set("scope", "openid email")

	authURL.RawQuery = params.Encode()
	return authURL.String(), err
}

func fetchToken(logger *zap.Logger, code string, loginGovSecretKey string, loginGovClientID string) (*openidConnect.Session, error) {
	// TODO: Get the token endpoint URL from Goth instead when
	// https://github.com/markbates/goth/pull/207 is resolved
	tokenURL := "https://idp.int.identitysandbox.gov/api/openid_connect/token"
	clientAssertion, err := createClientAssertionJWT(logger, tokenURL, loginGovSecretKey, loginGovClientID)
	if err != nil {
		return nil, err
	}

	params := url.Values{
		"client_assertion":      {clientAssertion},
		"client_assertion_type": {"urn:ietf:params:oauth:client-assertion-type:jwt-bearer"},
		"code":                  {code},
		"grant_type":            {"authorization_code"},
	}

	response, err := http.PostForm(tokenURL, params)
	if err != nil {
		logger.Error("Post to Login.gov token endpoint", zap.Error(err))
		return nil, err
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error("Reading Login.gov token response", zap.Error(err))
		return nil, err
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

func createClientAssertionJWT(logger *zap.Logger, tokenURL, loginGovSecretKey, loginGovClientID string) (string, error) {
	claims := &jwt.StandardClaims{
		Issuer:    loginGovClientID,
		Subject:   loginGovClientID,
		Audience:  tokenURL,
		Id:        generateNonce(),
		ExpiresAt: time.Now().Add(time.Minute * sessionExpiryInMinutes).Unix(),
	}

	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(loginGovSecretKey))
	if err != nil {
		logger.Error("JWT parse private key from PEM", zap.Error(err))
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	jwt, err := token.SignedString(rsaKey)
	if err != nil {
		logger.Error("Signing JWT", zap.Error(err))
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
