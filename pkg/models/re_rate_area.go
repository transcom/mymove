package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// ReRateArea model struct
type ReRateArea struct {
	ID        uuid.UUID `json:"id" db:"id"`
	IsOconus  bool      `json:"is_oconos" db:"is_oconos"`
	Code      string    `json:"code" db:"code"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ReRateAreas is not required by pop and may be deleted
type ReRateAreas []ReRateArea

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *ReRateArea) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: r.Code, Name: "Code"},
		&validators.StringIsPresent{Field: r.Name, Name: "Name"},
	), nil
}

// FetchReRateAreaItem returns an area for a matching code
func FetchReRateAreaItem(tx *pop.Connection, code string) (*ReRateArea, error) {
	var area ReRateArea
	query := `
		SELECT * from re_rate_area
		WHERE
			code = $1
	`
	err := tx.RawQuery(query, code).First(&area)

	if err != nil {
		return nil, err
	}

	if area.Code == code {
		return &area, err
	}

	return nil, err
}