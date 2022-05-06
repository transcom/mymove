package testdatagen

import (
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
		Name:                   makeRandomString(10),
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
	// Check if Yuma Duty Location exists, if not, create it.
	defaultLocation, err := models.FetchDutyLocationByName(db, "Yuma AFB")
	if err != nil {
		return MakeDutyLocation(db, Assertions{
			DutyLocation: models.DutyLocation{
				Name: "Yuma AFB",
			}})
	}
	return defaultLocation
}

// FetchOrMakeDefaultNewOrdersDutyLocation returns a default duty location - Yuma AFB
func FetchOrMakeDefaultNewOrdersDutyLocation(db *pop.Connection) models.DutyLocation {
	// Check if Fort Gordon exists, if not, create
	// Move date picker for this test case only works with an address of street name "Fort Gordon"
	fortGordon, err := models.FetchDutyLocationByName(db, "Fort Gordon")
	if err != nil {
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
	return fortGordon
}
