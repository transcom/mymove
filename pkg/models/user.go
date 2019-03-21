package models

import (
	"time"

	"strings"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// User is an entity with a registered uuid and email at login.gov
type User struct {
	ID            uuid.UUID `json:"id" db:"id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	LoginGovUUID  uuid.UUID `json:"login_gov_uuid" db:"login_gov_uuid"`
	LoginGovEmail string    `json:"login_gov_email" db:"login_gov_email"`
	Disabled      bool      `json:"disabled" db:"disabled"`
}

// Users is not required by pop and may be deleted
type Users []User

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: u.LoginGovUUID, Name: "LoginGovUUID"},
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

// GetUser loads the associated User from the DB
func GetUser(db *pop.Connection, userID uuid.UUID) (*User, error) {
	var user User
	err := db.Find(&user, userID)
	return &user, err
}

// CreateUser is called upon successful login.gov verification of a new user
func CreateUser(db *pop.Connection, loginGovID string, email string) (*User, error) {
	lgu, err := uuid.FromString(loginGovID)
	if err != nil {
		return nil, err
	}
	newUser := User{
		LoginGovUUID:  lgu,
		LoginGovEmail: strings.ToLower(email),
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

// UserIdentity is summary of the information about a user from the database
type UserIdentity struct {
	ID                     uuid.UUID  `db:"id"`
	Disabled               bool       `db:"disabled"`
	Email                  string     `db:"email"`
	ServiceMemberID        *uuid.UUID `db:"sm_id"`
	ServiceMemberFirstName *string    `db:"sm_fname"`
	ServiceMemberLastName  *string    `db:"sm_lname"`
	ServiceMemberMiddle    *string    `db:"sm_middle"`
	OfficeUserID           *uuid.UUID `db:"ou_id"`
	OfficeUserFirstName    *string    `db:"ou_fname"`
	OfficeUserLastName     *string    `db:"ou_lname"`
	OfficeUserMiddle       *string    `db:"ou_middle"`
	TspUserID              *uuid.UUID `db:"tu_id"`
	TspUserFirstName       *string    `db:"tu_fname"`
	TspUserLastName        *string    `db:"tu_lname"`
	TspUserMiddle          *string    `db:"tu_middle"`
	DpsUserID              *uuid.UUID `db:"du_id"`
}

// FetchUserIdentity queries the database for information about the logged in user
func FetchUserIdentity(db *pop.Connection, loginGovID string) (*UserIdentity, error) {
	var identities []UserIdentity
	query := `SELECT users.id,
				users.login_gov_email AS email,
				users.disabled AS disabled,
				sm.id AS sm_id,
				sm.first_name AS sm_fname,
				sm.last_name AS sm_lname,
				sm.middle_name AS sm_middle,
				ou.id AS ou_id,
				ou.first_name AS ou_fname,
				ou.last_name AS ou_lname,
				ou.middle_initials AS ou_middle,
				tu.id AS tu_id,
				tu.first_name AS tu_fname,
				tu.last_name AS tu_lname,
				tu.middle_initials AS tu_middle,
				du.id AS du_id
			FROM users
			LEFT OUTER JOIN service_members AS sm on sm.user_id = users.id
			LEFT OUTER JOIN office_users AS ou on ou.user_id = users.id
			LEFT OUTER JOIN tsp_users AS tu on tu.user_id = users.id
			LEFT OUTER JOIN dps_users AS du on du.login_gov_email = users.login_gov_email
			WHERE users.login_gov_uuid  = $1`
	err := db.RawQuery(query, loginGovID).All(&identities)
	if err != nil {
		return nil, err
	} else if len(identities) == 0 {
		return nil, ErrFetchNotFound
	}
	return &identities[0], nil
}

// FetchAllUserIdentities returns information for all users in the db
func FetchAllUserIdentities(db *pop.Connection) ([]UserIdentity, error) {
	var identities []UserIdentity
	query := `SELECT users.id,
				users.login_gov_email AS email,
				users.disabled AS disabled,
				sm.id AS sm_id,
				sm.first_name AS sm_fname,
				sm.last_name AS sm_lname,
				sm.middle_name AS sm_middle,
				ou.id AS ou_id,
				ou.first_name AS ou_fname,
				ou.last_name AS ou_lname,
				ou.middle_initials AS ou_middle,
				tu.id AS tu_id,
				tu.first_name AS tu_fname,
				tu.last_name AS tu_lname,
				tu.middle_initials AS tu_middle,
				du.id AS du_id
			FROM users
			LEFT OUTER JOIN service_members AS sm on sm.user_id = users.id
			LEFT OUTER JOIN office_users AS ou on ou.user_id = users.id
			LEFT OUTER JOIN tsp_users AS tu on tu.user_id = users.id
			LEFT OUTER JOIN dps_users AS du on du.login_gov_email = users.login_gov_email
			ORDER BY users.created_at`

	err := db.RawQuery(query).All(&identities)
	if err != nil {
		return nil, err
	}
	return identities, nil
}

// firstValue returns the first string value that is not nil
func firstValue(one *string, two *string, three *string) (value string) {

	if one != nil {
		value = *one
	} else if two != nil {
		value = *two
	} else if three != nil {
		value = *three
	}
	return
}

// FirstName gets the firstname of the user from either the ServiceMember or OfficeUser or TspUser identity
func (ui *UserIdentity) FirstName() string {
	return firstValue(ui.ServiceMemberFirstName, ui.OfficeUserFirstName, ui.TspUserFirstName)
}

// LastName gets the firstname of the user from either the ServiceMember or OfficeUser or TspUser identity
func (ui *UserIdentity) LastName() string {
	return firstValue(ui.ServiceMemberLastName, ui.OfficeUserLastName, ui.TspUserLastName)
}

// Middle gets the MiddleName or Initials from the ServiceMember or OfficeUser or TspUser Identity
func (ui *UserIdentity) Middle() string {
	return firstValue(ui.ServiceMemberMiddle, ui.OfficeUserMiddle, ui.TspUserMiddle)
}
