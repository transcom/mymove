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

// MakeWeightTicket creates a single WeightTicket and associated relationships with weights and documents
func MakeWeightTicket(db *pop.Connection, assertions Assertions) models.WeightTicket {
	assertions = EnsureServiceMemberIsSetUpInAssertionsForDocumentCreation(db, assertions)

	// Because this model points at multiple documents, it's not really good to point at the base assertions.Document,
	// so we'll look at assertions.WeightTicket.<Document>
	emptyDocument := GetOrCreateDocumentWithUploads(db, assertions.WeightTicket.EmptyDocument, assertions)
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
	return MakeWeightTicket(db, Assertions{})
}
