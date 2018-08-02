package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMovingExpenseDocument creates a single Moving Expense Document.
func MakeMovingExpenseDocument(db *pop.Connection, assertions Assertions) models.MovingExpenseDocument {
	moveDoc := assertions.MovingExpenseDocument.MoveDocument
	if isZeroUUID(assertions.MovingExpenseDocument.MoveDocumentID) {
		moveDoc = MakeMoveDocument(db, assertions)
	}

	reimbursement := assertions.MovingExpenseDocument.Reimbursement
	if isZeroUUID(assertions.MovingExpenseDocument.ReimbursementID) {
		reimbursement, _ = MakeDraftReimbursement(db)
	}

	movingExpenseDocument := models.MovingExpenseDocument{
		MoveDocumentID:    moveDoc.ID,
		MoveDocument:      moveDoc,
		MovingExpenseType: models.MovingExpenseTypeCONTRACTEDEXPENSE,
		ReimbursementID:   reimbursement.ID,
		Reimbursement:     reimbursement,
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
