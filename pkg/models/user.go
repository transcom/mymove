package models

import (
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models/roles"
)

// User is an entity with a registered profile ID and email in Okta
type User struct {
	ID                     uuid.UUID   `json:"id" db:"id"`
	CreatedAt              time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time   `json:"updated_at" db:"updated_at"`
	OktaID                 string      `json:"okta_id" db:"okta_id"`
	OktaEmail              string      `json:"okta_email" db:"okta_email"`
	Active                 bool        `json:"active" db:"active"`
	Roles                  roles.Roles `many_to_many:"users_roles"`
	Privileges             Privileges  `many_to_many:"users_privileges"`
	CurrentAdminSessionID  string      `json:"current_admin_session_id" db:"current_admin_session_id"`
	CurrentOfficeSessionID string      `json:"current_office_session_id" db:"current_office_session_id"`
	CurrentMilSessionID    string      `json:"current_mil_session_id" db:"current_mil_session_id"`
}

// TableName overrides the table name used by Pop.
func (u User) TableName() string {
	return "users"
}

type Users []User

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (u *User) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: u.OktaEmail, Name: "OktaEmail"},
	), nil
}

// GetUser loads the associated User from the DB using the user ID
func GetUser(db *pop.Connection, userID uuid.UUID) (*User, error) {
	var user User
	err := db.Find(&user, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to find user by id %s", userID)
	}
	return &user, nil
}

// GetUserFromEmail loads the associated User from the DB using the user email
func GetUserFromEmail(db *pop.Connection, email string) (*User, error) {
	users := []User{}
	downcasedEmail := strings.ToLower(email)
	err := db.Where("okta_email = $1", downcasedEmail).All(&users)
	if len(users) == 0 {
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to find user by email %s", downcasedEmail)
		}
		return nil, errors.Errorf("Unable to find user by email %s", downcasedEmail)
	}
	return &users[0], err
}

// CreateUser is called upon successful Okta verification of a new user
func CreateUser(db *pop.Connection, oktaID string, email string) (*User, error) {

	newUser := User{
		OktaID:    oktaID,
		OktaEmail: strings.ToLower(email),
		Active:    true,
	}
	verrs, err := db.ValidateAndCreate(&newUser)
	if verrs.HasAny() {
		return nil, verrs
	} else if err != nil {
		err = errors.Wrap(err, "Unable to create user")
		return nil, err
	}
	return &newUser, nil
}

// UpdateUserOktaID is called upon the first successful Okta verification of a new user
func UpdateUserOktaID(db *pop.Connection, user *User, oktaID string) error {

	user.OktaID = oktaID

	verrs, err := db.ValidateAndUpdate(user)
	if verrs.HasAny() {
		return verrs
	} else if err != nil {
		err = errors.Wrap(err, "Unable to update user")
		return err
	}

	return nil
}

// UserIdentity is summary of the information about a user from the database
type UserIdentity struct {
	ID                              uuid.UUID                       `db:"id"`
	Active                          bool                            `db:"active"`
	Email                           string                          `db:"email"`
	ServiceMemberID                 *uuid.UUID                      `db:"sm_id"`
	ServiceMemberFirstName          *string                         `db:"sm_fname"`
	ServiceMemberLastName           *string                         `db:"sm_lname"`
	ServiceMemberMiddle             *string                         `db:"sm_middle"`
	ServiceMemberCacValidated       *bool                           `db:"sm_cac_validated"`
	OfficeUserID                    *uuid.UUID                      `db:"ou_id"`
	OfficeUserFirstName             *string                         `db:"ou_fname"`
	OfficeUserLastName              *string                         `db:"ou_lname"`
	OfficeUserMiddle                *string                         `db:"ou_middle"`
	OfficeActive                    *bool                           `db:"ou_active"`
	AdminUserID                     *uuid.UUID                      `db:"au_id"`
	AdminUserRole                   *AdminRole                      `db:"au_role"`
	AdminUserFirstName              *string                         `db:"au_fname"`
	AdminUserLastName               *string                         `db:"au_lname"`
	AdminUserActive                 *bool                           `db:"au_active"`
	Roles                           roles.Roles                     `many_to_many:"users_roles" primary_id:"user_id"`
	Privileges                      Privileges                      `many_to_many:"users_privileges" primary_id:"user_id"`
	TransportationOfficeAssignments TransportationOfficeAssignments `many_to_many:"transportation_office_assignmentss" primary_id:"id"`
}

