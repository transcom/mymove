package models

import (
	"time"

	"github.com/transcom/mymove/pkg/models/roles"

	"strings"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
)

// User is an entity with a registered uuid and email at login.gov
type User struct {
	ID                     uuid.UUID   `json:"id" db:"id"`
	CreatedAt              time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time   `json:"updated_at" db:"updated_at"`
	LoginGovUUID           *uuid.UUID  `json:"login_gov_uuid" db:"login_gov_uuid"`
	LoginGovEmail          string      `json:"login_gov_email" db:"login_gov_email"`
	Active                 bool        `json:"active" db:"active"`
	Roles                  roles.Roles `many_to_many:"users_roles"`
	CurrentAdminSessionID  string      `json:"current_admin_session_id" db:"current_admin_session_id"`
	CurrentOfficeSessionID string      `json:"current_office_session_id" db:"current_office_session_id"`
	CurrentMilSessionID    string      `json:"current_mil_session_id" db:"current_mil_session_id"`
}

// Users is not required by pop and may be deleted
type Users []User

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: u.LoginGovEmail, Name: "LoginGovEmail"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// GetUser loads the associated User from the DB using the user ID
func GetUser(db *pop.Connection, userID uuid.UUID) (*User, error) {
	var user User
	err := db.Find(&user, userID.String())
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to find user by id %s", userID.String())
	}
	return &user, nil
}

// GetUserFromEmail loads the associated User from the DB using the user email
func GetUserFromEmail(db *pop.Connection, email string) (*User, error) {
	users := []User{}
	downcasedEmail := strings.ToLower(email)
	err := db.Where("login_gov_email = $1", downcasedEmail).All(&users)
	if len(users) == 0 {
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to find user by email %s", downcasedEmail)
		}
		return nil, errors.Errorf("Unable to find user by email %s", downcasedEmail)
	}
	return &users[0], err
}

// CreateUser is called upon successful login.gov verification of a new user
func CreateUser(db *pop.Connection, loginGovID string, email string) (*User, error) {
	lgu, err := uuid.FromString(loginGovID)
	if err != nil {
		return nil, err
	}
	newUser := User{
		LoginGovUUID:  &lgu,
		LoginGovEmail: strings.ToLower(email),
		Active:        true,
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

// UpdateUserLoginGovUUID is called upon the first successful login.gov verification of a new user
func UpdateUserLoginGovUUID(db *pop.Connection, user *User, loginGovID string) error {
	lgu, err := uuid.FromString(loginGovID)
	if err != nil {
		return err
	}

	user.LoginGovUUID = &lgu

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
	ID                     uuid.UUID   `db:"id"`
	Active                 bool        `db:"active"`
	Email                  string      `db:"email"`
	ServiceMemberID        *uuid.UUID  `db:"sm_id"`
	ServiceMemberFirstName *string     `db:"sm_fname"`
	ServiceMemberLastName  *string     `db:"sm_lname"`
	ServiceMemberMiddle    *string     `db:"sm_middle"`
	OfficeUserID           *uuid.UUID  `db:"ou_id"`
	OfficeUserFirstName    *string     `db:"ou_fname"`
	OfficeUserLastName     *string     `db:"ou_lname"`
	OfficeUserMiddle       *string     `db:"ou_middle"`
	OfficeActive           *bool       `db:"ou_active"`
	AdminUserID            *uuid.UUID  `db:"au_id"`
	AdminUserRole          *AdminRole  `db:"au_role"`
	AdminUserFirstName     *string     `db:"au_fname"`
	AdminUserLastName      *string     `db:"au_lname"`
	AdminUserActive        *bool       `db:"au_active"`
	Roles                  roles.Roles `many_to_many:"users_roles" primary_id:"user_id"`
}

// FetchUserIdentity queries the database for information about the logged in user
func FetchUserIdentity(db *pop.Connection, loginGovID string) (*UserIdentity, error) {
	var identities []UserIdentity
	query := `SELECT users.id,
				users.login_gov_email AS email,
				users.active AS active,
				sm.id AS sm_id,
				sm.first_name AS sm_fname,
				sm.last_name AS sm_lname,
				sm.middle_name AS sm_middle,
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
			WHERE users.login_gov_uuid  = $1`
	err := db.RawQuery(query, loginGovID).All(&identities)
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
				users.login_gov_email AS email,
				users.active AS active,
				ou.id AS ou_id,
				ou.first_name AS ou_fname,
				ou.last_name AS ou_lname,
				ou.middle_initials AS ou_middle
			FROM office_users as ou
			JOIN users on ou.user_id = users.id
			ORDER BY users.created_at LIMIT $1`
	case auth.AdminApp:
		query = `SELECT users.id,
				users.login_gov_email AS email,
				users.active AS active,
				au.id AS au_id,
				au.role AS au_role,
				au.first_name AS au_fname,
				au.last_name AS au_lname
			FROM admin_users as au
			JOIN users on au.user_id = users.id
			ORDER BY users.created_at LIMIT $1`
	default:
		query = `SELECT users.id,
				users.login_gov_email AS email,
				users.active AS active,
				sm.id AS sm_id,
				sm.first_name AS sm_fname,
				sm.last_name AS sm_lname,
				sm.middle_name AS sm_middle
			FROM service_members as sm
			JOIN users on sm.user_id = users.id
			WHERE users.login_gov_email != 'first.last@login.gov.test'
			ORDER BY users.created_at LIMIT $1`
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
