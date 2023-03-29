package factory

import (
	"log"

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
	tempAddressCustoms := customs
	result := findValidCustomization(customs, Addresses.DutyLocationAddress)
	if result != nil {
		tempAddressCustoms = convertCustomizationInList(tempAddressCustoms, Addresses.DutyLocationAddress, Address)
	}
	dlAddress := BuildAddress(db, tempAddressCustoms, []Trait{GetTraitAddress3})

	if db != nil {
		FindOrBuildPostalCodeToGBLOC(db, dlAddress.PostalCode, "KKFA")
	}

	// Find/create the transportationOffice Model
	tempTOAddressCustoms := customs
	dltoAddress := findValidCustomization(customs, Addresses.DutyLocationTOAddress)
	if dltoAddress != nil {
		tempTOAddressCustoms = convertCustomizationInList(tempTOAddressCustoms, Addresses.DutyLocationTOAddress, Address)
	}
	transportationOffice := BuildTransportationOfficeWithPhoneLine(db, tempTOAddressCustoms, traits)

	tarifCustoms := findValidCustomization(customs, Tariff400ngZip3)
	if tarifCustoms == nil {
		// Build the required Tariff 400 NG Zip3 to correspond with the
		// duty location address
		tarifCustoms = &Customization{
			Model: models.Tariff400ngZip3{
				Zip3:          "503",
				BasepointCity: "Des Moines",
				State:         "IA",
				ServiceArea:   "296",
				RateArea:      "US53",
				Region:        "7",
			},
		}
	}
	FetchOrBuildTariff400ngZip3(db, []Customization{*tarifCustoms}, nil)

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
	if db == nil {
		return BuildDutyLocation(nil, []Customization{
			{
				Model: models.DutyLocation{
					Name: "Yuma AFB",
				},
			},
		}, nil)
	}
	// Now that playwright tests create data on demand, it's possible
	// multiple tests will try to fetch or create the current duty
	// location simultaneously. If we do nothing, we can get failures
	// from the race condition of two different tests calling this at
	// the same time.
	//
	cleanupFunc := exclusiveDutyLocationLock(db)
	defer cleanupFunc()
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
	if db == nil {
		return BuildDutyLocation(nil, []Customization{
			{
				Model: models.DutyLocation{
					Name: "Fort Gordon",
				},
			},
		}, nil)
	}
	// Now that playwright tests create data on demand, it's possible
	// multiple tests will try to fetch or create the current duty
	// location simultaneously. If we do nothing, we can get failures
	// from the race condition of two different tests calling this at
	// the same time.
	//
	cleanupFunc := exclusiveDutyLocationLock(db)
	defer cleanupFunc()

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
	}
}

// exclusiveDutyLocationLock locks the duty_locations table in a savepoint
func exclusiveDutyLocationLock(db *pop.Connection) func() {
	// *sigh*, pop doesn't know about nested transactions, so manage
	// it ourselves.
	//
	// Assume we are in a transation so we can start a postgresql
	// SAVEPOINT (aka nested transaction)
	beginSavepoint := "SAVEPOINT duty_location"
	commitSavepoint := "RELEASE SAVEPOINT duty_location"
	err := db.RawQuery(beginSavepoint).Exec()
	if err != nil {
		log.Fatalf("Error starting duty location savepoint/txn: %s", err)
	}
	// lock the table exclusively, fetch to make sure no one has beat
	// us to it, and then create if necessary. This is not the most
	// performant way, but this is for tests and so being slightly
	// slower than theoritically optimal is ok.
	//
	// Use EXCLUSIVE lock so reads can happen, but not writes
	// https://www.postgresql.org/docs/current/explicit-locking.html
	err = db.RawQuery("LOCK TABLE duty_locations IN EXCLUSIVE MODE").Exec()
	if err != nil {
		log.Fatalf("Error locking duty location table: %s", err)
	}

	return func() {
		if err := db.RawQuery(commitSavepoint).Exec(); err != nil {
			log.Fatalf("Error commit duty location savepoint/tx: %s", err)
		}
	}
}
