package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMovingExpenseDocument creates a single Moving Expense Document.
func MakeMovingExpenseDocument(db *pop.Connection, assertions Assertions) models.MovingExpenseDocument {
	moveDoc := assertions.MovingExpenseDocument.MoveDocument
	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.MovingExpenseDocument.MoveDocumentID) {
		moveDoc = MakeMoveDocument(db, assertions)
	}

	movingExpenseDocument := models.MovingExpenseDocument{
		MoveDocumentID:       moveDoc.ID,
		MoveDocument:         moveDoc,
		MovingExpenseType:    models.MovingExpenseTypeCONTRACTEDEXPENSE,
		PaymentMethod:        "GTCC",
		RequestedAmountCents: unit.Cents(2589),
	}

	// Overwrite values with those from assertions
	mergeModels(&movingExpenseDocument, assertions.MovingExpenseDocument)

	mustCreate(db, &movingExpenseDocument)

	return movingExpenseDocument
}

// MakeDefaultMovingExpenseDocument returns a MoveDocument with default values
func MakeDefaultMovingExpenseDocument(db *pop.Connection) models.MovingExpenseDocument {
	return MakeMovingExpenseDocument(db, Assertions{})
}
