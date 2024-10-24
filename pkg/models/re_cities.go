package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type City struct {
	ID        uuid.UUID `json:"id" db:"id" rw:"r"`
	CityName  string    `json:"city_name" db:"city_name" rw:"r"`
	StateId   uuid.UUID `json:"state_id" db:"state_id" rw:"r"`
	State     State     `belongs_to:"re_states" fk_id:"state_id" rw:"r"`
	CountryId uuid.UUID `json:"country_id" db:"country_id" rw:"r"`
	Country   Country   `belongs_to:"re_countries" fk_id:"country_id" rw:"r"`
	IsOconus  bool      `json:"is_oconus" db:"is_oconus" rw:"r"`
	CreatedAt time.Time `json:"created_at" db:"created_at" rw:"r"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" rw:"r"`
}

// TableName overrides the table name used by Pop.
func (c City) TableName() string {
	return "re_cities"
}
