package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
)

// ShipmentRecalculateLog model/struct is a record used to log that a Shipment was
// recalculated
type ShipmentRecalculateLog struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	ShipmentID uuid.UUID `json:"shipment_id" db:"shipment_id"`
	Shipment   Shipment  `belongs_to:"shipments"`
}

// String is not required by pop and may be deleted
func (t ShipmentRecalculateLog) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// ShipmentRecalculateLogs store logs
type ShipmentRecalculateLogs []ShipmentRecalculateLog

// String is not required by pop and may be deleted
func (t ShipmentRecalculateLogs) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *ShipmentRecalculateLog) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: t.ShipmentID, Name: "ShipmentID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *ShipmentRecalculateLog) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *ShipmentRecalculateLog) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// SaveShipmentRecalculateLog validates and saves the log record
func (t *ShipmentRecalculateLog) SaveShipmentRecalculateLog(db *pop.Connection) (*validate.Errors, error) {
	verrs, err := db.ValidateAndSave(t)
	if verrs.HasAny() || err != nil {
		saveError := errors.Wrap(err, "Error saving re-calculate log")
		return verrs, saveError
	}
	return verrs, nil
}
