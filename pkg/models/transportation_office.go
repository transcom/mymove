package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// TransportationOffice is a PPPO, PPSO or JPPSO. If it is its own shipping office, ShippingOffice will be nil,
// otherwise its a pointer to the actual shipping office.
type TransportationOffice struct {
	ID               uuid.UUID             `json:"id" db:"id"`
	ShippingOfficeID *uuid.UUID            `json:"shipping_office_id" db:"shipping_office_id"`
	ShippingOffice   *TransportationOffice `belongs_to:"transportation_offices" fk_id:"shipping_office_id"`
	Name             string                `json:"name" db:"name"`
	Address          Address               `belongs_to:"address" fk_id:"address_id"`
	AddressID        uuid.UUID             `json:"address_id" db:"address_id"`
	Latitude         float32               `json:"latitude" db:"latitude"`
	Longitude        float32               `json:"longitude" db:"longitude"`
	PhoneLines       OfficePhoneLines      `has_many:"office_phone_lines" fk_id:"transportation_office_id"`
	Emails           OfficeEmails          `has_many:"office_emails" fk_id:"transportation_office_id"`
	Hours            *string               `json:"hours" db:"hours"`
	Services         *string               `json:"services" db:"services"`
	Note             *string               `json:"note" db:"note"`
	Gbloc            string                `json:"gbloc" db:"gbloc"`
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

// FetchNearestTransportationOffice fetches the nearest transportation office
func FetchNearestTransportationOffice(tx *pop.Connection, long float32, lat float32) (TransportationOffice, error) {
	var to TransportationOffice

	query := `
		select *
			from transportation_offices
		WHERE shipping_office_id IS NOT NULL
		order by ST_Distance(
  		ST_GeographyFromText(concat('point(',$1::text,' ',$2::text,')'))
  		, ST_GeographyFromText(concat('point(',longitude, ' ', latitude,')'))
 		) asc`

	err := tx.RawQuery(query, long, lat).First(&to)
	if err != nil {
		return to, errors.Wrap(err, "Fetch transportation office failed")
	}

	return to, nil
}
