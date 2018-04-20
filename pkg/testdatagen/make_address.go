package testdatagen

import (
	"log"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeAddress creates a single Address and associated service member.
func MakeAddress(db *pop.Connection) (models.Address, error) {
	address := models.Address{
		StreetAddress1: "123 Any Street",
		StreetAddress2: swag.String("P.O. Box 12345"),
		StreetAddress3: swag.String("c/o Some Person"),
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "90210",
		Country:        swag.String("US"),
	}

	verrs, err := db.ValidateAndSave(&address)
	if err != nil {
		log.Fatal(err)
	}
	if verrs.Count() != 0 {
		log.Fatal(verrs.Error())
	}

	return address, err
}
