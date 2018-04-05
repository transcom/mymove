package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// Tariff400ngFullUnpackRate describes the rates paid to unpack various weights of goods
type Tariff400ngFullUnpackRate struct {
	ID                 uuid.UUID `db:"id"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
	Schedule           int       `db:"schedule"`
	RateMillicents     int       `db:"rate_millicents"`
	EffectiveDateLower time.Time `db:"effective_date_lower"`
	EffectiveDateUpper time.Time `db:"effective_date_upper"`
}

// Tariff400ngFullUnpackRates is not required by pop and may be deleted
type Tariff400ngFullUnpackRates []Tariff400ngFullUnpackRate

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Tariff400ngFullUnpackRate) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsGreaterThan{Field: t.RateMillicents, Name: "RateMillicents", Compared: -1},
		&validators.TimeAfterTime{
			FirstTime: t.EffectiveDateUpper, FirstName: "EffectiveDateUpper",
			SecondTime: t.EffectiveDateLower, SecondName: "EffectiveDateLower"},
	), nil
}
