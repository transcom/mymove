package models

import (
	"time"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type MoveTaskOrder struct {
	ID                               uuid.UUID           `db:"id"`
	ActualWeight                     *unit.Pound         `json:"actual_weight" db:"actual_weight"`
	CreatedAt                        time.Time           `db:"created_at"`
	Customer                         ServiceMember       `belongs_to:"service_members"`
	CustomerID                       uuid.UUID           `db:"customer_id"`
	CustomerRemarks                  string              `db:"customer_remarks"`
	DestinationAddress               Address             `belongs_to:"addresses"`
	DestinationAddressID             uuid.UUID           `db:"destination_address_id"`
	DestinationDutyStation           DutyStation         `belongs_to:"duty_stations"`
	DestinationDutyStationID         uuid.UUID           `db:"destination_duty_station_id"`
	Entitlements                     GHCEntitlement      `has_one:"entitlements"`
	Move                             Move                `belongs_to:"moves"`
	MoveID                           uuid.UUID           `db:"move_id"`
	OriginDutyStation                DutyStation         `belongs_to:"duty_stations"`
	OriginDutyStationID              uuid.UUID           `db:"origin_duty_station_id"`
	PickupAddress                    Address             `belongs_to:"addresses"`
	PickupAddressID                  uuid.UUID           `db:"pickup_address_id"`
	PrimeEstimatedWeight             *unit.Pound         `db:"prime_estimated_weight"`
	PrimeEstimatedWeightRecordedDate *time.Time          `db:"prime_estimated_weight_recorded_date"`
	RequestedPickupDate              time.Time           `db:"requested_pickup_date"`
	Status                           MoveTaskOrderStatus `db:"status"`
	ServiceItems                     ServiceItems        `has_many:"service_items"`
	UpdatedAt                        time.Time           `db:"updated_at"`
}

type MoveTaskOrderStatus string

const (
	MoveTaskOrderStatusApproved  MoveTaskOrderStatus = "APPROVED"
	MoveTaskOrderStatusSubmitted MoveTaskOrderStatus = "SUBMITTED"
	MoveTaskOrderStatusRejected  MoveTaskOrderStatus = "REJECTED"
	MoveTaskOrderStatusDraft     MoveTaskOrderStatus = "DRAFT"
)

func (m *MoveTaskOrder) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: string(m.Status), Name: "Status"},
	), nil
}

type MoveTaskOrders []MoveTaskOrder
