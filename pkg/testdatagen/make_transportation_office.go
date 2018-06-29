package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeTransportationOffice creates a single ServiceMember and associated User.
func MakeTransportationOffice(db *pop.Connection) models.TransportationOffice {
	address := MakeDefaultAddress(db)

	office := models.TransportationOffice{
		Name:      "JPPSO Testy McTest",
		AddressID: address.ID,
		Address:   address,
		Latitude:  1.23445,
		Longitude: -23.34455,
	}

	mustSave(db, &office)

	var phoneLines []models.OfficePhoneLine
	phoneLine := models.OfficePhoneLine{
		TransportationOfficeID: office.ID,
		TransportationOffice:   office,
		Number:                 "(510) 555-5555",
		IsDsnNumber:            false,
		Type:                   "voice",
	}
	phoneLines = append(phoneLines, phoneLine)
	mustSave(db, &phoneLine)

	office.PhoneLines = phoneLines
	mustSave(db, &office)

	return office
}
