package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// Tariff400ngLinehaulRate describes the rate paids paid to transport various weights of goods
// various distances.
type Tariff400ngLinehaulRate struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	DistanceMilesLower int       `json:"distance_miles_lower" db:"distance_miles_lower"`
	DistanceMilesUpper int       `json:"distance_miles_upper" db:"distance_miles_upper"`
	Type               string    `json:"type" db:"type"`
	WeightLbsLower     int       `json:"weight_lbs_lower" db:"weight_lbs_lower"`
	WeightLbsUpper     int       `json:"weight_lbs_upper" db:"weight_lbs_upper"`
	RateCents          int       `json:"rate_cents" db:"rate_cents"`
	EffectiveDateLower time.Time `json:"effective_date_lower" db:"effective_date_lower"`
	EffectiveDateUpper time.Time `json:"effective_date_upper" db:"effective_date_upper"`
}

// String is not required by pop and may be deleted
func (t Tariff400ngLinehaulRate) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Tariff400ngLinehaulRates is not required by pop and may be deleted
type Tariff400ngLinehaulRates []Tariff400ngLinehaulRate

// String is not required by pop and may be deleted
func (t Tariff400ngLinehaulRates) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Tariff400ngLinehaulRate) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsGreaterThan{Field: t.RateCents, Name: "RateCents", Compared: -1},
		&validators.IntIsLessThan{Field: t.DistanceMilesLower, Name: "DistanceMilesLower",
			Compared: t.DistanceMilesUpper},
		&validators.IntIsLessThan{Field: t.WeightLbsLower, Name: "WeightLbsLower",
			Compared: t.WeightLbsUpper},
		&validators.TimeAfterTime{
			FirstTime: t.EffectiveDateUpper, FirstName: "EffectiveDateUpper",
			SecondTime: t.EffectiveDateLower, SecondName: "EffectiveDateLower"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Tariff400ngLinehaulRate) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Tariff400ngLinehaulRate) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// FetchBaseLinehaulRate takes a move's distance and weight and queries the tariff400ng_linehaul_rates table to find a move's base linehaul rate.
func FetchBaseLinehaulRate(tx *pop.Connection, mileage int, cwt int, date time.Time) (linehaulRate int, err error) {
	moveType := "ConusLinehaul" // TODO: change to a parameter once we're serving more move types
	fmt.Printf("FetchBaseLinehaulRate db: %v \n", tx)

	sql := `SELECT
		rate_cents
	FROM
		tariff400ng_linehaul_rates
	WHERE
		distance_miles_lower <= $1
	AND
		distance_miles_upper > $1
	AND
		weight_lbs_lower <= $2
	AND
		weight_lbs_upper > $2
	AND
		type = $3
	AND
		effective_date_lower <= $4
	AND
		effective_date_upper > $4;

	`
	err = tx.RawQuery(sql, mileage, (cwt * 100), moveType, date).First(&linehaulRate)

	fmt.Println("It's the linehaul rate int: ", linehaulRate)
	return linehaulRate, err

}
