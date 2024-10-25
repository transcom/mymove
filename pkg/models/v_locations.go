package models

import (
	"github.com/gofrs/uuid"
)

// VLocations is a read only view that represents postal region information retrieved from TRDM
/*
Column comments
uprc_id IS 'An id of a record in the us_post_region_cities table'
city_name IS 'A US postal region city name'
state IS 'A US postal region state name'
uspr_zip_id IS 'A US postal region zip identifier'
usprc_county_nm IS 'A name of the county or parish in which the UNITED-STATES-POSTAL-REGION-CITY resides'
country IS 'A name of the country'
cities_id IS 'An id of a record in the re_cities table'
state_id IS 'An id of a record in the re_states table'
us_post_regions_id IS 'An id of a record in the re_us_post_regions table'
country_id IS 'An id of a record in the re_countries table'
*/
type VLocation struct {
	UprcId           *uuid.UUID        `db:"uprc_id" json:"uprc_id" rw:"r"`
	UsPostRegionCity *UsPostRegionCity `belongs_to:"us_post_region_cities" fk_id:"uprc_id" rw:"r"`
	CityName         string            `db:"city_name" json:"city_name" rw:"r"`
	CityId           *uuid.UUID        `db:"cities_id" json:"cities_id" rw:"r"`
	City             *City             `belongs_to:"re_cities" fk_id:"cities_id" rw:"r"`
	StateName        string            `db:"state" json:"state" rw:"r"`
	StateId          *uuid.UUID        `db:"state_id" json:"state_id" rw:"r"`
	State            *State            `belongs_to:"re_states" fk_id:"state_id" rw:"r"`
	UsprZipID        string            `db:"uspr_zip_id" json:"uspr_zip_id" rw:"r"`
	UsprcCountyNm    string            `db:"usprc_county_nm" json:"usprc_county_nm" rw:"r"`
	UsPostRegionId   *uuid.UUID        `db:"us_post_regions_id" json:"us_post_regions_id" rw:"r"`
	UsPostRegion     *UsPostRegion     `belongs_to:"re_us_post_regions" fk_id:"us_post_regions_id" rw:"r"`
	CountryName      string            `db:"country" json:"country" rw:"r"`
	CountryId        *uuid.UUID        `db:"country_id" json:"country_id" rw:"r"`
	Country          *Country          `belongs_to:"re_countries" fk_id:"country_id" rw:"r"`
}

type VLocations []VLocation

// TableName overrides the table name used by Pop.
func (v VLocation) TableName() string {
	return "v_locations"
}
