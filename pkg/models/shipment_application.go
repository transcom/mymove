package models

import (
	"encoding/json"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
	"time"
)

// ShipmentApplication is the form 1299, submitted by service members creating or updating an application.
type ShipmentApplication struct {
	ID                    uuid.UUID `json:"id" db:"id"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
	NameOfPreparingOffice string    `json:"name_of_preparing_office" db:"name_of_preparing_office"`
}

// String is not required by pop and may be deleted
func (s ShipmentApplication) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// ShipmentApplications is not required by pop and may be deleted
type ShipmentApplications []ShipmentApplication

// String is not required by pop and may be deleted
func (s ShipmentApplications) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *ShipmentApplication) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: s.NameOfPreparingOffice, Name: "NameOfPreparingOffice"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *ShipmentApplication) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *ShipmentApplication) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
