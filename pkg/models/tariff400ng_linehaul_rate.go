package models

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// Tariff400ngLinehaulRate describes the rate paids paid to transport various weights of goods
// various distances.
type Tariff400ngLinehaulRate struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	DistanceMilesLower int        `json:"distance_miles_lower" db:"distance_miles_lower"`
	DistanceMilesUpper int        `json:"distance_miles_upper" db:"distance_miles_upper"`
	Type               string     `json:"type" db:"type"`
	WeightLbsLower     unit.Pound `json:"weight_lbs_lower" db:"weight_lbs_lower"`
	WeightLbsUpper     unit.Pound `json:"weight_lbs_upper" db:"weight_lbs_upper"`
	RateCents          unit.Cents `json:"rate_cents" db:"rate_cents"`
	EffectiveDateLower time.Time  `json:"effective_date_lower" db:"effective_date_lower"`
	EffectiveDateUpper time.Time  `json:"effective_date_upper" db:"effective_date_upper"`
}

// Tariff400ngLinehaulRates is not required by pop and may be deleted
type Tariff400ngLinehaulRates []Tariff400ngLinehaulRate

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Tariff400ngLinehaulRate) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: t.DistanceMilesLower, Name: "DistanceMilesLower"},
		&validators.IntIsPresent{Field: t.DistanceMilesUpper, Name: "DistanceMilesUpper"},
		&validators.IntIsLessThan{Field: t.DistanceMilesLower, Name: "DistanceMilesLower",
			Compared: t.DistanceMilesUpper},
		&validators.StringIsPresent{Field: t.Type, Name: "Type"},
		&validators.IntIsPresent{Field: t.WeightLbsLower.Int(), Name: "WeightLbsLower"},
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

// FetchBaseLinehaulRate takes a move's distance and weight and queries the tariff400ng_linehaul_rates table to find a move's base linehaul rate.
func FetchBaseLinehaulRate(tx *pop.Connection, mileage int, weight unit.Pound, date time.Time) (linehaulRate unit.Cents, err error) {
	// TODO: change to a parameter once we're serving more move types
	moveType := "ConusLinehaul"
	var linehaulRates []unit.Cents

	sql := `SELECT
		rate_cents
	FROM
		tariff400ng_linehaul_rates
	WHERE
		(distance_miles_lower <= $1 AND $1 < distance_miles_upper)
	AND
		(weight_lbs_lower <= $2 AND $2 < weight_lbs_upper)
	AND
		type = $3
	AND
		(effective_date_lower <= $4 AND $4 < effective_date_upper);`

	err = tx.RawQuery(sql, mileage, weight.Int(), moveType, date).All(&linehaulRates)

	if err != nil {
		return 0, fmt.Errorf("Error fetching linehaul rate: %s", err)
	}
	if len(linehaulRates) != 1 {
		return 0, fmt.Errorf("Wanted 1 rate, found %d rates for parameters: %v, %v, %v",
			len(linehaulRates), mileage, weight, date)
	}

	return linehaulRates[0], err
}
