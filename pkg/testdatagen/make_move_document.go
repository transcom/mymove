package testdatagen

import (
	"github.com/gobuffalo/pop"

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
	ppmId := assertions.MoveDocument.PersonallyProcuredMoveID
	ppm := assertions.MoveDocument.PersonallyProcuredMove
	shipmentId := assertions.MoveDocument.ShipmentID
	shipment := assertions.MoveDocument.Shipment

	moveDocumentType := models.MoveDocumentTypeOTHER
	if string(assertions.MoveDocument.MoveDocumentType) != "" {
		moveDocumentType = assertions.MoveDocument.MoveDocumentType
	}

	title := assertions.MoveDocument.Title
	if title == "" {
		title = "My-very-special-document.pdf"
	}

	moveDocument := models.MoveDocument{
		DocumentID: document.ID,
		Document:   document,
		MoveID:     move.ID,
		Move:       move,
		PersonallyProcuredMoveID: ppmId,
		PersonallyProcuredMove:   ppm,
		ShipmentID:               shipmentId,
		Shipment:                 shipment,
		Status:                   models.MoveDocumentStatusAWAITINGREVIEW,
		MoveDocumentType:         moveDocumentType,
		Title:                    title,
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
