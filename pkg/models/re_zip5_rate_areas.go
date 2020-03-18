package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// ReZip5RateArea model struct
type ReZip5RateArea struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Zip5       string    `json:"zip5" db:"zip5"`
	RateAreaID uuid.UUID `json:"rate_area_id" db:"rate_area_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`

	// Associations
	RateArea ReRateArea `belongs_to:"re_rate_areas"`
}

// ReZip5RateAreas is not required by pop and may be deleted
type ReZip5RateAreas []ReZip5RateArea

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *ReZip5RateArea) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.RateAreaID, Name: "RateAreaID"},
		&validators.StringLengthInRange{Field: r.Zip5, Name: "Zip5", Min: 5, Max: 5},
	), nil
}
