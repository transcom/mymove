package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// ShipmentRecalculate model/struct is a record of the date range to use to determine if a Shipment is a candidate
// to be re-calculated.
type ShipmentRecalculate struct {
	ID                    uuid.UUID `json:"id" db:"id"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
	ShipmentUpdatedBefore time.Time `json:"shipment_updated_before" db:"shipment_updated_before"`
	ShipmentUpdatedAfter  time.Time `json:"shipment_updated_after" db:"shipment_updated_after"`
	Active                bool      `json:"active" db:"active"`
}

// String is not required by pop and may be deleted
func (t ShipmentRecalculate) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// ShipmentRecalculates is not required by pop and may be deleted
type ShipmentRecalculates []ShipmentRecalculate

// String is not required by pop and may be deleted
func (t ShipmentRecalculates) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *ShipmentRecalculate) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.TimeIsPresent{Field: t.ShipmentUpdatedBefore, Name: "ShipmentUpdatedBefore"},
		&validators.TimeIsPresent{Field: t.ShipmentUpdatedAfter, Name: "ShipmentUpdatedAfter"},
		&validators.TimeAfterTime{
			FirstTime: t.ShipmentUpdatedBefore, FirstName: "ShipmentUpdatedBefore",
			SecondTime: t.ShipmentUpdatedAfter, SecondName: "ShipmentUpdatedAfter"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *ShipmentRecalculate) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *ShipmentRecalculate) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// FetchShipmentRecalculateDates returns the active recalculation date range
func FetchShipmentRecalculateDates(db *pop.Connection) (*ShipmentRecalculate, error) {
	type ShipmentRecalculateRecords []ShipmentRecalculate
	var recalculate ShipmentRecalculateRecords
	if err := db.Where("active = true").All(&recalculate); err != nil {
		return nil, errors.Wrap(err, "could not find re-calculate dates")
	}
	if len(recalculate) > 1 {
		return nil, errors.New("Too many active re-calculate date records")
	}
	if len(recalculate) == 1 {
		return &recalculate[0], nil
	}
	return nil, errors.New("No active records for re-calculate dates")
}
