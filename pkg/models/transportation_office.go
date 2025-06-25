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
// otherwise it's a pointer to the actual shipping office.
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
	ProvidesCloseout bool                  `json:"provides_ppm_closeout" db:"provides_ppm_closeout"`
}

// TableName overrides the table name used by Pop.
func (t TransportationOffice) TableName() string {
	return "transportation_offices"
}

type TransportationOffices []TransportationOffice

type GBLOCs []string

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (t *TransportationOffice) Validate(_ *pop.Connection) (*validate.Errors, error) {
	verrs := validate.Validate(
		&validators.StringIsPresent{Field: t.Name, Name: "Name"},
		&validators.UUIDIsPresent{Field: t.AddressID, Name: "AddressID"},
	)
	return verrs, nil
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

// GetCounselingOffices calls a db function that returns all the transportation offices in the GBLOC
// of the given duty location where provides_services_counseling = true
func GetCounselingOffices(db *pop.Connection, dutyLocationID uuid.UUID, serviceMemberID uuid.UUID) (TransportationOffices, error) {
	var officeList TransportationOffices

	err := db.RawQuery("SELECT * FROM get_counseling_offices($1, $2)", dutyLocationID, serviceMemberID).
		All(&officeList)
	if err != nil {
		return officeList, err
	}

	return officeList, nil
}

// FetchTransportationOfficeByID fetches an office user by ID
func FetchTransportationOfficeByID(tx *pop.Connection, id uuid.UUID) (*TransportationOffice, error) {
	var transportationOffice TransportationOffice
	err := tx.Find(&transportationOffice, id)
	return &transportationOffice, err
}
