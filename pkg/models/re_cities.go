package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type ReCities struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CityName  string    `json:"city_name" db:"city_name"`
	StateID   uuid.UUID `json:"state_id" db:"state_id"`
	CountryID uuid.UUID `json:"country_id" db:"country_id"`
	IsOconus  bool      `json:"is_oconus" db:"is_oconus"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (c ReCities) TableName() string {
	return "re_cities"
}
