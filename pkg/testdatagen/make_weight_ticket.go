package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// MakeMinimalWeightTicket creates a single WeightTicket and associated relationships with a minimal set of data
func MakeMinimalWeightTicket(db *pop.Connection, assertions Assertions) models.WeightTicket {
	assertions = EnsureServiceMemberIsSetUpInAssertionsForDocumentCreation(db, assertions)

	ppmShipment := checkOrCreatePPMShipment(db, assertions)

	// Because this model points at multiple documents, it's not really good to point at the base assertions.Document,
	// so we'll look at assertions.WeightTicket.<Document>
	emptyDocument := GetOrCreateDocument(db, assertions.WeightTicket.EmptyDocument, assertions)
	fullDocument := GetOrCreateDocument(db, assertions.WeightTicket.FullDocument, assertions)
	trailerDocument := GetOrCreateDocument(db, assertions.WeightTicket.ProofOfTrailerOwnershipDocument, assertions)

	newWeightTicket := models.WeightTicket{
		PPMShipmentID:                     ppmShipment.ID,
		PPMShipment:                       ppmShipment,
		EmptyDocumentID:                   emptyDocument.ID,
		EmptyDocument:                     emptyDocument,
		FullDocumentID:                    fullDocument.ID,
		FullDocument:                      fullDocument,
		ProofOfTrailerOwnershipDocumentID: trailerDocument.ID,
		ProofOfTrailerOwnershipDocument:   trailerDocument,
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

// makeWeightTicket creates a single WeightTicket and associated relationships with weights and documents
func makeWeightTicket(db *pop.Connection, assertions Assertions) models.WeightTicket {
	assertions = EnsureServiceMemberIsSetUpInAssertionsForDocumentCreation(db, assertions)

	assertionsHasFileToUse := false
	if assertions.File != nil {
		assertionsHasFileToUse = true
	}

	// Because this model points at multiple documents, it's not really good to point at the base assertions.Document,
	// so we'll look at assertions.WeightTicket.<Document>
	if !assertionsHasFileToUse {
		assertions.File = Fixture("empty-weight-ticket.png")
	}

	emptyDocument := GetOrCreateDocumentWithUploads(db, assertions.WeightTicket.EmptyDocument, assertions)

	if !assertionsHasFileToUse {
		assertions.File = Fixture("full-weight-ticket.png")
	}

	fullDocument := GetOrCreateDocumentWithUploads(db, assertions.WeightTicket.FullDocument, assertions)

	emptyWeight := unit.Pound(14500)
	fullWeight := emptyWeight + unit.Pound(4000)

	fullAssertions := Assertions{
		WeightTicket: models.WeightTicket{
			VehicleDescription:       models.StringPointer("2022 Honda CR-V Hybrid"),
			EmptyWeight:              &emptyWeight,
			MissingEmptyWeightTicket: models.BoolPointer(false),
			EmptyDocumentID:          emptyDocument.ID,
			EmptyDocument:            emptyDocument,
			FullWeight:               &fullWeight,
			MissingFullWeightTicket:  models.BoolPointer(false),
			FullDocumentID:           fullDocument.ID,
			FullDocument:             fullDocument,
			OwnsTrailer:              models.BoolPointer(false),
			TrailerMeetsCriteria:     models.BoolPointer(false),
		},
	}

	// Overwrite values with those from assertions
	mergeModels(&fullAssertions, assertions)

	return MakeMinimalWeightTicket(db, fullAssertions)
}

// MakeDefaultWeightTicket makes a WeightTicket with default values
func MakeDefaultWeightTicket(db *pop.Connection) models.WeightTicket {
	return makeWeightTicket(db, Assertions{})
}

// MakeWeightTicketWithConstructedWeight creates a single WeightTicket and associated relationships with weights and documents
func MakeWeightTicketWithConstructedWeight(db *pop.Connection, assertions Assertions) models.WeightTicket {
	assertions = EnsureServiceMemberIsSetUpInAssertionsForDocumentCreation(db, assertions)

	// If they don't have weight tickets, they'll be uploading a vehicle registration or a rental agreement
	assertions.File = Fixture("wa-vehicle-registration.pdf")

	// Because this model points at multiple documents, it's not really good to point at the base assertions.Document,
	// so we'll look at assertions.WeightTicket.<Document>
	assertions.WeightTicket.EmptyDocument = GetOrCreateDocumentWithUploads(db, assertions.WeightTicket.EmptyDocument, assertions)

	assertions.WeightTicket.MissingEmptyWeightTicket = models.BoolPointer(true)

	// If they don't have weight tickets, they'll be uploading a constructed weight spreadsheet for the
	// full document upload.
	assertions.File = Fixture("Weight Estimator.xls")

	assertions.WeightTicket.FullDocument = GetOrCreateDocumentWithUploads(db, assertions.WeightTicket.FullDocument, assertions)

	assertions.WeightTicket.MissingFullWeightTicket = models.BoolPointer(true)

	return makeWeightTicket(db, assertions)
}
