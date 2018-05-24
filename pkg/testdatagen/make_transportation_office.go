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
		Longitude: 23.34455,
	}

	verrs, err := db.ValidateAndSave(&office)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() != 0 {
		log.Panic(verrs.Error())
	}

	return office, nil
}
