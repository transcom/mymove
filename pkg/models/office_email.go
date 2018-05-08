package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"time"
)

// OfficeEmail is used to store Email addresses for the TransportationOffices
type OfficeEmail struct {
	ID                   uuid.UUID            `json:"id" db:"id"`
	TransportationOffice TransportationOffice `belongs_to:"transportation_office"`
	Email                string               `json:"email" db:"email"`
	Label                *string              `json:"label" db:"label"`
	CreatedAt            time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (o OfficeEmail) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

// OfficeEmails is not required by pop and may be deleted
type OfficeEmails []OfficeEmail

// String is not required by pop and may be deleted
func (o OfficeEmails) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (o *OfficeEmail) Validate(tx *pop.Connection) (*validate.Errors, error) {
	officeValidator := NewFieldValidator(tx, &o.TransportationOffice, "TransportationOffice")
	return validate.Validate(
		officeValidator,
		&validators.StringIsPresent{Field: o.Email, Name: "Email"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (o *OfficeEmail) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (o *OfficeEmail) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
