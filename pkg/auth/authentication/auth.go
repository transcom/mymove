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

	"github.com/alexedwards/scs/v2"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

// IsLoggedInMiddleware handles requests to is_logged_in endpoint by returning true if someone is logged in
func IsLoggedInMiddleware(logger Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"isLoggedIn": false,
		}

		session := auth.SessionFromRequestContext(r)
		if session != nil && session.UserID != uuid.Nil {
			data["isLoggedIn"] = true
		}

		newEncoderErr := json.NewEncoder(w).Encode(data)
		if newEncoderErr != nil {
			logger.Error("Failed encoding is_logged_in check response", zap.Error(newEncoderErr))
		}
	}
}

// UserAuthMiddleware enforces that the incoming request is tied to a user session
func UserAuthMiddleware(logger Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {

			session := auth.SessionFromRequestContext(r)

			// We must have a logged in session and a user
			if session == nil || session.UserID == uuid.Nil {
				logger.Info("unauthorized access, no session token or user id")
				http.Error(w, http.StatusText(401), http.StatusUnauthorized)
				return
			}

			// DO NOT CHECK MILMOVE SESSION BECAUSE NEW SERVICE MEMBERS WON'T HAVE AN ID RIGHT AWAY
			// This must be the right type of user for the application
			if session.IsOfficeApp() && !session.IsOfficeUser() {
				logger.Error("unauthorized user for office.move.mil", zap.String("email", session.Email))
				http.Error(w, http.StatusText(401), http.StatusUnauthorized)
				return
			} else if session.IsAdminApp() && !session.IsAdminUser() {
				logger.Error("unauthorized user for admin.move.mil", zap.String("email", session.Email))
				http.Error(w, http.StatusText(401), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(mw)
	}
}

func updateUserCurrentSessionID(session *auth.Session, sessionID string, db *pop.Connection, logger Logger) error {
	userID := session.UserID

	user, err := models.GetUser(db, userID)
	if err != nil {
		logger.Error("Fetching user", zap.String("user_id", userID.String()), zap.Error(err))
	}

	if session.IsAdminUser() {
		user.CurrentAdminSessionID = sessionID
	} else if session.IsOfficeUser() {
		user.CurrentOfficeSessionID = sessionID
	} else if session.IsServiceMember() {
		user.CurrentMilSessionID = sessionID
	}

	err = db.Save(user)
	if err != nil {
		logger.Error("Updating user's current_x_session_id", zap.String("email", session.Email), zap.Error(err))
		return err
	}

	return err
}

func resetUserCurrentSessionID(session *auth.Session, db *pop.Connection, logger Logger) error {
	userID := session.UserID
	user, err := models.GetUser(db, userID)
	if err != nil {
		logger.Error("Fetching user", zap.String("user_id", userID.String()), zap.Error(err))
	}

	if session.IsAdminUser() {
		user.CurrentAdminSessionID = ""
	} else if session.IsOfficeUser() {
		user.CurrentOfficeSessionID = ""
	} else if session.IsServiceMember() {
		user.CurrentMilSessionID = ""
	}
	err = db.Save(user)
	if err != nil {
		logger.Error("Updating user's current_x_session_id", zap.String("email", session.Email), zap.Error(err))
		return err
	}

	return err
}

func currentUser(session *auth.Session, db *pop.Connection) (*models.User, error) {
	userID := session.UserID
	user, err := models.GetUser(db, userID)
	if err != nil {
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

func authenticateUser(ctx context.Context, sessionManager *scs.SessionManager, session *auth.Session, logger Logger, db *pop.Connection) error {
	// The session token must be renewed during sign in to prevent
	// session fixation attacks
	err := sessionManager.RenewToken(ctx)
	if err != nil {
		logger.Error("Error renewing session token", zap.Error(err))
		return err
	}
	sessionID, _, err := sessionManager.Commit(ctx)
	if err != nil {
		logger.Error("Failed to write new user session to store", zap.Error(err))
		return err
	}
	sessionManager.Put(ctx, "session", session)

	user, err := currentUser(session, db)
	if err != nil {
		logger.Error("Fetching user", zap.String("user_id", session.UserID.String()), zap.Error(err))
		return err
	}
	// Check to see if sessionID is set on the user, presently
	existingSessionID := currentSessionID(session, user)
	if existingSessionID != "" {

		// Lookup the old session that wasn't logged out
		_, exists, err := sessionManager.Store.Find(existingSessionID)
		if err != nil {
			logger.Error("Error loading previous session", zap.Error(err))
			return err
		}

		if !exists {
			logger.Info("Session expired")
		} else {
			logger.Info("Concurrent session detected. Will delete previous session.")

			// We need to delete the concurrent session.
			err := sessionManager.Store.Delete(existingSessionID)
			if err != nil {
				logger.Error("Error deleting previous session", zap.Error(err))
				return err
			}
		}
	}

	updateErr := updateUserCurrentSessionID(session, sessionID, db, logger)
	if updateErr != nil {
		logger.Error("Updating user's current session ID", zap.Error(updateErr))
		return updateErr
	}
	logger.Info("Logged in", zap.Any("session", session))

	return nil
}

// AdminAuthMiddleware is middleware for admin authentication
func AdminAuthMiddleware(logger Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			session := auth.SessionFromRequestContext(r)

			if session == nil || !session.IsAdminUser() {
				http.Error(w, http.StatusText(403), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(mw)
	}
}

// PrimeAuthorizationMiddleware is the prime authorization middleware
func PrimeAuthorizationMiddleware(logger Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
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
	logger           Logger
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
func NewAuthContext(logger Logger, loginGovProvider LoginGovProvider, callbackProtocol string, callbackPort int, sessionManagers [3]*scs.SessionManager) Context {
	context := Context{
		logger:           logger,
		loginGovProvider: loginGovProvider,
		callbackTemplate: fmt.Sprintf("%s://%%s:%d/", callbackProtocol, callbackPort),
		sessionManagers:  sessionManagers,
	}
	return context
}

// LogoutHandler handles logging the user out of login.gov
type LogoutHandler struct {
	Context
	db *pop.Connection
}

// NewLogoutHandler creates a new LogoutHandler
func NewLogoutHandler(ac Context, db *pop.Connection) LogoutHandler {
	logoutHandler := LogoutHandler{
		Context: ac,
		db:      db,
	}
	return logoutHandler
}

func (h LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session := auth.SessionFromRequestContext(r)

	if session != nil {
		redirectURL := h.landingURL(session)
		if session.IDToken != "" {
			var logoutURL string
			// All users logged in via devlocal-auth will have this IDToken. We
			// don't want to make a call to login.gov for a logout URL as it will
			// fail for devlocal-auth'ed users.
			if session.IDToken == "devlocal" {
				logoutURL = redirectURL
			} else {
				logoutURL = h.loginGovProvider.LogoutURL(redirectURL, session.IDToken)
			}
			err := resetUserCurrentSessionID(session, h.db, h.logger)
			if err != nil {
				h.logger.Error("failed to reset user's current_x_session_id")
			}
			err = h.sessionManager(session).Destroy(r.Context())
			if err != nil {
				h.logger.Error("failed to destroy session")
			}
			auth.DeleteCSRFCookies(w)
			h.logger.Info("user logged out")
			fmt.Fprint(w, logoutURL)
		} else {
			// Can't log out of login.gov without a token, redirect and let them re-auth
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
	UseSecureCookie bool
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
	session := auth.SessionFromRequestContext(r)
	if session != nil && session.UserID != uuid.Nil {
		// User is already authenticated, redirect to landing page
		http.Redirect(w, r, h.landingURL(session), http.StatusTemporaryRedirect)
		return
	}

	loginData, err := h.loginGovProvider.AuthorizationURL(r)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// Hash the state/Nonce value sent to login.gov and set the result as an HttpOnly cookie
	// Check this when we return from login.gov
	if session == nil {
		h.logger.Error("Session is nil, so cannot get hostname for state Cookie")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	stateCookie := http.Cookie{
		Name:     StateCookieName(session),
		Value:    shaAsString(loginData.Nonce),
		Path:     "/",
		Expires:  time.Now().Add(time.Duration(loginStateCookieTTLInSecs) * time.Second),
		MaxAge:   loginStateCookieTTLInSecs,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.UseSecureCookie,
	}
	http.SetCookie(w, &stateCookie)
	http.Redirect(w, r, loginData.RedirectURL, http.StatusTemporaryRedirect)
}

// CallbackHandler processes a callback from login.gov
type CallbackHandler struct {
	Context
	db *pop.Connection
}

// NewCallbackHandler creates a new CallbackHandler
func NewCallbackHandler(ac Context, db *pop.Connection) CallbackHandler {
	handler := CallbackHandler{
		Context: ac,
		db:      db,
	}
	return handler
}

// AuthorizationCallbackHandler handles the callback from the Login.gov authorization flow
func (h CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session := auth.SessionFromRequestContext(r)

	if session == nil {
		h.logger.Error("Session missing")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	rawLandingURL := h.landingURL(session)

	landingURL, err := url.Parse(rawLandingURL)
	if err != nil {
		h.logger.Error("Error parsing landing URL")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	if err := r.URL.Query().Get("error"); len(err) > 0 {
		landingQuery := landingURL.Query()
		switch err {
		case "access_denied":
			// The user has either cancelled or declined to authorize the client
		case "invalid_request":
			h.logger.Error("INVALID_REQUEST error from login.gov")
			landingQuery.Add("error", "INVALID_REQUEST")
		default:
			h.logger.Error("unknown error from login.gov")
			landingQuery.Add("error", "UNKNOWN_ERROR")
		}
		landingURL.RawQuery = landingQuery.Encode()
		http.Redirect(w, r, landingURL.String(), http.StatusPermanentRedirect)
		return
	}

	// Check the state value sent back from login.gov with the value saved in the cookie
	returnedState := r.URL.Query().Get("state")
	stateCookie, err := r.Cookie(StateCookieName(session))
	if err != nil {
		h.logger.Error("Getting login.gov state cookie", zap.Error(err))
		http.Error(w, http.StatusText(403), http.StatusForbidden)
		return
	}

	hash := stateCookie.Value
	// case where user has 2 tabs open with different cookies
	if hash != shaAsString(returnedState) {
		h.logger.Error("State returned from Login.gov does not match state value stored in cookie",
			zap.String("state", returnedState),
			zap.String("cookie", hash),
			zap.String("hash", shaAsString(returnedState)))

		// Delete lg_state cookie
		auth.DeleteCookie(w, StateCookieName(session))

		// This operation will delete all cookies from the session
		err := h.sessionManager(session).Destroy(r.Context())
		if err != nil {
			h.logger.Error("Deleting login.gov state cookie", zap.Error(err))
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		// set error query
		landingQuery := landingURL.Query()
		landingQuery.Add("error", "SIGNIN_ERROR")
		landingURL.RawQuery = landingQuery.Encode()
		http.Redirect(w, r, landingURL.String(), http.StatusTemporaryRedirect)
		return
	}

	provider, err := getLoginGovProviderForRequest(r)
	if err != nil {
		h.logger.Error("Get Goth provider", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// TODO: validate the state is the same (pull from session)
	openIDSession, err := fetchToken(
		h.logger,
		r.URL.Query().Get("code"),
		provider.ClientKey,
		h.loginGovProvider)
	if err != nil {
		http.Error(w, http.StatusText(401), http.StatusUnauthorized)
		return
	}

	openIDUser, err := provider.FetchUser(openIDSession)
	if err != nil {
		h.logger.Error("Login.gov user info request", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	session.IDToken = openIDSession.IDToken
	session.Email = openIDUser.Email

	h.logger.Info("New Login", zap.String("OID_User", openIDUser.UserID), zap.String("OID_Email", openIDUser.Email), zap.String("Host", session.Hostname))

	userIdentity, err := models.FetchUserIdentity(h.db, openIDUser.UserID)
	if err == nil { // Someone we know already
		authorizeKnownUser(userIdentity, h, session, w, r, landingURL.String())
		return
	} else if err == models.ErrFetchNotFound { // Never heard of them so far
		authorizeUnknownUser(openIDUser, h, session, w, r, landingURL.String())
		return
	} else {
		h.logger.Error("Error loading Identity.", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
}

var authorizeKnownUser = func(userIdentity *models.UserIdentity, h CallbackHandler, session *auth.Session, w http.ResponseWriter, r *http.Request, lURL string) {
	if !userIdentity.Active {
		h.logger.Error("Active user requesting authentication",
			zap.String("application_name", string(session.ApplicationName)),
			zap.String("hostname", session.Hostname),
			zap.String("user_id", session.UserID.String()),
			zap.String("email", session.Email))
		http.Error(w, http.StatusText(403), http.StatusForbidden)
		return
	}
	for _, role := range userIdentity.Roles {
		session.Roles = append(session.Roles, role)
	}
	session.UserID = userIdentity.ID
	if userIdentity.ServiceMemberID != nil {
		session.ServiceMemberID = *(userIdentity.ServiceMemberID)
	}

	if userIdentity.DpsUserID != nil && (userIdentity.DpsActive != nil && *userIdentity.DpsActive) {
		session.DpsUserID = *(userIdentity.DpsUserID)
	}

	if session.IsOfficeApp() {
		if userIdentity.OfficeActive != nil && !*userIdentity.OfficeActive {
			h.logger.Error("Office user is deactivated", zap.String("email", session.Email))
			http.Error(w, http.StatusText(403), http.StatusForbidden)
			return
		}
		if userIdentity.OfficeUserID != nil {
			session.OfficeUserID = *(userIdentity.OfficeUserID)
		} else {
			// In case they managed to login before the office_user record was created
			officeUser, err := models.FetchOfficeUserByEmail(h.db, session.Email)
			if err == models.ErrFetchNotFound {
				h.logger.Error("Non-office user authenticated at office site", zap.String("email", session.Email))
				http.Error(w, http.StatusText(401), http.StatusUnauthorized)
				return
			} else if err != nil {
				h.logger.Error("Checking for office user", zap.String("email", session.Email), zap.Error(err))
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
			session.OfficeUserID = officeUser.ID
			officeUser.UserID = &userIdentity.ID
			err = h.db.Save(officeUser)
			if err != nil {
				h.logger.Error("Updating office user", zap.String("email", session.Email), zap.Error(err))
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
		}
	}

	if session.IsAdminApp() {
		if userIdentity.AdminUserActive != nil && !*userIdentity.AdminUserActive {
			h.logger.Error("Admin user is deactivated", zap.String("email", session.Email))
			http.Error(w, http.StatusText(403), http.StatusForbidden)
			return
		}
		if userIdentity.AdminUserID != nil {
			session.AdminUserID = *(userIdentity.AdminUserID)
			session.AdminUserRole = userIdentity.AdminUserRole.String()
		} else {
			// In case they managed to login before the admin_user record was created
			var adminUser models.AdminUser
			queryBuilder := query.NewQueryBuilder(h.db)
			filters := []services.QueryFilter{
				query.NewQueryFilter("email", "=", strings.ToLower(userIdentity.Email)),
			}
			err := queryBuilder.FetchOne(&adminUser, filters)
			if err == models.ErrFetchNotFound {
				h.logger.Error("Non-admin user authenticated at admin site", zap.String("email", session.Email))
				http.Error(w, http.StatusText(403), http.StatusForbidden)
				return
			} else if err != nil {
				h.logger.Error("Checking for admin user", zap.String("email", session.Email), zap.Error(err))
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}

			session.AdminUserID = adminUser.ID
			session.AdminUserRole = adminUser.Role.String()
			adminUser.UserID = &userIdentity.ID
			verrs, err := h.db.ValidateAndSave(&adminUser)
			if err != nil {
				h.logger.Error("Updating admin user", zap.String("email", session.Email), zap.Error(err))
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}

			if verrs != nil {
				h.logger.Error("Admin user validation errors", zap.String("email", session.Email), zap.Error(verrs))
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
		}
	}
	session.FirstName = userIdentity.FirstName()
	session.LastName = userIdentity.LastName()
	session.Middle = userIdentity.Middle()

	sessionManager := h.sessionManager(session)
	authError := authenticateUser(r.Context(), sessionManager, session, h.logger, h.db)
	if authError != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, lURL, http.StatusTemporaryRedirect)
}

var authorizeUnknownUser = func(openIDUser goth.User, h CallbackHandler, session *auth.Session, w http.ResponseWriter, r *http.Request, lURL string) {
	var officeUser *models.OfficeUser
	var user *models.User
	var err error

	// Loads the User and Roles associations of the office or admin user
	conn := h.db.Eager("User", "User.Roles")

	if session.IsOfficeApp() { // Look to see if we have OfficeUser with this email address
		officeUser, err = models.FetchOfficeUserByEmail(conn, session.Email)
		if err == models.ErrFetchNotFound {
			h.logger.Error("No Office user found", zap.String("email", session.Email))
			http.Error(w, http.StatusText(401), http.StatusUnauthorized)
			return
		} else if err != nil {
			h.logger.Error("Checking for office user", zap.String("email", session.Email), zap.Error(err))
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		if !officeUser.Active {
			h.logger.Error("Office user is deactivated", zap.String("email", session.Email))
			http.Error(w, http.StatusText(403), http.StatusForbidden)
			return
		}
		user = &officeUser.User
	}

	var adminUser models.AdminUser
	if session.IsAdminApp() {
		queryBuilder := query.NewQueryBuilder(conn)
		filters := []services.QueryFilter{
			query.NewQueryFilter("email", "=", session.Email),
		}
		err = queryBuilder.FetchOne(&adminUser, filters)

		if err != nil && errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			h.logger.Error("No admin user found", zap.String("email", session.Email))
			http.Error(w, http.StatusText(401), http.StatusUnauthorized)
			return
		} else if err != nil {
			h.logger.Error("Checking for admin user", zap.String("email", session.Email), zap.Error(err))
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		if !adminUser.Active {
			h.logger.Error("Admin user is deactivated", zap.String("email", session.Email))
			http.Error(w, http.StatusText(403), http.StatusForbidden)
			return
		}
		user = &adminUser.User
	}

	if session.IsMilApp() {
		user, err = models.CreateUser(h.db, openIDUser.UserID, openIDUser.Email)
	} else {
		err = models.UpdateUserLoginGovUUID(h.db, user, openIDUser.UserID)
	}

	if err != nil {
		h.logger.Error("Error updating/creating user", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	session.UserID = user.ID
	if session.IsOfficeApp() && officeUser != nil {
		session.OfficeUserID = officeUser.ID
	} else if session.IsAdminApp() && adminUser.ID != uuid.Nil {
		session.AdminUserID = adminUser.ID
	}

	for _, role := range user.Roles {
		session.Roles = append(session.Roles, role)
	}

	sessionManager := h.sessionManager(session)
	authError := authenticateUser(r.Context(), sessionManager, session, h.logger, h.db)
	if authError != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, h.landingURL(session), http.StatusTemporaryRedirect)
}

func fetchToken(logger Logger, code string, clientID string, loginGovProvider LoginGovProvider) (*openidConnect.Session, error) {
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
			logger.Error("Error in closing response", zap.Error(closeErr)))
		}
	}

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
func InitAuth(v *viper.Viper, logger Logger, appnames auth.ApplicationServername) (LoginGovProvider, error) {
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
