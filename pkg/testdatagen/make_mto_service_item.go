package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMTOServiceItem creates a single MTOServiceItem and associated set relationships
func MakeMTOServiceItem(db *pop.Connection, assertions Assertions) models.MTOServiceItem {
	var moveTaskOrder models.MoveTaskOrder
	if isZeroUUID(assertions.MoveTaskOrder.ID) {
		moveTaskOrder = MakeMoveTaskOrder(db, assertions)
	}
	var mtoShipment models.MTOShipment
	if isZeroUUID(assertions.MTOShipment.ID) {
		mtoShipment = MakeMTOShipment(db, assertions)
	}
	var reService models.ReService
	if isZeroUUID(assertions.ReService.ID) {
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
