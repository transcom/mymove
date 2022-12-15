package factory

import (
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildAddress creates a single Address and associated service member.
func BuildAddress(db *pop.Connection, customs []Customization, traits []Trait) models.Address {
	customs = setupCustomizations(customs, traits)

	// Find address assertion and convert to models address
	var cAddress models.Address
	if result := findValidCustomization(customs, Address); result != nil {
		cAddress = result.Model.(models.Address)
	}

	// Create default Address
	address := models.Address{
		StreetAddress1: "123 Any Street",
		StreetAddress2: swag.String("P.O. Box 12345"),
		StreetAddress3: swag.String("c/o Some Person"),
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "90210",
		Country:        swag.String("US"),
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&address, cAddress)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &address)
	}

	return address
}

// BuildDefaultAddress makes an Address with default values
func BuildDefaultAddress(db *pop.Connection) models.Address {
	FetchOrBuildTariff400ngZip3(db, nil, nil)
	return BuildAddress(db, nil, nil)
}

// BuildStubbedAddress returns a stubbed address without saving it to the DB
func BuildStubbedAddress(customs []Customization, traits []Trait) models.Address {
	return BuildAddress(nil, customs, traits)
}
