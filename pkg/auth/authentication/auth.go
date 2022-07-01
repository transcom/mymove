package authentication

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/notifications"

	"github.com/alexedwards/scs/v2"
	"github.com/gofrs/uuid"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

// IsLoggedInMiddleware handles requests to is_logged_in endpoint by returning true if someone is logged in
func IsLoggedInMiddleware(globalLogger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logging.FromContext(r.Context())
		data := map[string]interface{}{
			"isLoggedIn": false,
		}

		session := auth.SessionFromRequestContext(r)
		if session != nil && session.UserID != uuid.Nil {
			data["isLoggedIn"] = true
			logger.Info("Valid session, user logged in")
		}

		newEncoderErr := json.NewEncoder(w).Encode(data)
		if newEncoderErr != nil {
			logger.Error("Failed encoding is_logged_in check response", zap.Error(newEncoderErr))
		}
	}
}

type APIWithContext interface {
	Context() *middleware.Context
}

func PermissionsMiddleware(appCtx appcontext.AppContext, api APIWithContext) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {

			logger := logging.FromContext(r.Context())
			session := auth.SessionFromRequestContext(r)

			route, r, _ := api.Context().RouteInfo(r)
			if route == nil {
				// If we reach this error, something went wrong with the swagger router initialization, in reality will probably never be an issue except potentially in local testing
				logger.Error("Route not found on permission lookup")
				http.Error(w, http.StatusText(400), http.StatusBadRequest)
				return
			}

			permissionsRequired, exists := route.Operation.VendorExtensible.Extensions["x-permissions"]

			// no permissions defined on the route, we can move on
			if !exists {
				logger.Info("No permissions required on this route")
				next.ServeHTTP(w, r)
				return
			}

			// transform the object so we can iterate over permissions
			permissionsRequiredAsInterfaceArray := permissionsRequired.([]interface{})

			for _, v := range permissionsRequiredAsInterfaceArray {
				permission := v.(string)
				logger.Info("Permission required: ", zap.String("permission", permission))
				access, err := checkUserPermission(appCtx, session, permission)

				if err != nil {
					logger.Error("Unexpected error looking up permissions", zap.String("permission error", err.Error()))
					http.Error(w, http.StatusText(500), http.StatusInternalServerError)
					return
				}

				if !access {
					logger.Warn("Permission denied", zap.String("permission", permission))
					http.Error(w, http.StatusText(401), http.StatusUnauthorized)
					return
				}
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(mw)
	}
}

