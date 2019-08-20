package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// Tariff400ngItemRate contains pricing data for a Tariff400ngItem
type Tariff400ngItemRate struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	Code               string     `json:"code" db:"code"`
	Schedule           *int       `json:"schedule" db:"schedule"`
	WeightLbsLower     unit.Pound `json:"weight_lbs_lower" db:"weight_lbs_lower"`
	WeightLbsUpper     unit.Pound `json:"weight_lbs_upper" db:"weight_lbs_upper"`
	RateCents          unit.Cents `json:"rate_cents" db:"rate_cents"`
	EffectiveDateLower time.Time  `json:"effective_date_lower" db:"effective_date_lower"`
	EffectiveDateUpper time.Time  `json:"effective_date_upper" db:"effective_date_upper"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Tariff400ngItemRate) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: t.Code, Name: "Code"},
		&validators.IntIsGreaterThan{Field: t.RateCents.Int(), Name: "RateCents", Compared: -1},
		&validators.IntIsLessThan{Field: t.WeightLbsLower.Int(), Name: "WeightLbsLower",
			Compared: t.WeightLbsUpper.Int()},
		&validators.TimeAfterTime{
			FirstTime: t.EffectiveDateUpper, FirstName: "EffectiveDateUpper",
			SecondTime: t.EffectiveDateLower, SecondName: "EffectiveDateLower"},
	), nil
}

// FetchTariff400ngItemRate returns a rate for a matching code, schedule, weight, and ship date
func FetchTariff400ngItemRate(tx *pop.Connection, code string, schedule int, weight unit.Pound, shipDate time.Time) (Tariff400ngItemRate, error) {
	var rate Tariff400ngItemRate
	query := `
		SELECT * from tariff400ng_item_rates
		WHERE
			code = $1
			AND (schedule = $2 OR schedule IS NULL)
			AND weight_lbs_lower <= $3
			AND weight_lbs_upper >= $3
			AND effective_date_lower <= $4
			AND effective_date_upper > $4
	`

	err := tx.RawQuery(query, code, schedule, weight.Int(), shipDate).First(&rate)

	return rate, err
}
