package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
)

type OconusRateArea struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	RateAreaId         uuid.UUID `json:"rate_area_id" db:"rate_area_id"`
	CountryId          uuid.UUID `json:"country_id" db:"country_id"`
	UsPostRegionCityId uuid.UUID `json:"us_post_region_cities_id" db:"us_post_region_cities_id"`
	Active             bool      `json:"active" db:"active"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

func (o OconusRateArea) TableName() string {
	return "re_oconus_rate_areas"
}

func FetchOconusRateArea(db *pop.Connection, zip string) (*OconusRateArea, error) {
	var reOconusRateArea OconusRateArea
	err := db.Q().
		InnerJoin("re_rate_areas ra", "re_oconus_rate_areas.rate_area_id = ra.id").
		InnerJoin("us_post_region_cities upc", "upc.id = re_oconus_rate_areas.us_post_region_cities_id").
		Where("upc.uspr_zip_id = ?", zip).
		First(&reOconusRateArea)
	if err != nil {
		return nil, err
	}
	return &reOconusRateArea, nil
}

func FetchOconusRateAreaByCityId(db *pop.Connection, usprc string) (*OconusRateArea, error) {
	var reOconusRateArea OconusRateArea
	err := db.Q().
		Where("re_oconus_rate_areas.us_post_region_cities_id = ?", usprc).
		First(&reOconusRateArea)
	if err != nil {
		return nil, err
	}
	return &reOconusRateArea, nil
}
