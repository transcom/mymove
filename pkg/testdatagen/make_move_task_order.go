package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMoveTaskOrder creates a single MoveTaskOrder and associated set relationships
func MakeMoveTaskOrder(db *pop.Connection, assertions Assertions) models.MoveTaskOrder {
	var moveOrder models.MoveOrder
	if isZeroUUID(assertions.MoveOrder.ID) {
		moveOrder = MakeMoveOrder(db, assertions)
	}
	mtoStatus := assertions.MoveTaskOrder.Status
	if mtoStatus == "" {
		mtoStatus = models.MoveTaskOrderStatusApproved
	}
	var referenceID *string
	moveTaskOrder := models.MoveTaskOrder{
		MoveOrder:          moveOrder,
		MoveOrderID:        moveOrder.ID,
		ReferenceID:        referenceID,
		Status:             mtoStatus,
		IsAvailableToPrime: false,
		IsCancelled:        false,
	}

	// Overwrite values with those from assertions
	mergeModels(&moveTaskOrder, assertions.MoveTaskOrder)

	mustCreate(db, &moveTaskOrder)

	return moveTaskOrder
}

// MakeDefaultMoveTaskOrder makes an MoveTaskOrder with default values
func MakeDefaultMoveTaskOrder(db *pop.Connection) models.MoveTaskOrder {
	return MakeMoveTaskOrder(db, Assertions{})
}
