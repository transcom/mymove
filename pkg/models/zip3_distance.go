package models

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Zip3Distance model struct
type Zip3Distance struct {
	ID            uuid.UUID `json:"id" db:"id"`
	FromZip3      string    `json:"from_zip3" db:"from_zip3"`
	ToZip3        string    `json:"to_zip3" db:"to_zip3"`
	DistanceMiles int       `json:"distance_miles" db:"distance_miles"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Zip3Distances is not required by pop and may be deleted
type Zip3Distances []Zip3Distances

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (z *Zip3Distance) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringLengthInRange{Field: z.FromZip3, Name: "FromZip3", Min: 3, Max: 3},
		&validators.StringLengthInRange{Field: z.ToZip3, Name: "ToZip3", Min: 3, Max: 3},
		&validators.IntIsPresent{Field: z.DistanceMiles, Name: "DistanceMiles"},
	), nil
}
