package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// ReZip3 model struct
type ReZip3 struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	ContractID            uuid.UUID  `json:"contract_id" db:"contract_id"`
	Zip3                  string     `json:"zip3" db:"zip3"`
	BasePointCity         string     `json:"base_point_city" db:"base_point_city"`
	State                 string     `json:"state" db:"state"`
	DomesticServiceAreaID uuid.UUID  `json:"domestic_service_area_id" db:"domestic_service_area_id"`
	RateAreaID            *uuid.UUID `json:"rate_area_id" db:"rate_area_id"`
	HasMultipleRateAreas  bool       `json:"has_multiple_rate_areas" db:"has_multiple_rate_areas"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`

	// Associations
	Contract            ReContract            `belongs_to:"re_contract" fk_id:"contract_id"`
	DomesticServiceArea ReDomesticServiceArea `belongs_to:"re_domestic_service_areas" fk_id:"domestic_service_area_id"`
	RateArea            *ReRateArea           `belongs_to:"re_rate_areas" fk_id:"rate_area_id"`
}

// ReZip3s is not required by pop and may be deleted
type ReZip3s []ReZip3

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *ReZip3) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringLengthInRange{Field: r.Zip3, Name: "Zip3", Min: 3, Max: 3},
		&validators.UUIDIsPresent{Field: r.ContractID, Name: "ContractID"},
		&validators.StringIsPresent{Field: r.BasePointCity, Name: "BasePointCity"},
		&validators.StringIsPresent{Field: r.State, Name: "State"},
		&validators.UUIDIsPresent{Field: r.DomesticServiceAreaID, Name: "DomesticServiceAreaID"},
	), nil
}
