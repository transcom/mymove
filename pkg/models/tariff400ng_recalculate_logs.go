package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
)

// Tariff400ngRecalculateLog model/struct is a record used to log that a Shipment was
// recalculated
type Tariff400ngRecalculateLog struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	ShipmentID uuid.UUID `json:"shipment_id" db:"shipment_id"`
	Shipment   Shipment  `belongs_to:"shipments"`
}

// String is not required by pop and may be deleted
func (t Tariff400ngRecalculateLog) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Tariff400ngRecalculateLogs store logs
type Tariff400ngRecalculateLogs []Tariff400ngRecalculateLog

// String is not required by pop and may be deleted
func (t Tariff400ngRecalculateLogs) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Tariff400ngRecalculateLog) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Tariff400ngRecalculateLog) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Tariff400ngRecalculateLog) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// SaveTariff400ngRecalculateLog validates and saves the log record
func (t *Tariff400ngRecalculateLog) SaveTariff400ngRecalculateLog(db *pop.Connection) (*validate.Errors, error) {
	verrs, err := db.ValidateAndSave(t)
	if verrs.HasAny() || err != nil {
		saveError := errors.Wrap(err, "Error saving re-calculate log")
		return verrs, saveError
	}
	return verrs, nil
}
