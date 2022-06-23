package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// MakeMinimalWeightTicket creates a single WeightTicket and associated relationships with a minimal set of data
func MakeMinimalWeightTicket(db *pop.Connection, assertions Assertions) models.WeightTicket {
	ppmShipment := checkOrCreatePPMShipment(db, assertions)

	newWeightTicket := models.WeightTicket{
		PPMShipmentID: ppmShipment.ID,
		PPMShipment:   ppmShipment,
	}

	// Overwrite values with those from assertions
	mergeModels(&newWeightTicket, assertions.WeightTicket)

	mustCreate(db, &newWeightTicket, assertions.Stub)

	return newWeightTicket
}

// MakeMinimalDefaultWeightTicket makes a WeightTicket with minimal default values
func MakeMinimalDefaultWeightTicket(db *pop.Connection) models.WeightTicket {
	return MakeMinimalWeightTicket(db, Assertions{})
}

// MakeWeightTicket creates a single WeightTicket and associated relationships with weights and documents
func MakeWeightTicket(db *pop.Connection, assertions Assertions) models.WeightTicket {
	emptyWeight := unit.Pound(14500)
	fullWeight := emptyWeight + unit.Pound(4000)

	// Several of the downstream functions need a service member, but they don't always share assertions, look at the
	// same assertion, or create the service members in the same ways. We'll check now to see if we already have one
	// created, and if not, create one that we can place in the assertions for all the rest.
	if !assertions.Stub && assertions.ServiceMember.CreatedAt.IsZero() || assertions.ServiceMember.ID.IsNil() {
		serviceMember := MakeExtendedServiceMember(db, assertions)

		assertions.ServiceMember = serviceMember
		assertions.Order.ServiceMemberID = serviceMember.ID
		assertions.Order.ServiceMember = serviceMember
		assertions.Document.ServiceMemberID = serviceMember.ID
		assertions.Document.ServiceMember = serviceMember
	}

	// Because this model points at multiple documents, it's not really good to point at the base assertions.Document,
	// so we'll look at assertions.WeightTicket.<Document>
	emptyDocument := getOrCreateDocument(db, assertions.WeightTicket.EmptyDocument, assertions)
	fullDocument := getOrCreateDocument(db, assertions.WeightTicket.FullDocument, assertions)

	fullAssertions := Assertions{
		WeightTicket: models.WeightTicket{
			EmptyWeight:          &emptyWeight,
			HasEmptyWeightTicket: models.BoolPointer(true),
			EmptyDocumentID:      &emptyDocument.ID,
			EmptyDocument:        &emptyDocument,
			FullWeight:           &fullWeight,
			HasFullWeightTicket:  models.BoolPointer(true),
			FullDocumentID:       &fullDocument.ID,
			FullDocument:         &fullDocument,
			OwnsTrailer:          models.BoolPointer(false),
		},
	}

	// Overwrite values with those from assertions
	mergeModels(&fullAssertions, assertions)

	return MakeMinimalWeightTicket(db, fullAssertions)
}

// MakeDefaultWeightTicket makes a WeightTicket with default values
func MakeDefaultWeightTicket(db *pop.Connection) models.WeightTicket {
	return MakeWeightTicket(db, Assertions{})
}

// checkOrCreatePPMShipment checks PPMShipment in assertions, or creates one if none exists.
func checkOrCreatePPMShipment(db *pop.Connection, assertions Assertions) models.PPMShipment {
	ppmShipment := assertions.PPMShipment

	if !assertions.Stub && ppmShipment.CreatedAt.IsZero() || ppmShipment.ID.IsNil() {
		ppmShipment = MakeApprovedPPMShipmentWithActualInfo(db, assertions)
	}

	return ppmShipment
}

// getOrCreateDocument checks if a document exists. If it does, it returns it, otherwise, it creates it
func getOrCreateDocument(db *pop.Connection, document *models.Document, assertions Assertions) models.Document {
	if document == nil || assertions.Stub && document.CreatedAt.IsZero() || document.ID.IsNil() {
		upload := MakeUserUpload(db, assertions)

		return upload.Document
	}

	return *document
}
