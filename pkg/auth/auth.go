package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

const sessionExpiryInMinutes = 15

// This sets a small window during which tokens won't be renewed
const sessionRenewalTimeInMinutes = sessionExpiryInMinutes - 1

// UserSessionCookieName is the key at which we're storing our token cookie
const UserSessionCookieName = "user_session"

// Taken from answer here: https://stackoverflow.com/a/32620397
var maxPossibleTimeValue = time.Unix(1<<63-62135596801, 999999999)

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

// RequireAuthMiddleware enforces that the incoming request is tied to a user session
func RequireAuthMiddleware(next http.Handler) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {
		_, ok := context.GetUserID(r.Context())
		if !ok {
			http.Error(w, http.StatusText(401), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
		return
	}
	return http.HandlerFunc(mw)
}

// TokenParsingMiddleware attempts to populate user data onto request context
func TokenParsingMiddleware(logger *zap.Logger, secret string, noSessionTimeout bool) func(next http.Handler) http.Handler {
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
				// Never expire token if in development
				if noSessionTimeout {
					expiry = maxPossibleTimeValue
				}
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
	hostname         string
	logger           *zap.Logger
	loginGovProvider LoginGovProvider
}

// NewAuthContext creates an AuthorizationContext
func NewAuthContext(hostname string, logger *zap.Logger, loginGovProvider LoginGovProvider) AuthorizationContext {
	context := AuthorizationContext{
		hostname:         hostname,
		logger:           logger,
		loginGovProvider: loginGovProvider,
	}
	return context
}

// AuthorizationLogoutHandler handles logging the user out of login.gov
type AuthorizationLogoutHandler AuthorizationContext

func (h AuthorizationLogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	redirectURL := landingURL(h.hostname)

	idToken, ok := context.GetIDToken(r.Context())
	if !ok {
		// Can't log out of login.gov without a token, redirect and let them re-auth
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
		return
	}

	logoutURL := h.loginGovProvider.LogoutURL(redirectURL, idToken)

	// Also need to clear the cookie on the client
	deleteCookie(w, UserSessionCookieName)

	http.Redirect(w, r, logoutURL, http.StatusTemporaryRedirect)
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

	authURL, err := h.loginGovProvider.AuthorizationURL()
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// AuthorizationCallbackHandler processes a callback from login.gov
type AuthorizationCallbackHandler struct {
	db                  *pop.Connection
	clientAuthSecretKey string
	noSessionTimeout    bool
	hostname            string
	logger              *zap.Logger
	loginGovProvider    LoginGovProvider
}

// NewAuthorizationCallbackHandler creates a new AuthorizationCallbackHandler
func NewAuthorizationCallbackHandler(db *pop.Connection, clientAuthSecretKey string, noSessionTimeout bool, hostname string, logger *zap.Logger, loginGovProvider LoginGovProvider) AuthorizationCallbackHandler {
	handler := AuthorizationCallbackHandler{
		db:                  db,
		clientAuthSecretKey: clientAuthSecretKey,
		noSessionTimeout:    noSessionTimeout,
		hostname:            hostname,
		logger:              logger,
		loginGovProvider:    loginGovProvider,
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
	session, err := fetchToken(h.logger, code, h.loginGovProvider)
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
	// Never expire token if in development
	if h.noSessionTimeout {
		expiry = maxPossibleTimeValue
	}
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

func fetchToken(logger *zap.Logger, code string, loginGovProvider LoginGovProvider) (*openidConnect.Session, error) {
	tokenURL := loginGovProvider.TokenURL()
	expiry := getExpiryTimeFromMinutes(sessionExpiryInMinutes)
	params, err := loginGovProvider.TokenParams(code, expiry)
	if err != nil {
		logger.Error("Creating token endpoint params", zap.Error(err))
		return nil, err
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

	var parsedResponse LoginGovTokenResponse
	json.Unmarshal(responseBody, &parsedResponse)

	// TODO: get goth session from storage instead of constructing a new one
	session := openidConnect.Session{
		AccessToken: parsedResponse.AccessToken,
		ExpiresAt:   time.Now().Add(time.Second * time.Duration(parsedResponse.ExpiresIn)),
		IDToken:     parsedResponse.IDToken,
	}

	return &session, err
}
