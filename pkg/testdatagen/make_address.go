package testdatagen

import (
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeAddress creates a single Address and associated service member.
func MakeAddress(db *pop.Connection, assertions Assertions) models.Address {
	address := models.Address{
		StreetAddress1: "123 Any Street",
		StreetAddress2: swag.String("P.O. Box 12345"),
		StreetAddress3: swag.String("c/o Some Person"),
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "90210",
		Country:        swag.String("US"),
	}

	mergeModels(&address, assertions.Address)

	mustCreate(db, &address)

	return address
}

// MakeDefaultAddress makes an Address with default values
func MakeDefaultAddress(db *pop.Connection) models.Address {
	// Make associated lookup table records.
	MakeDefaultTariff400ngZip3(db)

	return MakeAddress(db, Assertions{})
}
