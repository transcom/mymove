package testdatagen

import (
	"log"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// MakeDutyLocation creates a single DutyLocation
func MakeDutyLocation(db *pop.Connection, assertions Assertions) models.DutyLocation {
	transportationOffice := assertions.DutyLocation.TransportationOffice
	if assertions.DutyLocation.TransportationOfficeID == nil {
		transportationOffice = MakeTransportationOffice(db, assertions)
	}

	address := assertions.DutyLocation.Address
	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.DutyLocation.AddressID) {
		address = MakeAddress3(db, assertions)

		// Make the required Tariff 400 NG Zip3 to correspond with the duty location address
		FetchOrMakeDefaultTariff400ngZip3(db)
		FetchOrMakeTariff400ngZip3(db, Assertions{
			Tariff400ngZip3: models.Tariff400ngZip3{
				Zip3:          "503",
				BasepointCity: "Des Moines",
				State:         "IA",
				ServiceArea:   "296",
				RateArea:      "US53",
				Region:        "7",
			},
		})
	}
	affiliation := internalmessages.AffiliationAIRFORCE
	location := models.DutyLocation{
		Name:                   MakeRandomString(10),
		Affiliation:            &affiliation,
		AddressID:              address.ID,
		Address:                address,
		TransportationOfficeID: &transportationOffice.ID,
		TransportationOffice:   transportationOffice,
	}
	mergeModels(&location, assertions.DutyLocation)

	mustCreate(db, &location, assertions.Stub)

	return location
}

// MakeDefaultDutyLocation makes a DutyLocation with default values
func MakeDefaultDutyLocation(db *pop.Connection) models.DutyLocation {
	return MakeDutyLocation(db, Assertions{})
}

// FetchOrMakeDefaultCurrentDutyLocation returns a default duty location - Yuma AFB
func FetchOrMakeDefaultCurrentDutyLocation(db *pop.Connection) models.DutyLocation {
	defaultLocation, err := models.FetchDutyLocationByName(db, "Yuma AFB")
	if err == nil {
		return defaultLocation
	}

	// Now that playwright tests create data on demand, it's possible
	// multiple tests will try to fetch or create the current duty
	// location simultaneously. If we do nothing, we can get failures
	// from the race condition of two different tests calling this at
	// the same time.
	//
	// *sigh*, pop doesn't know about nested transactions, so manage
	// it ourselves.
	//
	// Assume we are in a transation so we can start a postgresql
	// SAVEPOINT (aka nested transaction)
	beginSavepoint := "SAVEPOINT default_duty_location"
	commitSavepoint := "RELEASE SAVEPOINT default_duty_location"
	err = db.RawQuery(beginSavepoint).Exec()
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
	defaultLocation, err = models.FetchDutyLocationByName(db, "Yuma AFB")
	if err != nil {
		defaultLocation = MakeDutyLocation(db, Assertions{
			DutyLocation: models.DutyLocation{
				Name: "Yuma AFB",
			}})
	}
	err = db.RawQuery(commitSavepoint).Exec()
	if err != nil {
		log.Fatalf("Error commit duty location savepoint/tx: %s", err)
	}

	return defaultLocation
}

// FetchOrMakeDefaultNewOrdersDutyLocation returns a default duty location - Yuma AFB
func FetchOrMakeDefaultNewOrdersDutyLocation(db *pop.Connection) models.DutyLocation {
	// Check if Fort Gordon exists, if not, create
	// Move date picker for this test case only works with an address of street name "Fort Gordon"
	fortGordon, err := models.FetchDutyLocationByName(db, "Fort Gordon")
	if err == nil {
		fortGordon.TransportationOffice, err = models.FetchDutyLocationTransportationOffice(db, fortGordon.ID)
		if err == nil {
			return fortGordon
		}
	}
	fortGordonAssertions := Assertions{
		Address: models.Address{
			City:       "Augusta",
			State:      "GA",
			PostalCode: "30813",
		},
		DutyLocation: models.DutyLocation{
			Name: "Fort Gordon",
		},
	}
	FetchOrMakeTariff400ngZip3(db, Assertions{
		Tariff400ngZip3: models.Tariff400ngZip3{
			Zip3:          "308",
			BasepointCity: "Harlem",
			State:         "GA",
			ServiceArea:   "208",
			RateArea:      "US45",
			Region:        "12",
		},
	})
	return MakeDutyLocation(db, fortGordonAssertions)
}
