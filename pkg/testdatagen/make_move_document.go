package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMoveDocument creates a single Move Document.
func MakeMoveDocument(db *pop.Connection, assertions Assertions) models.MoveDocument {
	document := assertions.MoveDocument.Document
	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.MoveDocument.DocumentID) {
		document = MakeDocument(db, assertions)
	}

	move := assertions.MoveDocument.Move
	// See note above
	if isZeroUUID(assertions.MoveDocument.MoveID) {
		move = MakeMove(db, assertions)
	}

	// We can't know in advance if a move document is for a PPM or a Shipment
	// Better to force the user to choose and explicitly pass in the values
	// than to make default versions of these structs.
	ppmID := assertions.MoveDocument.PersonallyProcuredMoveID
	ppm := assertions.MoveDocument.PersonallyProcuredMove

	moveDocumentType := models.MoveDocumentTypeOTHER
	if string(assertions.MoveDocument.MoveDocumentType) != "" {
		moveDocumentType = assertions.MoveDocument.MoveDocumentType
	}

	title := assertions.MoveDocument.Title
	if title == "" {
		title = "My-very-special-document.pdf"
	}

	moveDocument := models.MoveDocument{
		DocumentID:               document.ID,
		Document:                 document,
		MoveID:                   move.ID,
		Move:                     move,
		PersonallyProcuredMoveID: ppmID,
		PersonallyProcuredMove:   ppm,
		Status:                   models.MoveDocumentStatusAWAITINGREVIEW,
		MoveDocumentType:         moveDocumentType,
		Title:                    title,
	}

	// Overwrite values with those from assertions
	mergeModels(&moveDocument, assertions.MoveDocument)

	mustCreate(db, &moveDocument, assertions.Stub)

	return moveDocument
}

// MakeMoveDocumentWeightTicketSet creates a single Move Document with a WeightTicketSetDocument.
func MakeMoveDocumentWeightTicketSet(db *pop.Connection, assertions Assertions) models.MoveDocument {
	moveDocumentType := models.MoveDocumentTypeWEIGHTTICKETSET

	moveDocumentDefaults := models.MoveDocument{
		MoveDocumentType: moveDocumentType,
	}

	// Overwrite values with those from assertions
	mergeModels(&moveDocumentDefaults, assertions.MoveDocument)

	assertions.MoveDocument = moveDocumentDefaults

	moveDocument := MakeMoveDocument(db, assertions)

	weightTicketSetAssertions := Assertions{
		WeightTicketSetDocument: models.WeightTicketSetDocument{
			MoveDocumentID: moveDocument.ID,
			MoveDocument:   moveDocument,
		},
	}

	weightTicketSetDocument := MakeWeightTicketSetDocument(db, weightTicketSetAssertions)
	moveDocument.WeightTicketSetDocument = &weightTicketSetDocument

	if !assertions.Stub {
		MustSave(db, &moveDocument)
	}

	return moveDocument
}

// MakeDefaultMoveDocument returns a MoveDocument with default values
func MakeDefaultMoveDocument(db *pop.Connection) models.MoveDocument {
	return MakeMoveDocument(db, Assertions{})
}
