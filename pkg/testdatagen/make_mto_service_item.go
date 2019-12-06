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
	mtoShipment := assertions.MTOShipment
	if isZeroUUID(mtoShipment.ID) {
		mtoShipment = MakeMTOShipment(db, assertions)
	}
	reService := assertions.ReService
	if isZeroUUID(reService.ID) {
		reService = MakeReService(db, assertions)
	}
	mtoServiceItem := models.MTOServiceItem{
		MoveTaskOrder:   moveTaskOrder,
		MoveTaskOrderID: moveTaskOrder.ID,
		MTOShipment:     mtoShipment,
		MTOShipmentID:   mtoShipment.ID,
		ReService:       reService,
		ReServiceID:     reService.ID,
	}
	// Overwrite values with those from assertions
	mergeModels(&mtoServiceItem, assertions.MTOServiceItem)

	mustCreate(db, &mtoServiceItem)

	return mtoServiceItem
}
