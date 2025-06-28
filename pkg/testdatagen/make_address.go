package testdatagen

import (
	"log"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakeAddress creates a single Address and associated service member.
func MakeAddress(db *pop.Connection, assertions Assertions) (models.Address, error) {
	address := models.Address{
		StreetAddress1: "123 Any Street",
		StreetAddress2: models.StringPointer("P.O. Box 12345"),
		StreetAddress3: models.StringPointer("c/o Some Person"),
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "90210",
		County:         models.StringPointer("LOS ANGELES"),
		IsOconus:       models.BoolPointer(false),
	}

	mergeModels(&address, assertions.Address)

	usPostRegionCity, err := models.FindByZipCodeAndCity(db, address.PostalCode, address.City)
	if err != nil {
		log.Panic(err)
	}
	address.UsPostRegionCity = usPostRegionCity
	address.UsPostRegionCityID = &usPostRegionCity.ID

	mustCreate(db, &address, assertions.Stub)

	return address, nil
}

// MakeAddress2 creates a different single Address and associated service member.
func MakeAddress2(db *pop.Connection, assertions Assertions) (models.Address, error) {
	address := models.Address{
		StreetAddress1: "987 Any Avenue",
		StreetAddress2: models.StringPointer("P.O. Box 9876"),
		StreetAddress3: models.StringPointer("c/o Some Person"),
		City:           "Fairfield",
		State:          "CA",
		PostalCode:     "94535",
		IsOconus:       models.BoolPointer(false),
		County:         models.StringPointer("SOLANO"),
	}

	mergeModels(&address, assertions.Address)

	usPostRegionCity, err := models.FindByZipCodeAndCity(db, address.PostalCode, address.City)
	if err != nil {
		log.Panic(err)
	}
	address.UsPostRegionCity = usPostRegionCity
	address.UsPostRegionCityID = &usPostRegionCity.ID

	mustCreate(db, &address, assertions.Stub)

	return address, nil
}

// MakeAddress3 creates a different single Address and associated service member.
func MakeAddress3(db *pop.Connection, assertions Assertions) (models.Address, error) {
	address := models.Address{
		StreetAddress1: "987 Other Avenue",
		StreetAddress2: models.StringPointer("P.O. Box 1234"),
		StreetAddress3: models.StringPointer("c/o Another Person"),
		City:           "Des Moines",
		State:          "IA",
		PostalCode:     "50309",
		County:         models.StringPointer("POLK"),
		IsOconus:       models.BoolPointer(false),
	}

	mergeModels(&address, assertions.Address)

	usPostRegionCity, err := models.FindByZipCodeAndCity(db, address.PostalCode, address.City)
	if err != nil {
		log.Panic(err)
	}
	address.UsPostRegionCity = usPostRegionCity
	address.UsPostRegionCityID = &usPostRegionCity.ID

	mustCreate(db, &address, assertions.Stub)

	return address, nil
}

// MakeDefaultAddress makes an Address with default values
func MakeDefaultAddress(db *pop.Connection) (models.Address, error) {
	return MakeAddress(db, Assertions{})
}
