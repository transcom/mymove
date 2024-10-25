package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type ReServiceItems struct {
	ID           uuid.UUID       `db:"id" rw:"r"`
	ServiceId    uuid.UUID       `db:"service_id" rw:"r"`
	Code         ReServiceCode   `belongs_to:"re_services" fk_id:"service_id"`
	ShipmentType MTOShipmentType `db:"shipment_type"`
	MarketCode   MarketCode      `db:"market_code"`
	AutoApproved bool            `db:"auto_approved" rw:"r"`
	CreatedAt    time.Time       `db:"created_at" rw:"r"`
	UpdatedAt    time.Time       `db:"updated_at" rw:"r"`
}

func (g ReServiceItems) TableName() string {
	return "re_service_items"
}
