package authentication

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

// UserListHandler handles redirection
type UserListHandler struct {
	db *pop.Connection
	Context
}

// NewUserListHandler returns a new UserListHandler
func NewUserListHandler(ac Context, db *pop.Connection) UserListHandler {
	handler := UserListHandler{
		Context: ac,
		db:      db,
	}
	return handler
}

// UserListHandler lists users in the local database for local login
func (h UserListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session := auth.SessionFromRequestContext(r)
	if session != nil && session.UserID != uuid.Nil {
		// User is already authenticated, redirect to landing page
		http.Redirect(w, r, h.landingURL(session), http.StatusTemporaryRedirect)
		return
	}
	identities, err := models.FetchAllUserIdentities(h.db)
	if err != nil {
		h.logger.Error("Could not load list of users", zap.Error(err))
		http.Error(w,
			fmt.Sprintf("%s - Could not load list of users, try migrating the DB", http.StatusText(500)),
			http.StatusInternalServerError)
		return
	}

	// Grab the CSRF token from cookies set by the middleware
	csrfCookie, err := auth.GetCookie(auth.MaskedGorillaCSRFToken, r)
	if err != nil {
		h.logger.Error("CSRF Cookie was not set via middleware")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	csrfToken := csrfCookie.Value

	t := template.Must(template.New("users").Parse(`
		<h1>Select an existing user</h1>
		{{range .}}
			<form method="post" action="/devlocal-auth/login">
				<p id="{{.ID}}">
					<input type="hidden" name="gorilla.csrf.Token" value="` + csrfToken + `">
					{{.Email}}
					({{if .DpsUserID}}dps{{else if .TspUserID}}tsp{{else if .OfficeUserID}}office{{else}}milmove{{end}})
					<input type="hidden" name="id" value="{{.ID}}" />
					<button type="submit" value="{{.ID}}" data-hook="existing-user-login">Login</button>
				</p>
			</form>
		{{else}}
			<p><em>No users in the system!</em></p>
		{{end}}

		<h1>Create a new user</h1>
		<form method="post" action="/devlocal-auth/new">
			<p>
				<input type="hidden" name="gorilla.csrf.Token" value="` + csrfToken + `">
				<button type="submit" data-hook="new-user-login">Login as New User</button>
			</p>
		</form>
	`))
	err = t.Execute(w, identities)
	if err != nil {
		h.logger.Error("Could not render template", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
}

type devlocalAuthHandler struct {
	Context
	db                  *pop.Connection
	clientAuthSecretKey string
	noSessionTimeout    bool
	useSecureCookie     bool
}

// AssignUserHandler logs a user in directly
type AssignUserHandler devlocalAuthHandler

// NewAssignUserHandler creates a new AssignUserHandler
func NewAssignUserHandler(ac Context, db *pop.Connection, clientAuthSecretKey string, noSessionTimeout bool, useSecureCookie bool) AssignUserHandler {
	handler := AssignUserHandler{
		Context:             ac,
		db:                  db,
		clientAuthSecretKey: clientAuthSecretKey,
		noSessionTimeout:    noSessionTimeout,
		useSecureCookie:     useSecureCookie,
	}
	return handler
}

// AssignUserHandler logs in a user locally
func (h AssignUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID := r.PostFormValue("id")
	if userID == "" {
		h.logger.Error("No user id specified")
		http.Redirect(w, r, "/devlocal-auth/login", http.StatusTemporaryRedirect)
		return
	}

	h.logger.Info("New Devlocal Login", zap.String("userID", userID))

	userUUID := uuid.Must(uuid.FromString(userID))
	user, err := models.GetUser(h.db, userUUID)
	if err != nil {
		h.logger.Error("Could not load user", zap.String("userID", userID), zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}

	session := loginUser(devlocalAuthHandler(h), user, w, r)
	if session == nil {
		return
	}
	http.Redirect(w, r, h.landingURL(session), http.StatusSeeOther)
}

// CreateUserHandler creates a new user
type CreateUserHandler devlocalAuthHandler

// NewCreateUserHandler creates a new CreateUserHandler
func NewCreateUserHandler(ac Context, db *pop.Connection, clientAuthSecretKey string, noSessionTimeout bool, useSecureCookie bool) CreateUserHandler {
	handler := CreateUserHandler{
		Context:             ac,
		db:                  db,
		clientAuthSecretKey: clientAuthSecretKey,
		noSessionTimeout:    noSessionTimeout,
		useSecureCookie:     useSecureCookie,
	}
	return handler
}

// CreateUserHandler creates a user, primarily used in automated testing
func (h CreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := createUser(devlocalAuthHandler(h), w, r)
	if user == nil {
		return
	}
	session := loginUser(devlocalAuthHandler(h), user, w, r)
	if session == nil {
		return
	}
	jsonOut, _ := json.Marshal(user)
	fmt.Fprintf(w, string(jsonOut))
}

// CreateAndLoginUserHandler creates and then logs in a new user
type CreateAndLoginUserHandler devlocalAuthHandler

// NewCreateAndLoginUserHandler creates a new CreateAndLoginUserHandler
func NewCreateAndLoginUserHandler(ac Context, db *pop.Connection, clientAuthSecretKey string, noSessionTimeout bool, useSecureCookie bool) CreateAndLoginUserHandler {
	handler := CreateAndLoginUserHandler{
		Context:             ac,
		db:                  db,
		clientAuthSecretKey: clientAuthSecretKey,
		noSessionTimeout:    noSessionTimeout,
		useSecureCookie:     useSecureCookie,
	}
	return handler
}

// CreateAndLoginUserHandler creates a user and logs them in
func (h CreateAndLoginUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := createUser(devlocalAuthHandler(h), w, r)
	if user == nil {
		return
	}
	session := loginUser(devlocalAuthHandler(h), user, w, r)
	if session == nil {
		return
	}
	http.Redirect(w, r, h.landingURL(session), http.StatusSeeOther)
}

// createUser creates a user
func createUser(h devlocalAuthHandler, w http.ResponseWriter, r *http.Request) *models.User {
	id := uuid.Must(uuid.NewV4())

	now := time.Now()
	email := fmt.Sprintf("%s@example.com", now.Format("20060102150405"))

	user := models.User{
		LoginGovUUID:  id,
		LoginGovEmail: email,
	}

	verrs, err := h.db.ValidateAndCreate(&user)
	if err != nil {
		h.logger.Error("could not create user", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return nil
	}
	if verrs.Count() != 0 {
		h.logger.Error("validation errors creating user", zap.Stringer("errors", verrs))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return nil
	}
	return &user
}

// createSession creates a new session for the user
func createSession(h devlocalAuthHandler, user *models.User, w http.ResponseWriter, r *http.Request) (*auth.Session, error) {
	session := auth.SessionFromRequestContext(r)
	if session == nil {
		return nil, errors.New("Unable to create session from request context")
	}

	lgUUID := user.LoginGovUUID.String()
	userIdentity, err := models.FetchUserIdentity(h.db, lgUUID)

	if err != nil {
		return nil, errors.Wrapf(err, "Unable to fetch user identity from LoginGovUUID %s", lgUUID)
	}

	// Assign user identity to session
	session.IDToken = "devlocal"
	session.UserID = userIdentity.ID
	session.Email = userIdentity.Email
	session.Disabled = userIdentity.Disabled

	if userIdentity.ServiceMemberID != nil {
		session.ServiceMemberID = *(userIdentity.ServiceMemberID)
	}

	if userIdentity.OfficeUserID != nil {
		session.OfficeUserID = *(userIdentity.OfficeUserID)
	}

	if userIdentity.TspUserID != nil {
		session.TspUserID = *(userIdentity.TspUserID)
	}

	if userIdentity.DpsUserID != nil {
		session.DpsUserID = *(userIdentity.DpsUserID)
	}

	session.FirstName = userIdentity.FirstName()
	session.LastName = userIdentity.LastName()
	session.Middle = userIdentity.Middle()

	// Writing out the session cookie logs in the user
	h.logger.Info("logged in", zap.Any("session", session))
	auth.WriteSessionCookie(w, session, h.clientAuthSecretKey, h.noSessionTimeout, h.logger, h.useSecureCookie)

	return session, nil
}

// verifySessionWithApp returns an error if the user id for a specific app is not available
func verifySessionWithApp(session *auth.Session) error {

	// TODO: Should this be a check that we do? Or will all office and tsp users also be service members?
	// if (session.ServiceMemberID == uuid.UUID{}) && session.IsMilApp() {
	// 	return errors.Errorf("Non-service member user %s authenticated at service member site", session.Email)
	// }

	if (session.OfficeUserID == uuid.UUID{}) && session.IsOfficeApp() {
		return errors.Errorf("Non-office user %s authenticated at office site", session.Email)
	}

	if (session.TspUserID == uuid.UUID{}) && session.IsTspApp() {
		return errors.Errorf("Non-TSP user %s authenticated at TSP site", session.Email)
	}

	return nil
}

// loginUser creates a session for the user and verifies the session against the app
func loginUser(h devlocalAuthHandler, user *models.User, w http.ResponseWriter, r *http.Request) *auth.Session {
	session, err := createSession(devlocalAuthHandler(h), user, w, r)
	if err != nil {
		h.logger.Error("Could not create session", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return nil
	}

	if session.Disabled {
		h.logger.Info("Disabled user requesting authentication", zap.Error(err), zap.String("email", session.Email))
		http.Error(w, http.StatusText(403), http.StatusForbidden)
		return nil
	}

	err = verifySessionWithApp(session)
	if err != nil {
		h.logger.Error("User unauthorized", zap.Error(err))
		http.Error(w, http.StatusText(401), http.StatusUnauthorized)
		return nil
	}
	return session
}
