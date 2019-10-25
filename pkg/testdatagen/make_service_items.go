package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeServiceItem creates a single ServiceItems and associated set relationships
func MakeServiceItem(db *pop.Connection, assertions Assertions) models.ServiceItem {

	// Create new Orders if not provided
	// ID is required because it must be populated for Eager saving to work.
	moveTaskOrder := assertions.MoveTaskOrder
	if isZeroUUID(assertions.MoveTaskOrder.ID) {
		moveTaskOrder = MakeMoveTaskOrder(db, assertions)
	}

	serviceItem := models.ServiceItem{
		MoveTaskOrder:   moveTaskOrder,
		MoveTaskOrderID: moveTaskOrder.ID,
	}

	// Overwrite values with those from assertions
	mergeModels(&serviceItem, assertions.ServiceItem)

	mustCreate(db, &serviceItem)

	return serviceItem
}
