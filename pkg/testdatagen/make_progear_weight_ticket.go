package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func MakeMinimalProgearWeightTicket(db *pop.Connection, assertions Assertions) models.ProgearWeightTicket {
	ppmShipment := checkOrCreatePPMShipment(db, assertions)

	emptyDocument := GetOrCreateDocument(db, assertions.ProgearWeightTicket.EmptyDocument, assertions)
	fullDocument := GetOrCreateDocument(db, assertions.ProgearWeightTicket.FullDocument, assertions)
	constructedWeightDocument := GetOrCreateDocument(db, assertions.ProgearWeightTicket.ConstructedWeightDocument, assertions)

	newProgearWeightTicket := models.ProgearWeightTicket{
		PPMShipmentID:               ppmShipment.ID,
		PPMShipment:                 ppmShipment,
		EmptyDocumentID:             emptyDocument.ID,
		EmptyDocument:               emptyDocument,
		FullDocumentID:              fullDocument.ID,
		FullDocument:                fullDocument,
		ConstructedWeightDocumentID: constructedWeightDocument.ID,
		ConstructedWeightDocument:   constructedWeightDocument,
	}

	// Overwrites model with data from assertions
	mergeModels(&newProgearWeightTicket, assertions.ProgearWeightTicket)

	mustCreate(db, &newProgearWeightTicket, assertions.Stub)

	return newProgearWeightTicket
}

func MakeMinimalDefaultProgearWeightTicket(db *pop.Connection) models.ProgearWeightTicket {
	return MakeMinimalProgearWeightTicket(db, Assertions{})
}

func MakeProgearWeightTicket(db *pop.Connection, assertions Assertions) models.ProgearWeightTicket {
	emptyDocument := GetOrCreateDocumentWithUploads(db, assertions.ProgearWeightTicket.EmptyDocument, assertions)
	fullDocument := GetOrCreateDocumentWithUploads(db, assertions.ProgearWeightTicket.FullDocument, assertions)

	description := "professional equipment"

	emptyWeight := unit.Pound(4500)
	fullWeight := emptyWeight + unit.Pound(500)

	fullAssertions := Assertions{
		ProgearWeightTicket: models.ProgearWeightTicket{
			EmptyDocumentID:  emptyDocument.ID,
			EmptyDocument:    emptyDocument,
			FullDocumentID:   fullDocument.ID,
			FullDocument:     fullDocument,
			BelongsToSelf:    models.BoolPointer(true),
			Description:      &description,
			HasWeightTickets: models.BoolPointer(true),
			EmptyWeight:      &emptyWeight,
			FullWeight:       &fullWeight,
		},
	}
	mergeModels(&fullAssertions, assertions)

	return MakeMinimalProgearWeightTicket(db, fullAssertions)
}

func MakeDefaultProgearWeightTicket(db *pop.Connection) models.ProgearWeightTicket {
	return MakeProgearWeightTicket(db, Assertions{})
}
