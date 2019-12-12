package testdatagen

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMTOServiceItem creates a single MtoServiceItem and associated set relationships
func MakeMTOServiceItem(db *pop.Connection, assertions Assertions) models.MtoServiceItem {
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
	MTOServiceItem := models.MtoServiceItem{
		MoveTaskOrder:   moveTaskOrder,
		MoveTaskOrderID: moveTaskOrder.ID,
		MTOShipment:     MTOShipment,
		MTOShipmentID:   MTOShipment.ID,
		ReService:       reService,
		ReServiceID:     reService.ID,
		MetaID:          uuid.Must(uuid.NewV4()),
		MetaType:        "unknown",
	}
	// Overwrite values with those from assertions
	mergeModels(&MTOServiceItem, assertions.MTOServiceItem)

	mustCreate(db, &MTOServiceItem)

	return MTOServiceItem
}
