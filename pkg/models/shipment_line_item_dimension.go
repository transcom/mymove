package models

import (
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/unit"
)

// ShipmentLineItemDimension is an object representing dimensions of a shipment line item
type ShipmentLineItemDimension struct {
	ID     uuid.UUID `json:"id" db:"id"`
	Length unit.Inch `json:"length" db:"length"`
	Width  unit.Inch `json:"width" db:"width"`
	Height unit.Inch `json:"height" db:"height"`
}
