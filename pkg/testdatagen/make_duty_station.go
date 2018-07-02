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
		transportationOffice = MakeTransportationOffice(db /*, assertions*/)
	}

	address := assertions.DutyStation.Address
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

// MakeDefaultDutyStation returns a duty station with default info
func MakeDefaultDutyStation(db *pop.Connection) models.DutyStation {
	return MakeDutyStation(db, Assertions{})
}
