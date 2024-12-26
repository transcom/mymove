package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// DutyLocation represents a military duty location for a specific affiliation
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

// TableName overrides the table name used by Pop.
func (d DutyLocation) TableName() string {
	return "duty_locations"
}

type DutyLocations []DutyLocation

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (d *DutyLocation) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: d.Name, Name: "Name"},
		&validators.UUIDIsPresent{Field: d.AddressID, Name: "AddressID"},
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
	err := tx.Q().Eager("Address", "Address.Country").Find(&dutyLocation, id)
	return dutyLocation, err
}

// FetchDutyLocationByName returns a DutyLocation for a given unique name
func FetchDutyLocationByName(tx *pop.Connection, name string) (DutyLocation, error) {
	var dutyLocation DutyLocation
	err := tx.Where("name = ?", name).Eager("Address", "Address.Country", "TransportationOffice",
		"TransportationOffice.Address", "TransportationOffice.PhoneLines").First(&dutyLocation)
	return dutyLocation, err
}

// FetchDutyLocationWithTransportationOffice returns a DutyLocation for a given id
// with the associated transportation office eagerly loaded
func FetchDutyLocationWithTransportationOffice(tx *pop.Connection, id uuid.UUID) (DutyLocation, error) {
	var dutyLocation DutyLocation
	err := tx.Q().Eager("Address", "Address.Country", "TransportationOffice", "TransportationOffice.Address",
		"TransportationOffice.PhoneLines").Find(&dutyLocation, id)
	return dutyLocation, err
}

// FindDutyLocations returns all duty locations matching a search query while excluding certain location by specified states.
func FindDutyLocationsExcludingStates(tx *pop.Connection, search string, exclusionStateFilters []string) (DutyLocations, error) {
	var locations DutyLocations

	// The % operator filters out strings that are below this similarity threshold
	err := tx.Q().RawQuery("SET pg_trgm.similarity_threshold = 0.03").Exec()
	if err != nil {
		return locations, err
	}

	sql_builder := strings.Builder{}
	sql_builder.WriteString(`with names as
	(
		(
	        -- search against duty_locations table
			(
				select
					id as duty_location_id,
					name,
					similarity(name, $1) as sim
				from
					duty_locations
				where
					name % $1
				order by sim desc limit 5
			)
			-- exclude OCONUS locations that are not active
			except
			(
				select
					d.id as duty_location_id,
					d.name,
					similarity(d.name, $1) as sim
				from
					duty_locations d,
					addresses a,
					re_oconus_rate_areas o,
					re_rate_areas r
				where d.name % $1
				and d.address_id = a.id
				and a.us_post_region_cities_id = o.us_post_region_cities_id
				and o.rate_area_id = r.id
				and o.active = false
			)
		)
		union
		(
		    -- search against duty_location_names table for alternative names
			(
				select
					duty_location_id,
					name,
					similarity(name, $1) as sim
				from
					duty_location_names
				where
					name % $1
				order by sim desc limit 5
			)
			-- exclude OCONUS locations that are not active
			except
			(
				select
					dn.duty_location_id,
					dn.name,
					similarity(dn.name, $1) as sim
				from
					duty_location_names dn,
					duty_locations d,
					addresses a,
					re_oconus_rate_areas o,
					re_rate_areas r
				where
					dn.name % $1
					and dn.duty_location_id = d.id
					and d.address_id = a.id
					and a.us_post_region_cities_id = o.us_post_region_cities_id
					and o.rate_area_id = r.id
					and o.active = false
				order by sim desc limit 5
			)
		)
		union
		(
		    -- search against duty_locations table if postal code
			(
				select
					dl.id as duty_location_id,
					dl.name as name,
					1 as sim
				from
					duty_locations as dl
					inner join addresses a2 on dl.address_id = a2.id
					and dl.affiliation is null
				where
					a2.postal_code ILIKE $1
				limit 5
			)
			-- exclude OCONUS locations that are not active
			except
			(
				select
					dl.id as duty_location_id,
					dl.name as name,
					1 as sim
				from
					duty_locations dl,
					addresses a,
					re_oconus_rate_areas o,
					re_rate_areas r
				where
					dl.address_id = a.id
					and a.us_post_region_cities_id = o.us_post_region_cities_id
					and o.rate_area_id = r.id
					and o.active = false
					and a.postal_code ILIKE $1
					and dl.affiliation is null
				limit 5
			)
		)
	)
	select
		dl.*
	from
		names n
		inner join duty_locations dl on n.duty_location_id = dl.id
		inner join addresses a on dl.address_id = a.id`)

	// apply filter to exclude specific states if provided
	if len(exclusionStateFilters) > 0 {
		exclusionStateParams := make([]string, 0)
		for _, value := range exclusionStateFilters {
			exclusionStateParams = append(exclusionStateParams, fmt.Sprintf("'%s'", value))
		}
		sql_builder.WriteString(fmt.Sprintf(" and a.state not in (%s)", strings.Join(exclusionStateParams, ",")))
	}

	// filter to exclude po box postal codes
	sql_builder.WriteString(`
	where
	not exists (
		select 1
		from re_us_post_regions upr
		where upr.uspr_zip_id = a.postal_code
			and upr.is_po_box = true
	)`)

	sql_builder.WriteString(`
	group by dl.id, dl.name, dl.affiliation, dl.address_id, dl.created_at, dl.updated_at, dl.transportation_office_id, dl.provides_services_counseling
	order by max(n.sim) desc, dl.name
	limit 7`)

	query := tx.Q().RawQuery(sql_builder.String(), search)
	if err := query.All(&locations); err != nil {
		if errors.Cause(err).Error() != RecordNotFoundErrorString {
			return locations, err
		}
	}

	return locations, nil

}

// FindDutyLocations returns all duty locations matching a search query
func FindDutyLocations(tx *pop.Connection, search string) (DutyLocations, error) {
	return FindDutyLocationsExcludingStates(tx, search, []string{})
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
		Where("addresses.postal_code like $1", postalCode).
		LeftJoin("addresses", "duty_location.address_id = addresses.id").
		LeftJoin("re_countries", "addresses.country_id = re_countries.id")

	err := query.All(&locations)
	if err != nil {
		return DutyLocations{}, err
	}

	return locations, nil
}
