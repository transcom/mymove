package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildServiceRequestDocument creates ServiceRequestDocument.
// Also creates, if not provided
// - MTOServiceItem and associated set relationships
//
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildServiceRequestDocument(db *pop.Connection, customs []Customization, traits []Trait) models.ServiceRequestDocument {
	customs = setupCustomizations(customs, traits)

	// Find upload assertion and convert to models upload
	var cServiceRequestDocument models.ServiceRequestDocument
	if result := findValidCustomization(customs, ServiceRequestDocument); result != nil {
		cServiceRequestDocument = result.Model.(models.ServiceRequestDocument)

		if result.LinkOnly {
			return cServiceRequestDocument
		}
	}

	mtoServiceItem := BuildMTOServiceItemBasic(db, customs, traits)

	ServiceRequestDocument := models.ServiceRequestDocument{
		MTOServiceItem:   mtoServiceItem,
		MTOServiceItemID: mtoServiceItem.ID,
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&ServiceRequestDocument, ServiceRequestDocument)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &ServiceRequestDocument)
	}

	return ServiceRequestDocument
}
