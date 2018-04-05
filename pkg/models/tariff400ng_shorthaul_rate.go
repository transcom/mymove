package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// Tariff400ngShorthaulRate describes the rate paid for a certain weight range for a certain
// schedule
type Tariff400ngShorthaulRate struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	CwtMilesLower      int       `json:"cwt_miles_lower" db:"cwt_miles_lower"`
	CwtMilesUpper      int       `json:"cwt_miles_upper" db:"cwt_miles_upper"`
	RateCents          int       `json:"rate_cents" db:"rate_cents"`
	EffectiveDateLower time.Time `json:"effective_date_lower" db:"effective_date_lower"`
	EffectiveDateUpper time.Time `json:"effective_date_upper" db:"effective_date_upper"`
}

// String is not required by pop and may be deleted
func (t Tariff400ngShorthaulRate) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Tariff400ngShorthaulRates is not required by pop and may be deleted
type Tariff400ngShorthaulRates []Tariff400ngShorthaulRate

// String is not required by pop and may be deleted
func (t Tariff400ngShorthaulRates) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

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

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Tariff400ngShorthaulRate) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Tariff400ngShorthaulRate) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
