package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// DutyStation represents a military duty station for a specific affiliation
type DutyStation struct {
	ID                     uuid.UUID                    `json:"id" db:"id"`
	CreatedAt              time.Time                    `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time                    `json:"updated_at" db:"updated_at"`
	Name                   string                       `json:"name" db:"name"`
	Affiliation            internalmessages.Affiliation `json:"affiliation" db:"affiliation"`
	AddressID              uuid.UUID                    `json:"address_id" db:"address_id"`
	Address                Address                      `belongs_to:"address"`
	TransportationOfficeID *uuid.UUID                   `json:"transportation_office_id" db:"transportation_office_id"`
	TransportationOffice   *TransportationOffice        `belongs_to:"transportation_offices"`
}

// DutyStations is not required by pop and may be deleted
type DutyStations []DutyStation

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (d *DutyStation) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: d.Name, Name: "Name"},
		&AffiliationIsPresent{Field: d.Affiliation, Name: "Affiliation"},
		&validators.UUIDIsPresent{Field: d.AddressID, Name: "AddressID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (d *DutyStation) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (d *DutyStation) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// FetchDutyStation returns a station for a given id
func FetchDutyStation(tx *pop.Connection, id uuid.UUID) (DutyStation, error) {
	var station DutyStation
	err := tx.Q().Eager().Find(&station, id)
	return station, err
}

// FindDutyStations returns all duty stations matching a search query and military affiliation
func FindDutyStations(tx *pop.Connection, search string) (DutyStations, error) {
	var stations DutyStations

	// ILIKE does case-insensitive pattern matching, "%" matches any string
	// We build a query by inserting '%' between each letter in the search string.
	// This allows matching substrings as well as abbreviations.
	// It would probably be worth ordering the results by similarity to the search string, one day.
	searchQuery := []rune("%")

	for _, runeChar := range search {
		searchQuery = append(searchQuery, runeChar)
		searchQuery = append(searchQuery, '%')
	}
	queryString := string(searchQuery)

	query := tx.Q().Eager().Where("name ILIKE $1", queryString)

	if err := query.All(&stations); err != nil {
		if errors.Cause(err).Error() != recordNotFoundErrorString {
			return stations, err
		}
	}

	return stations, nil
}

// FetchDutyStationTransportationOffice returns a transportation office for a duty station
func FetchDutyStationTransportationOffice(db *pop.Connection, dutyStationID uuid.UUID) (*TransportationOffice, error) {
	var dutyStation DutyStation

	err := db.Q().Eager("TransportationOffice").Find(&dutyStation, dutyStationID)
	if err != nil {
		return nil, err
	}

	if dutyStation.TransportationOffice == nil {
		return nil, ErrFetchNotFound
	}

	if err = db.Load(dutyStation.TransportationOffice, "Address"); err != nil {
		return nil, err
	}

	if err = db.Load(dutyStation.TransportationOffice, "PhoneLines"); err != nil {
		return nil, err
	}

	return dutyStation.TransportationOffice, nil
}
