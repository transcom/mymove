package factory

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildDutyLocation creates a single DutyLocation
func BuildDutyLocation(db *pop.Connection, customs []Customization, traits []Trait) models.DutyLocation {
	customs = setupCustomizations(customs, traits)

	// Find dutyLocation customization and extract the custom dutyLocation
	var cDutyLocation models.DutyLocation
	var address models.Address
	var transportationOffice models.TransportationOffice
	if result := findValidCustomization(customs, DutyLocation); result != nil {
		cDutyLocation = result.Model.(models.DutyLocation)
		if result.LinkOnly {
			return cDutyLocation
		}
		if cDutyLocation.TransportationOfficeID == nil {
			// BuildTransportationOffice Func goes here
			//transportationOffice = BuildTransportationOffice(db, customs, traits)
		}
	}

	//// Find a transportationOffice customization if available to BuildTransportationOffice
	//var transportationOffice models.TransportationOffice
	//if result := findValidCustomization(customs, TransportationOffice); result != nil {
	//	// Build Transportation Office func goes here
	//}

	// Create default Duty Location
	affiliation := internalmessages.AffiliationAIRFORCE
	location := models.DutyLocation{
		Name:        makeRandomString(10),
		Affiliation: &affiliation,
		AddressID:   address.ID,
		Address:     address,
		//TransportationOfficeID: &transportationOffice.ID,
		//TransportationOffice:   transportationOffice,
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&location, cDutyLocation)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &location)
	}

	return location

}
