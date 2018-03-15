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

	"github.com/dgrijalva/jwt-go"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth/context"
	"github.com/transcom/mymove/pkg/models"
)

const gothProviderType = "openid-connect"
const sessionExpiryInMinutes = 15

// This sets a small window during which tokens won't be renewed
const sessionRenewalTimeInMinutes = sessionExpiryInMinutes - 1

// UserSessionCookieName is the key at which we're storing our token cookie
const UserSessionCookieName = "user_session"

// UserClaims wraps StandardClaims with some user info we care about
type UserClaims struct {
	UserID  uuid.UUID `json:"user_id"`
	Email   string    `json:"email"`
	IDToken string    `json:"id_token"`
	jwt.StandardClaims
}

func landingURL(hostname string) string {
	return fmt.Sprintf("%s", hostname)
}

func getExpiryTimeFromMinutes(min int64) time.Time {
	return time.Now().Add(time.Minute * time.Duration(min))
}

// Returns true if the expiration time is inside the renewal window
func shouldRenewForClaims(claims UserClaims) bool {
	exp := claims.StandardClaims.ExpiresAt
	renewal := getExpiryTimeFromMinutes(sessionRenewalTimeInMinutes).Unix()
	return exp < renewal
}

func signTokenStringWithUserInfo(userID uuid.UUID, email string, idToken string, expiry time.Time, secret string) (ss string, err error) {
	claims := UserClaims{
		userID,
		email,
		idToken,
		jwt.StandardClaims{
			ExpiresAt: expiry.Unix(),
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

	return ss, err
}

func deleteCookie(w http.ResponseWriter, name string) {
	// Not all browsers support MaxAge, so set Expires too
	cookie := http.Cookie{
		Name:    name,
		Value:   "blank",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	}
	http.SetCookie(w, &cookie)
}

func getUserClaimsFromRequest(logger *zap.Logger, secret string, r *http.Request) (claims *UserClaims, ok bool) {
	cookie, err := r.Cookie(UserSessionCookieName)
	if err != nil {
		// No cookie set on client
		return
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(secret))
		return &rsaKey.PublicKey, err
	})

	if token == nil || !token.Valid {
		logger.Error("Failed token validation", zap.Error(err))
		return
	}

	// The token actually just stores a Claims interface, so we need to explicitly
	// cast back to UserClaims
	claims, ok = token.Claims.(*UserClaims)
	if !ok {
		logger.Error("Failed getting claims from token", zap.Error(err))
		return
	}

	return claims, ok
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

// UserAuthMiddleware attempts to populate user data onto request context
func UserAuthMiddleware(logger *zap.Logger, secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			claims, ok := getUserClaimsFromRequest(logger, secret, r)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			if shouldRenewForClaims(*claims) {
				// Renew the token
				expiry := getExpiryTimeFromMinutes(sessionExpiryInMinutes)
				ss, err := signTokenStringWithUserInfo(claims.UserID, claims.Email, claims.IDToken, expiry, secret)
				if err != nil {
					logger.Error("Generating signed token string", zap.Error(err))
				}
				cookie := http.Cookie{
					Name:    UserSessionCookieName,
					Value:   ss,
					Path:    "/",
					Expires: expiry,
				}
				http.SetCookie(w, &cookie)
			}

			// And put the user info on the request context
			ctx := context.PopulateAuthContext(r.Context(), claims.UserID, claims.IDToken)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(mw)
	}
}

// AuthorizationContext is the common handler type for auth handlers
type AuthorizationContext struct {
	hostname string
	logger   *zap.Logger
}

// NewAuthContext creates an AuthorizationContext
func NewAuthContext(hostname string, logger *zap.Logger) AuthorizationContext {
	context := AuthorizationContext{
		hostname: hostname,
		logger:   logger,
	}
	return context
}

// AuthorizationLogoutHandler handles logging the user out of login.gov
type AuthorizationLogoutHandler AuthorizationContext

func (h AuthorizationLogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logoutURL := "https://idp.int.identitysandbox.gov/openid_connect/logout"
	redirectURL := landingURL(h.hostname)

	idToken, ok := context.GetIDToken(r.Context())
	if !ok {
		// Can't log out of login.gov without a token, redirect and let them re-auth
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
		return
	}

	parsedURL, err := url.Parse(logoutURL)
	if err != nil {
		h.logger.Error("Parse logout URL", zap.Error(err))
	}

	// Parameters taken from https://developers.login.gov/oidc/#logout
	params := parsedURL.Query()
	params.Add("id_token_hint", idToken)
	params.Add("post_logout_redirect_uri", redirectURL)
	params.Set("state", generateNonce())
	parsedURL.RawQuery = params.Encode()

	// Also need to clear the cookie on the client
	deleteCookie(w, UserSessionCookieName)

	http.Redirect(w, r, parsedURL.String(), http.StatusTemporaryRedirect)
}

// AuthorizationRedirectHandler handles redirection
type AuthorizationRedirectHandler AuthorizationContext

// AuthorizationRedirectHandler constructs the Login.gov authentication URL and redirects to it
func (h AuthorizationRedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, ok := context.GetIDToken(r.Context())
	if ok {
		// User is already authed, redirect to landing page
		http.Redirect(w, r, landingURL(h.hostname), http.StatusTemporaryRedirect)
		return
	}

	url, err := getAuthorizationURL(h.logger)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// AuthorizationCallbackHandler processes a callback from login.gov
type AuthorizationCallbackHandler struct {
	db                  *pop.Connection
	clientAuthSecretKey string
	loginGovSecretKey   string
	loginGovClientID    string
	hostname            string
	logger              *zap.Logger
}

// NewAuthorizationCallbackHandler creates a new AuthorizationCallbackHandler
func NewAuthorizationCallbackHandler(db *pop.Connection, clientAuthSecretKey string, loginGovSecretKey string, loginGovClientID string, hostname string, logger *zap.Logger) AuthorizationCallbackHandler {
	handler := AuthorizationCallbackHandler{
		db:                  db,
		clientAuthSecretKey: clientAuthSecretKey,
		loginGovSecretKey:   loginGovSecretKey,
		loginGovClientID:    loginGovClientID,
		hostname:            hostname,
		logger:              logger,
	}
	return handler
}

// AuthorizationCallbackHandler handles the callback from the Login.gov authorization flow
func (h AuthorizationCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	authError := r.URL.Query().Get("error")

	// The user has either cancelled or declined to authorize the client
	if authError == "access_denied" {
		http.Redirect(w, r, landingURL(h.hostname), http.StatusTemporaryRedirect)
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

	openIDuser, err := provider.FetchUser(session)
	if err != nil {
		h.logger.Error("Login.gov user info request", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	user, err := models.GetOrCreateUser(h.db, openIDuser)
	if err != nil {
		h.logger.Error("Unable to create user.", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// Sign a token and save it as a cookie on the client
	expiry := getExpiryTimeFromMinutes(sessionExpiryInMinutes)
	ss, err := signTokenStringWithUserInfo(user.ID, user.LoginGovEmail, session.IDToken, expiry, h.clientAuthSecretKey)
	if err != nil {
		h.logger.Error("Generating signed token string", zap.Error(err))
	}
	cookie := http.Cookie{
		Name:    UserSessionCookieName,
		Value:   ss,
		Path:    "/",
		Expires: expiry,
	}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, landingURL(h.hostname), http.StatusTemporaryRedirect)
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
