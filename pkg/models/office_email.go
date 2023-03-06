package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// OfficeEmail is used to store Email addresses for the TransportationOffices
type OfficeEmail struct {
	ID                     uuid.UUID            `json:"id" db:"id"`
	TransportationOfficeID uuid.UUID            `json:"transportation_office_id" db:"transportation_office_id"`
	TransportationOffice   TransportationOffice `belongs_to:"transportation_office" fk_id:"transportation_office_id"`
	Email                  string               `json:"email" db:"email"`
	Label                  *string              `json:"label" db:"label"`
	CreatedAt              time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time            `json:"updated_at" db:"updated_at"`
}

type OfficeEmails []OfficeEmail

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (o *OfficeEmail) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: o.TransportationOfficeID, Name: "TransportationOfficeID"},
		&validators.StringIsPresent{Field: o.Email, Name: "Email"},
	), nil
}
