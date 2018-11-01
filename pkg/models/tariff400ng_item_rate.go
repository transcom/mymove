package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/transcom/mymove/pkg/unit"
)

// Tariff400ngItemRate contains pricing data for a Tariff400ngItem
type Tariff400ngItemRate struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	Code               string     `json:"code" db:"code"`
	ServicesSchedule   int        `json:"services_schedule" db:"services_schedule"`
	WeightLbsLower     unit.Pound `json:"weight_lbs_lower" db:"weight_lbs_lower"`
	WeightLbsUpper     unit.Pound `json:"weight_lbs_upper" db:"weight_lbs_upper"`
	RateMillicents     int        `json:"rate_millicents" db:"rate_millicents"`
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
		&validators.IntIsPresent{Field: t.ServicesSchedule, Name: "ServicesSchedule"},
		&validators.IntIsGreaterThan{Field: t.RateMillicents, Name: "RateMillicents", Compared: -1},
		&validators.IntIsLessThan{Field: t.WeightLbsLower.Int(), Name: "WeightLbsLower",
			Compared: t.WeightLbsUpper.Int()},
		&validators.TimeAfterTime{
			FirstTime: t.EffectiveDateUpper, FirstName: "EffectiveDateUpper",
			SecondTime: t.EffectiveDateLower, SecondName: "EffectiveDateLower"},
	), nil
}
