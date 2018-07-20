package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"strings"
)

// User is an entity with a registered uuid and email at login.gov
type User struct {
	ID            uuid.UUID `json:"id" db:"id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	LoginGovUUID  uuid.UUID `json:"login_gov_uuid" db:"login_gov_uuid"`
	LoginGovEmail string    `json:"login_gov_email" db:"login_gov_email"`
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
	Email                  string     `db:"email"`
	ServiceMemberID        *uuid.UUID `db:"sm_id"`
	ServiceMemberFirstName *string    `db:"sm_fname"`
	ServiceMemberLastName  *string    `db:"sm_lname"`
	ServiceMemberMiddle    *string    `db:"sm_middle"`
	OfficeUserID           *uuid.UUID `db:"ou_id"`
	OfficeUserFirstName    *string    `db:"ou_fname"`
	OfficeUserLastName     *string    `db:"ou_lname"`
	OfficeUserMiddle       *string    `db:"ou_middle"`
	OfficeUserEntityID     *uuid.UUID `db:"ou_entity_id"`
	TspUserID              *uuid.UUID `db:"tu_id"`
	TspUserFirstName       *string    `db:"tu_fname"`
	TspUserLastName        *string    `db:"tu_lname"`
	TspUserMiddle          *string    `db:"tu_middle"`
	TspUserEntityID        *uuid.UUID `db:"tu_entity_id"`
}

// FetchUserIdentity queries the database for information about the logged in user
func FetchUserIdentity(db *pop.Connection, loginGovID string) (*UserIdentity, error) {
	var identities []UserIdentity
	query := `SELECT users.id,
				users.login_gov_email as email,
				sm.id as sm_id,
				sm.first_name as sm_fname,
				sm.last_name as sm_lname,
				sm.middle_name as sm_middle,
				ou.id as ou_id,
				ou.first_name as ou_fname,
				ou.last_name as ou_lname,
				ou.middle_initials as ou_middle,
				ou.transportation_office_id as ou_entity_id,
				tu.id as tu_id,
				tu.first_name as tu_fname,
				tu.last_name as tu_lname,
				tu.middle_initials as tu_middle,
				tu.transportation_service_provider_id as tu_entity_id
			FROM users
			LEFT OUTER JOIN service_members as sm on sm.user_id = users.id
			LEFT OUTER JOIN office_users as ou on ou.user_id = users.id
			LEFT OUTER JOIN tsp_users as tu on tu.user_id = users.id
			WHERE users.login_gov_uuid  = $1`
	err := db.RawQuery(query, loginGovID).All(&identities)
	if err != nil {
		return nil, err
	} else if len(identities) == 0 {
		return nil, ErrFetchNotFound
	}
	return &identities[0], nil
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
