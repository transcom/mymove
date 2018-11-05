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
		MakeDefaultTariff400ngZip3(db)
		MakeTariff400ngZip3(db, Assertions{
			Tariff400ngZip3: models.Tariff400ngZip3{
				Zip3:          "945",
				BasepointCity: "Walnut Creek",
				State:         "CA",
				ServiceArea:   "80",
				RateArea:      "US87",
				Region:        "2",
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

// MakeDefaultDutyStation returns a duty station with default info
func MakeDefaultDutyStation(db *pop.Connection) models.DutyStation {
	return MakeDutyStation(db, Assertions{})
}
