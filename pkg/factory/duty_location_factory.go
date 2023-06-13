package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type dutyLocationBuildType byte

const (
	dutyLocationBuildStandard dutyLocationBuildType = iota
	dutyLocationBuildWithoutTransportationOffice
)

// buildDutyLocationWithBuildType does the actual work
// if buildType is standard, it builds
//   - DutyLocation
//   - TransportationOffice
//   - Address of the TO
//
// if buildType is withoutTransportationOffice, it builds
//   - DutyLocation
func buildDutyLocationWithBuildType(db *pop.Connection, customs []Customization, traits []Trait, buildType dutyLocationBuildType) models.DutyLocation {
	customs = setupCustomizations(customs, traits)

	// Find dutyLocation customization and extract the custom dutyLocation
	var cDutyLocation models.DutyLocation
	if result := findValidCustomization(customs, DutyLocation); result != nil {
		cDutyLocation = result.Model.(models.DutyLocation)
		if result.LinkOnly {
			return cDutyLocation
		}
	}

	// Create default Duty Location
	affiliation := internalmessages.AffiliationAIRFORCE

	// Prior to putting the duty location address information in the
	// duty location model, the duty location factory would create a
	// duty location address that contained street address 1, 2, and
	// 3. However, in practice, none of our duty locations have a
	// street address 2 or 3 and only 5 have a street address 1.
	//
	// Change the default to not include a street address
	location := models.DutyLocation{
		Name:           makeRandomString(10),
		Affiliation:    &affiliation,
		StreetAddress1: "n/a",
		City:           "Des Moines",
		State:          "IA",
		PostalCode:     "50309",
		Country:        "United States",
	}
	if buildType == dutyLocationBuildStandard {
		// Find/create the transportationOffice Model
		tempTOAddressCustoms := customs
		dltoAddress := findValidCustomization(customs, Addresses.DutyLocationTOAddress)
		if dltoAddress != nil {
			tempTOAddressCustoms = convertCustomizationInList(tempTOAddressCustoms, Addresses.DutyLocationTOAddress, Address)
		}
		transportationOffice := BuildTransportationOfficeWithPhoneLine(db, tempTOAddressCustoms, traits)
		location.TransportationOffice = transportationOffice
		location.TransportationOfficeID = &transportationOffice.ID
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&location, cDutyLocation)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		FetchOrBuildPostalCodeToGBLOC(db, location.PostalCode, "KKFA")
		mustCreate(db, &location)
	}

	return location
}

// BuildDutyLocation creates a single DutyLocation
// Also creates:
//   - TransportationOffice
//   - Address of the TO (use Addresses.DutyLocationTOAddress)
//
// Params:
//   - customs is a slice that will be modified by the factory
//   - db can be set to nil to create a stubbed model that is not stored in DB.
//
// Example:
//
//	dutyLocation := BuildDutyLocation(suite.DB(), []Customization{
//	       {Model: customDutyLocation},
//	       {Model: customTransportationOfficeAddress, Type: &Addresses.DutyLocationTOAddress},
//	       }, nil)
func BuildDutyLocation(db *pop.Connection, customs []Customization, traits []Trait) models.DutyLocation {
	return buildDutyLocationWithBuildType(db, customs, traits, dutyLocationBuildStandard)
}

// BuildDutyLocationWithoutTransportationOffice returns a duty location without a transportation office.
// Also creates:
//   - Address of the DL (use Addresses.DutyLocationAddress)
//   - Will not create a Transportation Office even if one is supplied in the customizations or traits
//
// Params:
//   - customs is a slice that will be modified by the factory
//   - db can be set to nil to create a stubbed model that is not stored in DB.
//
// Example:
//
//	dutyLocation := BuildDutyLocationWithoutTransportationOffice(suite.DB(), []Customization{
//	       {Model: customDutyLocation},
//	       }, nil)
func BuildDutyLocationWithoutTransportationOffice(db *pop.Connection, customs []Customization, traits []Trait) models.DutyLocation {
	return buildDutyLocationWithBuildType(db, customs, traits, dutyLocationBuildWithoutTransportationOffice)
}

// FetchOrBuildCurrentDutyLocation returns a default duty location
// It always fetches or builds a Yuma AFB Duty Location
func FetchOrBuildCurrentDutyLocation(db *pop.Connection) models.DutyLocation {
	if db == nil {
		return BuildDutyLocation(nil, []Customization{
			{
				Model: models.DutyLocation{
					Name: "Yuma AFB",
				},
			},
		}, nil)
	}
	// Check if Yuma Duty Location exists, if not, create it.
	defaultLocation, err := models.FetchDutyLocationByName(db, "Yuma AFB")
	if err != nil {
		return BuildDutyLocation(db, []Customization{
			{
				Model: models.DutyLocation{
					Name: "Yuma AFB",
				},
			},
		}, nil)
	}

	return defaultLocation
}

// FetchOrBuildOrdersDutyLocation returns a default orders duty location
// It always fetches or builds a Fort Gordon duty location with the
// specified city/state/postal code
// Some tests rely on the duty location being in 30813
func FetchOrBuildOrdersDutyLocation(db *pop.Connection) models.DutyLocation {
	if db == nil {
		return BuildDutyLocation(nil, GetTraitDefaultOrdersDutyLocation(), nil)
	}

	// Check if we already have a Fort Gordon Duty Location, return it if so
	fortGordon, err := models.FetchDutyLocationByName(db, "Fort Gordon")
	if err == nil {
		return fortGordon
	}

	return BuildDutyLocation(db, nil, []Trait{GetTraitDefaultOrdersDutyLocation})
}

func GetTraitDefaultOrdersDutyLocation() []Customization {
	return []Customization{
		{
			// Update the DutyLocationAddress (but not TO address) to Augusta, Georgia
			Model: models.DutyLocation{
				Name:           "Fort Gordon",
				StreetAddress1: "Fort Gordon",
				City:           "Augusta",
				State:          "GA",
				PostalCode:     "30813",
				Country:        "US",
			},
		},
	}
}
