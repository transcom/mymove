package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeDocument creates a single Document.
func MakeDocument(db *pop.Connection, assertions Assertions) models.Document {
	sm := assertions.Document.ServiceMember
	if isZeroUUID(assertions.Document.ServiceMemberID) {
		sm = MakeServiceMember(db, assertions)
	}

	document := models.Document{
		ServiceMemberID: sm.ID,
		ServiceMember:   sm,
		Name:            "Default name",
	}

	// Overwrite values with those from assertions
	mergeModels(&document, assertions.Document)

	mustCreate(db, &document)

	return document
}

// MakeDefaultDocument returns a Document with default values
func MakeDefaultDocument(db *pop.Connection) models.Document {
	return MakeDocument(db, Assertions{})
}
