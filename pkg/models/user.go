package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/markbates/goth"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/auth"
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

// GetFullServiceMemberProfile returns a service member profile if one is associated with this user, otherwise returns nil
func GetFullServiceMemberProfile(db *pop.Connection, session *auth.Session) (*ServiceMember, error) {
	serviceMembers := ServiceMembers{}
	err := db.Where("id = $1", session.ServiceMemberID).Eager("DutyStation", "ResidentialAddress", "BackupMailingAddress", "Orders.NewDutyStation.Address", "Orders.UploadedOrders", "Orders.Moves.PersonallyProcuredMoves", "User").All(&serviceMembers)
	if err != nil {
		return nil, err
	}

	// There can only ever be one service_member for a given user
	if len(serviceMembers) == 1 {
		return &serviceMembers[0], nil
	}

	return nil, nil

}

// GetOrCreateUser is called upon successful login.gov verification
func GetOrCreateUser(db *pop.Connection, gothUser goth.User) (*User, error) {

	// Check if user already exists
	query := db.Where("login_gov_uuid = $1", gothUser.UserID)
	var user User
	err := query.First(&user)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return nil, errors.Wrap(err, "Failed to load user")
		}
		// No user found, creating new user
		loginGovUUID, _ := uuid.FromString(gothUser.UserID)
		newUser := User{
			LoginGovUUID:  loginGovUUID,
			LoginGovEmail: gothUser.Email,
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
	// Return found user
	return &user, nil
}

// UserIdentity is summary of the information about a user from the database
type UserIdentity struct {
	ID                     uuid.UUID `db:"id"`
	ServiceMemberID        uuid.UUID `db:"sm_id"`
	ServiceMemberFirstName string    `db:"sm_fname"`
	ServiceMemberLastName  string    `db:"sm_lname"`
	ServiceMemberMiddle    string    `db:"sm_middle"`
	OfficeUserID           uuid.UUID `db:"ou_id"`
	OfficeUserFirstName    string    `db:"ou_fname"`
	OfficeUserLastName     string    `db:"ou_lname"`
	OfficeUserMiddle       string    `db:"ou_middle"`
}

// FetchUserIdentity queries the database for information about the logged in user
func FetchUserIdentity(db *pop.Connection, userID uuid.UUID) (*UserIdentity, error) {
	var identities []UserIdentity
	query := `SELECT users.id,
				sm.id as sm_id,
				sm.first_name as sm_fname,
				sm.last_name as sm_lname,
				sm.middle_name as sm_middle,
				ou.id as ou_id,
				ou.first_name as ou_fname,
				ou.last_name as ou_lname,
				ou.middle_initials as ou_middle
			FROM users
			OUTER JOIN service_members as sm on sm.user_id = users.id
			OUTER JOIN office_users as ou on ou.user_id = users.id
			WHERE users.id  = $1`
	err := db.RawQuery(query, userID).All(&identities)
	if err != nil {
		return nil, err
	} else if len(identities) == 0 {
		return nil, ErrFetchNotFound
	}
	return &identities[0], nil
}

func firstNotEmpty(one string, two string) string {
	if one != "" {
		return one
	}
	return two
}

// FirstName gets the firstname of the user from either the ServiceMember or OfficeUser identity
func (ui *UserIdentity) FirstName() string {
	return firstNotEmpty(ui.ServiceMemberFirstName, ui.OfficeUserFirstName)
}

// LastName gets the firstname of the user from either the ServiceMember or OfficeUser identity
func (ui *UserIdentity) LastName() string {
	return firstNotEmpty(ui.ServiceMemberLastName, ui.OfficeUserLastName)
}

// Middle gets the MiddleName or Initials from the ServiceMember or OfficeUserIdentity
func (ui *UserIdentity) Middle() string {
	return firstNotEmpty(ui.ServiceMemberMiddle, ui.OfficeUserMiddle)
}
