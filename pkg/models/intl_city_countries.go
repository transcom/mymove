package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type IntlCityCountries struct {
	ID                   uuid.UUID             `json:"id" db:"id" rw:"r"`
	CountryId            uuid.UUID             `json:"country_id" db:"country_id"`
	Country              Country               `belongs_to:"re_countries" fk_id:"country_id"`
	IntlCitiesId         uuid.UUID             `json:"intl_cities_id" db:"intl_cities_id" rw:"r"`
	IntlCities           IntlCity              `belongs_to:"re_intl_cities" fk_id:"intl_cities_id" rw:"r"`
	CountryPrnDivisionId uuid.UUID             `json:"country_prn_division_id" db:"country_prn_division_id" rw:"r"`
	CountryPrnDivision   ReCountryPrnDivisions `belongs_to:"re_country_prn_divisions" fk_id:"country_prn_division_id" rw:"r"`
	CreatedAt            time.Time             `json:"created_at" db:"created_at" rw:"r"`
	UpdatedAt            time.Time             `json:"updated_at" db:"updated_at" rw:"r"`
}

// TableName overrides the table name used by Pop.
func (c IntlCityCountries) TableName() string {
	return "intl_city_countries"
}
