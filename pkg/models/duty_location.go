package models

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// DutyLocation represents a military duty station for a specific affiliation
type DutyLocation struct {
	ID                         uuid.UUID                     `json:"id" db:"id"`
	CreatedAt                  time.Time                     `json:"created_at" db:"created_at"`
	UpdatedAt                  time.Time                     `json:"updated_at" db:"updated_at"`
	Name                       string                        `json:"name" db:"name"`
	Affiliation                *internalmessages.Affiliation `json:"affiliation" db:"affiliation"`
	AddressID                  uuid.UUID                     `json:"address_id" db:"address_id"`
	Address                    Address                       `belongs_to:"address" fk_id:"address_id"`
	TransportationOfficeID     *uuid.UUID                    `json:"transportation_office_id" db:"transportation_office_id"`
	TransportationOffice       TransportationOffice          `belongs_to:"transportation_offices" fk_id:"transportation_office_id"`
	ProvidesServicesCounseling bool                          `json:"provides_services_counseling" db:"provides_services_counseling"`
}

// DutyLocations is not required by pop and may be deleted
type DutyLocations []DutyLocation

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (d *DutyLocation) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: d.Name, Name: "Name"},
		&validators.UUIDIsPresent{Field: d.AddressID, Name: "AddressID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (d *DutyLocation) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (d *DutyLocation) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// DutyLocationTransportInfo contains all info needed for notifications emails
type DutyLocationTransportInfo struct {
	Name      string `db:"name"`
	PhoneLine string `db:"number"`
}

// FetchDLContactInfo loads a duty station's associated transportation office and its first listed office phone number.
func FetchDLContactInfo(db *pop.Connection, dutyLocationID *uuid.UUID) (*DutyLocationTransportInfo, error) {
	if dutyLocationID == nil {
		return nil, ErrFetchNotFound
	}
	DLTransportInfo := DutyLocationTransportInfo{}
	query := `SELECT d.name, opl.number
		FROM duty_locations as d
		JOIN office_phone_lines as opl
		ON d.transportation_office_id = opl.transportation_office_id
		WHERE d.id = $1
		LIMIT 1`
	err := db.RawQuery(query, *dutyLocationID).First(&DLTransportInfo)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// Non-installation duty locations do not have transportation offices
			// so we can't look up their contact information. This isn't an error.
			return nil, nil
		default:
			return nil, err
		}
	} else if DLTransportInfo.Name == "" || DLTransportInfo.PhoneLine == "" {
		return nil, ErrFetchNotFound
	}
	return &DLTransportInfo, nil
}

// FetchDutyLocation returns a DutyLocation for a given id
func FetchDutyLocation(tx *pop.Connection, id uuid.UUID) (DutyLocation, error) {
	var station DutyLocation
	err := tx.Q().Eager("Address").Find(&station, id)
	return station, err
}

// FetchDutyLocationByName returns a DutyLocation for a given unique name
func FetchDutyLocationByName(tx *pop.Connection, name string) (DutyLocation, error) {
	var station DutyLocation
	err := tx.Where("name = ?", name).Eager("Address").First(&station)
	return station, err
}

// FindDutyLocations returns all duty locations matching a search query
func FindDutyLocations(tx *pop.Connection, search string) (DutyLocations, error) {
	var locations DutyLocations

	// The % operator filters out strings that are below this similarity threshold
	err := tx.Q().RawQuery("SET pg_trgm.similarity_threshold = 0.03").Exec()
	if err != nil {
		return locations, err
	}

	sqlQuery := `
with names as (
(select id as duty_location_id, name, similarity(name, $1) as sim
from duty_locations
where name % $1
order by sim desc
limit 5)
union
(select duty_location_id, name, similarity(name, $1) as sim
from duty_location_names
where name % $1
order by sim desc
limit 5)
union
(select dl.id as duty_location_id, dl.name as name, 1 as sim
from duty_locations as dl
inner join addresses a2 on dl.address_id = a2.id  and dl.affiliation is null
where a2.postal_code ILIKE $1
limit 5)
)
select dl.*
from names n
inner join duty_locations dl on n.duty_location_id = dl.id
group by dl.id, dl.name, dl.affiliation, dl.address_id, dl.created_at, dl.updated_at, dl.transportation_office_id, dl.provides_services_counseling
order by max(n.sim) desc, dl.name
limit 7`

	query := tx.Q().RawQuery(sqlQuery, search)
	if err := query.All(&locations); err != nil {
		if errors.Cause(err).Error() != RecordNotFoundErrorString {
			return locations, err
		}
	}

	return locations, nil
}

// FetchDutyLocationTransportationOffice returns a transportation office for a duty station
func FetchDutyLocationTransportationOffice(db *pop.Connection, dutyLocationID uuid.UUID) (TransportationOffice, error) {
	var dutyLocation DutyLocation

	err := db.Q().Eager("TransportationOffice.Address", "TransportationOffice.PhoneLines").Find(&dutyLocation, dutyLocationID)
	if err != nil {
		return TransportationOffice{}, err
	}

	if dutyLocation.TransportationOfficeID == nil {
		return TransportationOffice{}, ErrFetchNotFound
	}

	return dutyLocation.TransportationOffice, nil
}

// FetchDutyLocationsByPostalCode returns a station for a given postal code
func FetchDutyLocationsByPostalCode(tx *pop.Connection, postalCode string) (DutyLocations, error) {
	var locations DutyLocations
	query := tx.
		Where("addresses.postal_code like $1", postalCode).
		LeftJoin("addresses", "duty_location.address_id = addresses.id")

	err := query.All(&locations)
	if err != nil {
		return DutyLocations{}, err
	}

	return locations, nil
}