// FetchUserIdentity queries the database for information about the logged in user
func FetchUserIdentity(db *pop.Connection, oktaID string) (*UserIdentity, error) {
	var identities []UserIdentity
	query := `SELECT users.id,
				users.okta_email AS email,
				users.active AS active,
				sm.id AS sm_id,
				sm.first_name AS sm_fname,
				sm.last_name AS sm_lname,
				sm.middle_name AS sm_middle,
				sm.cac_validated AS sm_cac_validated,
				ou.id AS ou_id,
				ou.first_name AS ou_fname,
				ou.last_name AS ou_lname,
				ou.middle_initials AS ou_middle,
				ou.active AS ou_active,
				au.id AS au_id,
				au.role AS au_role,
				au.first_name AS au_fname,
				au.last_name AS au_lname,
				au.active AS au_active
			FROM users
			LEFT OUTER JOIN service_members AS sm on sm.user_id = users.id
			LEFT OUTER JOIN office_users AS ou on ou.user_id = users.id
			LEFT OUTER JOIN admin_users AS au on au.user_id = users.id
			WHERE users.okta_id  = $1`
	err := db.RawQuery(query, oktaID).All(&identities)
	if err != nil {
		return nil, err
	} else if len(identities) == 0 {
		return nil, ErrFetchNotFound
	}
	identity := &identities[0]
	roleError := db.RawQuery(`SELECT * FROM roles
									WHERE id in (select role_id from users_roles
										where deleted_at is null and user_id = ?)`, identity.ID).All(&identity.Roles)
	if roleError != nil {
		return nil, roleError
	}
	privilegeError := db.RawQuery(`SELECT * FROM privileges
									WHERE id in (select privilege_id from users_privileges
										where deleted_at is null and user_id = ?)`, identity.ID).All(&identity.Privileges)
	if privilegeError != nil {
		return nil, privilegeError
	}
	transportationOfficeAssignmentError := db.EagerPreload("TransportationOffice").
		Join("transportation_offices", "transportation_office_assignments.transportation_office_id = transportation_offices.id").
		Where("transportation_office_assignments.id = ?", identity.OfficeUserID).All(&identity.TransportationOfficeAssignments)
	if transportationOfficeAssignmentError != nil {
		return nil, transportationOfficeAssignmentError
	}
	return identity, nil
}

// FetchAppUserIdentities returns a limited set of user records based on application
func FetchAppUserIdentities(db *pop.Connection, appname auth.Application, limit int) ([]UserIdentity, error) {
	var identities []UserIdentity

	var query string
	switch appname {
	case auth.OfficeApp:
		query = `SELECT
		        users.id,
				users.okta_email AS email,
				users.active AS active,
				ou.id AS ou_id,
				ou.first_name AS ou_fname,
				ou.last_name AS ou_lname,
				ou.middle_initials AS ou_middle
			FROM office_users as ou
			JOIN users on ou.user_id = users.id
			ORDER BY users.created_at DESC LIMIT $1`
	case auth.AdminApp:
		query = `SELECT users.id,
				users.okta_email AS email,
				users.active AS active,
				au.id AS au_id,
				au.role AS au_role,
				au.first_name AS au_fname,
				au.last_name AS au_lname
			FROM admin_users as au
			JOIN users on au.user_id = users.id
			ORDER BY users.created_at DESC LIMIT $1`
	default:
		query = `SELECT users.id,
				users.okta_email AS email,
				users.active AS active,
				sm.id AS sm_id,
				sm.first_name AS sm_fname,
				sm.last_name AS sm_lname,
				sm.middle_name AS sm_middle
			FROM service_members as sm
			JOIN users on sm.user_id = users.id
			WHERE users.okta_email != 'first.last@okta.mil'
			ORDER BY users.created_at DESC LIMIT $1`
	}

	err := db.RawQuery(query, limit).All(&identities)
	if err != nil {
		return nil, err
	}
	return identities, nil
}

// firstValue returns the first string value that is not nil
func firstValue(vals ...*string) string {
	for _, val := range vals {
		if val != nil {
			return *val
		}
	}
	return ""
}

// FirstName gets the firstname of the user from either the ServiceMember or OfficeUser identity
func (ui *UserIdentity) FirstName() string {
	return firstValue(ui.ServiceMemberFirstName, ui.OfficeUserFirstName, ui.AdminUserFirstName)
}

// LastName gets the firstname of the user from either the ServiceMember or OfficeUser or TspUser identity
func (ui *UserIdentity) LastName() string {
	return firstValue(ui.ServiceMemberLastName, ui.OfficeUserLastName, ui.AdminUserLastName)
}

// Middle gets the MiddleName or Initials from the ServiceMember or OfficeUser or TspUser Identity
func (ui *UserIdentity) Middle() string {
	return firstValue(ui.ServiceMemberMiddle, ui.OfficeUserMiddle)
}
