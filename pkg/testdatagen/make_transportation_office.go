package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeTransportationOffice creates a single ServiceMember and associated User.
func MakeTransportationOffice(db *pop.Connection, assertions Assertions) models.TransportationOffice {
	address := MakeDefaultAddress(db)

	office := models.TransportationOffice{
		Name:      "JPPSO Testy McTest",
		AddressID: address.ID,
		Address:   address,
		Gbloc:     "LKNQ",
		Latitude:  1.23445,
		Longitude: -23.34455,
	}

	mergeModels(&office, assertions.TransportationOffice)

	mustCreate(db, &office)

	var phoneLines []models.OfficePhoneLine
	phoneLine := models.OfficePhoneLine{
		TransportationOfficeID: office.ID,
		TransportationOffice:   office,
		Number:                 "(510) 555-5555",
		IsDsnNumber:            false,
		Type:                   "voice",
	}
	phoneLines = append(phoneLines, phoneLine)
	mustCreate(db, &phoneLine)

	office.PhoneLines = phoneLines

	return office
}

// MakeDefaultTransportationOffice makes a default TransportationOffice
func MakeDefaultTransportationOffice(db *pop.Connection) models.TransportationOffice {
	return MakeTransportationOffice(db, Assertions{})
}
