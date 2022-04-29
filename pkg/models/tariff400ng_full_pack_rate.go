package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/unit"
)

// Tariff400ngFullPackRate describes the rates paid to pack various weights of goods
type Tariff400ngFullPackRate struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	Schedule           int        `json:"schedule" db:"schedule"`
	WeightLbsLower     unit.Pound `json:"weight_lbs_lower" db:"weight_lbs_lower"`
	WeightLbsUpper     unit.Pound `json:"weight_lbs_upper" db:"weight_lbs_upper"`
	RateCents          unit.Cents `json:"rate_cents" db:"rate_cents"`
	EffectiveDateLower time.Time  `json:"effective_date_lower" db:"effective_date_lower"`
	EffectiveDateUpper time.Time  `json:"effective_date_upper" db:"effective_date_upper"`
}

// Tariff400ngFullPackRates is not required by pop and may be deleted
type Tariff400ngFullPackRates []Tariff400ngFullPackRate

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Tariff400ngFullPackRate) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: t.Schedule, Name: "Schedule"},
		&validators.IntIsPresent{Field: t.WeightLbsUpper.Int(), Name: "WeightLbsUpper"},
		&validators.IntIsLessThan{Field: t.WeightLbsLower.Int(), Name: "WeightLbsLower",
			Compared: t.WeightLbsUpper.Int()},
		&validators.IntIsGreaterThan{Field: t.RateCents.Int(), Name: "RateCents", Compared: -1},
		&validators.TimeIsPresent{Field: t.EffectiveDateLower, Name: "EffectiveDateLower"},
		&validators.TimeIsPresent{Field: t.EffectiveDateUpper, Name: "EffectiveDateUpper"},
		&validators.TimeAfterTime{
			FirstTime: t.EffectiveDateUpper, FirstName: "EffectiveDateUpper",
			SecondTime: t.EffectiveDateLower, SecondName: "EffectiveDateLower"},
	), nil
}

// FetchTariff400ngFullPackRateCents returns the full unpack rate for a service
// schedule and weight.
func FetchTariff400ngFullPackRateCents(tx *pop.Connection, weight unit.Pound, schedule int, date time.Time) (unit.Cents, error) {
	rate := Tariff400ngFullPackRate{}

	sql := `SELECT
			*
		FROM
			tariff400ng_full_pack_rates
		WHERE
			schedule = $1
		AND
			weight_lbs_lower <= $2 AND $2 < weight_lbs_upper
		AND
			effective_date_lower <= $3 AND $3 < effective_date_upper
		;
		`

	err := tx.RawQuery(sql, schedule, weight, date).First(&rate)
	if err != nil {
		return 0, errors.Wrap(err, "could not find a matching Tariff400ngFullPackRate")
	}
	return rate.RateCents, nil
}
