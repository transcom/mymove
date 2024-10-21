package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type City struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CityName  string    `json:"city_name" db:"city_name"`
	StateId   uuid.UUID `json:"state_id" db:"state_id"`
	State     State     `belongs_to:"re_states" fk_id:"state_id"`
	CountryId uuid.UUID `json:"country_id" db:"country_id"`
	Country   Country   `belongs_to:"re_countries" fk_id:"country_id"`
	IsOconus  bool      `json:"is_oconus" db:"is_oconus"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (c City) TableName() string {
	return "re_cities"
}
