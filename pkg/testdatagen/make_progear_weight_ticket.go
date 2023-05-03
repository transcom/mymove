package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func makeMinimalProgearWeightTicket(db *pop.Connection, assertions Assertions) models.ProgearWeightTicket {
	ppmShipment := checkOrCreatePPMShipment(db, assertions)

	document := GetOrCreateDocument(db, assertions.ProgearWeightTicket.Document, assertions)

	newProgearWeightTicket := models.ProgearWeightTicket{
		PPMShipmentID: ppmShipment.ID,
		PPMShipment:   ppmShipment,
		DocumentID:    document.ID,
		Document:      document,
	}

	// Overwrites model with data from assertions
	mergeModels(&newProgearWeightTicket, assertions.ProgearWeightTicket)

	mustCreate(db, &newProgearWeightTicket, assertions.Stub)

	return newProgearWeightTicket
}

func makeProgearWeightTicket(db *pop.Connection, assertions Assertions) models.ProgearWeightTicket {
	document := GetOrCreateDocumentWithUploads(db, assertions.ProgearWeightTicket.Document, assertions)

	description := "professional equipment"

	fullAssertions := Assertions{
		ProgearWeightTicket: models.ProgearWeightTicket{
			DocumentID:       document.ID,
			Document:         document,
			BelongsToSelf:    models.BoolPointer(true),
			Description:      &description,
			HasWeightTickets: models.BoolPointer(true),
			Weight:           models.PoundPointer(unit.Pound(500)),
		},
	}
	mergeModels(&fullAssertions, assertions)

	return makeMinimalProgearWeightTicket(db, fullAssertions)
}
