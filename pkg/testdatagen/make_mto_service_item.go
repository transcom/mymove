package testdatagen

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMtoServiceItem creates a single MtoServiceItem and associated set relationships
func MakeMtoServiceItem(db *pop.Connection, assertions Assertions) models.MtoServiceItem {
	moveTaskOrder := assertions.MoveTaskOrder
	if isZeroUUID(moveTaskOrder.ID) {
		moveTaskOrder = MakeMoveTaskOrder(db, assertions)
	}
	mtoShipment := assertions.MtoShipment
	if isZeroUUID(mtoShipment.ID) {
		mtoShipment = MakeMtoShipment(db, assertions)
	}
	reService := assertions.ReService
	if isZeroUUID(reService.ID) {
		reService = MakeReService(db, assertions)
	}
	MtoServiceItem := models.MtoServiceItem{
		MoveTaskOrder:   moveTaskOrder,
		MoveTaskOrderID: moveTaskOrder.ID,
		MtoShipment:     mtoShipment,
		MtoShipmentID:   mtoShipment.ID,
		ReService:       reService,
		ReServiceID:     reService.ID,
		MetaID:          uuid.Must(uuid.NewV4()),
		MetaType:        "unknown",
	}
	// Overwrite values with those from assertions
	mergeModels(&MtoServiceItem, assertions.MtoServiceItem)

	mustCreate(db, &MtoServiceItem)

	return MtoServiceItem
}
