package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"

	"github.com/gofrs/uuid"
)

type MoveTaskOrder struct {
	ID                       uuid.UUID           `db:"id"`
	CreatedAt                time.Time           `db:"created_at"`
	Customer                 ServiceMember       `belongs_to:"service_members"`
	CustomerID               uuid.UUID           `db:"customer_id"`
	CustomerRemarks          string              `db:"customer_remarks"`
	DestinationAddress       Address             `belongs_to:"addresses"`
	DestinationAddressID     uuid.UUID           `db:"destination_address_id"`
	DestinationDutyStation   DutyStation         `belongs_to:"duty_stations"`
	DestinationDutyStationID uuid.UUID           `db:"destination_duty_station_id"`
	Move                     Move                `belongs_to:"moves"`
	MoveID                   uuid.UUID           `db:"move_id"`
	NTSEntitlement           bool                `db:"nts_entitlement"`
	OriginDutyStation        DutyStation         `belongs_to:"duty_stations"`
	OriginDutyStationID      uuid.UUID           `db:"origin_duty_station_id"`
	POVEntitlement           bool                `db:"pov_entitlement"`
	PickupAddress            Address             `belongs_to:"addresses"`
	PickupAddressID          uuid.UUID           `db:"pickup_address_id"`
	RequestedPickupDates     time.Time           `db:"requested_pickup_dates"`
	SitEntitlement           int64               `db:"sit_entitlement"`
	Status                   MoveTaskOrderStatus `db:"status"`
	ServiceItems             ServiceItems        `has_many:"service_items"`
	UpdatedAt                time.Time           `db:"updated_at"`
	//TODO This was in ticket but looks different from API defn
	WeightEntitlement int64 `db:"weight_entitlement"`
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
