package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// Tariff400ngFullPackRate describes the rates paid to pack various weights of goods
type Tariff400ngFullPackRate struct {
	ID                 uuid.UUID `db:"id"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
	Schedule           int       `db:"schedule"`
	WeightLbsLower     int       `db:"weight_lbs_lower"`
	WeightLbsUpper     int       `db:"weight_lbs_upper"`
	RateCents          int       `db:"rate_cents"`
	EffectiveDateLower time.Time `db:"effective_date_lower"`
	EffectiveDateUpper time.Time `db:"effective_date_upper"`
}

// Tariff400ngFullPackRates is not required by pop and may be deleted
type Tariff400ngFullPackRates []Tariff400ngFullPackRate

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Tariff400ngFullPackRate) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsGreaterThan{Field: t.RateCents, Name: "RateCents", Compared: -1},
		&validators.IntIsLessThan{Field: t.WeightLbsLower, Name: "WeightLbsLower",
			Compared: t.WeightLbsUpper},
		&validators.TimeAfterTime{
			FirstTime: t.EffectiveDateUpper, FirstName: "EffectiveDateUpper",
			SecondTime: t.EffectiveDateLower, SecondName: "EffectiveDateLower"},
	), nil
}
