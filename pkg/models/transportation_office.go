package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"time"
)

// TransportationOffice is a PPO, PPSO or JPPSO. If it is its own shipping office, ShippingOffice will be nil,
// otherwise its a pointer to the actual shipping office.
type TransportationOffice struct {
	ID             uuid.UUID             `json:"id" db:"id"`
	ShippingOffice *TransportationOffice `belongs_to:"transportation_offices"`
	Name           string                `json:"name" db:"name"`
	Address        Address               `belongs_to:"address"`
	Latitude       float32               `json:"latitude" db:"latitude"`
	Longitude      float32               `json:"longitude" db:"longitude"`
	PhoneLines     OfficePhoneLines      `belongs_to:"office_phone_lines"`
	Emails         OfficeEmails          `belongs_to:"office_emails"`
	Hours          *string               `json:"hour" db:"hour"`
	Services       *string               `json:"service" db:"service"`
	Note           *string               `json:"note" db:"note"`
	CreatedAt      time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (t TransportationOffice) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// TransportationOffices is not required by pop and may be deleted
type TransportationOffices []TransportationOffice

// String is not required by pop and may be deleted
func (t TransportationOffices) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *TransportationOffice) Validate(tx *pop.Connection) (*validate.Errors, error) {
	addressValidator := NewFieldValidator(tx, &t.Address, "Address")
	verrs := validate.Validate(
		addressValidator,
		&validators.StringIsPresent{Field: t.Name, Name: "Name"},
	)
	return verrs, addressValidator.Error
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *TransportationOffice) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *TransportationOffice) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
