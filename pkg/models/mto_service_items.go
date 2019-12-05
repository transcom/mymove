package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type MTOServiceItem struct {
	ID              uuid.UUID     `db:"id"`
	MoveTaskOrder   MoveTaskOrder `belongs_to:"move_task_orders"`
	MoveTaskOrderID uuid.UUID     `db:"move_task_order_id"`
	MTOShipment     MTOShipment   `belongs_to:"mto_shipments"`
	MTOShipmentID   uuid.UUID     `db:"mto_shipment_id"`
	ReService       ReService     `belongs_to:"re_services"`
	ReServiceID     uuid.UUID     `db:"re_service_id"`
	CreatedAt       time.Time     `db:"created_at"`
	UpdatedAt       time.Time     `db:"updated_at"`
}
