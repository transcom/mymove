package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMoveDocument creates a single Move Document.
func MakeMoveDocument(db *pop.Connection, assertions Assertions) models.MoveDocument {
	document := assertions.MoveDocument.Document
	// In a real scenario, ID will have to be populated for the model
	// To be populated by Eager, which is why ID is required
	if isZeroUUID(assertions.MoveDocument.DocumentID) {
		document = MakeDocument(db, assertions)
	}

	move := assertions.MoveDocument.Move
	// See note above
	if isZeroUUID(assertions.MoveDocument.MoveID) {
		move = MakeMove(db, assertions)
	}

	moveDocument := models.MoveDocument{
		DocumentID: document.ID,
		Document:   document,
		MoveID:     move.ID,
		Move:       move,
		PersonallyProcuredMoveID: nil,
		Status:           models.MoveDocumentStatusAWAITINGREVIEW,
		MoveDocumentType: models.MoveDocumentTypeOTHER,
		Title:            "My-very-special-document.pdf",
	}

	// Overwrite values with those from assertions
	mergeModels(&moveDocument, assertions.MoveDocument)

	mustCreate(db, &moveDocument)

	return moveDocument
}

// MakeDefaultMoveDocument returns a MoveDocument with default values
func MakeDefaultMoveDocument(db *pop.Connection) models.MoveDocument {
	return MakeMoveDocument(db, Assertions{})
}
