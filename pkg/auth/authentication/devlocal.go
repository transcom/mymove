package authentication

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/gorilla/csrf"
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

	// get list of users in system
	var users []models.User
	err := h.db.All(&users)
	if err != nil {
		h.logger.Error("Could not load list of users", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// load user identities
	var identities []*models.UserIdentity
	for _, user := range users {
		uuid := user.LoginGovUUID.String()
		identity, err := models.FetchUserIdentity(h.db, uuid)
		if err != nil {
			h.logger.Error("Could not get user identity", zap.String("userID", uuid), zap.Error(err))
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		identities = append(identities, identity)
	}

	t := template.Must(template.New("users").Parse(`
		<h1>Select an existing user</h1>
		{{range .}}
			<form method="post" action="/devlocal-auth/login">
				<p id="{{.ID}}">
					<input type="hidden" name="gorilla.csrf.Token" value="` + csrf.Token(r) + `">
					{{.Email}}
					({{if .TspUserID}}tsp{{else if .OfficeUserID}}office{{else}}mymove{{end}})
					<button name="id" value="{{.ID}}" data-hook="existing-user-login">Login</button>
				</p>
			</form>
		{{else}}
			<p><em>No users in the system!</em></p>
		{{end}}

		<h1>Create a new user</h1>
		<form method="post" action="/devlocal-auth/new">
			<p>
				<input type="hidden" name="gorilla.csrf.Token" value="` + csrf.Token(r) + `">
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
}

// AssignUserHandler logs a user in directly
type AssignUserHandler devlocalAuthHandler

// NewAssignUserHandler creates a new AssignUserHandler
func NewAssignUserHandler(ac Context, db *pop.Connection, clientAuthSecretKey string, noSessionTimeout bool) AssignUserHandler {
	handler := AssignUserHandler{
		Context:             ac,
		db:                  db,
		clientAuthSecretKey: clientAuthSecretKey,
		noSessionTimeout:    noSessionTimeout,
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

	loginUser(devlocalAuthHandler(h), user, w, r)
}

// CreateUserHandler creates a new user
type CreateUserHandler devlocalAuthHandler

// NewCreateUserHandler creates a new CreateUserHandler
func NewCreateUserHandler(ac Context, db *pop.Connection, clientAuthSecretKey string, noSessionTimeout bool) CreateUserHandler {
	handler := CreateUserHandler{
		Context:             ac,
		db:                  db,
		clientAuthSecretKey: clientAuthSecretKey,
		noSessionTimeout:    noSessionTimeout,
	}
	return handler
}

// CreateUserHandler creates a user, primarily used in automated testing
func (h CreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := createUser(devlocalAuthHandler(h), w, r)
	jsonOut, _ := json.Marshal(user)
	fmt.Fprintf(w, string(jsonOut))
}

// createUser creates a user
func createUser(h devlocalAuthHandler, w http.ResponseWriter, r *http.Request) models.User {
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
	}
	if verrs.Count() != 0 {
		h.logger.Error("validation errors creating user", zap.Stringer("errors", verrs))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	return user
}

// CreateAndLoginUserHandler creates and then logs in a new user
type CreateAndLoginUserHandler devlocalAuthHandler

// NewCreateAndLoginUserHandler creates a new CreateAndLoginUserHandler
func NewCreateAndLoginUserHandler(ac Context, db *pop.Connection, clientAuthSecretKey string, noSessionTimeout bool) CreateAndLoginUserHandler {
	handler := CreateAndLoginUserHandler{
		Context:             ac,
		db:                  db,
		clientAuthSecretKey: clientAuthSecretKey,
		noSessionTimeout:    noSessionTimeout,
	}
	return handler
}

// CreateAndLoginUserHandler creates a user and logs them in
func (h CreateAndLoginUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	user := createUser(devlocalAuthHandler(h), w, r)
	loginUser(devlocalAuthHandler(h), &user, w, r)
}

func loginUser(handler devlocalAuthHandler, user *models.User, w http.ResponseWriter, r *http.Request) {
	session := auth.SessionFromRequestContext(r)
	if session == nil {
		handler.logger.Error("Session missing")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	userIdentity, err := models.FetchUserIdentity(handler.db, user.LoginGovUUID.String())
	if err == nil { // Someone we know already
		session.IDToken = "devlocal"
		session.UserID = userIdentity.ID
		session.Email = userIdentity.Email
		if userIdentity.ServiceMemberID != nil {
			session.ServiceMemberID = *(userIdentity.ServiceMemberID)
		}

		if userIdentity.OfficeUserID != nil {
			session.OfficeUserID = *(userIdentity.OfficeUserID)
		} else if session.IsOfficeApp() {
			handler.logger.Error("Non-office user authenticated at office site", zap.String("email", session.Email))
			http.Error(w, http.StatusText(401), http.StatusUnauthorized)
			return
		}

		if userIdentity.TspUserID != nil {
			session.TspUserID = *(userIdentity.TspUserID)
		} else if session.IsTspApp() {
			handler.logger.Error("Non-TSP user authenticated at TSP site", zap.String("email", session.Email))
			http.Error(w, http.StatusText(401), http.StatusUnauthorized)
			return
		}

		session.FirstName = userIdentity.FirstName()
		session.LastName = userIdentity.LastName()
		session.Middle = userIdentity.Middle()
	} else {
		handler.logger.Error("Error loading Identity.", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	handler.logger.Info("logged in", zap.Any("session", session))
	auth.WriteSessionCookie(w, session, handler.clientAuthSecretKey, handler.noSessionTimeout, handler.logger)

	lURL := handler.landingURL(session)
	http.Redirect(w, r, lURL, http.StatusSeeOther)
}
