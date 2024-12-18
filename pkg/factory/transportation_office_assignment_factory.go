package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildPrimaryTransportationOfficeAssignment creates a Transportation Assignment, and a transportation office and officer user if either doesn't exist
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
// Notes:
//   - Marks the transportation office assignment as the office user's primary transportation office,
//     use BuildAlternateTransportationOfficeAssignment for non-primary transportation office.
func BuildPrimaryTransportationOfficeAssignment(db *pop.Connection, customs []Customization, traits []Trait) models.TransportationOfficeAssignment {
	customs = setupCustomizations(customs, traits)

	// Find TransportationAssignment assertion and convert to models.TransportationOfficeAssignment
	var cTransportationOfficeAssignment models.TransportationOfficeAssignment
	if result := findValidCustomization(customs, TransportationOfficeAssignment); result != nil {
		cTransportationOfficeAssignment = result.Model.(models.TransportationOfficeAssignment)
		if result.LinkOnly {
			return cTransportationOfficeAssignment
		}
	}

	// Find/Create the associated transportation office model
	transportationOffice := BuildTransportationOffice(db, customs, traits)

	// Find/Create the associated office user model
	officeUser := BuildOfficeUser(db, customs, traits)

	// Create transportationOffice
	transportationOfficeAssignment := models.TransportationOfficeAssignment{
		ID:                     officeUser.ID,
		TransportationOfficeID: transportationOffice.ID,
		TransportationOffice:   transportationOffice,
		PrimaryOffice:          models.BoolPointer(true),
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&transportationOfficeAssignment, cTransportationOfficeAssignment)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &transportationOfficeAssignment)
	}
	return transportationOfficeAssignment
}

// BuildAlternateTransportationAssignment creates a Transportation Assignment, and a transportation office and officer user if either doesn't exist
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
// Notes:
//   - Marks the transportation office assignment as a non-primary transportation office assignment,
//     use BuildPrimaryTransportationOfficeAssignment for primary transportation office assignments.
func BuildAlternateTransportationOfficeAssignment(db *pop.Connection, customs []Customization, traits []Trait) models.TransportationOfficeAssignment {
	customs = setupCustomizations(customs, traits)

	// Find TransportationAssignment assertion and convert to models.TransportationAssignment
	var cTransportationOfficeAssignment models.TransportationOfficeAssignment
	if result := findValidCustomization(customs, TransportationOfficeAssignment); result != nil {
		cTransportationOfficeAssignment = result.Model.(models.TransportationOfficeAssignment)
		if result.LinkOnly {
			return cTransportationOfficeAssignment
		}
	}

	// Find/Create the associated office user model
	officeUser := BuildOfficeUser(db, customs, traits)

	// Find/Create the associated transportation office model
	transportationOffice := BuildTransportationOffice(db, customs, traits)

	// Create transportationOffice
	transportationOfficeAssignment := models.TransportationOfficeAssignment{
		ID:                     officeUser.ID,
		TransportationOfficeID: transportationOffice.ID,
		TransportationOffice:   transportationOffice,
		PrimaryOffice:          models.BoolPointer(false),
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&transportationOfficeAssignment, cTransportationOfficeAssignment)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &transportationOfficeAssignment)
	}
	return transportationOfficeAssignment
}
