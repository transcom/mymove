package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMTOServiceItem creates a single MTOServiceItem and associated set relationships
func MakeMTOServiceItem(db *pop.Connection, assertions Assertions) models.MTOServiceItem {
	moveTaskOrder := assertions.MoveTaskOrder
	if isZeroUUID(moveTaskOrder.ID) {
		moveTaskOrder = MakeMoveTaskOrder(db, assertions)
	}
	MTOShipment := assertions.MTOShipment
	if isZeroUUID(MTOShipment.ID) {
		MTOShipment = MakeMTOShipment(db, assertions)
	}
	reService := assertions.ReService
	if isZeroUUID(reService.ID) {
		reService = MakeReService(db, assertions)
	}

	MTOServiceItem := models.MTOServiceItem{
		MoveTaskOrder:   moveTaskOrder,
		MoveTaskOrderID: moveTaskOrder.ID,
		MTOShipment:     MTOShipment,
		MTOShipmentID:   &MTOShipment.ID,
		ReService:       reService,
		ReServiceID:     reService.ID,
	}
	// Overwrite values with those from assertions
	mergeModels(&MTOServiceItem, assertions.MTOServiceItem)

	mustCreate(db, &MTOServiceItem)

	return MTOServiceItem
}
