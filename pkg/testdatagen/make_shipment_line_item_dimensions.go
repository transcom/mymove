package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// MakeShipmentLineItemDimensions creates a ShipmentLineItemDimensions record
func MakeShipmentLineItemDimensions(db *pop.Connection, assertions Assertions) models.ShipmentLineItemDimensions {
	dimensions := models.ShipmentLineItemDimensions{
		Length: unit.ThousandthInches(1000),
		Width:  unit.ThousandthInches(1000),
		Height: unit.ThousandthInches(1000),
	}

	// Overwrite values with those from assertions
	mergeModels(&dimensions, assertions.ShipmentLineItemDimensions)

	mustCreate(db, &dimensions)

	return dimensions
}

// MakeDefaultShipmentLineItemDimensions makes a ShipmentLineItemDimensions with default values
func MakeDefaultShipmentLineItemDimensions(db *pop.Connection) models.ShipmentLineItemDimensions {
	return MakeShipmentLineItemDimensions(db, Assertions{})
}
