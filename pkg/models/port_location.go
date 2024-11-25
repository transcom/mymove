package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

type PortLocation struct {
	ID                   uuid.UUID `json:"id" db:"id"`
	PortId               uuid.UUID `json:"port_id" db:"port_id"`
	Port                 Port      `belongs_to:"port_locations" fk_id:"port_id"`
	CitiesId             uuid.UUID `json:"cities_id" db:"cities_id"`
	UsPostRegionCitiesId uuid.UUID `json:"us_post_region_cities_id" db:"us_post_region_cities_id"`
	CountryId            uuid.UUID `json:"country_id" db:"country_id"`
	IsActive             *bool     `json:"is_active" db:"is_active"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

func (l PortLocation) TableName() string {
	return "port_locations"
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *PortLocation) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: p.PortId, Name: "PortID"},
		&validators.UUIDIsPresent{Field: p.CitiesId, Name: "CitiesID"},
		&validators.UUIDIsPresent{Field: p.UsPostRegionCitiesId, Name: "UsPostRegionCitiesID"},
		&validators.UUIDIsPresent{Field: p.CountryId, Name: "CountryID"},
	), nil
}
