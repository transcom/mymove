package authentication

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	beeline "github.com/honeycombio/beeline-go"
	"github.com/honeycombio/beeline-go/trace"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

// UserAuthMiddleware enforces that the incoming request is tied to a user session
func UserAuthMiddleware(logger Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			ctx, span := beeline.StartSpan(r.Context(), "UserAuthMiddleware")
			defer span.Send()

			session := auth.SessionFromRequestContext(r)
			// We must have a logged in session and a user
			if session == nil || session.UserID == uuid.Nil {
				logger.Error("unauthorized access")
				http.Error(w, http.StatusText(401), http.StatusUnauthorized)
				return
			}
			// TODO: Add support for BackupContacts
			// And this must be the right type of user for the application
			if session.IsOfficeApp() && session.OfficeUserID == uuid.Nil {
				logger.Error("unauthorized user for office.move.mil", zap.String("email", session.Email))
				http.Error(w, http.StatusText(401), http.StatusUnauthorized)
				return
			}
			// This must be the right type of user for the application
			if session.IsTspApp() && session.TspUserID == uuid.Nil {
				logger.Error("unauthorized user for tsp.move.mil", zap.String("email", session.Email))
				http.Error(w, http.StatusText(401), http.StatusUnauthorized)
				return
			}

			// Include session office, service member, tsp and user IDs to the beeline event
			span.AddTraceField("auth.office_user_id", session.OfficeUserID)
			span.AddTraceField("auth.service_member_id", session.ServiceMemberID)
			span.AddTraceField("auth.tsp_user_id", session.TspUserID)
			span.AddTraceField("auth.user_id", session.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		return http.HandlerFunc(mw)
	}
}

func (context Context) landingURL(session *auth.Session) string {
	return fmt.Sprintf(context.callbackTemplate, session.Hostname)
}

// Context is the common handler type for auth handlers
type Context struct {
	logger           Logger
	loginGovProvider LoginGovProvider
	callbackTemplate string
}

// NewAuthContext creates an Context
func NewAuthContext(logger Logger, loginGovProvider LoginGovProvider, callbackProtocol string, callbackPort int) Context {
	context := Context{
		logger:           logger,
		loginGovProvider: loginGovProvider,
		callbackTemplate: fmt.Sprintf("%s://%%s:%d/", callbackProtocol, callbackPort),
	}
	return context
}

// LogoutHandler handles logging the user out of login.gov
type LogoutHandler struct {
	Context
	clientAuthSecretKey string
	noSessionTimeout    bool
	useSecureCookie     bool
}

// NewLogoutHandler creates a new LogoutHandler
func NewLogoutHandler(ac Context, clientAuthSecretKey string, noSessionTimeout bool, useSecureCookie bool) LogoutHandler {
	handler := LogoutHandler{
		Context:             ac,
		clientAuthSecretKey: clientAuthSecretKey,
		noSessionTimeout:    noSessionTimeout,
		useSecureCookie:     useSecureCookie,
	}
	return handler
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
			// This operation will delete all cookies from the session
			session.IDToken = ""
			session.UserID = uuid.Nil
			auth.WriteSessionCookie(w, session, h.clientAuthSecretKey, h.noSessionTimeout, h.logger, h.useSecureCookie)
			fmt.Fprint(w, logoutURL)
		} else {
			// Can't log out of login.gov without a token, redirect and let them re-auth
			fmt.Fprint(w, redirectURL)
		}
	}
}

// RedirectHandler handles redirection
type RedirectHandler struct {
	Context
}

