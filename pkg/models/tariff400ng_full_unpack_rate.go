package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// Tariff400ngFullUnpackRate describes the rates paid to unpack various weights of goods
type Tariff400ngFullUnpackRate struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	Schedule           int       `json:"schedule" db:"schedule"`
	RateMillicents     int       `json:"rate_millicents" db:"rate_millicents"`
	EffectiveDateLower time.Time `json:"effective_date_lower" db:"effective_date_lower"`
	EffectiveDateUpper time.Time `json:"effective_date_upper" db:"effective_date_upper"`
}

// Tariff400ngFullUnpackRates is not required by pop and may be deleted
type Tariff400ngFullUnpackRates []Tariff400ngFullUnpackRate

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Tariff400ngFullUnpackRate) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: t.Schedule, Name: "Schedule"},
		&validators.IntIsGreaterThan{Field: t.RateMillicents, Name: "RateMillicents", Compared: -1},
		&validators.TimeIsPresent{Field: t.EffectiveDateLower, Name: "EffectiveDateLower"},
		&validators.TimeIsPresent{Field: t.EffectiveDateUpper, Name: "EffectiveDateUpper"},
		&validators.TimeAfterTime{
			FirstTime: t.EffectiveDateUpper, FirstName: "EffectiveDateUpper",
			SecondTime: t.EffectiveDateLower, SecondName: "EffectiveDateLower"},
	), nil
}

// FetchTariff400ngFullUnpackRateMillicents returns the full unpack rate for a service
// schedule.
func FetchTariff400ngFullUnpackRateMillicents(tx *pop.Connection, serviceSchedule int, date time.Time) (int, error) {
	rate := Tariff400ngFullUnpackRate{}

	sql := `SELECT *
		FROM
			tariff400ng_full_unpack_rates
		WHERE
			schedule = $1
		AND
			effective_date_lower <= $2 AND $2 < effective_date_upper
		;`

	err := tx.RawQuery(sql, serviceSchedule, date).First(&rate)

	if err != nil {
		return 0, errors.Wrap(err, "could not find a matching Tariff400ngFullUnpackRate")
	}
	return rate.RateMillicents, nil
}
