package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// Tariff400ngShorthaulRate describes the rates paid for shorthaul shipments
type Tariff400ngShorthaulRate struct {
	ID                 uuid.UUID `db:"id"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
	CwtMilesLower      int       `db:"cwt_miles_lower"`
	CwtMilesUpper      int       `db:"cwt_miles_upper"`
	RateCents          int       `db:"rate_cents"`
	EffectiveDateLower time.Time `db:"effective_date_lower"`
	EffectiveDateUpper time.Time `db:"effective_date_upper"`
}

// Tariff400ngShorthaulRates is not required by pop and may be deleted
type Tariff400ngShorthaulRates []Tariff400ngShorthaulRate

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Tariff400ngShorthaulRate) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsGreaterThan{Field: t.RateCents, Name: "ServiceChargeCents", Compared: -1},
		&validators.IntIsGreaterThan{Field: t.CwtMilesUpper, Name: "CwtMilesUpper",
			Compared: t.CwtMilesLower},
		&validators.TimeAfterTime{
			FirstTime: t.EffectiveDateUpper, FirstName: "EffectiveDateUpper",
			SecondTime: t.EffectiveDateLower, SecondName: "EffectiveDateLower"},
	), nil
}
