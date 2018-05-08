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
	ID                   uuid.UUID            `json:"id" db:"id"`
	User                 User                 `belongs_to:"user"`
	FamilyName           string               `json:"family_name" db:"family_name"`
	GivenName            string               `json:"given_name" db:"given_name"`
	MiddleInitials       *string              `json:"middle_initial" db:"middle_initial"`
	Email                string               `json:"email" db:"email"`
	Telephone            string               `json:"telephone" db:"telephone"`
	TransportationOffice TransportationOffice `belongs_to:"transportation_office"`
	CreatedAt            time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at" db:"updated_at"`
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
	userValidator := NewFieldValidator(tx, &o.User, "User")
	officeValidator := NewFieldValidator(tx, &o.TransportationOffice, "TransportationOffice")

	return validate.Validate(
		userValidator,
		&validators.StringIsPresent{Field: o.FamilyName, Name: "FamilyName"},
		&validators.StringIsPresent{Field: o.GivenName, Name: "GivenName"},
		&validators.StringIsPresent{Field: o.Email, Name: "Email"},
		&validators.StringIsPresent{Field: o.Telephone, Name: "Telephone"},
		officeValidator,
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
