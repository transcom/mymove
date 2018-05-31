package testdatagen

import (
	"log"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeTransportationOffice creates a single ServiceMember and associated User.
func MakeTransportationOffice(db *pop.Connection) (models.TransportationOffice, error) {
	address, err := MakeAddress(db)
	if err != nil {
		log.Panic(err)
	}
	office := models.TransportationOffice{
		Name:      "JPPSO Testy McTest",
		AddressID: address.ID,
		Address:   address,
		Latitude:  1.23445,
		Longitude: -23.34455,
	}

	verrs, err := db.ValidateAndSave(&office)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() != 0 {
		log.Panic(verrs.Error())
	}

	var phoneLines []models.OfficePhoneLine
	phoneLine := models.OfficePhoneLine{
		TransportationOfficeID: office.ID,
		TransportationOffice:   office,
		Number:                 "(510) 555-5555",
		IsDsnNumber:            false,
		Type:                   "voice",
	}
	phoneLines = append(phoneLines, phoneLine)
	phoneVerrs, phoneErr := db.ValidateAndSave(&phoneLine)
	if phoneErr != nil {
		log.Panic(phoneErr)
	}
	if phoneVerrs.Count() != 0 {
		log.Panic(phoneVerrs.Error())
	}

	office.PhoneLines = phoneLines
	Office1Verrs, Office1Err := db.ValidateAndSave(&office)
	if Office1Err != nil {
		log.Panic(Office1Err)
	}
	if Office1Verrs.Count() != 0 {
		log.Panic(Office1Verrs.Error())
	}
	return office, nil
}
