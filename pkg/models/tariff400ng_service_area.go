package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// Tariff400ngServiceArea describes the rate paid for a certain weight range for a certain
// schedule
type Tariff400ngServiceArea struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	Name               string    `json:"name" db:"name"`
	ServiceArea        int       `json:"service_area" db:"service_area"`
	LinehaulFactor     int       `json:"linehaul_factor" db:"linehaul_factor"`
	ServiceChargeCents int       `json:"service_charge_cents" db:"service_charge_cents"`
	EffectiveDateLower time.Time `json:"effective_date_lower" db:"effective_date_lower"`
	EffectiveDateUpper time.Time `json:"effective_date_upper" db:"effective_date_upper"`
}

// String is not required by pop and may be deleted
func (t Tariff400ngServiceArea) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Tariff400ngServiceAreas is not required by pop and may be deleted
type Tariff400ngServiceAreas []Tariff400ngServiceArea

// String is not required by pop and may be deleted
func (t Tariff400ngServiceAreas) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Tariff400ngServiceArea) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsGreaterThan{Field: t.ServiceChargeCents, Name: "ServiceChargeCents", Compared: -1},
		&validators.TimeAfterTime{
			FirstTime: t.EffectiveDateUpper, FirstName: "EffectiveDateUpper",
			SecondTime: t.EffectiveDateLower, SecondName: "EffectiveDateLower"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Tariff400ngServiceArea) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Tariff400ngServiceArea) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
