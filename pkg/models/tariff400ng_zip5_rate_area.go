package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Tariff400ngZip5RateArea represents the mapping from a full 5-digit zipcode to a
// specific rate area. This is only needed for a small subset of zip3s.
type Tariff400ngZip5RateArea struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Zip5      string    `json:"zip5" db:"zip5"`
	RateArea  string    `json:"rate_area" db:"rate_area"`
}

// Tariff400ngZip5RateAreas is not required by pop and may be deleted
type Tariff400ngZip5RateAreas []Tariff400ngZip5RateArea

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (t *Tariff400ngZip5RateArea) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringLengthInRange{Field: t.Zip5, Name: "Zip5", Min: 5, Max: 5},
		&validators.StringIsPresent{Field: t.RateArea, Name: "RateArea"},
		&validators.RegexMatch{Field: t.RateArea, Name: "RateArea", Expr: "^(ZIP|US[0-9]+)$"},
	), nil
}
