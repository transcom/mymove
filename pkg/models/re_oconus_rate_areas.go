package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type OconusRateArea struct {
	ID                 uuid.UUID `json:"id" db:"id" rw:"r"`
	RateAreaId         uuid.UUID `json:"rate_area_id" db:"rate_area_id" rw:"r"`
	CountryId          uuid.UUID `json:"country_id" db:"country_id" rw:"r"`
	UsPostRegionCityId uuid.UUID `json:"us_post_region_city_id" db:"us_post_region_city_id" rw:"r"`
	Active             bool      `json:"active" db:"active" rw:"r"`
	CreatedAt          time.Time `json:"created_at" db:"created_at" rw:"r"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at" rw:"r"`
}

func (o OconusRateArea) TableName() string {
	return "re_oconus_rate_areas"
}