// RedirectHandler constructs the Login.gov authentication URL and redirects to it
func (h RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session := auth.SessionFromRequestContext(r)
	if session != nil && session.UserID != uuid.Nil {
		// User is already authenticated, redirect to landing page
		http.Redirect(w, r, h.landingURL(session), http.StatusTemporaryRedirect)
		return
	}

	authURL, err := h.loginGovProvider.AuthorizationURL(r)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// CallbackHandler processes a callback from login.gov
type CallbackHandler struct {
	Context
	db                  *pop.Connection
	clientAuthSecretKey string
	noSessionTimeout    bool
	useSecureCookie     bool
}

// NewCallbackHandler creates a new CallbackHandler
func NewCallbackHandler(ac Context, db *pop.Connection, clientAuthSecretKey string, noSessionTimeout bool, useSecureCookie bool) CallbackHandler {
	handler := CallbackHandler{
		Context:             ac,
		db:                  db,
		clientAuthSecretKey: clientAuthSecretKey,
		noSessionTimeout:    noSessionTimeout,
		useSecureCookie:     useSecureCookie,
	}
	return handler
}

// AuthorizationCallbackHandler handles the callback from the Login.gov authorization flow
func (h CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	_, span := beeline.StartSpan(r.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

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
		authorizeKnownUser(userIdentity, h, session, w, span, r, landingURL.String())
		return
	} else if err == models.ErrFetchNotFound { // Never heard of them so far
		authorizeUnknownUser(openIDUser, h, session, w, span, r, landingURL.String())
		return
	} else {
		h.logger.Error("Error loading Identity.", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
}

func authorizeKnownUser(userIdentity *models.UserIdentity, h CallbackHandler, session *auth.Session, w http.ResponseWriter, span *trace.Span, r *http.Request, lURL string) {
	if userIdentity.Disabled {
		h.logger.Error("Disabled user requesting authentication", zap.String("email", session.Email))
		http.Error(w, http.StatusText(403), http.StatusForbidden)
		return
	}
	session.UserID = userIdentity.ID
	span.AddField("session.user_id", session.UserID)
	if userIdentity.ServiceMemberID != nil {
		session.ServiceMemberID = *(userIdentity.ServiceMemberID)
		span.AddField("session.service_member_id", session.ServiceMemberID)
	}

	if userIdentity.DpsUserID != nil {
		session.DpsUserID = *(userIdentity.DpsUserID)
	}

	if userIdentity.OfficeUserID != nil {
		session.OfficeUserID = *(userIdentity.OfficeUserID)
	} else if session.IsOfficeApp() {
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
		span.AddField("session.office_user_id", session.OfficeUserID)
		officeUser.UserID = &userIdentity.ID
		err = h.db.Save(officeUser)
		if err != nil {
			h.logger.Error("Updating office user", zap.String("email", session.Email), zap.Error(err))
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
	}

	if userIdentity.TspUserID != nil {
		session.TspUserID = *(userIdentity.TspUserID)
	} else if session.IsTspApp() {
		// In case they managed to login before the tsp_user record was created
		tspUser, err := models.FetchTspUserByEmail(h.db, session.Email)
		if err == models.ErrFetchNotFound {
			h.logger.Error("Non-TSP user authenticated at tsp site", zap.String("email", session.Email))
			http.Error(w, http.StatusText(401), http.StatusUnauthorized)
			return
		} else if err != nil {
			h.logger.Error("Checking for TSP user", zap.String("email", session.Email), zap.Error(err))
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		session.TspUserID = tspUser.ID
		span.AddField("session.tsp_user_id", session.TspUserID)
		tspUser.UserID = &userIdentity.ID
		err = h.db.Save(tspUser)
		if err != nil {
			h.logger.Error("Updating TSP user", zap.String("email", session.Email), zap.Error(err))
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
	}
	session.FirstName = userIdentity.FirstName()
	session.LastName = userIdentity.LastName()
	session.Middle = userIdentity.Middle()

	h.logger.Info("logged in", zap.Any("session", session))

	auth.WriteSessionCookie(w, session, h.clientAuthSecretKey, h.noSessionTimeout, h.logger, h.useSecureCookie)
	http.Redirect(w, r, lURL, http.StatusTemporaryRedirect)
}

func authorizeUnknownUser(openIDUser goth.User, h CallbackHandler, session *auth.Session, w http.ResponseWriter, span *trace.Span, r *http.Request, lURL string) {
	var officeUser *models.OfficeUser
	var err error
	if session.IsOfficeApp() { // Look to see if we have OfficeUser with this email address
		officeUser, err = models.FetchOfficeUserByEmail(h.db, session.Email)
		if err == models.ErrFetchNotFound {
			h.logger.Error("No Office user found", zap.String("email", session.Email))
			http.Error(w, http.StatusText(401), http.StatusUnauthorized)
			return
		} else if err != nil {
			h.logger.Error("Checking for office user", zap.String("email", session.Email), zap.Error(err))
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
	}

	var tspUser *models.TspUser
	if session.IsTspApp() { // Look to see if we have TspUser with this email address
		tspUser, err = models.FetchTspUserByEmail(h.db, session.Email)
		if err == models.ErrFetchNotFound {
			h.logger.Error("No TSP user found", zap.String("email", session.Email))
			http.Error(w, http.StatusText(401), http.StatusUnauthorized)
			return
		} else if err != nil {
			h.logger.Error("Checking for TSP user", zap.String("email", session.Email), zap.Error(err))
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
	}

	user, err := models.CreateUser(h.db, openIDUser.UserID, openIDUser.Email)
	if err == nil { // Successfully created the user
		session.UserID = user.ID
		span.AddField("session.user_id", session.UserID)
		if officeUser != nil {
			session.OfficeUserID = officeUser.ID
			span.AddField("session.office_user_id", session.OfficeUserID)
			officeUser.UserID = &user.ID
			err = h.db.Save(officeUser)
		} else if tspUser != nil {
			session.TspUserID = tspUser.ID
			span.AddField("session.tsp_user_id", session.TspUserID)
			tspUser.UserID = &user.ID
			err = h.db.Save(tspUser)
		}
	}
	if err != nil {
		h.logger.Error("Error creating user", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	h.logger.Info("logged in", zap.Any("session", session))

	auth.WriteSessionCookie(w, session, h.clientAuthSecretKey, h.noSessionTimeout, h.logger, h.useSecureCookie)
	http.Redirect(w, r, h.landingURL(session), http.StatusTemporaryRedirect)
}

func fetchToken(logger Logger, code string, clientID string, loginGovProvider LoginGovProvider) (*openidConnect.Session, error) {
	tokenURL := loginGovProvider.TokenURL()
	expiry := auth.GetExpiryTimeFromMinutes(auth.SessionExpiryInMinutes)
	params, err := loginGovProvider.TokenParams(code, clientID, expiry)
	if err != nil {
		logger.Error("Creating token endpoint params", zap.Error(err))
		return nil, err
	}

	/* #nosec G107 */
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
