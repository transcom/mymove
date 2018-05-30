package models

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"time"
)

// TransportationOffice is a PPPO, PPSO or JPPSO. If it is its own shipping office, ShippingOffice will be nil,
// otherwise its a pointer to the actual shipping office.
type TransportationOffice struct {
	ID               uuid.UUID             `json:"id" db:"id"`
	ShippingOfficeID *uuid.UUID            `json:"shipping_office_id" db:"shipping_office_id"`
	ShippingOffice   *TransportationOffice `belongs_to:"transportation_offices"`
	Name             string                `json:"name" db:"name"`
	Address          Address               `belongs_to:"address"`
	AddressID        uuid.UUID             `json:"address_id" db:"address_id"`
	Latitude         float32               `json:"latitude" db:"latitude"`
	Longitude        float32               `json:"longitude" db:"longitude"`
	PhoneLines       OfficePhoneLines      `has_many:"office_phone_lines"`
	Emails           OfficeEmails          `has_many:"office_emails"`
	Hours            *string               `json:"hours" db:"hours"`
	Services         *string               `json:"services" db:"services"`
	Note             *string               `json:"note" db:"note"`
	CreatedAt        time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at" db:"updated_at"`
}

// TransportationOffices is not required by pop and may be deleted
type TransportationOffices []TransportationOffice

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *TransportationOffice) Validate(tx *pop.Connection) (*validate.Errors, error) {
	verrs := validate.Validate(
		&validators.StringIsPresent{Field: t.Name, Name: "Name"},
		&validators.UUIDIsPresent{Field: t.AddressID, Name: "AddressID"},
	)
	return verrs, nil
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
