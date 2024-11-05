package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type ReServiceItem struct {
	ID             uuid.UUID       `db:"id" rw:"r"`
	ServiceId      uuid.UUID       `db:"service_id" rw:"r"`
	Code           ReServiceCode   `belongs_to:"re_services" fk_id:"service_id" rw:"r"`
	ShipmentType   MTOShipmentType `db:"shipment_type" rw:"r"`
	MarketCode     MarketCode      `db:"market_code" rw:"r"`
	IsAutoApproved bool            `db:"is_auto_approved" rw:"r"`
	CreatedAt      time.Time       `db:"created_at" rw:"r"`
	UpdatedAt      time.Time       `db:"updated_at" rw:"r"`
}

func (g ReServiceItem) TableName() string {
	return "re_service_items"
}

// ReServiceItems is a slice of ReServiceItem
type ReServiceItems []ReServiceItem
