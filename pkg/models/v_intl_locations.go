package models

import (
	"github.com/gofrs/uuid"
)

type VIntlLocation struct {
	IntlCityCountriesID    *uuid.UUID             `db:"icc_id" json:"icc_id" rw:"r"`
	IntlCityCountries      *IntlCityCountries     `belongs_to:"intl_city_countries" fk_id:"icc_id" rw:"r"`
	CityName               *string                `db:"city_name" json:"city_name" rw:"r"`
	CountryPrnDivID        *string                `db:"country_prn_dv_id" json:"country_prn_dv_id" rw:"r"`
	CountryPrnDivName      *string                `db:"country_prn_dv_nm" json:"country_prn_dv_nm" rw:"r"`
	CountryPrnDivCode      *string                `db:"country_prn_dv_cd" json:"country_prn_dv_cd" rw:"r"`
	CountryCode            *string                `db:"country" json:"country" rw:"r"`
	IntlCityId             *uuid.UUID             `db:"intl_cities_id" json:"intl_cities_id" rw:"r"`
	IntlCity               *IntlCity              `belongs_to:"re_intl_cities" fk_id:"intl_cities_id" rw:"r"`
	ReCountryPrnDivisionID *uuid.UUID             `db:"re_country_prn_division_id" json:"re_country_prn_division_id" rw:"r"`
	ReCountryPrnDivision   *ReCountryPrnDivisions `belongs_to:"re_country_prn_divisions" fk_id:"country_prn_division_id" rw:"r"`
	CountryId              *uuid.UUID             `db:"country_id" json:"country_id" rw:"r"`
	Country                *Country               `belongs_to:"re_countries" fk_id:"country_id" rw:"r"`
}

type VIntlLocations []VIntlLocation

// TableName overrides the table name used by Pop.
func (v VIntlLocation) TableName() string {
	return "v_intl_locations"
}
