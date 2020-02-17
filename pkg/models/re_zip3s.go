package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// ReZip3 model struct
type ReZip3 struct {
	ID                    uuid.UUID `json:"id" db:"id"`
	ContractID            uuid.UUID `json:"contract_id" db:"contract_id"`
	Zip3                  string    `json:"zip3" db:"zip3"`
	DomesticServiceAreaID uuid.UUID `json:"domestic_service_area_id" db:"domestic_service_area_id"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`

	// Associations
	Contract            ReContract            `belongs_to:"re_contract"`
	DomesticServiceArea ReDomesticServiceArea `belongs_to:"re_domestic_service_area"`
}

// ReZip3s is not required by pop and may be deleted
type ReZip3s []ReZip3

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *ReZip3) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ContractID, Name: "ContractID"},
		&validators.UUIDIsPresent{Field: r.DomesticServiceAreaID, Name: "DomesticServiceAreaID"},
		&validators.StringLengthInRange{Field: r.Zip3, Name: "Zip3", Min: 3, Max: 3},
	), nil
}
