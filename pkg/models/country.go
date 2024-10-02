package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// Country is a model representing a country
type Country struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Country     string    `json:"country" db:"country"`
	CountryName string    `json:"country_name" db:"country_name"`
}

// TableName overrides the table name used by Pop.
func (b Country) TableName() string {
	return "re_countries"
}
