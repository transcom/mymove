package models

import (
	"encoding/json"
	"github.com/markbates/goth"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"time"
)

// User is an entity with a registered uuid and email at login.gov
type User struct {
	ID            uuid.UUID `json:"id" db:"id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	LoginGovUUID  uuid.UUID `json:"login_gov_uuid" db:"login_gov_uuid"`
	LoginGovEmail string    `json:"login_gov_email" db:"login_gov_email"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

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

// GetOrCreateUser is called upon successful login.gov verification
func GetOrCreateUser(db *pop.Connection, gothUser goth.User) (User, error) {

	// Check if user already exists
	query := db.Where("login_gov_uuid = $1", gothUser.UserID)
	var users []User
	err := query.All(&users)
	if err != nil {
		err = errors.Wrap(err, "DB Query Error")
		return (User{}), err
	}

	// If user is not in DB, create it
	if len(users) == 0 {
		loginGovUUID, _ := uuid.FromString(gothUser.UserID)
		newUser := User{
			LoginGovUUID:  loginGovUUID,
			LoginGovEmail: gothUser.Email,
		}
		verrs, err := db.ValidateAndCreate(&newUser)
		if verrs.HasAny() {
			return (User{}), verrs
		} else if err != nil {
			err = errors.Wrap(err, "Unable to create user")
			return (User{}), err
		}
		// Create new move for user
		newMove := Move{
			UserID: newUser.ID,
		}
		moveVerrs, moveErr := db.ValidateAndCreate(&newMove)
		if moveVerrs.HasAny() {
			return (User{}), moveVerrs
		} else if moveErr != nil {
			moveErr = errors.Wrap(err, "Unable to create move")
			return (User{}), moveErr
		}
		return newUser, nil
	}
	// one user was found, return it
	return users[0], nil
}