// UserAuthMiddleware enforces that the incoming request is tied to a user session
func UserAuthMiddleware(globalLogger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {

			logger := logging.FromContext(r.Context())
			session := auth.SessionFromRequestContext(r)

			// We must have a logged in session and a user
			if session == nil || session.UserID == uuid.Nil {
				logger.Error("unauthorized access, no session token or user id")
				http.Error(w, http.StatusText(401), http.StatusUnauthorized)
				return
			}

			// DO NOT CHECK MILMOVE SESSION BECAUSE NEW SERVICE MEMBERS WON'T HAVE AN ID RIGHT AWAY
			// This must be the right type of user for the application
			if session.IsOfficeApp() && !session.IsOfficeUser() {
				logger.Error("unauthorized user for office.move.mil", zap.String("user_id", session.UserID.String()))
				http.Error(w, http.StatusText(403), http.StatusForbidden)
				return
			} else if session.IsAdminApp() && !session.IsAdminUser() {
				logger.Error("unauthorized user for admin.move.mil", zap.String("user_id", session.UserID.String()))
				http.Error(w, http.StatusText(403), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(mw)
	}
}

func updateUserCurrentSessionID(appCtx appcontext.AppContext, sessionID string) error {
	userID := appCtx.Session().UserID

	user, err := models.GetUser(appCtx.DB(), userID)
	if err != nil {
		appCtx.Logger().Error("Fetching user", zap.String("user_id", userID.String()), zap.Error(err))
	}

	if appCtx.Session().IsAdminUser() {
		user.CurrentAdminSessionID = sessionID
		appCtx.Logger().Info("User is an Admin user")
	} else if appCtx.Session().IsOfficeUser() {
		user.CurrentOfficeSessionID = sessionID
		appCtx.Logger().Info("User is an Office user")
	} else if appCtx.Session().IsServiceMember() {
		user.CurrentMilSessionID = sessionID
		appCtx.Logger().Info("User is an Service Member user")
	}

	err = appCtx.DB().Save(user)
	if err != nil {
		appCtx.Logger().Error("Updating user's current_x_session_id", zap.String("user_id", appCtx.Session().UserID.String()), zap.Error(err))
		return err
	}

	return err
}

func resetUserCurrentSessionID(appCtx appcontext.AppContext) error {
	userID := appCtx.Session().UserID
	user, err := models.GetUser(appCtx.DB(), userID)
	if err != nil {
		appCtx.Logger().Error("Fetching user", zap.String("user_id", userID.String()), zap.Error(err))
	}

	if appCtx.Session().IsAdminUser() {
		user.CurrentAdminSessionID = ""
	} else if appCtx.Session().IsOfficeUser() {
		user.CurrentOfficeSessionID = ""
	} else if appCtx.Session().IsServiceMember() {
		user.CurrentMilSessionID = ""
	}
	err = appCtx.DB().Save(user)
	if err != nil {
		appCtx.Logger().Error("Updating user's current_x_session_id", zap.String("user_id", appCtx.Session().UserID.String()), zap.Error(err))
		return err
	}

	return err
}

func currentUser(appCtx appcontext.AppContext) (*models.User, error) {
	userID := appCtx.Session().UserID
	user, err := models.GetUser(appCtx.DB(), userID)
	if err != nil {
		appCtx.Logger().Error("Getting the user", zap.String("user_id", appCtx.Session().UserID.String()), zap.Error(err))
		return nil, err
	}

	return user, nil
}

func currentSessionID(session *auth.Session, user *models.User) string {
	if session.IsAdminUser() {
		return user.CurrentAdminSessionID
	} else if session.IsOfficeUser() {
		return user.CurrentOfficeSessionID
	} else if session.IsServiceMember() {
		return user.CurrentMilSessionID
	}

	return ""
}

func authenticateUser(ctx context.Context, appCtx appcontext.AppContext, sessionManager *scs.SessionManager) error {
	// The session token must be renewed during sign in to prevent
	// session fixation attacks
	err := sessionManager.RenewToken(ctx)
	if err != nil {
		appCtx.Logger().Error("Error renewing session token", zap.Error(err))
		return err
	}
	sessionID, _, err := sessionManager.Commit(ctx)
	if err != nil {
		appCtx.Logger().Error("Failed to write new user session to store", zap.Error(err))
		return err
	}
	sessionManager.Put(ctx, "session", appCtx.Session())

	user, err := currentUser(appCtx)
	if err != nil {
		appCtx.Logger().Error("Fetching user", zap.String("user_id", appCtx.Session().UserID.String()), zap.Error(err))
		return err
	}
	// Check to see if sessionID is set on the user, presently
	existingSessionID := currentSessionID(appCtx.Session(), user)
	if existingSessionID != "" {
		appCtx.Logger().Info("SessionID is not set on the current user", zap.String("user_id", appCtx.Session().UserID.String()))

		// Lookup the old session that wasn't logged out
		_, exists, err := sessionManager.Store.Find(existingSessionID)
		if err != nil {
			appCtx.Logger().Error("Error loading previous session", zap.Error(err))
			return err
		}

		if !exists {
			appCtx.Logger().Info("Session expired", zap.String("user_id", appCtx.Session().UserID.String()))
		} else {
			appCtx.Logger().Info("Concurrent session detected. Will delete previous session.", zap.String("user_id", appCtx.Session().UserID.String()))

			// We need to delete the concurrent session.
			err := sessionManager.Store.Delete(existingSessionID)
			if err != nil {
				appCtx.Logger().Error("Error deleting previous session", zap.Error(err))
				return err
			}
		}
	}

	updateErr := updateUserCurrentSessionID(appCtx, sessionID)
	if updateErr != nil {
		appCtx.Logger().Error("Updating user's current session ID", zap.Error(updateErr))
		return updateErr
	}
	appCtx.Logger().Info("Logged in", zap.Any("session", appCtx.Session()))

	return nil
}

// AdminAuthMiddleware is middleware for admin authentication
func AdminAuthMiddleware(globalLogger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			logger := logging.FromContext(r.Context())
			session := auth.SessionFromRequestContext(r)

			if session == nil || !session.IsAdminUser() {
				logger.Error("unauthorized user for admin.move.mil")
				http.Error(w, http.StatusText(403), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(mw)
	}
}

// PrimeAuthorizationMiddleware is the prime authorization middleware
func PrimeAuthorizationMiddleware(globalLogger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			logger := logging.FromContext(r.Context())
			clientCert := ClientCertFromContext(r.Context())
			if clientCert == nil {
				logger.Error("unauthorized user for ghc prime")
				http.Error(w, http.StatusText(401), http.StatusUnauthorized)
				return
			}
			if !clientCert.AllowPrime {
				logger.Error("forbidden user for ghc prime")
				http.Error(w, http.StatusText(403), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(mw)
	}
}

// PrimeSimulatorAuthorizationMiddleware ensures only users with the
// prime simulator role can access the simulator
func PrimeSimulatorAuthorizationMiddleware(globalLogger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			logger := logging.FromContext(r.Context())
			session := auth.SessionFromRequestContext(r)
			if session == nil || !session.Roles.HasRole(roles.RoleTypePrimeSimulator) {
				logger.Error("forbidden user for prime simulator")
				http.Error(w, http.StatusText(403), http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(mw)
	}
}

func (context Context) landingURL(session *auth.Session) string {
	return fmt.Sprintf(context.callbackTemplate, session.Hostname)
}

// SetFeatureFlag sets a feature flag in the context
func (context *Context) SetFeatureFlag(flag FeatureFlag) {
	if context.featureFlags == nil {
		context.featureFlags = make(map[string]bool)
	}

	context.featureFlags[flag.Name] = flag.Active
}

// GetFeatureFlag gets a feature flag from the context
func (context *Context) GetFeatureFlag(flag string) bool {
	if value, ok := context.featureFlags[flag]; ok {
		return value
	}
	return false
}

// sessionManager returns the session manager corresponding to the current app.
// A user can be signed in at the same time across multiple apps.
func (context Context) sessionManager(session *auth.Session) *scs.SessionManager {
	if session.IsMilApp() {
		return context.sessionManagers[0]
	} else if session.IsAdminApp() {
		return context.sessionManagers[1]
	} else if session.IsOfficeApp() {
		return context.sessionManagers[2]
	}

	return nil
}

// Context is the common handler type for auth handlers
type Context struct {
	loginGovProvider LoginGovProvider
	callbackTemplate string
	featureFlags     map[string]bool
	sessionManagers  [3]*scs.SessionManager
}

// FeatureFlag holds the name of a feature flag and if it is enabled
type FeatureFlag struct {
	Name   string
	Active bool
}

// NewAuthContext creates an Context
func NewAuthContext(logger *zap.Logger, loginGovProvider LoginGovProvider, callbackProtocol string, callbackPort int, sessionManagers [3]*scs.SessionManager) Context {
	context := Context{
		loginGovProvider: loginGovProvider,
		callbackTemplate: fmt.Sprintf("%s://%%s:%d/", callbackProtocol, callbackPort),
		sessionManagers:  sessionManagers,
	}
	return context
}

// LogoutHandler handles logging the user out of login.gov
type LogoutHandler struct {
	Context
	handlers.HandlerConfig
}

// NewLogoutHandler creates a new LogoutHandler
func NewLogoutHandler(ac Context, hc handlers.HandlerConfig) LogoutHandler {
	logoutHandler := LogoutHandler{
		Context:       ac,
		HandlerConfig: hc,
	}
	return logoutHandler
}

func (h LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	appCtx := h.AppContextFromRequest(r)
	if appCtx.Session() != nil {
		redirectURL := h.landingURL(appCtx.Session())
		if appCtx.Session().IDToken != "" {
			var logoutURL string
			// All users logged in via devlocal-auth will have this IDToken. We
			// don't want to make a call to login.gov for a logout URL as it will
			// fail for devlocal-auth'ed users.
			if appCtx.Session().IDToken == "devlocal" {
				logoutURL = redirectURL
			} else {
				logoutURL = h.loginGovProvider.LogoutURL(redirectURL, appCtx.Session().IDToken)
			}
			err := resetUserCurrentSessionID(appCtx)
			if err != nil {
				appCtx.Logger().Error("failed to reset user's current_x_session_id")
			}
			err = h.sessionManager(appCtx.Session()).Destroy(r.Context())
			if err != nil {
				appCtx.Logger().Error("failed to destroy session")
			}
			auth.DeleteCSRFCookies(w)
			appCtx.Logger().Info("user logged out")
			fmt.Fprint(w, logoutURL)
		} else {
			// Can't log out of login.gov without a token, redirect and let them re-auth
			appCtx.Logger().Info("session exists but has an empty IDToken")

			if appCtx.Session().UserID != uuid.Nil {
				err := resetUserCurrentSessionID(appCtx)
				if err != nil {
					appCtx.Logger().Error("failed to reset user's current_x_session_id")
				}
			}

			err := h.sessionManager(appCtx.Session()).Destroy(r.Context())
			if err != nil {
				appCtx.Logger().Error("failed to destroy session", zap.Error(err))
			}

			auth.DeleteCSRFCookies(w)
			fmt.Fprint(w, redirectURL)
		}
	}
}

// loginStateCookieName is the name given to the cookie storing the encrypted Login.gov state nonce.
const loginStateCookieName = "lg_state"
const loginStateCookieTTLInSecs = 1800 // 30 mins to transit through login.gov.

// RedirectHandler handles redirection
type RedirectHandler struct {
	Context
	handlers.HandlerConfig
	UseSecureCookie bool
}

func NewRedirectHandler(ac Context, hc handlers.HandlerConfig, useSecureCookie bool) RedirectHandler {
	return RedirectHandler{
		Context:         ac,
		HandlerConfig:   hc,
		UseSecureCookie: useSecureCookie,
	}
}

func shaAsString(nonce string) string {
	s := sha256.Sum256([]byte(nonce))
	return hex.EncodeToString(s[:])
}

// StateCookieName returns the login.gov state cookie name
func StateCookieName(session *auth.Session) string {
	return fmt.Sprintf("%s_%s", string(session.ApplicationName), loginStateCookieName)
}

// RedirectHandler constructs the Login.gov authentication URL and redirects to it
func (h RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	appCtx := h.AppContextFromRequest(r)

	if appCtx.Session() != nil && appCtx.Session().UserID != uuid.Nil {
		// User is already authenticated, redirect to landing page
		http.Redirect(w, r, h.landingURL(appCtx.Session()), http.StatusTemporaryRedirect)
		return
	}

	loginData, err := h.loginGovProvider.AuthorizationURL(r)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// Hash the state/Nonce value sent to login.gov and set the result as an HttpOnly cookie
	// Check this when we return from login.gov
	if appCtx.Session() == nil {
		appCtx.Logger().Error("Session is nil, so cannot get hostname for state Cookie")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	stateCookie := http.Cookie{
		Name:     StateCookieName(appCtx.Session()),
		Value:    shaAsString(loginData.Nonce),
		Path:     "/",
		Expires:  time.Now().Add(time.Duration(loginStateCookieTTLInSecs) * time.Second),
		MaxAge:   loginStateCookieTTLInSecs,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.UseSecureCookie,
	}

	http.SetCookie(w, &stateCookie)
	appCtx.Logger().Info("Cookie has been set", zap.Any("stateCookie", stateCookie))
	http.Redirect(w, r, loginData.RedirectURL, http.StatusTemporaryRedirect)
	appCtx.Logger().Info("User has been redirected", zap.Any("redirectURL", loginData.RedirectURL))
}

// CallbackHandler processes a callback from login.gov
type CallbackHandler struct {
	Context
	handlers.HandlerConfig
	sender notifications.NotificationSender
}

// NewCallbackHandler creates a new CallbackHandler
func NewCallbackHandler(ac Context, hc handlers.HandlerConfig, sender notifications.NotificationSender) CallbackHandler {
	handler := CallbackHandler{
		Context:       ac,
		HandlerConfig: hc,
		sender:        sender,
	}
	return handler
}

// AuthorizationCallbackHandler handles the callback from the Login.gov authorization flow
func (h CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	appCtx := h.AppContextFromRequest(r)

	if appCtx.Session() == nil {
		appCtx.Logger().Error("Session missing")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	rawLandingURL := h.landingURL(appCtx.Session())

	landingURL, err := url.Parse(rawLandingURL)
	if err != nil {
		appCtx.Logger().Error("Error parsing landing URL")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	if err := r.URL.Query().Get("error"); len(err) > 0 {
		landingQuery := landingURL.Query()
		switch err {
		case "access_denied":
			// The user has either cancelled or declined to authorize the client
			appCtx.Logger().Error("ACCESS_DENIED error from login.gov")
		case "invalid_request":
			appCtx.Logger().Error("INVALID_REQUEST error from login.gov")
			landingQuery.Add("error", "INVALID_REQUEST")
		default:
			appCtx.Logger().Error("unknown error from login.gov")
			landingQuery.Add("error", "UNKNOWN_ERROR")
		}
		landingURL.RawQuery = landingQuery.Encode()
		http.Redirect(w, r, landingURL.String(), http.StatusPermanentRedirect)
		appCtx.Logger().Info("User redirected from login.gov", zap.String("landingURL", landingURL.String()))

		return
	}

	// Check the state value sent back from login.gov with the value saved in the cookie
	returnedState := r.URL.Query().Get("state")
	stateCookieName := StateCookieName(appCtx.Session())
	stateCookie, err := r.Cookie(stateCookieName)
	if err != nil {
		appCtx.Logger().Error("Getting login.gov state cookie",
			zap.String("stateCookieName", stateCookieName),
			zap.String("sessionUserId", appCtx.Session().UserID.String()),
			zap.Error(err))
		http.Error(w, http.StatusText(403), http.StatusForbidden)
		return
	}

	hash := stateCookie.Value
	// case where user has 2 tabs open with different cookies
	if hash != shaAsString(returnedState) {
		appCtx.Logger().Error("State returned from Login.gov does not match state value stored in cookie",
			zap.String("state", returnedState),
			zap.String("cookie", hash),
			zap.String("hash", shaAsString(returnedState)))

		// Delete lg_state cookie
		auth.DeleteCookie(w, StateCookieName(appCtx.Session()))
		appCtx.Logger().Info("lg_state cookie deleted")

		// This operation will delete all cookies from the session
		err = h.sessionManager(appCtx.Session()).Destroy(r.Context())
		if err != nil {
			appCtx.Logger().Error("Deleting login.gov state cookie", zap.Error(err))
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		// set error query
		landingQuery := landingURL.Query()
		landingQuery.Add("error", "SIGNIN_ERROR")
		landingURL.RawQuery = landingQuery.Encode()
		http.Redirect(w, r, landingURL.String(), http.StatusTemporaryRedirect)
		appCtx.Logger().Info("User redirected", zap.String("landingURL", landingURL.String()))

		return
	}

	provider, err := getLoginGovProviderForRequest(r)
	if err != nil {
		appCtx.Logger().Error("Get Goth provider", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// TODO: validate the state is the same (pull from session)
	openIDSession, err := fetchToken(
		appCtx.Logger(),
		r.URL.Query().Get("code"),
		provider.ClientKey,
		h.loginGovProvider)
	if err != nil {
		appCtx.Logger().Error("Reading openIDSession from login.gov", zap.Error(err))
		http.Error(w, http.StatusText(401), http.StatusUnauthorized)
		return
	}

	openIDUser, err := provider.FetchUser(openIDSession)
	if err != nil {
		appCtx.Logger().Error("Login.gov user info request", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	appCtx.Session().IDToken = openIDSession.IDToken
	appCtx.Session().Email = openIDUser.Email

	appCtx.Logger().Info("New Login", zap.String("OID_User", openIDUser.UserID), zap.String("OID_Email", openIDUser.Email), zap.String("Host", appCtx.Session().Hostname))

	userIdentity, err := models.FetchUserIdentity(appCtx.DB(), openIDUser.UserID)
	if err == nil { // Someone we know already
		authorizeKnownUser(appCtx, userIdentity, h, w, r, landingURL.String())
		appCtx.Logger().Info("Authorized and known user detected", zap.String("OID_User", openIDUser.UserID), zap.String("OID_Email", openIDUser.Email))
		return
	} else if err == models.ErrFetchNotFound { // Never heard of them so far
		authorizeUnknownUser(appCtx, openIDUser, h, w, r, landingURL.String())
		appCtx.Logger().Error("Unknown user detected", zap.Error(err))
		return
	} else {
		appCtx.Logger().Error("Error loading Identity.", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
}

var authorizeKnownUser = func(appCtx appcontext.AppContext, userIdentity *models.UserIdentity, h CallbackHandler, w http.ResponseWriter, r *http.Request, lURL string) {
	if !userIdentity.Active {
		appCtx.Logger().Error("Inactive user requesting authentication",
			zap.String("application_name", string(appCtx.Session().ApplicationName)),
			zap.String("hostname", appCtx.Session().Hostname),
			zap.String("user_id", appCtx.Session().UserID.String()))
		http.Error(w, http.StatusText(403), http.StatusForbidden)
		return
	}
	appCtx.Session().Roles = append(appCtx.Session().Roles, userIdentity.Roles...)
	appCtx.Session().Permissions = getPermissionsForUser(appCtx, userIdentity.ID)

	appCtx.Session().UserID = userIdentity.ID
	if appCtx.Session().IsMilApp() && userIdentity.ServiceMemberID != nil {
		appCtx.Session().ServiceMemberID = *(userIdentity.ServiceMemberID)
	}

	if appCtx.Session().IsOfficeApp() {
		if userIdentity.OfficeActive != nil && !*userIdentity.OfficeActive {
			appCtx.Logger().Error("Office user is deactivated", zap.String("userID", appCtx.Session().UserID.String()))
			http.Error(w, http.StatusText(403), http.StatusForbidden)
			return
		}
		if userIdentity.OfficeUserID != nil {
			appCtx.Session().OfficeUserID = *(userIdentity.OfficeUserID)
		} else {
			// In case they managed to login before the office_user record was created
			officeUser, err := models.FetchOfficeUserByEmail(appCtx.DB(), appCtx.Session().Email)
			if err == models.ErrFetchNotFound {
				appCtx.Logger().Error("Non-office user authenticated at office site", zap.String("userID", appCtx.Session().UserID.String()))
				http.Error(w, http.StatusText(403), http.StatusForbidden)
				return
			} else if err != nil {
				appCtx.Logger().Error("Checking for office user", zap.String("userID", appCtx.Session().UserID.String()), zap.Error(err))
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
			appCtx.Session().OfficeUserID = officeUser.ID
			officeUser.UserID = &userIdentity.ID
			err = appCtx.DB().Save(officeUser)
			if err != nil {
				appCtx.Logger().Error("Updating office user", zap.String("userID", appCtx.Session().UserID.String()), zap.Error(err))
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
		}
	}

	if appCtx.Session().IsAdminApp() {
		if userIdentity.AdminUserActive != nil && !*userIdentity.AdminUserActive {
			appCtx.Logger().Error("Admin user is deactivated", zap.String("userID", appCtx.Session().UserID.String()))
			http.Error(w, http.StatusText(403), http.StatusForbidden)
			return
		}
		if userIdentity.AdminUserID != nil {
			appCtx.Session().AdminUserID = *(userIdentity.AdminUserID)
			appCtx.Session().AdminUserRole = userIdentity.AdminUserRole.String()
		} else {
			// In case they managed to login before the admin_user record was created
			var adminUser models.AdminUser
			queryBuilder := query.NewQueryBuilder()
			filters := []services.QueryFilter{
				query.NewQueryFilter("email", "=", strings.ToLower(userIdentity.Email)),
			}
			err := queryBuilder.FetchOne(appCtx, &adminUser, filters)

			if err != nil && errors.Cause(err).Error() == models.RecordNotFoundErrorString {
				appCtx.Logger().Error("No admin user found", zap.String("userID", appCtx.Session().UserID.String()))
				http.Error(w, http.StatusText(403), http.StatusForbidden)
				return
			} else if err != nil {
				appCtx.Logger().Error("Checking for admin user", zap.String("userID", appCtx.Session().UserID.String()), zap.Error(err))
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}

			appCtx.Session().AdminUserID = adminUser.ID
			appCtx.Session().AdminUserRole = adminUser.Role.String()
			adminUser.UserID = &userIdentity.ID
			verrs, err := appCtx.DB().ValidateAndSave(&adminUser)
			if err != nil {
				appCtx.Logger().Error("Updating admin user", zap.String("userID", appCtx.Session().UserID.String()), zap.Error(err))
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}

			if verrs != nil {
				appCtx.Logger().Error("Admin user validation errors", zap.String("userID", appCtx.Session().UserID.String()), zap.Error(verrs))
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
		}
	}
	appCtx.Session().FirstName = userIdentity.FirstName()
	appCtx.Session().LastName = userIdentity.LastName()
	appCtx.Session().Middle = userIdentity.Middle()

	sessionManager := h.sessionManager(appCtx.Session())
	authError := authenticateUser(r.Context(), appCtx, sessionManager)
	if authError != nil {
		appCtx.Logger().Error("Authenticating user", zap.Error(authError))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, lURL, http.StatusTemporaryRedirect)
}

var authorizeUnknownUser = func(appCtx appcontext.AppContext, openIDUser goth.User, h CallbackHandler, w http.ResponseWriter, r *http.Request, lURL string) {
	var officeUser *models.OfficeUser
	var user *models.User
	var err error

	// Loads the User and Roles associations of the office or admin user
	conn := appCtx.DB().Eager("User", "User.Roles")

	if appCtx.Session().IsOfficeApp() { // Look to see if we have OfficeUser with this email address
		officeUser, err = models.FetchOfficeUserByEmail(conn, appCtx.Session().Email)
		if err == models.ErrFetchNotFound {
			appCtx.Logger().Error("No Office user found", zap.String("userID", appCtx.Session().UserID.String()))
			http.Error(w, http.StatusText(403), http.StatusForbidden)
			return
		} else if err != nil {
			appCtx.Logger().Error("Checking for office user", zap.String("userID", appCtx.Session().UserID.String()), zap.Error(err))
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		if !officeUser.Active {
			appCtx.Logger().Error("Office user is deactivated", zap.String("userID", appCtx.Session().UserID.String()))
			http.Error(w, http.StatusText(403), http.StatusForbidden)
			return
		}
		user = &officeUser.User
	}

	var adminUser models.AdminUser
	if appCtx.Session().IsAdminApp() {
		queryBuilder := query.NewQueryBuilder()
		filters := []services.QueryFilter{
			query.NewQueryFilter("email", "=", appCtx.Session().Email),
		}
		err = queryBuilder.FetchOne(appCtx, &adminUser, filters)

		if err != nil && errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			appCtx.Logger().Error("No admin user found", zap.String("userID", appCtx.Session().UserID.String()))
			http.Error(w, http.StatusText(403), http.StatusForbidden)
			return
		} else if err != nil {
			appCtx.Logger().Error("Checking for admin user", zap.String("userID", appCtx.Session().UserID.String()), zap.Error(err))
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		if !adminUser.Active {
			appCtx.Logger().Error("Admin user is deactivated", zap.String("userID", appCtx.Session().UserID.String()))
			http.Error(w, http.StatusText(403), http.StatusForbidden)
			return
		}
		user = &adminUser.User
	}

	if appCtx.Session().IsMilApp() {
		user, err = models.CreateUser(appCtx.DB(), openIDUser.UserID, openIDUser.Email)
		if err == nil {
			sysAdminEmail := notifications.GetSysAdminEmail(h.sender)
			appCtx.Logger().Info(
				"New user account created through Login.gov",
				zap.String("newUserID", user.ID.String()),
			)
			email, emailErr := notifications.NewUserAccountCreated(appCtx, sysAdminEmail, user.ID, user.UpdatedAt)
			if emailErr == nil {
				sendErr := h.sender.SendNotification(appCtx, email)
				if sendErr != nil {
					appCtx.Logger().Error("Error sending user creation email", zap.Error(sendErr))
				}
			} else {
				appCtx.Logger().Error("Error creating user creation email", zap.Error(emailErr))
			}
		}
		// Create the user's service member now and add the ServiceMemberID to
		// the session to allow the user's `CurrentMilSessionId` field to be
		// populated. This field is only populated if `session.IsServiceMember()`
		// returns true, and it only returns true if the user has a service
		// member associated with it. Previously, the service member was created
		// after the auth flow was over, when the user was redirected to the
		// onboarding home page (via /src/sagas/onboarding.js). This meant that
		// on the very first sign in, a user's `CurrentMilSessionId` would be
		// empty, which was misleading and prevented us from revoking their session.
		newServiceMember := models.ServiceMember{
			UserID: user.ID,
		}
		smVerrs, smErr := models.SaveServiceMember(appCtx, &newServiceMember)
		if smVerrs.HasAny() || smErr != nil {
			appCtx.Logger().Error("Error creating service member for user", zap.Error(smErr))
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		appCtx.Session().ServiceMemberID = newServiceMember.ID
	} else {
		err = models.UpdateUserLoginGovUUID(appCtx.DB(), user, openIDUser.UserID)
	}

	if err != nil {
		appCtx.Logger().Error("Error updating/creating user", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	appCtx.Session().UserID = user.ID
	if appCtx.Session().IsOfficeApp() && officeUser != nil {
		appCtx.Session().OfficeUserID = officeUser.ID
	} else if appCtx.Session().IsAdminApp() && adminUser.ID != uuid.Nil {
		appCtx.Session().AdminUserID = adminUser.ID
	}

	appCtx.Session().Roles = append(appCtx.Session().Roles, user.Roles...)

	sessionManager := h.sessionManager(appCtx.Session())
	authError := authenticateUser(r.Context(), appCtx, sessionManager)
	if authError != nil {
		appCtx.Logger().Error("Authenticate user", zap.Error(authError))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, lURL, http.StatusTemporaryRedirect)
}

func fetchToken(logger *zap.Logger, code string, clientID string, loginGovProvider LoginGovProvider) (*openidConnect.Session, error) {
	expiry := auth.GetExpiryTimeFromMinutes(auth.SessionExpiryInMinutes)
	params, err := loginGovProvider.TokenParams(code, clientID, expiry)
	if err != nil {
		logger.Error("Creating token endpoint params", zap.Error(err))
		return nil, err
	}

	response, err := http.PostForm(loginGovProvider.TokenURL(), params)
	if err != nil {
		logger.Error("Post to Login.gov token endpoint", zap.Error(err))
		return nil, err
	}

	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			logger.Error("Error in closing response", zap.Error(closeErr))
		}
	}()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error("Reading Login.gov token response", zap.Error(err))
		return nil, err
	}

	var parsedResponse LoginGovTokenResponse
	err = json.Unmarshal(responseBody, &parsedResponse)
	if err != nil {
		logger.Error("Parsing login.gov token", zap.Error(err))
		return nil, errors.Wrap(err, "parsing login.gov")
	}
	if parsedResponse.Error != "" {
		logger.Error("Error in Login.gov token response", zap.String("error", parsedResponse.Error))
		return nil, errors.New(parsedResponse.Error)
	}

	// TODO: get goth session from storage instead of constructing a new one
	session := openidConnect.Session{
		AccessToken: parsedResponse.AccessToken,
		ExpiresAt:   time.Now().Add(time.Second * time.Duration(parsedResponse.ExpiresIn)),
		IDToken:     parsedResponse.IDToken,
	}
	return &session, err
}

// InitAuth initializes the Login.gov provider
func InitAuth(v *viper.Viper, logger *zap.Logger, appnames auth.ApplicationServername) (LoginGovProvider, error) {
	loginGovCallbackProtocol := v.GetString(cli.LoginGovCallbackProtocolFlag)
	loginGovCallbackPort := v.GetInt(cli.LoginGovCallbackPortFlag)
	loginGovSecretKey := v.GetString(cli.LoginGovSecretKeyFlag)
	loginGovHostname := v.GetString(cli.LoginGovHostnameFlag)

	loginGovProvider := NewLoginGovProvider(loginGovHostname, loginGovSecretKey, logger)
	err := loginGovProvider.RegisterProvider(
		appnames.MilServername,
		v.GetString(cli.LoginGovMyClientIDFlag),
		appnames.OfficeServername,
		v.GetString(cli.LoginGovOfficeClientIDFlag),
		appnames.AdminServername,
		v.GetString(cli.LoginGovAdminClientIDFlag),
		loginGovCallbackProtocol,
		loginGovCallbackPort)
	return loginGovProvider, err
}
