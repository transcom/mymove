package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
)

type ReServiceItem struct {
	ID             uuid.UUID       `db:"id" rw:"r"`
	ServiceId      uuid.UUID       `db:"service_id" rw:"r"`
	ReService      ReService       `belongs_to:"re_services" fk_id:"service_id" rw:"r"`
	ShipmentType   MTOShipmentType `db:"shipment_type" rw:"r"`
	MarketCode     MarketCode      `db:"market_code" rw:"r"`
	IsAutoApproved bool            `db:"is_auto_approved" rw:"r"`
	Sort           *string         `db:"sort" rw:"r"`
	CreatedAt      time.Time       `db:"created_at" rw:"r"`
	UpdatedAt      time.Time       `db:"updated_at" rw:"r"`
}

func (r ReServiceItem) TableName() string {
	return "re_service_items"
}

// ReServiceItems is a slice of ReServiceItem
type ReServiceItems []ReServiceItem

func FetchReServiceByCode(db *pop.Connection, code ReServiceCode) (*ReService, error) {
	reService := ReService{}
	err := db.Where("code = ?", code).First(&reService)
	if err != nil {
		return nil, apperror.NewQueryError("ReService", err, "")
	}
	return &reService, err
}
