package testdatagen

import (
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// GetAddress constructs a single Address object
func GetAddress(assertions Assertions) *models.Address {
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
	return &address
}

// MakeAddress creates a single Address
func MakeAddress(db *pop.Connection, assertions Assertions) models.Address {
	addr := GetAddress(assertions)
	mustCreate(db, addr)
	return *addr
}

// MakeAddress2 creates a different single Address
func MakeAddress2(db *pop.Connection, assertions Assertions) models.Address {
	address := models.Address{
		StreetAddress1: "987 Any Avenue",
		StreetAddress2: swag.String("P.O. Box 9876"),
		StreetAddress3: swag.String("c/o Some Person"),
		City:           "Fairfield",
		State:          "CA",
		PostalCode:     "94535",
		Country:        swag.String("US"),
	}

	mergeModels(&address, assertions.Address)

	mustCreate(db, &address)

	return address
}

// MakeAddress3 creates a different single Address and associated service member.
func MakeAddress3(db *pop.Connection, assertions Assertions) models.Address {
	address := models.Address{
		StreetAddress1: "987 Other Avenue",
		StreetAddress2: swag.String("P.O. Box 1234"),
		StreetAddress3: swag.String("c/o Another Person"),
		City:           "Des Moines",
		State:          "IA",
		PostalCode:     "50309",
		Country:        swag.String("US"),
	}

	mergeModels(&address, assertions.Address)

	mustCreate(db, &address)

	return address
}

// MakeDefaultAddress makes an Address with default values
func MakeDefaultAddress(db *pop.Connection) models.Address {
	// Make associated lookup table records.
	FetchOrMakeDefaultTariff400ngZip3(db)

	return MakeAddress(db, Assertions{})
}

// GetDefaultAddress returns an address with default values
func GetDefaultAddress() *models.Address {
	return GetAddress(Assertions{})
}
