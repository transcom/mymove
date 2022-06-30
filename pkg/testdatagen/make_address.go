package testdatagen

import (
	"fmt"
	"runtime"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeAddress creates a single Address and associated service member.
func MakeAddress(db *pop.Connection, assertions Assertions) models.Address {
	_, file, no, ok := runtime.Caller(1)
	if ok {
		fmt.Printf("called from %s#%d\n", file, no)
	}

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

	mustCreate(db, &address, assertions.Stub)

	return address
}

// MakeAddress2 creates a different single Address and associated service member.
func MakeAddress2(db *pop.Connection, assertions Assertions) models.Address {
	_, file, no, ok := runtime.Caller(1)
	if ok {
		fmt.Printf("called from %s#%d\n", file, no)
	}
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

	mustCreate(db, &address, assertions.Stub)

	return address
}

// MakeAddress3 creates a different single Address and associated service member.
func MakeAddress3(db *pop.Connection, assertions Assertions) models.Address {
	_, file, no, ok := runtime.Caller(1)
	if ok {
		fmt.Printf("called from %s#%d\n", file, no)
	}
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

	mustCreate(db, &address, assertions.Stub)

	return address
}

// MakeAddress4 creates a different single Address and associated service member.
func MakeAddress4(db *pop.Connection, assertions Assertions) models.Address {
	_, file, no, ok := runtime.Caller(1)
	if ok {
		fmt.Printf("called from %s#%d\n", file, no)
	}
	address := models.Address{
		StreetAddress1: "987 Over There Avenue",
		StreetAddress2: swag.String("P.O. Box 1234"),
		StreetAddress3: swag.String("c/o Another Person"),
		City:           "Houston",
		State:          "TX",
		PostalCode:     "77083",
		Country:        swag.String("US"),
	}

	mergeModels(&address, assertions.Address)

	mustCreate(db, &address, assertions.Stub)

	return address
}

// MakeDefaultAddress makes an Address with default values
func MakeDefaultAddress(db *pop.Connection) models.Address {
	// Make associated lookup table records.
	FetchOrMakeDefaultTariff400ngZip3(db)
	_, file, no, ok := runtime.Caller(1)
	if ok {
		fmt.Printf("called from %s#%d\n", file, no)
	}
	return MakeAddress(db, Assertions{})
}

// MakeStubbedAddress returns a stubbed address without saving it to the DB
func MakeStubbedAddress(db *pop.Connection) models.Address {
	return MakeAddress(db, Assertions{
		Address: models.Address{
			ID: uuid.Must(uuid.NewV4()),
		},
		Stub: true,
	})
}
