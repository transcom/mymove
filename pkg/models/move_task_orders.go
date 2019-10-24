package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type MoveTaskOrder struct {
	CreatedAt                time.Time     `db:"created_at"`
	Customer                 ServiceMember `belongs_to:"service_members"`
	CustomerID               uuid.UUID     `db:"customer_id"`
	CustomerRemarks          string        `db:"customer_remarks"`
	DestinationAddress       Address       `belongs_to:"addresses"`
	DestinationAddressID     uuid.UUID     `db:"destination_address_id"`
	DestinationDutyStation   DutyStation   `belongs_to:"duty_stations"`
	DestinationDutyStationID uuid.UUID     `db:"destination_duty_station_id"`
	ID                       uuid.UUID     `db:"id"`
	Move                     Move          `belongs_to:"moves"`
	MoveID                   uuid.UUID     `db:"move_id"`
	NTSEntitlement           bool          `db:"nts_entitlement"`
	OriginDutyStation        DutyStation   `belongs_to:"duty_stations"`
	OriginDutyStationID      uuid.UUID     `db:"origin_duty_station_id"`
	POVEntitlement           bool          `db:"pov_entitlement"`
	PickupAddress            Address       `belongs_to:"addresses"`
	PickupAddressID          uuid.UUID     `db:"pickup_address_id"`
	RequestedPickupDates     time.Time     `db:"requested_pickup_dates"`
	SitEntitlement           int           `db:"sit_entitlement"`
	UpdatedAt                time.Time     `db:"updated_at"`
	WeightEntitlement        int           `db:"weight_entitlement"`
}
