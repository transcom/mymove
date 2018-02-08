package models

import (
	"encoding/json"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
	"time"
)

// TransportationServiceProvider models moving companies used to move
// Shipments.
type TransportationServiceProvider struct {
	ID                       uuid.UUID `json:"id" db:"id"`
	CreatedAt                time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time `json:"updated_at" db:"updated_at"`
	StandardCarrierAlphaCode string    `json:"standard_carrier_alpha_code" db:"standard_carrier_alpha_code"`
	Name                     string    `json:"name" db:"name"`
}

// String is not required by pop and may be deleted
func (t TransportationServiceProvider) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// TransportationServiceProviders is not required by pop and may be deleted
type TransportationServiceProviders []TransportationServiceProvider

// String is not required by pop and may be deleted
func (t TransportationServiceProviders) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *TransportationServiceProvider) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: t.StandardCarrierAlphaCode, Name: "StandardCarrierAlphaCode"},
		&validators.StringIsPresent{Field: t.Name, Name: "Name"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *TransportationServiceProvider) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *TransportationServiceProvider) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
