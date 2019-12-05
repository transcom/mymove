package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/unit"
)

type MTOShipment struct {
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
