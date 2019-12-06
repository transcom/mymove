package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// MtoShipment is an object representing data for a move task order shipment
type MtoShipment struct {
	ID                               uuid.UUID     `db:"id"`
	MoveTaskOrder                    MoveTaskOrder `belongs_to:"move_task_orders"`
	MoveTaskOrderID                  uuid.UUID     `db:"move_task_order_id"`
	ScheduledPickupDate              *time.Time    `db:"scheduled_pickup_date"`
	RequestedPickupDate              *time.Time    `db:"requested_pickup_date"`
	CustomerRemarks                  *string       `db:"customer_remarks"`
	PickupAddress                    Address       `belongs_to:"addresses"`
	PickupAddressID                  uuid.UUID     `db:"pickup_address_id"`
	DestinationAddress               Address       `belongs_to:"addresses"`
	DestinationAddressID             uuid.UUID     `db:"destination_address_id"`
	SecondaryPickupAddress           *Address      `belongs_to:"addresses"`
	SecondaryPickupAddressID         *uuid.UUID    `db:"secondary_pickup_address_id"`
	SecondaryDeliveryAddress         *Address      `belongs_to:"addresses"`
	SecondaryDeliveryAddressID       *uuid.UUID    `db:"secondary_delivery_address_id"`
	PrimeEstimatedWeight             *unit.Pound   `db:"prime_estimated_weight"`
	PrimeEstimatedWeightRecordedDate *time.Time    `db:"prime_estimated_weight_recorded_date"`
	PrimeActualWeight                *unit.Pound   `db:"prime_actual_weight"`
	CreatedAt                        time.Time     `db:"created_at"`
	UpdatedAt                        time.Time     `db:"updated_at"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (m *MtoShipment) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs, &validators.UUIDIsPresent{Field: m.MoveTaskOrderID, Name: "MoveTaskOrderID"})
	vs = append(vs, &validators.UUIDIsPresent{Field: m.PickupAddressID, Name: "PickupAddressID"})
	vs = append(vs, &validators.UUIDIsPresent{Field: m.DestinationAddressID, Name: "DestinationAddressID"})
	if m.PrimeEstimatedWeight != nil {
		vs = append(vs, &validators.IntIsGreaterThan{Field: m.PrimeEstimatedWeight.Int(), Compared: -1, Name: "PrimeEstimatedWeight"})
	}
	if m.PrimeActualWeight != nil {
		vs = append(vs, &validators.IntIsGreaterThan{Field: m.PrimeActualWeight.Int(), Compared: -1, Name: "PrimeActualWeight"})
	}
	return validate.Validate(vs...), nil
}
