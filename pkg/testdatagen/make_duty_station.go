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
		address = MakeAddress(db, assertions)
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

// FetchOrMakeDefaultDutyStation returns a default duty station - Yuma AFB
func FetchOrMakeDefaultDutyStation(db *pop.Connection) models.DutyStation {
	// Check if Yuma Duty Station exists, if not, create it.
	defaultStation, err := models.FetchDutyStationByName(db, "Yuma AFB")
	if err != nil {
		return MakeDutyStation(db, Assertions{})
	}
	return defaultStation
}
