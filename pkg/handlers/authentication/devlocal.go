package authentication

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/csrf"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

const (
	// MilMoveUserType is the type of user for a Service Member
	MilMoveUserType string = "milmove"
	// TOOOfficeUserType is the type of user for an Office user
	TOOOfficeUserType string = "TOO office"
	// TIOOfficeUserType is the type of user for an Office user
	TIOOfficeUserType string = "TIO office"
	// ServicesCounselorOfficeUserType is the type of user for an Office User
	ServicesCounselorOfficeUserType string = "Services Counselor office"
	// PrimeSimulatorOfficeUserType is the type of user for an Office user
	PrimeSimulatorOfficeUserType string = "Prime Simulator office"
	// QaeOfficeUserType is a type of user for an Office user
	QaeOfficeUserType string = "QAE office"
	// CustomerServiceRepresentativeOfficeUserType is the Customer Service Representative type of user for an Office user
	CustomerServiceRepresentativeOfficeUserType string = "CSR office"
	// MultiRoleOfficeUserType has all the Office user roles
	MultiRoleOfficeUserType string = "Multi role office"
	// AdminUserType is the type of user for an admin user
	AdminUserType string = "admin"
)

// UserListHandler handles redirection
type UserListHandler struct {
	Context
	handlers.HandlerConfig
}

// NewUserListHandler returns a new UserListHandler
func NewUserListHandler(ac Context, hc handlers.HandlerConfig) UserListHandler {
	handler := UserListHandler{
		Context:       ac,
		HandlerConfig: hc,
	}
	return handler
}

