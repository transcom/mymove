package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// Tariff400ngFullPackRate describes the rate paid for a certain weight range for a certain
// schedule
type Tariff400ngFullPackRate struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	Schedule           int       `json:"schedule" db:"schedule"`
	WeightLbsLower     int       `json:"weight_lbs_lower" db:"weight_lbs_lower"`
	WeightLbsUpper     int       `json:"weight_lbs_upper" db:"weight_lbs_upper"`
	RateCents          int       `json:"rate_cents" db:"rate_cents"`
	EffectiveDateLower time.Time `json:"effective_date_lower" db:"effective_date_lower"`
	EffectiveDateUpper time.Time `json:"effective_date_upper" db:"effective_date_upper"`
}

// String is not required by pop and may be deleted
func (t Tariff400ngFullPackRate) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Tariff400ngFullPackRates is not required by pop and may be deleted
type Tariff400ngFullPackRates []Tariff400ngFullPackRate

// String is not required by pop and may be deleted
func (t Tariff400ngFullPackRates) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Tariff400ngFullPackRate) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsGreaterThan{Field: t.RateCents, Name: "RateCents", Compared: -1},
		&validators.IntIsLessThan{Field: t.WeightLbsLower, Name: "WeightLbsLower",
			Compared: t.WeightLbsUpper + 1},
		&validators.TimeAfterTime{
			FirstTime: t.EffectiveDateUpper, FirstName: "EffectiveDateUpper",
			SecondTime: t.EffectiveDateLower, SecondName: "EffectiveDateLower"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Tariff400ngFullPackRate) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Tariff400ngFullPackRate) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
