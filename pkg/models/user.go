package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/markbates/goth"
	"github.com/pkg/errors"
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

// GetUserByID fetches a user model by their database ID
func GetUserByID(db *pop.Connection, id uuid.UUID) (User, error) {
	user := User{}
	err := db.Find(&user, id)
	return user, err
}

// GetServiceMemberProfile returns a service member profile if one is associated with this user, otherwise returns nil
func (u User) GetServiceMemberProfile(db *pop.Connection) (*ServiceMember, error) {
	serviceMembers := ServiceMembers{}
	err := db.Where("user_id = $1", u.ID).Eager("DutyStation", "ResidentialAddress", "BackupMailingAddress", "Orders.NewDutyStation.Address", "Orders.UploadedOrders.Uploads", "Orders.Moves.PersonallyProcuredMoves", "BackupContacts").All(&serviceMembers)
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
