package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type IntlCity struct {
	ID        uuid.UUID `json:"id" db:"id" rw:"r"`
	CityName  string    `json:"city_name" db:"city_name" rw:"r"`
	CountryId uuid.UUID `json:"country_id" db:"country_id"`
	Country   Country   `belongs_to:"re_countries" fk_id:"country_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at" rw:"r"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" rw:"r"`
}

// TableName overrides the table name used by Pop.
func (c IntlCity) TableName() string {
	return "re_intl_cities"
}
