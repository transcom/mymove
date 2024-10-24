package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type PortLocation struct {
	ID                   uuid.UUID `json:"id" db:"id"`
	PortId               uuid.UUID `json:"port_id" db:"port_id"`
	CitiesId             uuid.UUID `json:"cities_id" db:"cities_id"`
	UsPostRegionCitiesId uuid.UUID `json:"us_post_region_cities_id" db:"us_post_region_cities_id"`
	CountryId            uuid.UUID `json:"country_id" db:"country_id"`
	IsActive             bool      `json:"is_active" db:"is_active"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

func (l PortLocation) TableName() string {
	return "port_locations"
}
