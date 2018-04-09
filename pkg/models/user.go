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
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// User is an entity with a registered uuid and email at login.gov
type User struct {
	ID            uuid.UUID                 `json:"id" db:"id"`
	CreatedAt     time.Time                 `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time                 `json:"updated_at" db:"updated_at"`
	LoginGovUUID  uuid.UUID                 `json:"login_gov_uuid" db:"login_gov_uuid"`
	LoginGovEmail string                    `json:"login_gov_email" db:"login_gov_email"`
	Type          internalmessages.UserType `json:"type" db:"type"`
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
	validTypes := []string{string(internalmessages.UserTypeUNKNOWN), string(internalmessages.UserTypeSERVICEMEMBER), string(internalmessages.UserTypeTRUSTEDAGENT)}
	return validate.Validate(
		&validators.UUIDIsPresent{Field: u.LoginGovUUID, Name: "LoginGovUUID"},
		&validators.StringIsPresent{Field: u.LoginGovEmail, Name: "LoginGovEmail"},
		&validators.StringInclusion{Field: string(u.Type), Name: "Type", List: validTypes},
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

// GetServiceMemberProfile returns a service member profile if one is associated with this user, otherwise returns nil
func (u User) GetServiceMemberProfile(db *pop.Connection) (*ServiceMember, error) {
	if u.Type != internalmessages.UserTypeSERVICEMEMBER {
		return nil, nil
	}

	serviceMembers := ServiceMembers{}
	err := db.Where("user_id = $1", u.ID).All(&serviceMembers)
	if err != nil {
		return nil, err
	}

	if len(serviceMembers) > 1 {
		return nil, errors.New("Should not ever have more than one service member profile for a user")
	}

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
			Type:          internalmessages.UserTypeUNKNOWN,
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
