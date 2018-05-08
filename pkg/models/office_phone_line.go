package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"time"
)

// OfficePhoneLine is used to store Phone lines (voice or fax) for the TransportationOffices
type OfficePhoneLine struct {
	ID                   uuid.UUID            `json:"id" db:"id"`
	TransportationOffice TransportationOffice `belongs_to:"transportation_office"`
	Number               string               `json:"number" db:"number"`
	Label                *string              `json:"label" db:"label"`
	IsDsnNumber          bool                 `json:"is_dsn_number" db:"is_dsn_number"`
	Type                 string               `json:"type" db:"type"`
	CreatedAt            time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (o OfficePhoneLine) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

// OfficePhoneLines is not required by pop and may be deleted
type OfficePhoneLines []OfficePhoneLine

// String is not required by pop and may be deleted
func (o OfficePhoneLines) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (o *OfficePhoneLine) Validate(tx *pop.Connection) (*validate.Errors, error) {
	officeValidator := NewFieldValidator(tx, &o.TransportationOffice, "TransportationOffice")
	validLineTypes := []string{"voice", "fax"}
	return validate.Validate(
		&validators.StringIsPresent{Field: o.Number, Name: "Number"},
		officeValidator,
		&validators.StringInclusion{Field: o.Type, Name: "Type", List: validLineTypes},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (o *OfficePhoneLine) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (o *OfficePhoneLine) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
