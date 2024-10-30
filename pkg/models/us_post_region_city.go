package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// UsPostRegionCity represents postal region information retrieved from TRDM
/*
Column comments
uspr_zip_id IS 'A US postal region zip identifier.';
u_s_post_region_city_nm IS 'A US postal region city name.';
usprc_county_nm IS 'A name of the county or parish in which the UNITED-STATES- POSTAL-REGION-CITY resides.';
ctry_genc_dgph_cd IS 'A 2-digit Geopolitical Entities, Names, and Codes (GENC) Standard.';
*/
type UsPostRegionCity struct {
	ID                 uuid.UUID    `db:"id" json:"id"`
	UsprZipID          string       `db:"uspr_zip_id" json:"uspr_zip_id"`
	USPostRegionCityNm string       `db:"u_s_post_region_city_nm" json:"u_s_post_region_city_nm"`
	UsprcCountyNm      string       `db:"usprc_county_nm" json:"usprc_county_nm"`
	CtryGencDgphCd     string       `db:"ctry_genc_dgph_cd" json:"ctry_genc_dgph_cd"`
	CreatedAt          time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time    `db:"updated_at" json:"updated_at"`
	UsPostRegionId     uuid.UUID    `db:"us_post_regions_id" json:"us_post_regions_id"`
	UsPostRegion       UsPostRegion `belongs_to:"re_us_post_regions" fk_id:"us_post_regions_id"`
	CityId             uuid.UUID    `db:"cities_id" json:"cities_id"`
	City               City         `belongs_to:"re_cities" fk_id:"cities_id"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (usprc *UsPostRegionCity) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringLengthInRange{Field: usprc.UsprZipID, Name: "UsprZipID", Min: 5, Max: 5},
		&validators.StringLengthInRange{Field: usprc.CtryGencDgphCd, Name: "CtryGencDgphCd", Min: 2, Max: 2},
		&validators.StringIsPresent{Field: usprc.USPostRegionCityNm, Name: "USPostRegionCityNm"},
		&validators.StringIsPresent{Field: usprc.UsprcCountyNm, Name: "UsprcCountyNm"},
	), nil
}

// Find a corresponding county for a provided zip code from the USPRC table
func FindCountyByZipCode(db *pop.Connection, zipCode string) (string, error) {
	var usprc UsPostRegionCity
	err := db.Where("uspr_zip_id = ?", zipCode).First(&usprc)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", fmt.Errorf("No county found for provided zip code %s", zipCode)
		default:
			return "", err
		}
	}
	return usprc.UsprcCountyNm, nil
}
