package factory

import (
	"fmt"

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
//   - Address of the DL
//   - TransportationOffice
//   - Address of the TO
//
// if buildType is withoutTransportationOffice, it builds
//   - DutyLocation
//   - Address of the DL
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

	// Find/create the DutyLocationAddress
	tempAddressCustoms := customs
	result := findValidCustomization(customs, Addresses.DutyLocationAddress)
	if result != nil {
		tempAddressCustoms = convertCustomizationInList(tempAddressCustoms, Addresses.DutyLocationAddress, Address)
	}
	dlAddress := BuildAddress(db, tempAddressCustoms, []Trait{GetTraitAddress3})

	if db != nil {
		FetchOrBuildPostalCodeToGBLOC(db, dlAddress.PostalCode, "KKFA")
	}

	// Create default Duty Location
	affiliation := internalmessages.AffiliationAIRFORCE

	location := models.DutyLocation{
		// Make test duty location name consistent with how TRDM duty location is formatted
		Name:        fmt.Sprintf("%s, %s %s", MakeRandomString(10), dlAddress.State, dlAddress.PostalCode),
		Affiliation: &affiliation,
		AddressID:   dlAddress.ID,
		Address:     dlAddress,
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
		mustCreate(db, &location)
	}

	return location
}

// BuildDutyLocation creates a single DutyLocation
// Also creates:
//   - Address of the DL (use Addresses.DutyLocationAddress)
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
//	       {Model: customDutyLocationAddress, Type: &Addresses.DutyLocationAddress},
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
//	       {Model: customDutyLocationAddress, Type: &Addresses.DutyLocationAddress},
//	       }, nil)
func BuildDutyLocationWithoutTransportationOffice(db *pop.Connection, customs []Customization, traits []Trait) models.DutyLocation {
	return buildDutyLocationWithBuildType(db, customs, traits, dutyLocationBuildWithoutTransportationOffice)
}

// FetchOrBuildOtherDutyLocation returns a duty location other
// than the default. It always fetches or builds a Bessemer, AL Duty Location
func FetchOrBuildOtherDutyLocation(db *pop.Connection) models.DutyLocation {
	bessemerAddress := models.Address{
		StreetAddress1: "123 Main St",
		City:           "Bessemer",
		State:          "AL",
		PostalCode:     "35023",
	}
	dutyLoc := models.DutyLocation{
		Name:    "Bessemer, AL 35023",
		Address: bessemerAddress,
	}

	if db == nil {
		return BuildDutyLocation(nil, []Customization{
			{
				Model: dutyLoc,
			},
		}, nil)
	}
	// Check if Bessemer, AL Location exists, if not, create it.
	defaultLocation, err := models.FetchDutyLocationByName(db, "Bessemer, AL 35023")
	if err != nil {
		return BuildDutyLocation(db, []Customization{
			{
				Model: dutyLoc,
			},
		}, nil)
	}

	return defaultLocation
}

// FetchOrBuildCurrentDutyLocation returns a default duty location
// It always fetches or builds a Yuma AFB Duty Location
func FetchOrBuildCurrentDutyLocation(db *pop.Connection) models.DutyLocation {
	if db == nil {
		return BuildDutyLocation(nil, []Customization{
			{
				Model: models.DutyLocation{
					Name: "Yuma AFB, IA 50309",
				},
			},
		}, nil)
	}
	// Check if Yuma Duty Location exists, if not, create it.
	defaultLocation, err := models.FetchDutyLocationByName(db, "Yuma AFB, IA 50309")
	if err != nil {
		return BuildDutyLocation(db, []Customization{
			{
				Model: models.DutyLocation{
					Name: "Yuma AFB, IA 50309",
				},
			},
		}, nil)
	}

	return defaultLocation
}

// FetchOrBuildOrdersDutyLocation returns a default orders duty location
// It always fetches or builds a Fort Eisenhower duty location with the specified city/state/postal code
// Some tests rely on the duty location being in 30813
func FetchOrBuildOrdersDutyLocation(db *pop.Connection) models.DutyLocation {
	if db == nil {
		return BuildDutyLocation(nil, []Customization{
			{
				Model: models.DutyLocation{
					Name: "Fort Eisenhower, GA 30813",
				},
			},
			{
				Model: models.Address{
					City:       "Fort Eisenhower",
					State:      "GA",
					PostalCode: "30813",
				},
				Type: &Addresses.DutyLocationAddress,
			},
		}, nil)
	}

	// Check if we already have a Fort Eisenhower Duty Location, return it if so
	fortEisenhower, err := models.FetchDutyLocationByName(db, "Fort Eisenhower, GA 30813")
	if err == nil {
		return fortEisenhower
	}

	return BuildDutyLocation(db, nil, []Trait{GetTraitDefaultOrdersDutyLocation})
}

func GetTraitDefaultOrdersDutyLocation() []Customization {
	return []Customization{
		{
			Model: models.DutyLocation{
				Name: "Fort Eisenhower, GA 30813",
			},
		},
		{
			// Update the DutyLocationAddress (but not TO address) to Fort Eisenhower, Georgia
			Type: &Addresses.DutyLocationAddress,
			Model: models.Address{
				City:       "Fort Eisenhower",
				State:      "GA",
				PostalCode: "30813",
			},
		},
	}
}
