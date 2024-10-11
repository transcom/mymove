package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type ReInternationalTransitTime struct {
	ID                    uuid.UUID `json:"id" db:"id"`
	OriginRateAreaId      uuid.UUID `json:"origin_rate_area_id" db:"origin_rate_area_id"`
	DestinationRateAreaId uuid.UUID `json:"destination_rate_area_id" db:"destination_rate_area_id"`
	HhgTransitTime        *int      `json:"hhg_transit_time" db:"hhg_transit_time"`
	UbTransitTime         *int      `json:"ub_transit_time" db:"ub_transit_time"`
	Active                *bool     `json:"active" db:"active"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}

func (ReInternationalTransitTime) TableName() string {
	return "re_international_transit_times"
}
