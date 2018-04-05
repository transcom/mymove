package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"time"
)

// Tariff400ngZip3 is the first 3 numbers of a zip for Tariff400NG calculations
type Tariff400ngZip3 struct {
	ID            uuid.UUID `json:"id" db:"id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	Zip3          string    `json:"zip3" db:"zip3"`
	BasepointCity string    `json:"basepoint_city" db:"basepoint_city"`
	State         string    `json:"state" db:"state"`
	ServiceArea   int       `json:"service_area" db:"service_area"`
	RateArea      int       `json:"rate_area" db:"rate_area"`
	Region        int       `json:"region" db:"region"`
}

// String is not required by pop and may be deleted
func (t Tariff400ngZip3) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Tariff400ngZip3s is not required by pop and may be deleted
type Tariff400ngZip3s []Tariff400ngZip3

// String is not required by pop and may be deleted
func (t Tariff400ngZip3s) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Tariff400ngZip3) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: t.Zip3, Name: "Zip3"},
		&validators.StringIsPresent{Field: t.BasepointCity, Name: "BasepointCity"},
		&validators.StringIsPresent{Field: t.State, Name: "State"},
		&validators.IntIsPresent{Field: t.ServiceArea, Name: "ServiceArea"},
		&validators.IntIsPresent{Field: t.RateArea, Name: "RateArea"},
		&validators.IntIsPresent{Field: t.Region, Name: "Region"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Tariff400ngZip3) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Tariff400ngZip3) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
