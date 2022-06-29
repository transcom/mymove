package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakeDocument creates a single Document.
func MakeDocument(db *pop.Connection, assertions Assertions) models.Document {
	sm := assertions.ServiceMember
	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.ServiceMember.ID) {
		sm = MakeServiceMember(db, assertions)
	}

	document := models.Document{
		ServiceMemberID: sm.ID,
		ServiceMember:   sm,
	}

	// Overwrite values with those from assertions
	mergeModels(&document, assertions.Document)

	mustCreate(db, &document, assertions.Stub)

	return document
}

// MakeDefaultDocument returns a Document with default values
func MakeDefaultDocument(db *pop.Connection) models.Document {
	return MakeDocument(db, Assertions{})
}