// UserListHandler lists users in the local database for local login
func (h UserListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	appCtx := h.AppContextFromRequest(r)
	limit := 100
	identities, err := models.FetchAppUserIdentities(appCtx.DB(), appCtx.Session().ApplicationName, limit)
	if err != nil {
		appCtx.Logger().Error("Could not load list of users", zap.Error(err))
		http.Error(w,
			fmt.Sprintf("%s - Could not load list of users, try migrating the DB", http.StatusText(500)),
			http.StatusInternalServerError)
		return
	}

	// get roles for all identities
	for i := range identities {
		rerr := appCtx.DB().RawQuery(`SELECT * FROM roles
									WHERE id in (select role_id from users_roles
										where deleted_at is null and user_id = ?)`, identities[i].ID).All(&identities[i].Roles)
		if rerr != nil {
			appCtx.Logger().Error("Could not load identity roles", zap.Error(rerr))
			http.Error(w,
				fmt.Sprintf("%s - Could not load list of users, try migrating the DB", http.StatusText(500)),
				http.StatusInternalServerError)
			return
		}
	}

	// grab all GBLOCs
	var gblocList []string
	err = appCtx.DB().RawQuery("select distinct gbloc from transportation_offices").All(&gblocList)
	if err != nil {
		appCtx.Logger().Error("Could not load gblocs", zap.Error(err))
		http.Error(w,
			fmt.Sprintf("%s - Could not load gblocs, try migrating the DB", http.StatusText(500)),
			http.StatusInternalServerError)
		return
	}

	type TemplateData struct {
		Identities                                  []models.UserIdentity
		Gblocs                                      []string
		GblocDefault                                string
		IsMilApp                                    bool
		MilMoveUserType                             string
		IsOfficeApp                                 bool
		TOOOfficeUserType                           string
		TIOOfficeUserType                           string
		ServicesCounselorOfficeUserType             string
		PrimeSimulatorOfficeUserType                string
		QaeOfficeUserType                           string
		CustomerServiceRepresentativeOfficeUserType string
		MultiRoleOfficeUserType                     string
		IsAdminApp                                  bool
		AdminUserType                               string
		CsrfToken                                   string
		QueryLimit                                  int
	}

	templateData := TemplateData{
		Identities:                      identities,
		Gblocs:                          gblocList,
		GblocDefault:                    "KKFA", // Most seed data is tied to this
		IsMilApp:                        auth.MilApp == appCtx.Session().ApplicationName,
		MilMoveUserType:                 MilMoveUserType,
		IsOfficeApp:                     auth.OfficeApp == appCtx.Session().ApplicationName,
		TOOOfficeUserType:               TOOOfficeUserType,
		TIOOfficeUserType:               TIOOfficeUserType,
		ServicesCounselorOfficeUserType: ServicesCounselorOfficeUserType,
		PrimeSimulatorOfficeUserType:    PrimeSimulatorOfficeUserType,
		QaeOfficeUserType:               QaeOfficeUserType,
		CustomerServiceRepresentativeOfficeUserType: CustomerServiceRepresentativeOfficeUserType,
		MultiRoleOfficeUserType:                     MultiRoleOfficeUserType,
		IsAdminApp:                                  auth.AdminApp == appCtx.Session().ApplicationName,
		AdminUserType:                               AdminUserType,
		// Build CSRF token instead of grabbing from middleware. Otherwise throws errors when accessed directly.
		CsrfToken:  csrf.Token(r),
		QueryLimit: limit,
	}

	gblocSelectHTML := `
		<label for="gblocSelect">Select GBLOC:</label>
		<select id="gblocSelect" name="gbloc">
			{{ range $index, $element := .Gblocs }}
				{{if eq $element $.GblocDefault}}
					<option value="{{$element}}" selected="">{{$element}}</option>
				{{else}}
					<option value="{{$element}}">{{$element}}</option>
				{{end}}
			{{ end }}
		</select>`
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
			  <h4>Login with user email:</h4>
				<form method="post" action="/devlocal-auth/login">
					<p>
						<input type="hidden" name="gorilla.csrf.Token" value="{{$.CsrfToken}}">
						<input type="hidden" name="userType" value="{{if $.IsOfficeApp}}{{$.TOOOfficeUserType}}{{else}}{{$.MilMoveUserType}}{{end}}">
						<label for="email">User Email</label>
						<input type="text" name="email" size="60">
						<button type="submit" data-hook="existing-user-login">Login</button>
					</p>
				</form>
			  <h4>Showing the first {{$.QueryLimit}} users by creation date in descending order:</h4>
			  {{range .Identities}}
				<form method="post" action="/devlocal-auth/login">
					<p id="{{.ID}}">
						<input type="hidden" name="gorilla.csrf.Token" value="{{$.CsrfToken}}">
						{{.Email}}
						{{if .AdminUserID}}
						  ({{$.AdminUserType}})
						  <input type="hidden" name="userType" value="{{$.AdminUserType}}">
						{{else if .OfficeUserID}}

						  <input type="hidden" name="userType" value="{{$.TOOOfficeUserType}}">
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
			  {{ if $.IsMilApp }}
				  <form method="post" action="/devlocal-auth/new">
					<p>
					  <input type="hidden" name="gorilla.csrf.Token" value="{{.CsrfToken}}">
					  <input type="hidden" name="userType" value="{{.MilMoveUserType}}">
					  <label for="email">User Email</label>
					  <input type="text" name="email" value="">
					  <button type="submit" data-hook="new-user-login-{{.MilMoveUserType}}">Create a New {{.MilMoveUserType}} User</button>
					</p>
				  </form>
			  {{else if $.IsAdminApp }}
				  <form method="post" action="/devlocal-auth/new">
					<p>
					  <input type="hidden" name="gorilla.csrf.Token" value="{{.CsrfToken}}">
					  <input type="hidden" name="userType" value="{{.AdminUserType}}">
					  <button type="submit" data-hook="new-user-login-{{.AdminUserType}}">Create a New {{.AdminUserType}} User</button>
					</p>
				  </form>
			  {{else if $.IsOfficeApp }}
				  <form method="post" action="/devlocal-auth/new">
					<p>
					  <input type="hidden" name="gorilla.csrf.Token" value="{{.CsrfToken}}">
					  <input type="hidden" name="userType" value="{{.TOOOfficeUserType}}">
					  ` + gblocSelectHTML + `
					  <button type="submit" data-hook="new-user-login-{{.TOOOfficeUserType}}">Create a New {{.TOOOfficeUserType}} User</button>
					</p>
				  </form>

				  <form method="post" action="/devlocal-auth/new">
					<p>
					  <input type="hidden" name="gorilla.csrf.Token" value="{{.CsrfToken}}">
					  <input type="hidden" name="userType" value="{{.TIOOfficeUserType}}">
					  ` + gblocSelectHTML + `
					  <button type="submit" data-hook="new-user-login-{{.TIOOfficeUserType}}">Create a New {{.TIOOfficeUserType}} User</button>
					</p>
				  </form>

				  <form method="post" action="/devlocal-auth/new">
				  <p>
					<input type="hidden" name="gorilla.csrf.Token" value="{{.CsrfToken}}">
					<input type="hidden" name="userType" value="{{.ServicesCounselorOfficeUserType}}">
					` + gblocSelectHTML + `
					<button type="submit" data-hook="new-user-login-{{.ServicesCounselorOfficeUserType}}">Create a New {{.ServicesCounselorOfficeUserType}} User</button>
				  </p>
				  </form>

				  <form method="post" action="/devlocal-auth/new">
				  <p>
					<input type="hidden" name="gorilla.csrf.Token" value="{{.CsrfToken}}">
					<input type="hidden" name="userType" value="{{.PrimeSimulatorOfficeUserType}}">
					` + gblocSelectHTML + `
					<button type="submit" data-hook="new-user-login-{{.PrimeSimulatorOfficeUserType}}">Create a New {{.PrimeSimulatorOfficeUserType}} User</button>
				  </p>
				</form>

				<form method="post" action="/devlocal-auth/new">
					<p>
					  <input type="hidden" name="gorilla.csrf.Token" value="{{.CsrfToken}}">
					  <input type="hidden" name="userType" value="{{.QaeOfficeUserType}}">
					  ` + gblocSelectHTML + `
					  <button type="submit" data-hook="new-user-login-{{.QaeOfficeUserType}}">Create a New {{.QaeOfficeUserType}} User</button>
					</p>
				  </form>

				<form method="post" action="/devlocal-auth/new">
				  <p>
					<input type="hidden" name="gorilla.csrf.Token" value="{{.CsrfToken}}">
					<input type="hidden" name="userType" value="{{.CustomerServiceRepresentativeOfficeUserType}}">
					` + gblocSelectHTML + `
					<button type="submit" data-hook="new-user-login-{{.CustomerServiceRepresentativeOfficeUserType}}">Create a New {{.CustomerServiceRepresentativeOfficeUserType}} User</button>
				  </p>
				</form>

				<form method="post" action="/devlocal-auth/new">
					<p>
					  <input type="hidden" name="gorilla.csrf.Token" value="{{.CsrfToken}}">
					  <input type="hidden" name="userType" value="{{.MultiRoleOfficeUserType}}">
					  ` + gblocSelectHTML + `
					  <button type="submit" data-hook="new-user-login-{{.MultiRoleOfficeUserType}}">Create a New {{.MultiRoleOfficeUserType}} User</button>
					</p>
				  </form>

			  {{end}}
			</div>
		  </div>
		</div> <!-- container -->
	  </body>
	  </html>
	`))
	err = t.Execute(w, templateData)
	if err != nil {
		appCtx.Logger().Error("Could not render template", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
}

type devlocalAuthHandler struct {
	Context
	handlers.HandlerConfig
}

// AssignUserHandler logs a user in directly
type AssignUserHandler devlocalAuthHandler

// NewAssignUserHandler creates a new AssignUserHandler
func NewAssignUserHandler(ac Context, hc handlers.HandlerConfig) AssignUserHandler {
	handler := AssignUserHandler{
		Context:       ac,
		HandlerConfig: hc,
	}
	return handler
}

// AssignUserHandler logs in a user locally using a user id or email
func (h AssignUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	appCtx := h.AppContextFromRequest(r)
	userID := r.PostFormValue("id")
	email := r.PostFormValue("email")
	if userID == "" && email == "" {
		appCtx.Logger().Error("No user id or email specified")
		http.Redirect(w, r, "/devlocal-auth/login", http.StatusTemporaryRedirect)
		return
	}

	appCtx.Logger().Info("New Devlocal Login",
		zap.String("userID", userID),
		zap.String("email", email))

	var user *models.User
	if userID != "" {
		userUUID := uuid.Must(uuid.FromString(userID))
		var err error
		user, err = models.GetUser(appCtx.DB(), userUUID)
		if err != nil {
			appCtx.Logger().Error("Could not load user from user id", zap.String("userID", userID), zap.Error(err))
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		}
	} else if email != "" {
		var err error
		user, err = models.GetUserFromEmail(appCtx.DB(), email)
		if err != nil {
			appCtx.Logger().Error("Could not load user from email", zap.String("email", email), zap.Error(err))
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		}
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
func NewCreateUserHandler(ac Context, hc handlers.HandlerConfig) CreateUserHandler {
	handler := CreateUserHandler{
		Context:       ac,
		HandlerConfig: hc,
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
	fmt.Fprint(w, string(jsonOut))
}

// CreateAndLoginUserHandler creates and then logs in a new user
type CreateAndLoginUserHandler devlocalAuthHandler

// NewCreateAndLoginUserHandler creates a new CreateAndLoginUserHandler
func NewCreateAndLoginUserHandler(ac Context, hc handlers.HandlerConfig) CreateAndLoginUserHandler {
	handler := CreateAndLoginUserHandler{
		Context:       ac,
		HandlerConfig: hc,
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
	appCtx := h.HandlerConfig.AppContextFromRequest(r)
	id := uuid.Must(uuid.NewV4()).String()

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
	userType := r.PostFormValue("userType")
	email := r.PostFormValue("email")
	if email == "" {
		// Time alone doesn't guarantee uniqueness if a system is being automated
		// To add some more uniqueness without making the email unreadable a UUID adds a nonce
		now := time.Now()
		guid, _ := uuid.NewV4()
		nonce := strings.Split(guid.String(), "-")[4]
		prefix := now.Format("20060102150405")
		i := strings.Index(userType, " office")
		if i > 0 {
			// turn e.g. Prime Simulator office into "prime-simulator"
			re := regexp.MustCompile(`[^\w]+`)
			prefix = re.ReplaceAllString(strings.ToLower(userType[0:i]), "-") + "-" + prefix
		}
		email = fmt.Sprintf("%s-%s@example.com", prefix, nonce)
	}
	gbloc := r.PostFormValue("gbloc")
	if gbloc == "" {
		gbloc = "KKFA" // most seed data uses this
	}

	// Create the User (which is the basis of all Service Members)
	user := models.User{
		OktaID:    id,
		OktaEmail: email,
		Active:    true,
	}

	verrs, err := appCtx.DB().ValidateAndCreate(&user)
	if err != nil {
		appCtx.Logger().Error("could not create user", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return nil, userType
	}
	if verrs.Count() != 0 {
		appCtx.Logger().Error("validation errors creating user", zap.Stringer("errors", verrs))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return nil, userType
	}

	switch userType {
	case MilMoveUserType:
		newServiceMember := models.ServiceMember{
			UserID:       user.ID,
			CacValidated: true,
		}
		smVerrs, smErr := models.SaveServiceMember(appCtx, &newServiceMember)
		if smVerrs.HasAny() || smErr != nil {
			appCtx.Logger().Error("Error creating service member for user", zap.Error(smErr))
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		}
	case TOOOfficeUserType:
		// Now create the Truss JPPSO
		address := models.Address{
			StreetAddress1: "1333 Minna St",
			City:           "San Francisco",
			State:          "CA",
			PostalCode:     "94115",
			County:         "SAINT CLAIR",
		}

		verrs, err := appCtx.DB().ValidateAndSave(&address)
		if err != nil {
			appCtx.Logger().Error("could not create address", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating address", zap.Stringer("errors", verrs))
		}

		role := roles.Role{}
		err = appCtx.DB().Where("role_type = $1", roles.RoleTypeTOO).First(&role)
		if err != nil {
			appCtx.Logger().Error("could not fetch role task_ordering_officer", zap.Error(err))
		}

		usersRole := models.UsersRoles{
			UserID: user.ID,
			RoleID: role.ID,
		}

		verrs, err = appCtx.DB().ValidateAndSave(&usersRole)
		if err != nil {
			appCtx.Logger().Error("could not create user role", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating user role", zap.Stringer("errors", verrs))
		}

		office := models.TransportationOffice{
			Name:      "Truss",
			AddressID: address.ID,
			Latitude:  37.7678355,
			Longitude: -122.4199298,
			Hours:     models.StringPointer("0900-1800 Mon-Sat"),
			Gbloc:     gbloc,
		}

		verrs, err = appCtx.DB().ValidateAndSave(&office)
		if err != nil {
			appCtx.Logger().Error("could not create office", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating office", zap.Stringer("errors", verrs))
		}

		officeUser := models.OfficeUser{
			FirstName:              firstName,
			LastName:               lastName,
			Telephone:              telephone,
			TransportationOfficeID: office.ID,
			Email:                  email,
			Active:                 true,
		}
		if user.ID != uuid.Nil {
			officeUser.UserID = &user.ID
		}

		verrs, err = appCtx.DB().ValidateAndSave(&officeUser)
		if err != nil {
			appCtx.Logger().Error("could not create office user", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating office user", zap.Stringer("errors", verrs))
		}
	case TIOOfficeUserType:
		// Now create the Truss JPPSO
		address := models.Address{
			StreetAddress1: "1333 Minna St",
			City:           "San Francisco",
			State:          "CA",
			PostalCode:     "94115",
			County:         "SAINT CLAIR",
		}

		verrs, err := appCtx.DB().ValidateAndSave(&address)
		if err != nil {
			appCtx.Logger().Error("could not create address", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating address", zap.Stringer("errors", verrs))
		}

		role := roles.Role{}
		err = appCtx.DB().Where("role_type = $1", roles.RoleTypeTIO).First(&role)
		if err != nil {
			appCtx.Logger().Error("could not fetch role transporation_invoicing_officer", zap.Error(err))
		}
		usersRole := models.UsersRoles{
			UserID: user.ID,
			RoleID: role.ID,
		}

		verrs, err = appCtx.DB().ValidateAndSave(&usersRole)
		if err != nil {
			appCtx.Logger().Error("could not create user role", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating user role", zap.Stringer("errors", verrs))
		}

		office := models.TransportationOffice{
			Name:      "Truss",
			AddressID: address.ID,
			Latitude:  37.7678355,
			Longitude: -122.4199298,
			Hours:     models.StringPointer("0900-1800 Mon-Sat"),
			Gbloc:     gbloc,
		}

		verrs, err = appCtx.DB().ValidateAndSave(&office)
		if err != nil {
			appCtx.Logger().Error("could not create office", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating office", zap.Stringer("errors", verrs))
		}

		officeUser := models.OfficeUser{
			FirstName:              firstName,
			LastName:               lastName,
			Telephone:              telephone,
			TransportationOfficeID: office.ID,
			Email:                  email,
			Active:                 true,
		}
		if user.ID != uuid.Nil {
			officeUser.UserID = &user.ID
		}

		verrs, err = appCtx.DB().ValidateAndSave(&officeUser)
		if err != nil {
			appCtx.Logger().Error("could not create office user", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating office user", zap.Stringer("errors", verrs))
		}
	case ServicesCounselorOfficeUserType:
		// Now create the Truss JPPSO
		address := models.Address{
			StreetAddress1: "1333 Minna St",
			City:           "San Francisco",
			State:          "CA",
			PostalCode:     "94115",
			County:         "SAINT CLAIR",
		}

		verrs, err := appCtx.DB().ValidateAndSave(&address)
		if err != nil {
			appCtx.Logger().Error("could not create address", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating address", zap.Stringer("errors", verrs))
		}

		role := roles.Role{}
		err = appCtx.DB().Where("role_type = $1", "services_counselor").First(&role)
		if err != nil {
			appCtx.Logger().Error("could not fetch role services_counselor", zap.Error(err))
		}
		usersRole := models.UsersRoles{
			UserID: user.ID,
			RoleID: role.ID,
		}

		verrs, err = appCtx.DB().ValidateAndSave(&usersRole)
		if err != nil {
			appCtx.Logger().Error("could not create user role", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating user role", zap.Stringer("errors", verrs))
		}

		office := models.TransportationOffice{
			Name:      "Truss",
			AddressID: address.ID,
			Latitude:  37.7678355,
			Longitude: -122.4199298,
			Hours:     models.StringPointer("0900-1800 Mon-Sat"),
			Gbloc:     gbloc,
		}

		verrs, err = appCtx.DB().ValidateAndSave(&office)
		if err != nil {
			appCtx.Logger().Error("could not create office", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating office", zap.Stringer("errors", verrs))
		}

		officeUser := models.OfficeUser{
			FirstName:              firstName,
			LastName:               lastName,
			Telephone:              telephone,
			TransportationOfficeID: office.ID,
			Email:                  email,
			Active:                 true,
		}
		if user.ID != uuid.Nil {
			officeUser.UserID = &user.ID
		}

		verrs, err = appCtx.DB().ValidateAndSave(&officeUser)
		if err != nil {
			appCtx.Logger().Error("could not create office user", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating office user", zap.Stringer("errors", verrs))
		}
	case PrimeSimulatorOfficeUserType:
		// Now create the Truss JPPSO
		address := models.Address{
			StreetAddress1: "1333 Minna St",
			City:           "San Francisco",
			State:          "CA",
			PostalCode:     "94115",
			County:         "SAINT CLAIR",
		}

		verrs, err := appCtx.DB().ValidateAndSave(&address)
		if err != nil {
			appCtx.Logger().Error("could not create address", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating address", zap.Stringer("errors", verrs))
		}

		role := roles.Role{}
		err = appCtx.DB().Where("role_type = $1", "prime_simulator").First(&role)
		if err != nil {
			appCtx.Logger().Error("could not fetch role prime_simulator", zap.Error(err))
		}
		usersRole := models.UsersRoles{
			UserID: user.ID,
			RoleID: role.ID,
		}

		verrs, err = appCtx.DB().ValidateAndSave(&usersRole)
		if err != nil {
			appCtx.Logger().Error("could not create user role", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating user role", zap.Stringer("errors", verrs))
		}

		office := models.TransportationOffice{
			Name:      "Truss",
			AddressID: address.ID,
			Latitude:  37.7678355,
			Longitude: -122.4199298,
			Hours:     models.StringPointer("0900-1800 Mon-Sat"),
			Gbloc:     gbloc,
		}

		verrs, err = appCtx.DB().ValidateAndSave(&office)
		if err != nil {
			appCtx.Logger().Error("could not create office", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating office", zap.Stringer("errors", verrs))
		}

		officeUser := models.OfficeUser{
			FirstName:              firstName,
			LastName:               lastName,
			Telephone:              telephone,
			TransportationOfficeID: office.ID,
			Email:                  email,
			Active:                 true,
		}
		if user.ID != uuid.Nil {
			officeUser.UserID = &user.ID
		}

		verrs, err = appCtx.DB().ValidateAndSave(&officeUser)
		if err != nil {
			appCtx.Logger().Error("could not create office user", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating office user", zap.Stringer("errors", verrs))
		}
	case QaeOfficeUserType:
		// Now create the Truss JPPSO
		address := models.Address{
			StreetAddress1: "1333 Minna St",
			City:           "San Francisco",
			State:          "CA",
			PostalCode:     "94115",
			County:         "SAINT CLAIR",
		}

		verrs, err := appCtx.DB().ValidateAndSave(&address)
		if err != nil {
			appCtx.Logger().Error("could not create address", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating address", zap.Stringer("errors", verrs))
		}

		role := roles.Role{}
		err = appCtx.DB().Where("role_type = $1", roles.RoleTypeQae).First(&role)
		if err != nil {
			appCtx.Logger().Error("could not fetch role qae", zap.Error(err))
		}

		usersRole := models.UsersRoles{
			UserID: user.ID,
			RoleID: role.ID,
		}

		verrs, err = appCtx.DB().ValidateAndSave(&usersRole)
		if err != nil {
			appCtx.Logger().Error("could not create user role", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating user role", zap.Stringer("errors", verrs))
		}

		office := models.TransportationOffice{
			Name:      "Truss",
			AddressID: address.ID,
			Latitude:  37.7678355,
			Longitude: -122.4199298,
			Hours:     models.StringPointer("0900-1800 Mon-Sat"),
			Gbloc:     gbloc,
		}

		verrs, err = appCtx.DB().ValidateAndSave(&office)
		if err != nil {
			appCtx.Logger().Error("could not create office", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating office", zap.Stringer("errors", verrs))
		}

		officeUser := models.OfficeUser{
			FirstName:              firstName,
			LastName:               lastName,
			Telephone:              telephone,
			TransportationOfficeID: office.ID,
			Email:                  email,
			Active:                 true,
		}
		if user.ID != uuid.Nil {
			officeUser.UserID = &user.ID
		}

		verrs, err = appCtx.DB().ValidateAndSave(&officeUser)
		if err != nil {
			appCtx.Logger().Error("could not create office user", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating office user", zap.Stringer("errors", verrs))
		}
	case CustomerServiceRepresentativeOfficeUserType:
		address := models.Address{
			StreetAddress1: "1333 Minna St",
			City:           "San Francisco",
			State:          "CA",
			PostalCode:     "94115",
			County:         "SAINT CLAIR",
		}

		verrs, err := appCtx.DB().ValidateAndSave(&address)
		if err != nil {
			appCtx.Logger().Error("could not create address", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating address", zap.Stringer("errors", verrs))
		}

		role := roles.Role{}
		err = appCtx.DB().Where("role_type = $1", roles.RoleTypeCustomerServiceRepresentative).First(&role)
		if err != nil {
			appCtx.Logger().Error("could not fetch role customer_service_representative", zap.Error(err))
		}

		usersRole := models.UsersRoles{
			UserID: user.ID,
			RoleID: role.ID,
		}

		verrs, err = appCtx.DB().ValidateAndSave(&usersRole)
		if err != nil {
			appCtx.Logger().Error("could not create user role", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating user role", zap.Stringer("errors", verrs))
		}

		office := models.TransportationOffice{
			Name:      "Transcom",
			AddressID: address.ID,
			Latitude:  37.7678355,
			Longitude: -122.4199298,
			Hours:     models.StringPointer("0900-1800 Mon-Sat"),
			Gbloc:     gbloc,
		}

		verrs, err = appCtx.DB().ValidateAndSave(&office)
		if err != nil {
			appCtx.Logger().Error("could not create office", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating office", zap.Stringer("errors", verrs))
		}

		officeUser := models.OfficeUser{
			FirstName:              firstName,
			LastName:               lastName,
			Telephone:              telephone,
			TransportationOfficeID: office.ID,
			Email:                  email,
			Active:                 true,
		}
		if user.ID != uuid.Nil {
			officeUser.UserID = &user.ID
		}

		verrs, err = appCtx.DB().ValidateAndSave(&officeUser)
		if err != nil {
			appCtx.Logger().Error("could not create office user", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating office user", zap.Stringer("errors", verrs))
		}
	case MultiRoleOfficeUserType:
		// Now create the Truss JPPSO
		address := models.Address{
			StreetAddress1: "1333 Minna St",
			City:           "San Francisco",
			State:          "CA",
			PostalCode:     "94115",
			County:         "SAINT CLAIR",
		}

		verrs, err := appCtx.DB().ValidateAndSave(&address)
		if err != nil {
			appCtx.Logger().Error("could not create address", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating address", zap.Stringer("errors", verrs))
		}

		officeUserRoleTypes := []roles.RoleType{roles.RoleTypeTOO,
			roles.RoleTypeTIO,
			roles.RoleTypeServicesCounselor,
			roles.RoleTypePrimeSimulator,
		}
		var userRoles roles.Roles
		err = appCtx.DB().Where("role_type IN (?)", officeUserRoleTypes).All(&userRoles)
		if err != nil {
			appCtx.Logger().Error("could not fetch multi user roles", zap.Error(err))
		}

		for i := range userRoles {
			usersRole := models.UsersRoles{
				UserID: user.ID,
				RoleID: userRoles[i].ID,
			}

			verrs, err = appCtx.DB().ValidateAndSave(&usersRole)
			if err != nil {
				appCtx.Logger().Error("could not create user role", zap.Error(err))
			}
			if verrs.HasAny() {
				appCtx.Logger().Error("validation errors creating user role", zap.Stringer("errors", verrs))
			}
		}

		office := models.TransportationOffice{
			Name:      "Truss",
			AddressID: address.ID,
			Latitude:  37.7678355,
			Longitude: -122.4199298,
			Hours:     models.StringPointer("0900-1800 Mon-Sat"),
			Gbloc:     gbloc,
		}

		verrs, err = appCtx.DB().ValidateAndSave(&office)
		if err != nil {
			appCtx.Logger().Error("could not create office", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating office", zap.Stringer("errors", verrs))
		}

		officeUser := models.OfficeUser{
			FirstName:              firstName,
			LastName:               lastName,
			Telephone:              telephone,
			TransportationOfficeID: office.ID,
			Email:                  email,
			Active:                 true,
		}
		if user.ID != uuid.Nil {
			officeUser.UserID = &user.ID
		}

		verrs, err = appCtx.DB().ValidateAndSave(&officeUser)
		if err != nil {
			appCtx.Logger().Error("could not create office user", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating office user", zap.Stringer("errors", verrs))
		}
	case AdminUserType:

		adminUser := models.AdminUser{
			UserID:    &user.ID,
			Email:     user.OktaEmail,
			FirstName: "Leo",
			LastName:  "Spaceman",
			Role:      models.SystemAdminRole,
		}
		verrs, err := appCtx.DB().ValidateAndSave(&adminUser)

		if err != nil {
			appCtx.Logger().Error("could not create admin user", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("validation errors creating admin user", zap.Stringer("errors", verrs))
		}
	}

	return &user, userType
}

// createSession creates a new session for the user
func createSession(h devlocalAuthHandler, user *models.User, userType string, _ http.ResponseWriter, r *http.Request) (*auth.Session, error) {
	appCtx := h.AppContextFromRequest(r)
	// Preference any session already in the request context. Otherwise just create a new empty session.
	session := auth.SessionFromRequestContext(r)
	if session == nil {
		session = &auth.Session{}
	}

	lgUUID := user.OktaID
	userIdentity, err := models.FetchUserIdentity(appCtx.DB(), lgUUID)

	if err != nil {
		return nil, errors.Wrapf(err, "Unable to fetch user identity from OktaID %s", lgUUID)
	}

	session.Roles = append(session.Roles, userIdentity.Roles...)
	session.Permissions = getPermissionsForUser(appCtx, userIdentity.ID)

	// Assign user identity to session
	session.IDToken = "devlocal"
	session.UserID = userIdentity.ID
	session.Email = userIdentity.Email

	// Set the app
	active := userIdentity.Active

	// Keep the logic for redirection separate from setting the session user ids
	switch userType {
	case TOOOfficeUserType, TIOOfficeUserType, ServicesCounselorOfficeUserType, PrimeSimulatorOfficeUserType, QaeOfficeUserType, CustomerServiceRepresentativeOfficeUserType, MultiRoleOfficeUserType:
		session.ApplicationName = auth.OfficeApp
		session.Hostname = h.AppNames().OfficeServername
		active = userIdentity.Active || (userIdentity.OfficeActive != nil && *userIdentity.OfficeActive)
	case AdminUserType:
		session.ApplicationName = auth.AdminApp
		session.Hostname = h.AppNames().AdminServername
		session.AdminUserID = *userIdentity.AdminUserID
		session.AdminUserRole = userIdentity.AdminUserRole.String()
	default:
		session.ApplicationName = auth.MilApp
		session.Hostname = h.AppNames().MilServername
	}

	// If the user is active they should be denied a session
	if !active {
		appCtx.Logger().Error("Deactivated user requesting authentication",
			zap.String("application_name", string(session.ApplicationName)),
			zap.String("hostname", session.Hostname),
			zap.String("user_id", session.UserID.String()),
			zap.String("email", session.Email))
		return nil, errors.New("Deactivated user requesting authentication")
	}

	if session.IsMilApp() && userIdentity.ServiceMemberID != nil {
		session.ServiceMemberID = *(userIdentity.ServiceMemberID)
	}

	if userIdentity.OfficeUserID != nil && (session.IsOfficeApp() || isOfficeUser(userType)) {
		session.OfficeUserID = *(userIdentity.OfficeUserID)
	}

	session.FirstName = userIdentity.FirstName()
	session.LastName = userIdentity.LastName()
	session.Middle = userIdentity.Middle()

	sessionManager := h.SessionManagers().
		SessionManagerForApplication(session.ApplicationName)
	if sessionManager == nil {
		appCtx.Logger().Error("Cannot determine session manager for session", zap.String("session.ApplicationName", string(session.ApplicationName)))
		return nil, errors.New("Cannot determine session manager")
	}

	sessionManager.Put(r.Context(), "session", session)
	// Writing out the session cookie logs in the user
	appCtx.Logger().Info("logged in", zap.Any("session", session))
	return session, nil
}

// verifySessionWithApp returns an error if the user id for a specific app is not available
func verifySessionWithApp(session *auth.Session) error {

	// TODO: Should this be a check that we do? Or will all office users also be service members?
	// if (session.ServiceMemberID == uuid.UUID{}) && session.IsMilApp() {
	// 	return errors.Errorf("Non-service member user %s authenticated at service member site", session.Email)
	// }

	if (session.OfficeUserID == uuid.UUID{}) && session.IsOfficeApp() {
		return errors.Errorf("Non-office user %s authenticated at office site", session.Email)
	}

	if !session.IsAdminUser() && session.IsAdminApp() {
		return errors.Errorf("Non-admin user %s authenticated at admin site", session.Email)
	}

	return nil
}

// loginUser creates a session for the user and verifies the session against the app
func loginUser(h devlocalAuthHandler, user *models.User, userType string, w http.ResponseWriter, r *http.Request) (*auth.Session, error) {
	appCtx := h.AppContextFromRequest(r)
	session, err := createSession(devlocalAuthHandler(h), user, userType, w, r)
	if err != nil {
		appCtx.Logger().Error("Could not create session", zap.Error(err))
		http.Error(w, http.StatusText(403), http.StatusForbidden)
		return nil, err
	}

	err = verifySessionWithApp(session)
	if err != nil {
		appCtx.Logger().Error("User unauthorized", zap.Error(err))
		http.Error(w, http.StatusText(401), http.StatusUnauthorized)
		return nil, err
	}
	return session, nil
}

func isOfficeUser(userType string) bool {
	if userType == TOOOfficeUserType || userType == TIOOfficeUserType || userType == ServicesCounselorOfficeUserType || userType == QaeOfficeUserType || userType == CustomerServiceRepresentativeOfficeUserType {
		return true
	}
	return false
}
