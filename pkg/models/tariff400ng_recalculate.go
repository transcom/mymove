package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"time"
)

// Tariff400ngRecalculate model/struct is a record of the date range to use to determine if a Shipment is a candidate
// to be re-calculated.
type Tariff400ngRecalculate struct {
	ID                    uuid.UUID `json:"id" db:"id"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
	ShipmentUdpatedBefore time.Time `json:"shipment_udpated_before" db:"shipment_udpated_before"`
	ShipmentUpdatedAfter  time.Time `json:"shipment_updated_after" db:"shipment_updated_after"`
	Active                bool      `json:"active" db:"active"`
}

// String is not required by pop and may be deleted
func (t Tariff400ngRecalculate) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Tariff400ngRecalculates is not required by pop and may be deleted
type Tariff400ngRecalculates []Tariff400ngRecalculate

// String is not required by pop and may be deleted
func (t Tariff400ngRecalculates) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Tariff400ngRecalculate) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.TimeIsPresent{Field: t.ShipmentUdpatedBefore, Name: "ShipmentUdpatedBefore"},
		&validators.TimeIsPresent{Field: t.ShipmentUpdatedAfter, Name: "ShipmentUpdatedAfter"},
		&validators.TimeAfterTime{
			FirstTime: t.ShipmentUdpatedBefore, FirstName: "ShipmentUdpatedBefore",
			SecondTime: t.ShipmentUpdatedAfter, SecondName: "ShipmentUpdatedAfter"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Tariff400ngRecalculate) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Tariff400ngRecalculate) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// FetchTariff400ngRecalculateDates returns the active recalculation date range
func FetchTariff400ngRecalculateDates(db *pop.Connection) (*Tariff400ngRecalculate, error) {
	type Tariff400ngRecalculateRecords []Tariff400ngRecalculate
	var recalculate Tariff400ngRecalculateRecords
	if err := db.Where("active = true").All(&recalculate); err != nil {
		return &Tariff400ngRecalculate{}, errors.Wrap(err, "could not find re-calculate dates")
	}
	if len(recalculate) > 1 {
		return &Tariff400ngRecalculate{}, errors.New("Too many active re-calculate date records")
	}
	if len(recalculate) == 1 {
		return &recalculate[0], nil
	}
	return &Tariff400ngRecalculate{}, errors.New("No active records for re-calculate dates")
}
