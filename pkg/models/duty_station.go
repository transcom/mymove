package models

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
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
	TransportationOffice   TransportationOffice         `belongs_to:"transportation_offices"`
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

// DutyStationTransportInfo contains all info needed for notifications emails
type DutyStationTransportInfo struct {
	Name      string `db:"name"`
	PhoneLine string `db:"number"`
}

// FetchDSContactInfo loads a duty station's associated transportation office and its first listed office phone number.
func FetchDSContactInfo(db *pop.Connection, dutyStationID *uuid.UUID) (*DutyStationTransportInfo, error) {
	if dutyStationID == nil {
		return nil, ErrFetchNotFound
	}
	DSTransportInfo := DutyStationTransportInfo{}
	query := `SELECT d.name, opl.number
		FROM duty_stations as d
		JOIN office_phone_lines as opl
		ON d.transportation_office_id = opl.transportation_office_id
		WHERE d.id = $1
		LIMIT 1`
	err := db.RawQuery(query, *dutyStationID).First(&DSTransportInfo)
	if err != nil {
		return nil, err
	} else if DSTransportInfo.Name == "" || DSTransportInfo.PhoneLine == "" {
		return nil, ErrFetchNotFound
	}
	return &DSTransportInfo, nil
}

// FetchDutyStation returns a station for a given id
func FetchDutyStation(tx *pop.Connection, id uuid.UUID) (DutyStation, error) {
	var station DutyStation
	err := tx.Q().Eager("Address").Find(&station, id)
	return station, err
}

// FetchDutyStationByName returns a station for a given unique name
func FetchDutyStationByName(tx *pop.Connection, name string) (DutyStation, error) {
	var station DutyStation
	err := tx.Where("name = ?", name).Eager("Address").First(&station)
	return station, err
}

// FindDutyStations returns all duty stations matching a search query
func FindDutyStations(tx *pop.Connection, search string) (DutyStations, error) {
	var stations DutyStations

	sql := `
with names as (
(select id as duty_station_id, name, similarity(name, $1) as sim
from duty_stations
where similarity(name, $1) > 0.03
order by sim desc
limit 5)
union
(select duty_station_id, name, similarity(name, $1) as sim
from duty_station_names
where similarity(name, $1) > 0.03
order by sim desc
limit 5)
)
select ds.*
from names n
inner join duty_stations ds on n.duty_station_id = ds.id
group by ds.id, ds.name, ds.affiliation, ds.address_id, ds.created_at, ds.updated_at, ds.transportation_office_id
order by max(n.sim) desc, ds.name
limit 7`

	query := tx.Q().RawQuery(sql, search)
	if err := query.All(&stations); err != nil {
		if errors.Cause(err).Error() != RecordNotFoundErrorString {
			return stations, err
		}
	}

	return stations, nil
}

// FetchDutyStationTransportationOffice returns a transportation office for a duty station
func FetchDutyStationTransportationOffice(db *pop.Connection, dutyStationID uuid.UUID) (TransportationOffice, error) {
	var dutyStation DutyStation

	err := db.Q().Eager("TransportationOffice.Address", "TransportationOffice.PhoneLines").Find(&dutyStation, dutyStationID)
	if err != nil {
		return TransportationOffice{}, err
	}

	if dutyStation.TransportationOfficeID == nil {
		return TransportationOffice{}, ErrFetchNotFound
	}

	return dutyStation.TransportationOffice, nil
}

// FetchDutyStationsByPostalCode returns a station for a given postal code
func FetchDutyStationsByPostalCode(tx *pop.Connection, postalCode string) (DutyStations, error) {
	var stations DutyStations
	query := tx.
		Where("addresses.postal_code like $1", postalCode).
		LeftJoin("addresses", "duty_stations.address_id = addresses.id")

	err := query.All(&stations)
	if err != nil {
		return DutyStations{}, err
	}

	return stations, nil
}
