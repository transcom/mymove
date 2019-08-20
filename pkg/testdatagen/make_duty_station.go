package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// MakeDutyStation creates a single DutyStation
func MakeDutyStation(db *pop.Connection, assertions Assertions) models.DutyStation {
	transportationOffice := assertions.DutyStation.TransportationOffice
	if assertions.DutyStation.TransportationOfficeID == nil {
		transportationOffice = MakeTransportationOffice(db, assertions)
	}

	address := assertions.DutyStation.Address
	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.DutyStation.AddressID) {
		address = MakeAddress3(db, assertions)

		// Make the required Tariff 400 NG Zip3 to correspond with the duty station address
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

	station := models.DutyStation{
		Name:                   "Yuma AFB",
		Affiliation:            internalmessages.AffiliationAIRFORCE,
		AddressID:              address.ID,
		Address:                address,
		TransportationOfficeID: &transportationOffice.ID,
		TransportationOffice:   transportationOffice,
	}

	mergeModels(&station, assertions.DutyStation)

	mustCreate(db, &station)

	return station
}

// FetchOrMakeDefaultCurrentDutyStation returns a default duty station - Yuma AFB
func FetchOrMakeDefaultCurrentDutyStation(db *pop.Connection) models.DutyStation {
	// Check if Yuma Duty Station exists, if not, create it.
	defaultStation, err := models.FetchDutyStationByName(db, "Yuma AFB")
	if err != nil {
		return MakeDutyStation(db, Assertions{})
	}
	return defaultStation
}

// FetchOrMakeDefaultNewOrdersDutyStation returns a default duty station - Yuma AFB
func FetchOrMakeDefaultNewOrdersDutyStation(db *pop.Connection) models.DutyStation {
	// Check if Fort Gordon exists, if not, create
	// Move date picker for this test case only works with an address of street name "Fort Gordon"
	fortGordon, err := models.FetchDutyStationByName(db, "Fort Gordon")
	if err != nil {
		fortGordonAssertions := Assertions{
			Address: models.Address{
				City:       "Augusta",
				State:      "GA",
				PostalCode: "30813",
			},
			DutyStation: models.DutyStation{
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
		return MakeDutyStation(db, fortGordonAssertions)
	}
	return fortGordon
}
