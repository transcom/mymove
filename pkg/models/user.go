package models

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/markbates/goth"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth/context"
)

// User is an entity with a registered uuid and email at login.gov
type User struct {
	ID            uuid.UUID `db:"id"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	LoginGovUUID  uuid.UUID `db:"login_gov_uuid"`
	LoginGovEmail string    `db:"login_gov_email"`
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

// GetUserFromRequest extracts the user model from the request context's user ID
func GetUserFromRequest(db *pop.Connection, r *http.Request) (user User, err error) {
	userID, ok := context.GetUserID(r.Context())
	if !ok {
		err = errors.New("Failed to fetch user_id from context")
		return
	}

	user, err = GetUserByID(db, userID)
	if err != nil {
		return
	}

	return user, err
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
