package testdatagen

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMoveTaskOrder creates a single MoveTaskOrder and associated set relationships
func MakeMoveTaskOrder(db *pop.Connection, assertions Assertions) models.Move {
	order := assertions.Order
	if isZeroUUID(order.ID) {
		order = MakeOrder(db, assertions)
	}

	var referenceID string
	assertedReferenceID := assertions.MoveTaskOrder.ReferenceID
	if assertedReferenceID == nil || *assertedReferenceID == "" {
		referenceID, _ = models.GenerateReferenceID(db)
	}

	var contractorID uuid.UUID

	if assertions.MoveTaskOrder.ContractorID == nil {
		contractor := MakeContractor(db, assertions)
		contractorID = contractor.ID
	}

	ppmType := assertions.MoveTaskOrder.PPMType
	if assertions.MoveTaskOrder.PPMType == nil {
		partialType := "PARTIAL"
		ppmType = &partialType
	}

	moveTaskOrder := models.Move{
		AvailableToPrimeAt: assertions.MoveTaskOrder.AvailableToPrimeAt,
		Orders:             order,
		OrdersID:           order.ID,
		ContractorID:       &contractorID,
		ReferenceID:        &referenceID,
		Locator:            models.GenerateLocator(),
		Status:             models.MoveStatusDRAFT,
		PPMType:            ppmType,
		Show:               setShow(assertions.Move.Show),
	}

	// Overwrite values with those from assertions
	mergeModels(&moveTaskOrder, assertions.MoveTaskOrder)

	mustCreate(db, &moveTaskOrder)

	return moveTaskOrder
}

// MakeDefaultMoveTaskOrder makes an MoveTaskOrder with default values
func MakeDefaultMoveTaskOrder(db *pop.Connection) models.Move {
	return MakeMoveTaskOrder(db, Assertions{})
}
