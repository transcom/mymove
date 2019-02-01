package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/unit"
)

// Dimensions is an object representing dimensions of a shipment line item
type Dimensions struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Length    unit.Inch `json:"length" db:"length"`
	Width     unit.Inch `json:"width" db:"width"`
	Height    unit.Inch `json:"height" db:"height"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
