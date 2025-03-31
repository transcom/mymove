package authentication

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/services"
	movelocker "github.com/transcom/mymove/pkg/services/lock_move"
	"github.com/transcom/mymove/pkg/services/query"
)

// used by authorizeKnownUser and authorizeUnknownUser
type AuthorizationResult byte

const (
	authorizationResultAuthorized AuthorizationResult = iota
	authorizationResultUnauthorized
	authorizationResultError
)

func (ar AuthorizationResult) String() string {
	return []string{
		"authorizationResultAuthorized",
		"authorizationResultUnauthorized",
		"authorizationResultError",
	}[ar]
}

// IsLoggedInMiddleware handles requests to is_logged_in endpoint by returning true if someone is logged in
func IsLoggedInMiddleware(_ *zap.Logger, maintenanceFlag bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logging.FromContext(r.Context())
		data := map[string]interface{}{
			"isLoggedIn":       false,
			"underMaintenance": maintenanceFlag,
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
func UserAuthMiddleware(_ *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {

			logger := logging.FromContext(r.Context())
			session := auth.SessionFromRequestContext(r)

			// We must have a logged in session and a user
			if session == nil {
				logger.Error("unauthorized access, no session token")
				http.Error(w, http.StatusText(401), http.StatusUnauthorized)
				return
			}
			if session.UserID == uuid.Nil {
				logger.Error("unauthorized access, no userid")
				http.Error(w, http.StatusText(401), http.StatusUnauthorized)
				return
			}

			// DO NOT CHECK MILMOVE SESSION BECAUSE WE'LL BE CHECKING THAT IN ANOTHER MIDDLEWARE
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

// CustomerAPIAuthMiddleware checks to see if the request matches one of the routes that should be allowed through with
// less strict authentication requirements. If it is on the allow list, it will allow the request to continue. If it
// is not, it will check to see if the user is a service member. Ideally, we will get rid of the allow list eventually
// and the service member check can be rolled into the UserAuthMiddleware.
func CustomerAPIAuthMiddleware(_ appcontext.AppContext, api APIWithContext) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			// Pulling logger & session from the request context instead of the app context because the app context
			// getting passed into here is the one from the server set up, not the one from the request that came in.
			logger := logging.FromContext(r.Context())
			session := auth.SessionFromRequestContext(r)

			route, r, _ := api.Context().RouteInfo(r)
			if route == nil {
				// If we reach this error, something went wrong with the swagger router initialization, in reality will probably never be an issue except potentially in local testing
				logger.Error("Route not found while checking authorization")
				http.Error(w, http.StatusText(400), http.StatusBadRequest)
				return
			}

			routeIsOnAllowList := checkIfRouteIsAllowed(route)

			if !routeIsOnAllowList && !session.IsServiceMember() {
				logger.Error("unauthorized user for my.move.mil", zap.String("user_id", session.UserID.String()))
				http.Error(w, http.StatusText(403), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(mw)
	}
}

// This is a list of routes (<tag>.<id>) that cannot be strictly checked. The goal is to get rid of this list entirely,
// so please try not to add routes to it.
var allowedRoutes = map[string]bool{
	"addresses.showAddress":                       true,
	"duty_locations.searchDutyLocations":          true,
	"featureFlags.booleanFeatureFlagForUser":      true,
	"featureFlags.variantFeatureFlagForUser":      true,
	"move_docs.createGenericMoveDocument":         true,
	"move_docs.deleteMoveDocument":                true,
	"move_docs.indexMoveDocuments":                true,
	"move_docs.updateMoveDocument":                true,
	"moves.showMove":                              true,
	"office.approveMove":                          true,
	"office.approveReimbursement":                 true,
	"office.cancelMove":                           true,
	"office.showOfficeOrders":                     true,
	"orders.showOrders":                           true,
	"orders.updateOrders":                         true,
	"postal_codes.validatePostalCodeWithRateData": true,
	"queues.showQueue":                            true,
	"uploads.deleteUpload":                        true,
	"users.showLoggedInUser":                      true,
	"okta_profile.showOktaInfo":                   true,
	"uploads.getUploadStatus":                     true,
}

// checkIfRouteIsAllowed checks to see if the route is one of the ones that should be allowed through without stricter
// checks. This is a temporary solution until we can implement robust permissions checks.
func checkIfRouteIsAllowed(route *middleware.MatchedRoute) bool {
	currentRouteTagID := fmt.Sprintf("%s.%s", route.Operation.OperationProps.Tags[0], route.Operation.OperationProps.ID)

	return allowedRoutes[currentRouteTagID]
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
	if userID == uuid.Nil {
		return nil, errors.New("No current user")
	}
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

func authenticateUser(ctx context.Context, appCtx appcontext.AppContext, sessionManager auth.SessionManager) error {
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
	appCtx.Logger().Info("User authenticated with new session", zap.String("new_session_id", sessionID))
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
		_, exists, err := sessionManager.Store().Find(existingSessionID)
		if err != nil {
			appCtx.Logger().Error("Error loading previous session", zap.Error(err))
			return err
		}

		if !exists {
			appCtx.Logger().Info("Session expired", zap.String("user_id", appCtx.Session().UserID.String()))
		} else {
			appCtx.Logger().Info("Concurrent session detected. Will delete previous session.", zap.String("user_id", appCtx.Session().UserID.String()))

			// We need to delete the concurrent session.
			err := sessionManager.Store().Delete(existingSessionID)
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
	appCtx.Logger().Info("Logged in",
		zap.Any("session.user_id", appCtx.Session().UserID),
		zap.Any("session.appname", appCtx.Session().ApplicationName))

	return nil
}

// AdminAuthMiddleware is middleware for admin authentication
func AdminAuthMiddleware(_ *zap.Logger) func(next http.Handler) http.Handler {
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
func PrimeAuthorizationMiddleware(_ *zap.Logger) func(next http.Handler) http.Handler {
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

// PPTASAuthorizationMiddleware is the PPTAS authorization middleware
func PPTASAuthorizationMiddleware(_ *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			logger := logging.FromContext(r.Context())
			clientCert := ClientCertFromContext(r.Context())
			if clientCert == nil {
				logger.Error("unauthorized user for PPTAS")
				http.Error(w, http.StatusText(401), http.StatusUnauthorized)
				return
			}
			if !clientCert.AllowPPTAS {
				logger.Error("forbidden user for PPTAS")
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
func PrimeSimulatorAuthorizationMiddleware(_ *zap.Logger) func(next http.Handler) http.Handler {
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

// Context is the common handler type for auth handlers
type Context struct {
	oktaProvider     okta.Provider
	callbackTemplate string
}

// FeatureFlag holds the name of a feature flag and if it is enabled
type FeatureFlag struct {
	Name   string
	Active bool
}

// NewAuthContext creates an Context
func NewAuthContext(_ *zap.Logger, oktaProvider okta.Provider, callbackProtocol string, callbackPort int) Context {
	context := Context{
		oktaProvider:     oktaProvider,
		callbackTemplate: fmt.Sprintf("%s://%%s:%d/", callbackProtocol, callbackPort),
	}
	return context
}

// LogoutHandler handles logging the user out of okta.mil
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

// logic for the /auth/logout endpoint
func (h LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	appCtx := h.AppContextFromRequest(r)
	provider, err := okta.GetOktaProviderForRequest(r)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	if appCtx.Session() != nil {
		// if the user is an office user, we need to unlock any moves that they have locked
		if appCtx.Session().IsOfficeApp() && appCtx.Session().OfficeUserID != uuid.Nil {
			moveUnlocker := movelocker.NewMoveUnlocker()
			officeUserID := appCtx.Session().OfficeUserID
			err := moveUnlocker.CheckForLockedMovesAndUnlock(appCtx, officeUserID)
			if err != nil {
				appCtx.Logger().Error(fmt.Sprintf("failed to unlock moves for office user ID: %s", officeUserID), zap.Error(err))
			}
		}

		sessionManager := h.SessionManagers().SessionManagerForApplication(appCtx.Session().ApplicationName)
		if sessionManager == nil {
			appCtx.Logger().Error("Authenticating user, cannot get session manager from request")
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		redirectURL := h.landingURL(appCtx.Session())
		if appCtx.Session().IDToken != "" {

			// storing ID token to use for /logout call to Okta
			userIDToken := appCtx.Session().IDToken

			// clearing okta.mil sessions by clearing Access Token & ID Token
			// this is shown in a sample app here: https://github.com/okta/samples-golang/blob/master/okta-hosted-login/main.go
			appCtx.Session().AccessToken = ""
			appCtx.Session().IDToken = ""

			// getting okta logout URL that will contain ID token and redirect
			oktaLogoutURL, err := logoutOktaUserURL(provider, userIDToken, redirectURL)
			if oktaLogoutURL == "" || err != nil {
				appCtx.Logger().Error("failed to get Okta Logout URL")
			}

			// Remember, UserID is UUID; however, the Okta ID is not.
			if appCtx.Session().UserID != uuid.Nil {
				err = resetUserCurrentSessionID(appCtx)
				if err != nil {
					appCtx.Logger().Error("failed to reset user's current_x_session_id")
				}
			}
			err = sessionManager.Destroy(r.Context())
			if err != nil {
				appCtx.Logger().Error("failed to destroy session")
			}
			auth.DeleteCSRFCookies(w)
			appCtx.Logger().Info("user logged out of application")
			fmt.Fprint(w, oktaLogoutURL)
		} else {
			// Can't log out of okta.mil without a token, redirect and let them re-auth
			appCtx.Logger().Info("session exists but has an empty IDToken")

			if appCtx.Session().UserID != uuid.Nil {
				err := resetUserCurrentSessionID(appCtx)
				if err != nil {
					appCtx.Logger().Error("failed to reset user's current_x_session_id")
				}
			}

			err := sessionManager.Destroy(r.Context())
			if err != nil {
				appCtx.Logger().Error("failed to destroy session", zap.Error(err))
			}

			auth.DeleteCSRFCookies(w)
			fmt.Fprint(w, redirectURL)
		}
	}
}

// LogoutOktaRedirectHandler handles logging the user out of okta.mil
// and then redirecting the user BACK to the sign in page
// this will be used for customers that are required to authenticate with CAC first
type LogoutOktaRedirectHandler struct {
	Context
	handlers.HandlerConfig
}

// NewLogoutOktaRedirectHandler creates a new NewLogoutOktaRedirectHandler
func NewLogoutOktaRedirectHandler(ac Context, hc handlers.HandlerConfig) LogoutOktaRedirectHandler {
	logoutHandler := LogoutOktaRedirectHandler{
		Context:       ac,
		HandlerConfig: hc,
	}
	return logoutHandler
}

// logic for the /auth/logoutOktaRedirect endpoint
func (h LogoutOktaRedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	appCtx := h.AppContextFromRequest(r)
	provider, err := okta.GetOktaProviderForRequest(r)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	if appCtx.Session() != nil {
		sessionManager := h.SessionManagers().SessionManagerForApplication(appCtx.Session().ApplicationName)
		if sessionManager == nil {
			appCtx.Logger().Error("Authenticating user, cannot get session manager from request")
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		redirectURL := h.landingURL(appCtx.Session())
		if appCtx.Session().IDToken != "" {

			// storing ID token to use for logging the user out of Okta
			// this contains the user's token that Okta needs to clear their session
			userIDToken := appCtx.Session().IDToken

			// clearing okta.mil sessions by clearing Access Token & ID Token
			// this is shown in a sample app here: https://github.com/okta/samples-golang/blob/master/okta-hosted-login/main.go
			appCtx.Session().AccessToken = ""
			appCtx.Session().IDToken = ""

			// getting okta logout URL that will contain ID token and redirect back to the Okta sign in page via redirect
			oktaLogoutURL, err := logoutOktaUserURLWithRedirect(provider, userIDToken, redirectURL)
			if oktaLogoutURL == "" || err != nil {
				appCtx.Logger().Error("failed to get Okta Logout URL")
			}

			// Remember, UserID is UUID; however, the Okta ID is not.
			if appCtx.Session().UserID != uuid.Nil {
				err = resetUserCurrentSessionID(appCtx)
				if err != nil {
					appCtx.Logger().Error("failed to reset user's current_x_session_id")
				}
			}
			err = sessionManager.Destroy(r.Context())
			if err != nil {
				appCtx.Logger().Error("failed to destroy session")
			}
			auth.DeleteCSRFCookies(w)
			appCtx.Logger().Info("user logged out of application")
			fmt.Fprint(w, oktaLogoutURL)
		} else {
			// Can't log out of okta.mil without a token, redirect and let them re-auth
			appCtx.Logger().Info("session exists but has an empty IDToken")

			if appCtx.Session().UserID != uuid.Nil {
				err := resetUserCurrentSessionID(appCtx)
				if err != nil {
					appCtx.Logger().Error("failed to reset user's current_x_session_id")
				}
			}

			err := sessionManager.Destroy(r.Context())
			if err != nil {
				appCtx.Logger().Error("failed to destroy session", zap.Error(err))
			}

			auth.DeleteCSRFCookies(w)
			fmt.Fprint(w, redirectURL)
		}
	}
}

// loginStateCookieName is the name given to the cookie storing the encrypted okta.mil state nonce.
const loginStateCookieName = "okta_state"
const loginStateCookieTTLInSecs = 1800 // 30 mins to transit through okta.mil.

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

// StateCookieName returns the okta.mil state cookie name
func StateCookieName(session *auth.Session) string {
	return fmt.Sprintf("%s_%s", string(session.ApplicationName), loginStateCookieName)
}

// RedirectHandler constructs the okta.mil authentication URL and redirects to it
// This will be called when logging in
func (h RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	appCtx := h.AppContextFromRequest(r)

	if appCtx.Session() != nil && appCtx.Session().UserID != uuid.Nil {
		// User is already authenticated, redirect to landing page
		http.Redirect(w, r, h.landingURL(appCtx.Session()), http.StatusTemporaryRedirect)
		return
	}

	loginData, err := h.oktaProvider.AuthorizationURL(r)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// Hash the state/Nonce value sent to okta.mil and set the result as an HttpOnly cookie
	// Check this when we return from okta.mil
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

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// CallbackHandler processes a callback from okta.mil
type CallbackHandler struct {
	Context
	handlers.HandlerConfig
	sender     notifications.NotificationSender
	HTTPClient HTTPClient
}

// NewCallbackHandler creates a new CallbackHandler
func NewCallbackHandler(ac Context, hc handlers.HandlerConfig, sender notifications.NotificationSender) CallbackHandler {
	handler := CallbackHandler{
		Context:       ac,
		HandlerConfig: hc,
		sender:        sender,
		HTTPClient:    &http.Client{},
	}
	return handler
}

// invalidPermissionsResponse generates an http response when invalid
// permissions are encountered. It *also* saves the session
// information. This is needed so we have the necessary info to create
// a redirect to logout of okta.mil
func invalidPermissionsResponse(appCtx appcontext.AppContext, handlerConfig handlers.HandlerConfig, authContext Context, w http.ResponseWriter, r *http.Request) {

	sessionManager := handlerConfig.SessionManagers().SessionManagerForApplication(appCtx.Session().ApplicationName)
	_, _, err := sessionManager.Commit(r.Context())
	if err != nil {
		appCtx.Logger().Error("Failed to write invalid permissions user session to store", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	sessionManager.Put(r.Context(), "session", appCtx.Session())
	if err != nil {
		appCtx.Logger().Error("Error authenticating user with invalid permissions", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	rawLandingURL := authContext.landingURL(appCtx.Session()) + "invalid-permissions"
	landingURL, err := url.Parse(rawLandingURL)
	if err != nil {
		appCtx.Logger().Error("Error parsing invalid permissions url", zap.Any("url", rawLandingURL))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	traceID := handlerConfig.GetTraceIDFromRequest(r)
	if !traceID.IsNil() {
		landingQuery := landingURL.Query()
		landingQuery.Add("traceId", traceID.String())
		landingURL.RawQuery = landingQuery.Encode()
	}

	// We need to redirect here because we got to this handler after a
	// redirect from okta.mil. Our client application did not make
	// this request, so we need to redirect to the client app so that
	// we can present a "pretty" error page to the user
	appCtx.Logger().Info("Redirect invalid permissions",
		zap.String("request_path", r.URL.Path),
		zap.String("redirect_url", landingURL.String()))
	http.Redirect(w, r, landingURL.String(), http.StatusTemporaryRedirect)
}

type MockHTTPClient struct {
	Response *http.Response
	Err      error
}

func (m *MockHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}

// AuthorizationCallbackHandler handles the callback from the Okta.mil authorization flow
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

	// if the user is still signed into okta and has an active session in the browser
	// but is being forced to authenticate/re-authenticate from MilMove, we need to handle logout or let the user know they need to log out
	// so they can re-use their authenticator (CAC)
	errDescription := r.URL.Query().Get("error_description")
	// this is the description okta sends when the user has used all of their authenticators
	if errDescription == "The resource owner or authorization server denied the request." {
		provider, providerErr := okta.GetOktaProviderForRequest(r)
		if providerErr != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		// if the user just closed their tab and appCtx is still holding the ID token, we can use it
		// MM will still have the active IDToken and we can use that to log them out and clear their session
		if appCtx.Session().IDToken != "" {
			oktaLogoutURL, logoutErr := logoutOktaUserURL(provider, appCtx.Session().IDToken, landingURL.String())
			if oktaLogoutURL == "" || logoutErr != nil {
				appCtx.Logger().Error("failed to get Okta Logout URL")
			}
			http.Redirect(w, r, oktaLogoutURL, http.StatusTemporaryRedirect)
			return
		}
		// if not, we will need the user to go to okta and sign out, adding these params will display a UI info banner
		redirectURL := landingURL.String() + "sign-in" + "?okta_error=true"
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
		return
	}

	if err := r.URL.Query().Get("error"); len(err) > 0 {
		landingQuery := landingURL.Query()
		switch err {
		case "access_denied":
			// The user has either cancelled or declined to authorize the client
			appCtx.Logger().Error("ACCESS_DENIED error from okta.mil")
		case "invalid_request":
			appCtx.Logger().Error("INVALID_REQUEST error from okta.mil")
			landingQuery.Add("error", "INVALID_REQUEST")
		default:
			appCtx.Logger().Error("unknown error from okta.mil")
			landingQuery.Add("error", "UNKNOWN_ERROR")
		}
		landingURL.RawQuery = landingQuery.Encode()
		http.Redirect(w, r, landingURL.String(), http.StatusTemporaryRedirect)
		appCtx.Logger().Info("User redirected from okta.mil", zap.String("landingURL", landingURL.String()))

		return
	}
	sessionManager := h.SessionManagers().SessionManagerForApplication(appCtx.Session().ApplicationName)
	if sessionManager == nil {
		appCtx.Logger().Error("Authenticating user, cannot get session manager from request")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// Check the state value sent back from okta.mil with the value saved in the cookie
	returnedState := r.URL.Query().Get("state")
	stateCookieName := StateCookieName(appCtx.Session())
	stateCookie, err := r.Cookie(stateCookieName)
	if err != nil {
		appCtx.Logger().Error("Getting okta.mil state cookie",
			zap.String("stateCookieName", stateCookieName),
			zap.String("sessionUserId", appCtx.Session().UserID.String()),
			zap.Error(err))
		landingQuery := landingURL.Query()
		landingQuery.Add("error", "STATE_COOKIE_MISSING")
		landingURL.RawQuery = landingQuery.Encode()
		http.Redirect(w, r, landingURL.String(), http.StatusTemporaryRedirect)
		appCtx.Logger().Info("User redirected from okta.mil", zap.String("landingURL", landingURL.String()))
		return
	}

	hash := stateCookie.Value
	// case where user has 2 tabs open with different cookies
	if hash != shaAsString(returnedState) {
		appCtx.Logger().Error("State returned from okta.mil does not match state value stored in cookie",
			zap.String("state", returnedState),
			zap.String("cookie", hash),
			zap.String("hash", shaAsString(returnedState)))

		// Delete okta_state cookie
		auth.DeleteCookie(w, StateCookieName(appCtx.Session()))
		appCtx.Logger().Info("okta_state cookie deleted")

		// This operation will delete all cookies from the session
		err = sessionManager.Destroy(r.Context())
		if err != nil {
			appCtx.Logger().Error("Deleting okta.mil state cookie", zap.Error(err))
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

	provider, providerErr := okta.GetOktaProviderForRequest(r)
	if providerErr != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	// Exchange code received from login for access token. This is used during the grant_type auth flow
	exchange, err := exchangeCode(r.URL.Query().Get("code"), r, appCtx, *provider, h.HTTPClient)
	// Double error check
	if exchange.Error != "" {
		fmt.Println(exchange.Error)
		fmt.Println(exchange.ErrorDescription)
		return
	} else if err != nil {
		appCtx.Logger().Error("exchange code for access token", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// Verify access token
	jwtResult, verificationError := verifyToken(exchange.IDToken, returnedState, *provider)

	if verificationError != nil {
		appCtx.Logger().Error("token exchange verification", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	// Assign token values to session
	appCtx.Session().IDToken = exchange.IDToken
	appCtx.Session().AccessToken = exchange.AccessToken

	// Retrieve user info
	profileData, err := getProfileData(appCtx, *provider)
	if err != nil {
		appCtx.Logger().Error("get profile data", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// adding Okta profile data with intent to use for Okta profile editing from MilMove app
	appCtx.Session().IDToken = exchange.IDToken
	appCtx.Session().Email = profileData.Email

	// checking to see if user signed in with smart card
	// this will be leveraged in the customer app to ensure they authenticate with SC at least once
	var loggedInWithSmartCard bool
	if didUserSignInWithSmartCard(jwtResult.Claims, "sc") {
		loggedInWithSmartCard = true
	}

	oktaInfo := auth.OktaSessionInfo{
		Login:                 profileData.PreferredUsername,
		Email:                 profileData.Email,
		FirstName:             profileData.GivenName,
		LastName:              profileData.FamilyName,
		Edipi:                 profileData.Edipi,
		Sub:                   profileData.Sub,
		SignedInWithSmartCard: loggedInWithSmartCard,
	}
	appCtx.Session().OktaSessionInfo = oktaInfo

	appCtx.Logger().Info("New Login", zap.String("Okta user", profileData.PreferredUsername), zap.String("Okta email", profileData.Email), zap.String("Host", appCtx.Session().Hostname))

	result := authorizeUser(r.Context(), appCtx, profileData, sessionManager, h.sender)
	switch result {
	case authorizationResultError:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	case authorizationResultUnauthorized:
		invalidPermissionsResponse(appCtx, h.HandlerConfig, h.Context, w, r)
	case authorizationResultAuthorized:
		http.Redirect(w, r, landingURL.String(), http.StatusTemporaryRedirect)
	}
}

// didUserSignInWithSmartCard checks if the given value is present in the "amr" claim of the JWT claims interface
func didUserSignInWithSmartCard(claims map[string]interface{}, value string) bool {
	// isolate amr claim
	amr, ok := claims["amr"].([]interface{})
	if !ok {
		return false
	}

	// sift through to find the passed in value
	for _, v := range amr {
		if str, ok := v.(string); ok && str == value {
			return true
		}
	}

	return false
}

func authorizeUser(ctx context.Context, appCtx appcontext.AppContext, oktaUser models.OktaUser, sessionManager auth.SessionManager, notificationSender notifications.NotificationSender) AuthorizationResult {
	userIdentity, err := models.FetchUserIdentity(appCtx.DB(), oktaUser.Sub)

	if err == nil {
		// In this case, we found an existing user associated with the
		// unique okta.mil UUID (aka OID_User, aka openIDUser.UserID,
		// aka models.User.okta_id)
		appCtx.Logger().Info("Known user: found by okta.mil OID_User, checking authorization", zap.String("OID_User", oktaUser.Sub), zap.String("OID_Email", oktaUser.Email), zap.String("user.id", userIdentity.ID.String()), zap.String("user.okta_email", userIdentity.Email))

		result := AuthorizeKnownUser(ctx, appCtx, userIdentity, sessionManager)
		appCtx.Logger().Info("Known user authorization",
			zap.Any("authorizedResult", result),
			zap.String("OID_User", oktaUser.Sub),
			zap.String("OID_Email", oktaUser.Email))
		return result
	} else if err == models.ErrFetchNotFound { // Never heard of them
		// so far In this case, we can't find an existing user
		// associated with the unique okta.mil UUID (aka OID_User,
		// aka openIDUser.UserID, models.User.okta_id).
		// The authorizeUnknownUser method tries to find a user record
		// with a matching email address
		appCtx.Logger().Info("Unknown user: not found by okta.mil OID_User, associating email and checking authorization", zap.String("OID_User", oktaUser.Sub), zap.String("OID_Email", oktaUser.Email))
		result := authorizeUnknownUser(ctx, appCtx, oktaUser, sessionManager, notificationSender)
		appCtx.Logger().Info("Unknown user authorization",
			zap.Any("authorizedResult", result),
			zap.String("OID_User", oktaUser.Sub),
			zap.String("OID_Email", oktaUser.Email))
		return result
	}

	appCtx.Logger().Error("Error loading Identity.", zap.Error(err))
	return authorizationResultError
}

func AuthorizeKnownUser(ctx context.Context, appCtx appcontext.AppContext, userIdentity *models.UserIdentity, sessionManager auth.SessionManager) AuthorizationResult {
	if !userIdentity.Active {
		appCtx.Logger().Error("Inactive user requesting authentication",
			zap.String("application_name", string(appCtx.Session().ApplicationName)),
			zap.String("hostname", appCtx.Session().Hostname),
			zap.String("user_id", appCtx.Session().UserID.String()))
		return authorizationResultUnauthorized
	}
	appCtx.Session().Roles = append(appCtx.Session().Roles, userIdentity.Roles...)
	appCtx.Session().Permissions = getPermissionsForUser(appCtx, userIdentity.ID)

	appCtx.Session().UserID = userIdentity.ID
	if appCtx.Session().IsMilApp() && userIdentity.ServiceMemberID != nil {
		appCtx.Session().ServiceMemberID = *(userIdentity.ServiceMemberID)
	}

	// we want to check if the service member is signing in with CAC for the first time
	// if they are, their account is now validated with CAC and this check won't happen again
	if appCtx.Session().IsMilApp() &&
		appCtx.Session().OktaSessionInfo.SignedInWithSmartCard &&
		!*(userIdentity.ServiceMemberCacValidated) {
		sm, err := models.FetchServiceMember(appCtx.DB(), *userIdentity.ServiceMemberID)
		if err != nil {
			appCtx.Logger().Error("Error fetching service member to update", zap.Error(err))
		}
		sm.CacValidated = true
		smVerrs, err := models.SaveServiceMember(appCtx, &sm)
		if err != nil {
			appCtx.Logger().Error("Error updating service member's cac_verified value", zap.Error(err))
		}
		if smVerrs.HasAny() {
			appCtx.Logger().Error("Error updating service member's cac_verified value", zap.Error(smVerrs))
		}
	}

	if appCtx.Session().IsOfficeApp() {
		if userIdentity.OfficeActive != nil && !*userIdentity.OfficeActive {
			appCtx.Logger().Error("Inactive office user requesting authorization",
				zap.String("application_name", string(appCtx.Session().ApplicationName)),
				zap.String("hostname", appCtx.Session().Hostname),
				zap.String("user_id", appCtx.Session().UserID.String()))
			return authorizationResultUnauthorized
		}
		if userIdentity.OfficeUserID != nil {
			appCtx.Session().OfficeUserID = *(userIdentity.OfficeUserID)
		} else {
			// In case they managed to login before the office_user record was created
			officeUser, err := models.FetchOfficeUserByEmail(appCtx.DB(), appCtx.Session().Email)
			if err == models.ErrFetchNotFound {
				appCtx.Logger().Error("Non-office user authenticated at office site",
					zap.String("application_name", string(appCtx.Session().ApplicationName)),
					zap.String("hostname", appCtx.Session().Hostname),
					zap.String("user_id", appCtx.Session().UserID.String()))
				return authorizationResultUnauthorized
			} else if err != nil {
				appCtx.Logger().Error("Checking for office user during authorization",
					zap.String("user_id", appCtx.Session().UserID.String()),
					zap.Error(err))
				return authorizationResultError
			}
			appCtx.Session().OfficeUserID = officeUser.ID
			officeUser.UserID = &userIdentity.ID
			err = appCtx.DB().Save(officeUser)
			if err != nil {
				appCtx.Logger().Error("Updating office user during authorization",
					zap.String("user_id", appCtx.Session().UserID.String()),
					zap.Error(err))
				return authorizationResultError
			}
		}
	}

	if appCtx.Session().IsAdminApp() {
		if userIdentity.AdminUserActive != nil && !*userIdentity.AdminUserActive {
			appCtx.Logger().Error("Inactive admin user requesting authorization",
				zap.String("application_name", string(appCtx.Session().ApplicationName)),
				zap.String("hostname", appCtx.Session().Hostname),
				zap.String("user_id", appCtx.Session().UserID.String()))
			return authorizationResultUnauthorized
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
				appCtx.Logger().Error("No admin user found during authorization",
					zap.String("user_id", appCtx.Session().UserID.String()))
				return authorizationResultUnauthorized
			} else if err != nil {
				appCtx.Logger().Error("Checking for admin user", zap.String("userID", appCtx.Session().UserID.String()), zap.Error(err))
				return authorizationResultError
			}

			appCtx.Session().AdminUserID = adminUser.ID
			appCtx.Session().AdminUserRole = adminUser.Role.String()
			adminUser.UserID = &userIdentity.ID
			verrs, err := appCtx.DB().ValidateAndSave(&adminUser)
			if err != nil {
				appCtx.Logger().Error("Updating admin user", zap.String("userID", appCtx.Session().UserID.String()), zap.Error(err))
				return authorizationResultError
			}

			if verrs != nil {
				appCtx.Logger().Error("Admin user validation errors", zap.String("userID", appCtx.Session().UserID.String()), zap.Error(verrs))
				return authorizationResultError
			}
		}
	}
	appCtx.Session().FirstName = userIdentity.FirstName()
	appCtx.Session().LastName = userIdentity.LastName()
	appCtx.Session().Middle = userIdentity.Middle()

	if sessionManager == nil {
		appCtx.Logger().Error("Authenticating user, cannot get session manager from request")
		return authorizationResultError
	}

	authError := authenticateUser(ctx, appCtx, sessionManager)
	if authError != nil {
		appCtx.Logger().Error("Authenticating user", zap.Error(authError))
		return authorizationResultError
	}

	return authorizationResultAuthorized
}

func authorizeUnknownUser(ctx context.Context, appCtx appcontext.AppContext, oktaUser models.OktaUser, sessionManager auth.SessionManager, notificationSender notifications.NotificationSender) AuthorizationResult {
	var officeUser *models.OfficeUser
	var user *models.User
	var err error

	// Loads the User and Roles associations of the office or admin user
	conn := appCtx.DB().Eager("User", "User.Roles")

	if appCtx.Session().IsOfficeApp() {
		// Look to see if we have OfficeUser with this email address
		officeUser, err = models.FetchOfficeUserByEmail(conn, appCtx.Session().Email)
		if err == models.ErrFetchNotFound {
			appCtx.Logger().Error("Unauthorized: No Office user found",
				zap.String("OID_User", oktaUser.Sub),
				zap.String("OID_Email", oktaUser.Email))
			return authorizationResultUnauthorized
		} else if err != nil {
			appCtx.Logger().Error("Authorization checking for office user",
				zap.String("OID_User", oktaUser.Sub),
				zap.String("OID_Email", oktaUser.Email),
				zap.Error(err))
			return authorizationResultError
		}
		if !officeUser.Active {
			appCtx.Logger().Error("Unauthorized: Office user deactivated",
				zap.String("OID_User", oktaUser.Sub),
				zap.String("OID_Email", oktaUser.Email))
			return authorizationResultUnauthorized
		}
		user = &officeUser.User
	}

	var adminUser models.AdminUser
	if appCtx.Session().IsAdminApp() {
		// Look to see if we have AdminUser with this email address
		queryBuilder := query.NewQueryBuilder()
		filters := []services.QueryFilter{
			query.NewQueryFilter("email", "=", appCtx.Session().Email),
		}
		err = queryBuilder.FetchOne(appCtx, &adminUser, filters)

		// Log error and return if no AdminUser found with this email
		if err != nil && errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			appCtx.Logger().Error("Unauthorized: No admin user found",
				zap.String("OID_User", oktaUser.Sub),
				zap.String("OID_Email", oktaUser.Email))
			return authorizationResultUnauthorized
		} else if err != nil {
			appCtx.Logger().Error("Authorization checking for admin user",
				zap.String("OID_User", oktaUser.Sub),
				zap.String("OID_Email", oktaUser.Email),
				zap.Error(err))
			return authorizationResultError
		}
		// Log error and return if adminUser was found but deactivated
		if !adminUser.Active {
			appCtx.Logger().Error("Unauthorized: Admin user deactivated",
				zap.String("OID_User", oktaUser.Sub),
				zap.String("OID_Email", oktaUser.Email))
			return authorizationResultUnauthorized
		}
		user = &adminUser.User
	}

	if appCtx.Session().IsMilApp() {
		user, err = models.CreateUser(appCtx.DB(), oktaUser.Sub, oktaUser.Email)
		if err == nil {
			sysAdminEmail := notifications.GetSysAdminEmail(notificationSender)
			appCtx.Logger().Info(
				"New user account created through Okta.mil",
				zap.String("newUserID", user.ID.String()),
			)
			email, emailErr := notifications.NewUserAccountCreated(appCtx, sysAdminEmail, user.ID, user.UpdatedAt)
			if emailErr == nil {
				sendErr := notificationSender.SendNotification(appCtx, email)
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

		// setting cac_verified to false due to initial registration not allowing for smart card authentication
		// this will let the user into the application, but show them an error page telling them to sign in with CAC
		newServiceMember := models.ServiceMember{
			UserID:       user.ID,
			CacValidated: false,
		}

		smVerrs, smErr := models.SaveServiceMember(appCtx, &newServiceMember)
		if smVerrs.HasAny() || smErr != nil {
			appCtx.Logger().Error("Error creating service member for user", zap.Error(smErr))
			return authorizationResultError
		}
		appCtx.Session().ServiceMemberID = newServiceMember.ID
	} else {
		// If in Office App or Admin App with valid user - update user's OktaID
		appCtx.Logger().Error("Authorization associating UUID with user",
			zap.String("OID_User", oktaUser.Sub),
			zap.String("OID_Email", oktaUser.Email),
			zap.String("user.id", user.ID.String()),
		)
		err = models.UpdateUserOktaID(appCtx.DB(), user, oktaUser.Sub)

	}

	if err != nil {
		appCtx.Logger().Error("Authorization error updating/creating user", zap.Error(err))
		return authorizationResultError
	}

	appCtx.Session().UserID = user.ID
	if appCtx.Session().IsOfficeApp() && officeUser != nil {
		appCtx.Session().OfficeUserID = officeUser.ID
	} else if appCtx.Session().IsAdminApp() && adminUser.ID != uuid.Nil {
		appCtx.Session().AdminUserID = adminUser.ID
	}

	appCtx.Session().Roles = append(appCtx.Session().Roles, user.Roles...)
	appCtx.Session().Permissions = getPermissionsForUser(appCtx, user.ID)

	if sessionManager == nil {
		appCtx.Logger().Error("Authenticating user, cannot get session manager from request")
		return authorizationResultError
	}

	authError := authenticateUser(ctx, appCtx, sessionManager)
	if authError != nil {
		appCtx.Logger().Error("Authenticate user", zap.Error(authError))
		return authorizationResultError
	}

	return authorizationResultAuthorized
}

// InitAuth initializes the Okta provider
func InitAuth(v *viper.Viper, logger *zap.Logger, _ auth.ApplicationServername) (*okta.Provider, error) {

	// Create a new Okta Provider. This will be used in the creation of the additional providers for each subdomain
	oktaProvider := okta.NewOktaProvider(logger)
	err := oktaProvider.RegisterProviders(v)
	if err != nil {
		logger.Error("Initializing auth", zap.Error(err))
		return nil, err
	}

	return oktaProvider, nil
}
