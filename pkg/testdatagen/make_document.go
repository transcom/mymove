package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// makeDocument creates a single Document.
//
// Deprecated: use factory.BuildDocument
func makeDocument(db *pop.Connection, assertions Assertions) (models.Document, error) {
	sm := assertions.Document.ServiceMember
	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.Document.ServiceMemberID) {
		var err error
		sm, err = makeServiceMember(db, assertions)
		if err != nil {
			return models.Document{}, err
		}
	}

	document := models.Document{
		ServiceMemberID: sm.ID,
		ServiceMember:   sm,
	}

	// Overwrite values with those from assertions
	mergeModels(&document, assertions.Document)

	mustCreate(db, &document, assertions.Stub)

	return document, nil
}
