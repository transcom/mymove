package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type MoveOrder struct {
	ID                       uuid.UUID   `db:"id"`
	CreatedAt                time.Time   `db:"created_at"`
	Customer                 Customer    `belongs_to:"customers"`
	CustomerID               uuid.UUID   `db:"customer_id"`
	Entitlement              Entitlement `belongs_to:"entitlements"`
	EntitlementID            uuid.UUID   `db:"entitlement_id"`
	DestinationAddress       Address     `belongs_to:"addresses"`
	DestinationAddressID     uuid.UUID   `db:"destination_address_id"`
	DestinationDutyStation   DutyStation `belongs_to:"duty_stations"`
	DestinationDutyStationID uuid.UUID   `db:"destination_duty_station_id"`
	OriginDutyStation        DutyStation `belongs_to:"duty_stations"`
	OriginDutyStationID      uuid.UUID   `db:"origin_duty_station_id"`
}
