package authentication

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

// UserAuthMiddleware enforces that the incoming request is tied to a user session
func UserAuthMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
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
			next.ServeHTTP(w, r)
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
	logger           *zap.Logger
	loginGovProvider LoginGovProvider
	callbackTemplate string
}

// NewAuthContext creates an Context
func NewAuthContext(logger *zap.Logger, loginGovProvider LoginGovProvider, callbackProtocol string, callbackPort string) Context {
	context := Context{
		logger:           logger,
		loginGovProvider: loginGovProvider,
		callbackTemplate: fmt.Sprintf("%s%%s:%s/", callbackProtocol, callbackPort),
	}
	return context
}

// LogoutHandler handles logging the user out of login.gov
type LogoutHandler struct {
	Context
	clientAuthSecretKey string
	noSessionTimeout    bool
}

// NewLogoutHandler creates a new LogoutHandler
func NewLogoutHandler(ac Context, clientAuthSecretKey string, noSessionTimeout bool) LogoutHandler {
	handler := LogoutHandler{
		Context:             ac,
		clientAuthSecretKey: clientAuthSecretKey,
		noSessionTimeout:    noSessionTimeout,
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
				logoutURL = "/"
			} else {
				logoutURL = h.loginGovProvider.LogoutURL(redirectURL, session.IDToken)
			}
			session.IDToken = ""
			session.UserID = uuid.Nil
			auth.WriteSessionCookie(w, session, h.clientAuthSecretKey, h.noSessionTimeout, h.logger)
			http.Redirect(w, r, logoutURL, http.StatusTemporaryRedirect)
		} else {
			// Can't log out of login.gov without a token, redirect and let them re-auth
			http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
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
	db                     *pop.Connection
	clientAuthSecretKey    string
	noSessionTimeout       bool
	loginGovMyClientID     string
	loginGovOfficeClientID string
}

// NewCallbackHandler creates a new CallbackHandler
func NewCallbackHandler(ac Context, db *pop.Connection, clientAuthSecretKey string, noSessionTimeout bool) CallbackHandler {
	handler := CallbackHandler{
		Context:             ac,
		db:                  db,
		clientAuthSecretKey: clientAuthSecretKey,
		noSessionTimeout:    noSessionTimeout,
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
	lURL := h.landingURL(session)

	authError := r.URL.Query().Get("error")
	// The user has either cancelled or declined to authorize the client
	if authError == "access_denied" {
		http.Redirect(w, r, lURL, http.StatusTemporaryRedirect)
		return
	}
	if authError == "invalid_request" {
		h.logger.Error("INVALID_REQUEST error from login.gov")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	provider, err := getLoginGovProviderForRequest(r)
	if err != nil {
		h.logger.Error("Get Goth provider", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// TODO: validate the state is the same (pull from session)
	code := r.URL.Query().Get("code")
	openIDSession, err := fetchToken(h.logger, code, provider.ClientKey, h.loginGovProvider)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
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

		session.UserID = userIdentity.ID
		if userIdentity.ServiceMemberID != nil {
			session.ServiceMemberID = *(userIdentity.ServiceMemberID)
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
			officeUser.UserID = &userIdentity.ID
			err = h.db.Save(officeUser)
			if err != nil {
				h.logger.Error("Updating office user", zap.String("email", session.Email), zap.Error(err))
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
		}
		session.FirstName = userIdentity.FirstName()
		session.LastName = userIdentity.LastName()
		session.Middle = userIdentity.Middle()

	} else if err == models.ErrFetchNotFound { // Never heard of them so far

		var officeUser *models.OfficeUser
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

		user, err := models.CreateUser(h.db, openIDUser.UserID, openIDUser.Email)
		if err == nil { // Successfully created the user
			session.UserID = user.ID
			if officeUser != nil {
				session.OfficeUserID = officeUser.ID
				officeUser.UserID = &user.ID
				err = h.db.Save(officeUser)
			}
		}
		if err != nil {
			h.logger.Error("Error creating user", zap.Error(err))
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
	} else {
		h.logger.Error("Error loading Identity.", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	h.logger.Info("logged in", zap.Any("session", session))
	auth.WriteSessionCookie(w, session, h.clientAuthSecretKey, h.noSessionTimeout, h.logger)
	http.Redirect(w, r, lURL, http.StatusTemporaryRedirect)
}

func fetchToken(logger *zap.Logger, code string, clientID string, loginGovProvider LoginGovProvider) (*openidConnect.Session, error) {
	tokenURL := loginGovProvider.TokenURL()
	expiry := auth.GetExpiryTimeFromMinutes(auth.SessionExpiryInMinutes)
	params, err := loginGovProvider.TokenParams(code, clientID, expiry)
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
