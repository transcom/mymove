package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMoveTaskOrder creates a single MoveTaskOrder and associated set relationships
func MakeMoveTaskOrder(db *pop.Connection, assertions Assertions) models.MoveTaskOrder {
	moveOrder := assertions.MoveOrder
	if isZeroUUID(moveOrder.ID) {
		moveOrder = MakeMoveOrder(db, assertions)
	}

	var referenceID string
	if assertions.MoveTaskOrder.ReferenceID == "" {
		referenceID, _ = models.GenerateReferenceID(db)
	}

	var ppmType string
	if assertions.MoveTaskOrder.PPMType == "" {
		ppmType = "PARTIAL"
	}

	moveTaskOrder := models.MoveTaskOrder{
		MoveOrder:          moveOrder,
		MoveOrderID:        moveOrder.ID,
		ReferenceID:        referenceID,
		IsAvailableToPrime: false,
		IsCanceled:         false,
		PPMType:            ppmType,
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
