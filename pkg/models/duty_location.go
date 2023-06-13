package models

import (
	"database/sql"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// DutyLocation represents a military duty location for a specific
// affiliation
//
// As of 2023-06-13, none of the duty locations had secondary street
// addresses. Only 5 had a street address at all other than 'n/a' or
// 'N/A'
type DutyLocation struct {
	ID                         uuid.UUID                     `json:"id" db:"id"`
	CreatedAt                  time.Time                     `json:"created_at" db:"created_at"`
	UpdatedAt                  time.Time                     `json:"updated_at" db:"updated_at"`
	Name                       string                        `json:"name" db:"name"`
	Affiliation                *internalmessages.Affiliation `json:"affiliation" db:"affiliation"`
	StreetAddress1             string                        `json:"street_address_1" db:"street_address_1"`
	City                       string                        `json:"city" db:"city"`
	State                      string                        `json:"state" db:"state"`
	PostalCode                 string                        `json:"postal_code" db:"postal_code"`
	Country                    string                        `json:"country" db:"country"`
	TransportationOfficeID     *uuid.UUID                    `json:"transportation_office_id" db:"transportation_office_id"`
	TransportationOffice       TransportationOffice          `belongs_to:"transportation_offices" fk_id:"transportation_office_id"`
	ProvidesServicesCounseling bool                          `json:"provides_services_counseling" db:"provides_services_counseling"`
}

// TableName overrides the table name used by Pop.
func (d DutyLocation) TableName() string {
	return "duty_locations"
}

type DutyLocations []DutyLocation

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (d *DutyLocation) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: d.Name, Name: "Name"},
		&validators.StringIsPresent{Field: d.City, Name: "City"},
		&validators.StringIsPresent{Field: d.State, Name: "State"},
		&validators.StringIsPresent{Field: d.PostalCode, Name: "PostalCode"},
		&validators.StringIsPresent{Field: d.Country, Name: "Country"},
	), nil
}

// DutyLocationTransportInfo contains all info needed for notifications emails
type DutyLocationTransportInfo struct {
	Name      string `db:"name"`
	PhoneLine string `db:"number"`
}

// FetchDLContactInfo loads a duty location's associated transportation office and its first listed office phone number.
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
	var dutyLocation DutyLocation
	err := tx.Q().Find(&dutyLocation, id)
	return dutyLocation, err
}

// FetchDutyLocationByName returns a DutyLocation for a given unique name
func FetchDutyLocationByName(tx *pop.Connection, name string) (DutyLocation, error) {
	var dutyLocation DutyLocation
	err := tx.Where("name = ?", name).Eager("TransportationOffice",
		"TransportationOffice.Address", "TransportationOffice.PhoneLines").First(&dutyLocation)
	return dutyLocation, err
}

// FetchDutyLocationWithTransportationOffice returns a DutyLocation for a given id
// with the associated transportation office eagerly loaded
func FetchDutyLocationWithTransportationOffice(tx *pop.Connection, id uuid.UUID) (DutyLocation, error) {
	var dutyLocation DutyLocation
	err := tx.Q().Eager("TransportationOffice", "TransportationOffice.Address",
		"TransportationOffice.PhoneLines").Find(&dutyLocation, id)
	return dutyLocation, err
}

// FindDutyLocations returns all duty locations matching a search query
func FindDutyLocations(tx *pop.Connection, search string) (DutyLocations, error) {
	var locations DutyLocations

	// The % operator filters out strings that are below this similarity threshold
	err := tx.Q().RawQuery("SET LOCAL pg_trgm.similarity_threshold = 0.03").Exec()
	if err != nil {
		return locations, err
	}

	// Because we are changing the schema to move address information
	// into the duty_locations table, the final SELECT query has to
	// explicitly name the columns we want so that the address_id
	// field is not included, since that field is no longer available
	// in the model
	sqlQuery := `
WITH names AS (
(SELECT id AS duty_location_id, name, similarity(name, $1) AS sim
FROM duty_locations
WHERE name % $1
ORDER BY sim DESC
LIMIT 5)
UNION
(SELECT duty_location_id, name, similarity(name, $1) AS sim
FROM duty_location_names
WHERE name % $1
ORDER BY sim desc
LIMIT 5)
UNION
(SELECT dl.id AS duty_location_id, dl.name AS name, 1 AS sim
FROM duty_locations AS dl
WHERE dl.postal_code LIKE $1 AND dl.affiliation IS NULL
LIMIT 5)
)
SELECT dl.id, dl.name, dl.affiliation, dl.created_at, dl.updated_at,
       dl.transportation_office_id, dl.provides_services_counseling,
       dl.street_address_1, dl.city, dl.state, dl.postal_code, dl.country
FROM names n
INNER JOIN duty_locations dl ON n.duty_location_id = dl.id
GROUP BY dl.id, dl.name, dl.affiliation, dl.postal_code, dl.created_at,
         dl.updated_at, dl.transportation_office_id, dl.provides_services_counseling
ORDER BY MAX(n.sim) DESC, dl.name
LIMIT 7`

	query := tx.Q().RawQuery(sqlQuery, search)
	if err := query.All(&locations); err != nil {
		if errors.Cause(err).Error() != RecordNotFoundErrorString {
			return locations, err
		}
	}

	return locations, nil
}

// FetchDutyLocationTransportationOffice returns a transportation office for a duty location
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

// FetchDutyLocationsByPostalCode returns a duty location for a given postal code
func FetchDutyLocationsByPostalCode(tx *pop.Connection, postalCode string) (DutyLocations, error) {
	var locations DutyLocations
	query := tx.
		Where("postal_code like $1", postalCode)

	err := query.All(&locations)
	if err != nil {
		return DutyLocations{}, err
	}

	return locations, nil
}
