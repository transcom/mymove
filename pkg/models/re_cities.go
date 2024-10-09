package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type ReCities struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CityName  string    `json:"city_name" db:"city_name"`
	StateID    uuid.UUID `json:"state_id" db:"state_id"`
	CountryID  uuid.UUID `json:"country_id" db:"country_id"`
	IsConus    bool      `json:"is_conus" db:"is_conus"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (d ReCities) TableName() string {
	return "re_cities"
}