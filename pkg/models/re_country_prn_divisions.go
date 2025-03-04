package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type ReCountryPrnDivisions struct {
	ID             uuid.UUID `json:"id" db:"id"`
	CountryID      uuid.UUID `json:"country_id" db:"country_id"`
	Country        Country   `belongs_to:"re_countries" fk_id:"country_id"`
	CountryPrnDvID string    `json:"country_prn_dv_id" db:"country_prn_dv_id"`
	CountryPrnDvNm string    `json:"country_prn_dv_nm" db:"country_prn_dv_nm"`
	CountryPrnDvCd string    `json:"country_prn_dv_cd" db:"country_prn_dv_cd"`
	CommandOrgCd   *string   `json:"command_org_cd" db:"command_org_cd"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (r ReCountryPrnDivisions) TableName() string {
	return "re_country_prn_divisions"
}
