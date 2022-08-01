package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// MakeMinimalWeightTicket creates a single WeightTicket and associated relationships with a minimal set of data
func MakeMinimalWeightTicket(db *pop.Connection, assertions Assertions) models.WeightTicket {
	assertions = ensureServiceMemberIsSetUpInAssertions(db, assertions)

	ppmShipment := checkOrCreatePPMShipment(db, assertions)

	// Because this model points at multiple documents, it's not really good to point at the base assertions.Document,
	// so we'll look at assertions.WeightTicket.<Document>
	emptyDocument := getOrCreateDocument(db, assertions.WeightTicket.EmptyDocument, assertions)
	fullDocument := getOrCreateDocument(db, assertions.WeightTicket.FullDocument, assertions)
	trailerDocument := getOrCreateDocument(db, assertions.WeightTicket.FullDocument, assertions)

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
	assertions = ensureServiceMemberIsSetUpInAssertions(db, assertions)

	// Because this model points at multiple documents, it's not really good to point at the base assertions.Document,
	// so we'll look at assertions.WeightTicket.<Document>
	emptyDocument := getOrCreateDocumentWithUploads(db, assertions.WeightTicket.EmptyDocument, assertions)
	fullDocument := getOrCreateDocumentWithUploads(db, assertions.WeightTicket.FullDocument, assertions)

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

// ensureServiceMemberIsSetUpInAssertions checks for ServiceMember in assertions, or creates one if none exists. Several
// of the downstream functions need a service member, but they don't always share assertions, look at the same
// assertion, or create the service members in the same ways. We'll check now to see if we already have one created,
// and if not, create one that we can place in the assertions for all the rest.
func ensureServiceMemberIsSetUpInAssertions(db *pop.Connection, assertions Assertions) Assertions {
	if !assertions.Stub && assertions.ServiceMember.CreatedAt.IsZero() || assertions.ServiceMember.ID.IsNil() {
		serviceMember := MakeExtendedServiceMember(db, assertions)

		assertions.ServiceMember = serviceMember
		assertions.Order.ServiceMemberID = serviceMember.ID
		assertions.Order.ServiceMember = serviceMember
		assertions.Document.ServiceMemberID = serviceMember.ID
		assertions.Document.ServiceMember = serviceMember
	} else {
		assertions.Order.ServiceMemberID = assertions.ServiceMember.ID
		assertions.Order.ServiceMember = assertions.ServiceMember
		assertions.Document.ServiceMemberID = assertions.ServiceMember.ID
		assertions.Document.ServiceMember = assertions.ServiceMember
	}

	return assertions
}

// getOrCreateDocument checks if a document exists. If it does, it returns it, otherwise, it creates it
func getOrCreateDocument(db *pop.Connection, document models.Document, assertions Assertions) models.Document {
	if assertions.Stub && document.CreatedAt.IsZero() || document.ID.IsNil() {
		// Ensure our doc is associated with the expected ServiceMember
		document.ServiceMemberID = assertions.ServiceMember.ID
		document.ServiceMember = assertions.ServiceMember
		// Set generic Document to have the specific assertions that were passed in
		assertions.Document = document

		return MakeDocument(db, assertions)
	}

	return document
}

// getOrCreateUpload checks if an upload exists. If it does, it returns it, otherwise, it creates it.
func getOrCreateUpload(db *pop.Connection, upload models.UserUpload, assertions Assertions) models.UserUpload {
	if assertions.Stub && upload.CreatedAt.IsZero() || upload.ID.IsNil() {
		// Set generic UserUpload to have the specific assertions that were passed in
		assertions.UserUpload = upload

		return MakeUserUpload(db, assertions)
	}

	return upload
}

// getOrCreateDocumentWithUploads checks if a document exists. If it doesn't, it creates it. Then checks if the document
// has any uploads. If not, creates an upload associated with the document. Returns the document at the end. This
// function expects to get a specific document assertion since we're dealing with multiple documents in this overall
// file.
//
// Usage example:
//
//     emptyDocument := getOrCreateDocumentWithUploads(db, assertions.WeightTicket.EmptyDocument, assertions)
//
func getOrCreateDocumentWithUploads(db *pop.Connection, document models.Document, assertions Assertions) models.Document {
	// hang on to UserUploads, if any, for later
	userUploads := document.UserUploads

	// Ensure our doc is associated with the expected ServiceMember
	document.ServiceMemberID = assertions.ServiceMember.ID
	document.ServiceMember = assertions.ServiceMember

	doc := getOrCreateDocument(db, document, assertions)

	// Clear out doc.UserUploads because we'll be looping over the assertions that were passed in and potentially
	// creating data from those. It's easier to start with a clean slate than to track which ones were already created
	// vs which ones are newly created.
	doc.UserUploads = nil

	// Try getting or creating any uploads that were passed in via specific assertions
	for _, userUpload := range userUploads {
		// In case these weren't already set, set them so that they point at the correct document.
		userUpload.DocumentID = &doc.ID
		userUpload.Document = doc

		upload := getOrCreateUpload(db, userUpload, assertions)

		doc.UserUploads = append(doc.UserUploads, upload)
	}

	// If at the end we still don't have an upload, we'll just create the default one.
	if len(doc.UserUploads) == 0 {
		// This will be overriding the assertions locally only because we have a copy rather than a pointer
		assertions.UserUpload.DocumentID = &doc.ID
		assertions.UserUpload.Document = doc

		upload := MakeUserUpload(db, assertions)

		doc.UserUploads = append(doc.UserUploads, upload)
	}

	return doc
}
