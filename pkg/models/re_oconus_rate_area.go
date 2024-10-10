package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type ReOconusRateArea struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	RateAreaId         uuid.UUID `json:"rate_area_id" db:"rate_area_id"`
	CountryId          uuid.UUID `json:"country_id" db:"country_id"`
	UsPostRegionCityId uuid.UUID `json:"us_post_region_city_id" db:"us_post_region_city_id"`
	InactiveFlag       bool      `json:"inactive_flag" db:"inactive_flag"`
}

func (ReOconusRateArea) TableName() string {
	return "re_oconus_rate_areas"
}
