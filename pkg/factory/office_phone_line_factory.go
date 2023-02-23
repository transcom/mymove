package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildOfficePhoneLine creates a OfficePhoneLine
// Also creates, if not provided
// - TransportationOffice
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildOfficePhoneLine(db *pop.Connection, customs []Customization, traits []Trait) models.OfficePhoneLine {
	customs = setupCustomizations(customs, traits)

	// Find OfficePhoneLine assertion and convert to models officephoneline
	var cOfficePhoneLine models.OfficePhoneLine
	if result := findValidCustomization(customs, OfficePhoneLine); result != nil {
		cOfficePhoneLine = result.Model.(models.OfficePhoneLine)
		if result.LinkOnly {
			return cOfficePhoneLine
		}
	}

	// Create the associated transportationOffice model
	office := BuildTransportationOffice(db, customs, nil)

	// create officephoneline
	phoneLine := models.OfficePhoneLine{
		TransportationOfficeID: office.ID,
		TransportationOffice:   office,
		Number:                 "(510) 555-5555",
		IsDsnNumber:            false,
		Type:                   "voice",
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&phoneLine, cOfficePhoneLine)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &phoneLine)
	}

	return phoneLine
}
