package testdatagen

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeTransportationOffice creates a single TransportationOffice.
func MakeTransportationOffice(db *pop.Connection, assertions Assertions) models.TransportationOffice {

	transportationOfficeID := assertions.TransportationOffice.ID
	if isZeroUUID(transportationOfficeID) {
		transportationOfficeID = uuid.Must(uuid.NewV4())
	}

	address := MakeDefaultAddress(db)

	office := models.TransportationOffice{
		ID:        transportationOfficeID,
		Name:      "JPPSO Testy McTest",
		AddressID: address.ID,
		Address:   address,
		Gbloc:     "LKNQ",
		Latitude:  1.23445,
		Longitude: -23.34455,
	}

	mergeModels(&office, assertions.TransportationOffice)

	if assertions.Stub != true {
		mustCreate(db, &office)
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
	if assertions.Stub != true {
		mustCreate(db, &phoneLine)
	}

	office.PhoneLines = phoneLines

	return office
}

// MakeDefaultTransportationOffice makes a default TransportationOffice
func MakeDefaultTransportationOffice(db *pop.Connection) models.TransportationOffice {
	return MakeTransportationOffice(db, Assertions{})
}
