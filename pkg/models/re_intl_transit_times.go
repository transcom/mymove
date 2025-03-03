package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
)

type InternationalTransitTime struct {
	ID                    uuid.UUID `json:"id" db:"id" rw:"r"`
	OriginRateAreaId      uuid.UUID `json:"origin_rate_area_id" db:"origin_rate_area_id" rw:"r"`
	DestinationRateAreaId uuid.UUID `json:"destination_rate_area_id" db:"destination_rate_area_id" rw:"r"`
	HhgTransitTime        *int      `json:"hhg_transit_time" db:"hhg_transit_time" rw:"r"`
	UbTransitTime         *int      `json:"ub_transit_time" db:"ub_transit_time" rw:"r"`
	Active                bool      `json:"active" db:"active" rw:"r"`
	CreatedAt             time.Time `json:"created_at" db:"created_at" rw:"r"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at" rw:"r"`
}

func (InternationalTransitTime) TableName() string {
	return "re_intl_transit_times"
}

// fetch the re_intl_transit_time record from the db. This function is bi-directional.
func FetchInternationalTransitTime(db *pop.Connection, originRateAreaId uuid.UUID, destinationRateAreaId uuid.UUID) (InternationalTransitTime, error) {
	var internationalTransitTime InternationalTransitTime
	err := db.
		Where("origin_rate_area_id = $1 and destination_rate_area_id = $2", originRateAreaId, destinationRateAreaId).
		First(&internationalTransitTime)

	if err != nil && internationalTransitTime.ID.IsNil() {
		err = db.
			Where("origin_rate_area_id = $1 and destination_rate_area_id = $2", destinationRateAreaId, originRateAreaId).
			First(&internationalTransitTime)
	}

	if err != nil && internationalTransitTime.ID.IsNil() {
		return internationalTransitTime, apperror.NewQueryError("InternationalTransitTime", err, "could not look up intl transit time")
	}

	return internationalTransitTime, nil
}
