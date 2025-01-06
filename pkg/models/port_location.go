package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type PortLocation struct {
	ID                   uuid.UUID        `json:"id" db:"id" rw:"r"`
	PortId               uuid.UUID        `json:"port_id" db:"port_id" rw:"r"`
	Port                 Port             `belongs_to:"port_locations" fk_id:"port_id" rw:"r"`
	CitiesId             uuid.UUID        `json:"cities_id" db:"cities_id" rw:"r"`
	City                 City             `belongs_to:"re_cities" fk_id:"cities_id" rw:"r"`
	UsPostRegionCitiesId uuid.UUID        `json:"us_post_region_cities_id" db:"us_post_region_cities_id" rw:"r"`
	UsPostRegionCity     UsPostRegionCity `belongs_to:"us_post_region_cities" fk_id:"us_post_region_cities_id" rw:"r"`
	CountryId            uuid.UUID        `json:"country_id" db:"country_id" rw:"r"`
	Country              Country          `belongs_to:"re_countries" fk_id:"country_id" rw:"r"`
	IsActive             *bool            `json:"is_active" db:"is_active" rw:"r"`
	CreatedAt            time.Time        `json:"created_at" db:"created_at" rw:"r"`
	UpdatedAt            time.Time        `json:"updated_at" db:"updated_at" rw:"r"`
}

func (l PortLocation) TableName() string {
	return "port_locations"
}
