package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"time"
)

// OfficeUser is someone who works in one of the TransportationOffices
type OfficeUser struct {
	ID                     uuid.UUID            `json:"id" db:"id"`
	UserID                 uuid.UUID            `json:"user_id" db:"user_id"`
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

// String is not required by pop and may be deleted
func (o OfficeUser) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

// OfficeUsers is not required by pop and may be deleted
type OfficeUsers []OfficeUser

// String is not required by pop and may be deleted
func (o OfficeUsers) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (o *OfficeUser) Validate(tx *pop.Connection) (*validate.Errors, error) {

	return validate.Validate(
		&validators.UUIDIsPresent{Field: o.UserID, Name: "UserID"},
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
