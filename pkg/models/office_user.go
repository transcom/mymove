package models

import (
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// OfficeUser is someone who works in one of the TransportationOffices
type OfficeUser struct {
	ID                     uuid.UUID            `json:"id" db:"id"`
	UserID                 *uuid.UUID           `json:"user_id" db:"user_id"`
	User                   User                 `belongs_to:"user"`
	LastName               string               `json:"last_name" db:"last_name"`
	FirstName              string               `json:"first_name" db:"first_name"`
	MiddleInitials         *string              `json:"middle_initials" db:"middle_initials"`
	Email                  string               `json:"email" db:"email"`
	Telephone              string               `json:"telephone" db:"telephone"`
	TransportationOfficeID uuid.UUID            `json:"transportation_office_id" db:"transportation_office_id"`
	TransportationOffice   TransportationOffice `belongs_to:"transportation_office"`
	CreatedAt              time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time            `json:"updated_at" db:"updated_at"`
}

// OfficeUsers is not required by pop and may be deleted
type OfficeUsers []OfficeUser

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (o *OfficeUser) Validate(tx *pop.Connection) (*validate.Errors, error) {

	return validate.Validate(
		&validators.StringIsPresent{Field: o.LastName, Name: "LastName"},
		&validators.StringIsPresent{Field: o.FirstName, Name: "FirstName"},
		&validators.StringIsPresent{Field: o.Email, Name: "Email"},
		&validators.StringIsPresent{Field: o.Telephone, Name: "Telephone"},
		&validators.UUIDIsPresent{Field: o.TransportationOfficeID, Name: "TransportationOfficeID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (o *OfficeUser) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (o *OfficeUser) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// FetchOfficeUserByEmail looks for an office user with a specific email
func FetchOfficeUserByEmail(tx *pop.Connection, email string) (*OfficeUser, error) {
	var users OfficeUsers
	err := tx.Where("LOWER(email) = $1", strings.ToLower(email)).All(&users)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, ErrFetchNotFound
	}
	return &users[0], nil
}

// FetchOfficeUserByID fetches an office user by ID
func FetchOfficeUserByID(tx *pop.Connection, id uuid.UUID) (*OfficeUser, error) {
	var user OfficeUser
	err := tx.Find(&user, id)
	return &user, err
}
