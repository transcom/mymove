package testdatagen

import (
	"log"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// MakeDutyStation creates a single DutyStation
func MakeDutyStation(db *pop.Connection, name string, affiliation internalmessages.Affiliation, address models.Address) (models.DutyStation, error) {
	transportationOffice, err := MakeTransportationOffice(db)
	if err != nil {
		log.Panic(err)
	}

	var station models.DutyStation
	verrs, err := db.ValidateAndSave(&address)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() != 0 {
		log.Panic(verrs.Error())
	}

	station = models.DutyStation{
		Name:                   name,
		Affiliation:            affiliation,
		AddressID:              address.ID,
		Address:                address,
		TransportationOfficeID: &transportationOffice.ID,
		TransportationOffice:   &transportationOffice,
	}

	verrs, err = db.ValidateAndSave(&station)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() != 0 {
		log.Panic(verrs.Error())
	}

	return station, err
}

// MakeAnyDutyStation returns a duty station with dummy info
func MakeAnyDutyStation(db *pop.Connection) models.DutyStation {
	station, _ := MakeDutyStation(db, "Air Station Yuma", internalmessages.AffiliationMARINES,
		models.Address{StreetAddress1: "duty station", City: "Yuma", State: "Arizona", PostalCode: "85364"})
	return station
}

// MakeDutyStationWithoutTransportationOffice returns a duty station with dummy info and no office
func MakeDutyStationWithoutTransportationOffice(db *pop.Connection, name string, affiliation internalmessages.Affiliation, address models.Address) (models.DutyStation, error) {
	var station models.DutyStation
	verrs, err := db.ValidateAndSave(&address)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() != 0 {
		log.Panic(verrs.Error())
	}

	station = models.DutyStation{
		Name:        name,
		Affiliation: affiliation,
		AddressID:   address.ID,
		Address:     address,
	}

	verrs, err = db.ValidateAndSave(&station)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() != 0 {
		log.Panic(verrs.Error())
	}

	return station, err
}
