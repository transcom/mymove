package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

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
	var dlAddress models.Address
	var tempAddressCustoms = customs
	result := findValidCustomization(customs, Addresses.DutyLocationAddress)
	if result != nil {
		dlAddress = result.Model.(models.Address)
		tempAddressCustoms = convertCustomizationInList(tempAddressCustoms, Addresses.DutyLocationAddress, Address)
	}
	dlAddress = BuildAddress(db, tempAddressCustoms, []Trait{GetTraitAddress3})

	// Find/create the transportationOffice Model
	var transportationOffice models.TransportationOffice
	var tempTOAddressCustoms = customs
	result = findValidCustomization(customs, TransportationOffice)
	dltoAddress := findValidCustomization(customs, Addresses.DutyLocationTOAddress)
	if result != nil {
		transportationOffice = result.Model.(models.TransportationOffice)
	}

	if dltoAddress != nil {
		tempTOAddressCustoms = convertCustomizationInList(tempTOAddressCustoms, Addresses.DutyLocationTOAddress, Address)
	}

	transportationOffice = BuildTransportationOfficeWithPhoneLine(db, tempTOAddressCustoms, traits)

	// Build the required Tariff 400 NG Zip3 to correspond with the duty location address
	FetchOrBuildTariff400ngZip3(db, []Customization{
		{
			Model: models.Tariff400ngZip3{
				Zip3:          "503",
				BasepointCity: "Des Moines",
				State:         "IA",
				ServiceArea:   "296",
				RateArea:      "US53",
				Region:        "7",
			},
		},
	}, nil)

	// Create default Duty Location
	affiliation := internalmessages.AffiliationAIRFORCE
	location := models.DutyLocation{
		Name:                   makeRandomString(10),
		Affiliation:            &affiliation,
		AddressID:              dlAddress.ID,
		Address:                dlAddress,
		TransportationOfficeID: &transportationOffice.ID,
		TransportationOffice:   transportationOffice,
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&location, cDutyLocation)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &location)
	}

	return location

}

// FetchOrBuildCurrentDutyLocation returns a default duty location
// It always fetches or builds a Yuma AFB Duty Location
func FetchOrBuildCurrentDutyLocation(db *pop.Connection) models.DutyLocation {
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
// It always fetches or builds a Fort Gordon duty location
// It also creates a GA 208 tariff
func FetchOrBuildOrdersDutyLocation(db *pop.Connection) models.DutyLocation {
	// Check if we already have a Fort Gordon Duty Location, return it if so
	fortGordon, err := models.FetchDutyLocationByName(db, "Fort Gordon")
	if err == nil {
		return fortGordon
	}

	// If not, build the Fort Gordon Duty location with the associated
	// address and tariff
	FetchOrBuildTariff400ngZip3(db, []Customization{
		{
			Model: models.Tariff400ngZip3{
				Zip3:          "308",
				BasepointCity: "Harlem",
				State:         "GA",
				ServiceArea:   "208",
				RateArea:      "US45",
				Region:        "12",
			},
		},
	}, nil)

	return BuildDutyLocation(db, []Customization{
		{
			Model: models.DutyLocation{
				Name: "Fort Gordon",
			},
		},
		{
			// Update the DutyLocationAddress (but not TO address) to Augusta, Georgia
			Type: &Addresses.DutyLocationAddress,
			Model: models.Address{
				City:       "Augusta",
				State:      "GA",
				PostalCode: "30813",
			},
		},
	}, nil)
}
