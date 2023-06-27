package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// OfficePhoneLine is used to store Phone lines (voice or fax) for the TransportationOffices
type OfficePhoneLine struct {
	ID                     uuid.UUID            `json:"id" db:"id"`
	TransportationOfficeID uuid.UUID            `json:"transportation_office_id" db:"transportation_office_id"`
	TransportationOffice   TransportationOffice `belongs_to:"transportation_office" fk_id:"transportation_office_id"`
	Number                 string               `json:"number" db:"number"`
	Label                  *string              `json:"label" db:"label"`
	IsDsnNumber            bool                 `json:"is_dsn_number" db:"is_dsn_number"`
	Type                   string               `json:"type" db:"type"`
	CreatedAt              time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time            `json:"updated_at" db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (o OfficePhoneLine) TableName() string {
	return "office_phone_lines"
}

type OfficePhoneLines []OfficePhoneLine

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (o *OfficePhoneLine) Validate(_ *pop.Connection) (*validate.Errors, error) {
	validLineTypes := []string{"voice", "fax"}
	return validate.Validate(
		&validators.StringIsPresent{Field: o.Number, Name: "Number"},
		&validators.UUIDIsPresent{Field: o.TransportationOfficeID, Name: "TransportationOfficeID"},
		&validators.StringInclusion{Field: o.Type, Name: "Type", List: validLineTypes},
	), nil
}
