package testdatagen

import (
	"log"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// makeDutyLocation creates a single DutyLocation
//
// Deprecated: use factory.BuildDutyLocation
func makeDutyLocation(db *pop.Connection, assertions Assertions) (models.DutyLocation, error) {
	transportationOffice := assertions.DutyLocation.TransportationOffice
	if assertions.DutyLocation.TransportationOfficeID == nil {
		var err error
		transportationOffice, err = MakeTransportationOffice(db, assertions)
		if err != nil {
			return models.DutyLocation{}, err
		}
	}

	address := assertions.DutyLocation.Address
	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.DutyLocation.AddressID) {
		var err error
		address, err = MakeAddress3(db, assertions)
		if err != nil {
			return models.DutyLocation{}, err
		}
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

	return location, nil
}

// fetchOrMakeDefaultCurrentDutyLocation returns a default duty
// location - Yuma AFB
//
// Deprecated: use factory.FetchOrMakeDefaultCurrentDutyLocation
func fetchOrMakeDefaultCurrentDutyLocation(db *pop.Connection) (models.DutyLocation, error) {
	defaultLocation, err := models.FetchDutyLocationByName(db, "Yuma AFB")
	if err == nil {
		return defaultLocation, nil
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
		var errResponse error
		defaultLocation, errResponse = makeDutyLocation(db, Assertions{
			DutyLocation: models.DutyLocation{
				Name: "Yuma AFB",
			},
		})

		if errResponse != nil {
			return models.DutyLocation{}, errResponse
		}
	}
	err = db.RawQuery(commitSavepoint).Exec()
	if err != nil {
		log.Fatalf("Error commit duty location savepoint/tx: %s", err)
	}

	return defaultLocation, nil
}

// fetchOrMakeDefaultNewOrdersDutyLocation returns a default duty
// location - Yuma AFB
//
// Deprecated: use factory.fetchOrMakeDefaultNewOrdersDutyLocation
func fetchOrMakeDefaultNewOrdersDutyLocation(db *pop.Connection) (models.DutyLocation, error) {
	// Check if Fort Eisenhower exists, if not, create
	// Move date picker for this test case only works with an address of street name "Fort Eisenhower"
	fortEisenhower, err := models.FetchDutyLocationByName(db, "Fort Eisenhower, GA 30813")
	if err == nil {
		fortEisenhower.TransportationOffice, err = models.FetchDutyLocationTransportationOffice(db, fortEisenhower.ID)
		if err == nil {
			return fortEisenhower, nil
		}
	}

	fortEisenhowerAssertions := Assertions{
		Address: models.Address{
			City:       "GROVETOWN",
			State:      "GA",
			PostalCode: "30813",
		},
		DutyLocation: models.DutyLocation{
			Name: "Fort Eisenhower, GA 30813",
		},
	}

	return makeDutyLocation(db, fortEisenhowerAssertions)
}
