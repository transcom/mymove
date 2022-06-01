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

// Tariff400ngShorthaulRate describes the rates paid for shorthaul shipments
type Tariff400ngShorthaulRate struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	CwtMilesLower      int        `json:"cwt_miles_lower" db:"cwt_miles_lower"`
	CwtMilesUpper      int        `json:"cwt_miles_upper" db:"cwt_miles_upper"`
	RateCents          unit.Cents `json:"rate_cents" db:"rate_cents"`
	EffectiveDateLower time.Time  `json:"effective_date_lower" db:"effective_date_lower"`
	EffectiveDateUpper time.Time  `json:"effective_date_upper" db:"effective_date_upper"`
}

// Tariff400ngShorthaulRates is not required by pop and may be deleted
type Tariff400ngShorthaulRates []Tariff400ngShorthaulRate

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Tariff400ngShorthaulRate) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsGreaterThan{Field: t.RateCents.Int(), Name: "RateCents", Compared: -1},
		&validators.IntIsGreaterThan{Field: t.CwtMilesUpper, Name: "CwtMilesUpper",
			Compared: t.CwtMilesLower},
		&validators.IntIsPresent{Field: t.CwtMilesUpper, Name: "CwtMilesUpper"},
		&validators.TimeIsPresent{Field: t.EffectiveDateLower, Name: "EffectiveDateLower"},
		&validators.TimeIsPresent{Field: t.EffectiveDateUpper, Name: "EffectiveDateUpper"},
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

// FetchShorthaulRateCents returns the shorthaul rate for a given Centumweight-Miles
// (cwtMiles is a unit capturing the movement of 100lbs by 1 mile.) The value returned
// is in cents of 1 USD.
func FetchShorthaulRateCents(tx *pop.Connection, cwtMiles int, date time.Time) (rateCents unit.Cents, err error) {
	sh := Tariff400ngShorthaulRates{}

	sql := `SELECT
		rate_cents
	FROM
		tariff400ng_shorthaul_rates
	WHERE
		cwt_miles_lower <= $1 AND $1 < cwt_miles_upper
	AND
		effective_date_lower <= $2 AND $2 < effective_date_upper`

	err = tx.RawQuery(sql, cwtMiles, date).All(&sh)
	if err != nil {
		return 0, errors.Wrapf(err, "error fetching shorthaul rate for %d cwtmiles on %s", cwtMiles, date)
	}
	if len(sh) != 1 {
		return 0, errors.Errorf("Wanted 1 shorthaul rate, found %d rates for parameters: %v cwtMiles, %v",
			len(sh), cwtMiles, date)
	}

	return sh[0].RateCents, nil
}
