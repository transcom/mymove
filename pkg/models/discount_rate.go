package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// DiscountRate describes how great a discount a TSP will provide for SIT and Linehaul
// for a given TDL/rate cycle.
type DiscountRate struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	RateCycle          string    `json:"rate_cycle" db:"rate_cycle"`
	Origin             string    `json:"origin" db:"origin"`
	Destination        string    `json:"destination" db:"destination"`
	CodeOfService      string    `json:"code_of_service" db:"code_of_service"`
	Scac               string    `json:"scac" db:"scac"`
	LhRate             float64   `json:"lh_rate" db:"lh_rate"`
	SitRate            float64   `json:"sit_rate" db:"sit_rate"`
	EffectiveDateLower time.Time `json:"effective_date_lower" db:"effective_date_lower"`
	EffectiveDateUpper time.Time `json:"effective_date_upper" db:"effective_date_upper"`
}

// String is not required by pop and may be deleted
func (d DiscountRate) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// DiscountRates is not required by pop and may be deleted
type DiscountRates []DiscountRate

// String is not required by pop and may be deleted
func (d DiscountRates) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (d *DiscountRate) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.TimeIsPresent{Field: d.EffectiveDateUpper, Name: "EffectiveDateUpper"},
		&validators.TimeAfterTime{
			FirstTime: d.EffectiveDateUpper, FirstName: "EffectiveDateUpper",
			SecondTime: d.EffectiveDateLower, SecondName: "EffectiveDateLower"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (d *DiscountRate) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (d *DiscountRate) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
