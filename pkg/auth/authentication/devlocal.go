package authentication

import (
	"html/template"
	"net/http"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
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
		<h1>Select a user to login</h1>
		{{range .}}
			<p id="{{.ID}}">
				<form method="post" action="/devlocal-auth/login">
				{{.Email}}
				({{if .OfficeUserID}}office{{else}}mymove{{end}})
				<button name="id" value="{{.ID}}">Login</button>
			</p>
		{{else}}
			<p><em>No users in the system!</em></p>
		{{end}}
	`))
	err = t.Execute(w, identities)
	if err != nil {
		h.logger.Error("Could not render template", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
}

// AssignUserHandler processes a callback from login.gov
type AssignUserHandler struct {
	Context
	db                  *pop.Connection
	clientAuthSecretKey string
	noSessionTimeout    bool
}

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

	session := auth.SessionFromRequestContext(r)
	if session == nil {
		h.logger.Error("Session missing")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	lURL := h.landingURL(session)

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

	userIdentity, err := models.FetchUserIdentity(h.db, user.LoginGovUUID.String())
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
			h.logger.Error("Non-office user authenticated at office site", zap.String("email", session.Email))
			http.Error(w, http.StatusText(401), http.StatusUnauthorized)
			return
		}

		session.FirstName = userIdentity.FirstName()
		session.LastName = userIdentity.LastName()
		session.Middle = userIdentity.Middle()
	} else {
		h.logger.Error("Error loading Identity.", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	h.logger.Info("logged in", zap.Any("session", session))
	auth.WriteSessionCookie(w, session, h.clientAuthSecretKey, h.noSessionTimeout, h.logger)
	http.Redirect(w, r, lURL, http.StatusSeeOther)
}
