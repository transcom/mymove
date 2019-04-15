package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// DpsUser are users who have permission to access MyMove - DPS integration resources
type DpsUser struct {
	ID            uuid.UUID `json:"id" db:"id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	LoginGovEmail string    `json:"login_gov_email" db:"login_gov_email"`
}

// String is not required by pop and may be deleted
func (d DpsUser) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// DpsUsers is not required by pop and may be deleted
type DpsUsers []DpsUser

// String is not required by pop and may be deleted
func (d DpsUsers) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (d *DpsUser) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: d.LoginGovEmail, Name: "LoginGovEmail"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (d *DpsUser) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (d *DpsUser) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// IsDPSUser checks if a user is a DPS user given email
func IsDPSUser(db *pop.Connection, email string) (bool, error) {
	count, err := db.Q().Where("LOWER(login_gov_email) = ?", strings.ToLower(email)).Count(DpsUser{})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FetchDPSUserByEmail looks for an DPS user with a specific email
func FetchDPSUserByEmail(tx *pop.Connection, email string) (*DpsUser, error) {
	var users DpsUsers
	err := tx.Where("LOWER(login_gov_email) = $1", strings.ToLower(email)).All(&users)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, ErrFetchNotFound
	}
	return &users[0], nil
}

// FetchDPSUserByID fetches an DPS user by ID
func FetchDPSUserByID(tx *pop.Connection, id uuid.UUID) (*DpsUser, error) {
	var user DpsUser
	err := tx.Find(&user, id)
	return &user, err
}
