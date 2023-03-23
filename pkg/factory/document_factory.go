package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildDocument creates a single Document.
// Also creates, if not provided
// - ServiceMember
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildDocument(db *pop.Connection, customs []Customization, traits []Trait) models.Document {
	customs = setupCustomizations(customs, traits)

	// Find Document assertion and convert to models.Document
	var cDocument models.Document
	if result := findValidCustomization(customs, Document); result != nil {
		cDocument = result.Model.(models.Document)
		if result.LinkOnly {
			return cDocument
		}
	}

	// Find/create the ServiceMember model
	serviceMember := BuildServiceMember(db, customs, traits)

	// Create document
	document := models.Document{
		ServiceMemberID: serviceMember.ID,
		ServiceMember:   serviceMember,
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&document, cDocument)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &document)
	}
	return document
}

func BuildDocumentLinkServiceMember(db *pop.Connection, serviceMember models.ServiceMember) models.Document {
	return BuildDocument(db, []Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
	}, nil)
}
