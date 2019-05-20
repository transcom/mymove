package authentication

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/gorilla/csrf"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

const (
	// MilMoveUserType is the type of user for a Service Member
	MilMoveUserType string = "milmove"
	// OfficeUserType is the type of user for an Office user
	OfficeUserType string = "office"
	// TspUserType is the type of user for a TSP user
	TspUserType string = "tsp"
	// DpsUserType is the type of user for a DPS user
	DpsUserType string = "dps"
	// AdminUserType is the type of user for an admin user
	AdminUserType string = "admin"
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
	// Truncate the list if larger than 25
	if len(identities) > 25 {
		identities = identities[:25]
	}

	type TemplateData struct {
		Identities      []models.UserIdentity
		MilMoveUserType string
		OfficeUserType  string
		TspUserType     string
		DpsUserType     string
		AdminUserType   string
		CsrfToken       string
	}

	templateData := TemplateData{
		Identities:      identities,
		MilMoveUserType: MilMoveUserType,
		OfficeUserType:  OfficeUserType,
		TspUserType:     TspUserType,
		DpsUserType:     DpsUserType,
		AdminUserType:   AdminUserType,
		// Build CSRF token instead of grabbing from middleware. Otherwise throws errors when accessed directly.
		CsrfToken: csrf.Token(r),
	}

	t := template.Must(template.New("users").Parse(`
	  <html>
	  <head>
		<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
	  </head>
	  <body class="py-4">
		<div class="container">
		  <div class="row mb-3">
			<div class="col-md-8">
			  <h2 class="mt-4">Select an Existing User</h1>
			  <p>Showing the first 25 users:</p>
			  {{range .Identities}}
				<form method="post" action="/devlocal-auth/login">
					<p id="{{.ID}}">
						<input type="hidden" name="gorilla.csrf.Token" value="{{$.CsrfToken}}">
						{{.Email}}
						{{if .IsSuperuser}}
						  ({{$.AdminUserType}})
						  <input type="hidden" name="userType" value="{{$.AdminUserType}}">
						{{else if .DpsUserID}}
						  ({{$.DpsUserType}})
						  <input type="hidden" name="userType" value="{{$.DpsUserType}}">
						{{else if .TspUserID}}
						  ({{$.TspUserType}})
						  <input type="hidden" name="userType" value="{{$.TspUserType}}">
						{{else if .OfficeUserID}}
						  ({{$.OfficeUserType}})
						  <input type="hidden" name="userType" value="{{$.OfficeUserType}}">
						{{else}}
						  ({{$.MilMoveUserType}})
						  <input type="hidden" name="userType" value="{{$.MilMoveUserType}}">
						{{end}}
						<input type="hidden" name="id" value="{{.ID}}" />
						<button type="submit" value="{{.ID}}" data-hook="existing-user-login">Login</button>
					</p>
				</form>
			  {{else}}
				<p><em>No users in the system!</em></p>
			  {{end}}
			</div>

			<div class="col-md-4">
			  <h2 class="mt-4">Create a New User</h1>
			  <p>Creating new users for different sites will mean you need to redirect and log in yourself.</p>
			  <form method="post" action="/devlocal-auth/new">
				<p>
				  <input type="hidden" name="gorilla.csrf.Token" value="{{.CsrfToken}}">
				  <input type="hidden" name="userType" value="{{.MilMoveUserType}}">
				  <button type="submit" data-hook="new-user-login-{{.MilMoveUserType}}">Create a New {{.MilMoveUserType}} User</button>
				</p>
			  </form>
			  <form method="post" action="/devlocal-auth/new">
				<p>
				  <input type="hidden" name="gorilla.csrf.Token" value="{{.CsrfToken}}">
				  <input type="hidden" name="userType" value="{{.OfficeUserType}}">
				  <button type="submit" data-hook="new-user-login-{{.OfficeUserType}}">Create a New {{.OfficeUserType}} User</button>
				</p>
			  </form>
			  <form method="post" action="/devlocal-auth/new">
				<p>
				  <input type="hidden" name="gorilla.csrf.Token" value="{{.CsrfToken}}">
				  <input type="hidden" name="userType" value="{{.TspUserType}}">
				  <button type="submit" data-hook="new-user-login-{{.TspUserType}}">Create a New {{.TspUserType}} User</button>
				</p>
			  </form>
			  <form method="post" action="/devlocal-auth/new">
				<p>
				  <input type="hidden" name="gorilla.csrf.Token" value="{{.CsrfToken}}">
				  <input type="hidden" name="userType" value="{{.DpsUserType}}">
				  <button type="submit" data-hook="new-user-login-{{.DpsUserType}}">Create a New {{.DpsUserType}} User</button>
				</p>
			  </form>
			  <form method="post" action="/devlocal-auth/new">
				<p>
				  <input type="hidden" name="gorilla.csrf.Token" value="{{.CsrfToken}}">
				  <input type="hidden" name="userType" value="{{.AdminUserType}}">
				  <button type="submit" data-hook="new-user-login-{{.AdminUserType}}">Create a New {{.AdminUserType}} User</button>
				</p>
			  </form>
			</div>
		  </div>
		</div> <!-- container -->
	  </body>
	  </html>
	`))
	err = t.Execute(w, templateData)
	if err != nil {
		h.logger.Error("Could not render template", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
}

type devlocalAuthHandler struct {
	Context
	db                  *pop.Connection
	appnames            auth.ApplicationServername
	clientAuthSecretKey string
	noSessionTimeout    bool
	useSecureCookie     bool
}

// AssignUserHandler logs a user in directly
type AssignUserHandler devlocalAuthHandler

// NewAssignUserHandler creates a new AssignUserHandler
func NewAssignUserHandler(ac Context, db *pop.Connection, appnames auth.ApplicationServername, clientAuthSecretKey string, noSessionTimeout bool, useSecureCookie bool) AssignUserHandler {
	handler := AssignUserHandler{
		Context:             ac,
		db:                  db,
		appnames:            appnames,
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

	userType := r.PostFormValue("userType")
	session, err := loginUser(devlocalAuthHandler(h), user, userType, w, r)
	if err != nil {
		return
	}
	if session == nil {
		return
	}
	http.Redirect(w, r, h.landingURL(session), http.StatusSeeOther)
}

// CreateUserHandler creates a new user
type CreateUserHandler devlocalAuthHandler

// NewCreateUserHandler creates a new CreateUserHandler
func NewCreateUserHandler(ac Context, db *pop.Connection, appnames auth.ApplicationServername, clientAuthSecretKey string, noSessionTimeout bool, useSecureCookie bool) CreateUserHandler {
	handler := CreateUserHandler{
		Context:             ac,
		db:                  db,
		appnames:            appnames,
		clientAuthSecretKey: clientAuthSecretKey,
		noSessionTimeout:    noSessionTimeout,
		useSecureCookie:     useSecureCookie,
	}
	return handler
}

// CreateUserHandler creates a user, primarily used in automated testing
func (h CreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, userType := createUser(devlocalAuthHandler(h), w, r)
	if user == nil {
		return
	}
	session, err := loginUser(devlocalAuthHandler(h), user, userType, w, r)
	if err != nil {
		return
	}
	if session == nil {
		return
	}
	jsonOut, _ := json.Marshal(user)
	fmt.Fprintf(w, string(jsonOut))
}

// CreateAndLoginUserHandler creates and then logs in a new user
type CreateAndLoginUserHandler devlocalAuthHandler

// NewCreateAndLoginUserHandler creates a new CreateAndLoginUserHandler
func NewCreateAndLoginUserHandler(ac Context, db *pop.Connection, appnames auth.ApplicationServername, clientAuthSecretKey string, noSessionTimeout bool, useSecureCookie bool) CreateAndLoginUserHandler {
	handler := CreateAndLoginUserHandler{
		Context:             ac,
		db:                  db,
		appnames:            appnames,
		clientAuthSecretKey: clientAuthSecretKey,
		noSessionTimeout:    noSessionTimeout,
		useSecureCookie:     useSecureCookie,
	}
	return handler
}

// CreateAndLoginUserHandler creates a user and logs them in
func (h CreateAndLoginUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, userType := createUser(devlocalAuthHandler(h), w, r)
	if user == nil {
		return
	}
	session, err := loginUser(devlocalAuthHandler(h), user, userType, w, r)
	if err != nil {
		return
	}
	if session == nil {
		return
	}
	http.Redirect(w, r, h.landingURL(session), http.StatusSeeOther)
}

// createUser creates a user
func createUser(h devlocalAuthHandler, w http.ResponseWriter, r *http.Request) (*models.User, string) {
	id := uuid.Must(uuid.NewV4())

	// Set up some defaults that we can pass in from a form
	firstName := r.PostFormValue("firstName")
	if firstName == "" {
		firstName = "Alice"
	}
	lastName := r.PostFormValue("lastName")
	if lastName == "" {
		lastName = "Bob"
	}
	telephone := r.PostFormValue("telephone")
	if telephone == "" {
		telephone = "333-333-3333"
	}
	email := r.PostFormValue("email")
	if email == "" {
		// Time alone doesn't guarantee uniqueness if a system is being automated
		// To add some more uniqueness without making the email unreadable a UUID adds a nonce
		now := time.Now()
		guid, _ := uuid.NewV4()
		nonce := strings.Split(guid.String(), "-")[4]
		email = fmt.Sprintf("%s-%s@example.com", now.Format("20060102150405"), nonce)
	}

	// Create the User (which is the basis of all Service Members)
	user := models.User{
		LoginGovUUID:  id,
		LoginGovEmail: email,
		IsSuperuser:   false,
	}

	userType := r.PostFormValue("userType")
	verrs, err := h.db.ValidateAndCreate(&user)
	if err != nil {
		h.logger.Error("could not create user", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return nil, userType
	}
	if verrs.Count() != 0 {
		h.logger.Error("validation errors creating user", zap.Stringer("errors", verrs))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return nil, userType
	}

	switch userType {
	case OfficeUserType:
		// Now create the Truss JPPSO
		address := models.Address{
			StreetAddress1: "1333 Minna St",
			City:           "San Francisco",
			State:          "CA",
			PostalCode:     "94115",
		}

		verrs, err := h.db.ValidateAndSave(&address)
		if err != nil {
			h.logger.Error("could not create address", zap.Error(err))
		}
		if verrs.HasAny() {
			h.logger.Error("validation errors creating address", zap.Stringer("errors", verrs))
		}

		office := models.TransportationOffice{
			Name:      "Truss",
			AddressID: address.ID,
			Latitude:  37.7678355,
			Longitude: -122.4199298,
			Hours:     models.StringPointer("0900-1800 Mon-Sat"),
		}

		verrs, err = h.db.ValidateAndSave(&office)
		if err != nil {
			h.logger.Error("could not create office", zap.Error(err))
		}
		if verrs.HasAny() {
			h.logger.Error("validation errors creating office", zap.Stringer("errors", verrs))
		}

		officeUser := models.OfficeUser{
			FirstName:              firstName,
			LastName:               lastName,
			Telephone:              telephone,
			TransportationOfficeID: office.ID,
			Email:                  email,
		}
		if user.ID != uuid.Nil {
			officeUser.UserID = &user.ID
		}

		verrs, err = h.db.ValidateAndSave(&officeUser)
		if err != nil {
			h.logger.Error("could not create office user", zap.Error(err))
		}
		if verrs.HasAny() {
			h.logger.Error("validation errors creating office user", zap.Stringer("errors", verrs))
		}
	case TspUserType:
		var tsp models.TransportationServiceProvider
		h.db.Where("standard_carrier_alpha_code = $1", "TRS1").First(&tsp)
		if tsp.ID == uuid.Nil {
			// TSP not found, create one for Truss
			tsp = models.TransportationServiceProvider{
				StandardCarrierAlphaCode: "TRSS",
			}
			verrs, err := h.db.ValidateAndSave(&tsp)
			if err != nil {
				h.logger.Error("could not create TSP", zap.Error(err))
			}
			if verrs.HasAny() {
				h.logger.Error("validation errors creating TSP", zap.Stringer("errors", verrs))
			}
		}

		tspUser := models.TspUser{
			FirstName:                       firstName,
			LastName:                        lastName,
			Telephone:                       telephone,
			TransportationServiceProviderID: tsp.ID,
			Email:                           email,
		}
		if user.ID != uuid.Nil {
			tspUser.UserID = &user.ID
		}

		verrs, err := h.db.ValidateAndSave(&tspUser)
		if err != nil {
			h.logger.Error("could not create tsp user", zap.Error(err))
		}
		if verrs.HasAny() {
			h.logger.Error("validation errors creating tsp user", zap.Stringer("errors", verrs))
		}
	case DpsUserType:
		dpsUser := models.DpsUser{
			LoginGovEmail: email,
		}

		verrs, err := h.db.ValidateAndSave(&dpsUser)
		if err != nil {
			h.logger.Error("could not create dps user", zap.Error(err))
		}
		if verrs.HasAny() {
			h.logger.Error("validation errors creating dps user", zap.Stringer("errors", verrs))
		}
	case AdminUserType:
		user.IsSuperuser = true
		verrs, err := h.db.ValidateAndSave(&user)
		if err != nil {
			h.logger.Error("could not create admin user", zap.Error(err))
		}
		if verrs.HasAny() {
			h.logger.Error("validation errors creating admin user", zap.Stringer("errors", verrs))
		}
	}

	return &user, userType
}

// createSession creates a new session for the user
func createSession(h devlocalAuthHandler, user *models.User, userType string, w http.ResponseWriter, r *http.Request) (*auth.Session, error) {
	// Preference any session already in the request context. Otherwise just create a new empty session.
	session := auth.SessionFromRequestContext(r)
	if session == nil {
		session = &auth.Session{}
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
	session.IsSuperuser = userIdentity.IsSuperuser

	// Set the app

	// Keep the logic for redirection separate from setting the session user ids
	switch userType {
	case OfficeUserType:
		session.ApplicationName = auth.OfficeApp
		session.Hostname = h.appnames.OfficeServername
	case TspUserType:
		session.ApplicationName = auth.TspApp
		session.Hostname = h.appnames.TspServername
	case AdminUserType:
		session.ApplicationName = auth.AdminApp
		session.Hostname = h.appnames.AdminServername
	default:
		session.ApplicationName = auth.MilApp
		session.Hostname = h.appnames.MilServername
	}

	if userIdentity.ServiceMemberID != nil {
		session.ServiceMemberID = *(userIdentity.ServiceMemberID)
	}

	if userIdentity.OfficeUserID != nil && (session.IsOfficeApp() || userType == OfficeUserType) {
		session.OfficeUserID = *(userIdentity.OfficeUserID)
	}

	if userIdentity.TspUserID != nil && (session.IsTspApp() || userType == TspUserType) {
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

	if !session.IsSuperuser && session.IsAdminApp() {
		return errors.Errorf("Non-superuser %s authenticated at admin site", session.Email)
	}

	return nil
}

// loginUser creates a session for the user and verifies the session against the app
func loginUser(h devlocalAuthHandler, user *models.User, userType string, w http.ResponseWriter, r *http.Request) (*auth.Session, error) {
	session, err := createSession(devlocalAuthHandler(h), user, userType, w, r)
	if err != nil {
		h.logger.Error("Could not create session", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return nil, err
	}

	if session.Disabled {
		h.logger.Info("Disabled user requesting authentication", zap.Error(err), zap.String("email", session.Email))
		http.Error(w, http.StatusText(403), http.StatusForbidden)
		return nil, nil
	}

	err = verifySessionWithApp(session)
	if err != nil {
		h.logger.Error("User unauthorized", zap.Error(err))
		http.Error(w, http.StatusText(401), http.StatusUnauthorized)
		return nil, err
	}
	return session, nil
}
