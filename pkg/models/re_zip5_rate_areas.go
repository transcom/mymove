package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// ReZip5RateArea model struct
type ReZip5RateArea struct {
	ID         uuid.UUID `json:"id" db:"id"`
	ContractID uuid.UUID `json:"contract_id" db:"contract_id"`
	Zip5       string    `json:"zip5" db:"zip5"`
	RateAreaID uuid.UUID `json:"rate_area_id" db:"rate_area_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`

	// Associations
	Contract ReContract `belongs_to:"re_contract" fk_id:"contract_id"`
	RateArea ReRateArea `belongs_to:"re_rate_areas" fk_id:"rate_area_id"`
}

// TableName overrides the table name used by Pop.
func (r ReZip5RateArea) TableName() string {
	return "re_zip5_rate_areas"
}

type ReZip5RateAreas []ReZip5RateArea

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (r *ReZip5RateArea) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ContractID, Name: "ContractID"},
		&validators.UUIDIsPresent{Field: r.RateAreaID, Name: "RateAreaID"},
		&validators.StringLengthInRange{Field: r.Zip5, Name: "Zip5", Min: 5, Max: 5},
	), nil
}
