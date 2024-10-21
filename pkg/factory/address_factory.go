package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildAddress creates a single Address.
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildAddress(db *pop.Connection, customs []Customization, traits []Trait) models.Address {
	customs = setupCustomizations(customs, traits)

	// Find address customization and extract the custom address
	var cAddress models.Address
	if result := findValidCustomization(customs, Address); result != nil {
		cAddress = result.Model.(models.Address)
		if result.LinkOnly {
			return cAddress
		}
	}

	// Create default Address
	address := models.Address{
		StreetAddress1: "123 Any Street",
		StreetAddress2: models.StringPointer("P.O. Box 12345"),
		StreetAddress3: models.StringPointer("c/o Some Person"),
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "90210",
		County:         "LOS ANGELES",
		IsOconus:       models.BoolPointer(false),
	}

	// Find/create the Country if customization is provided
	var country models.Country
	if result := findValidCustomization(customs, Country); result != nil {
		country = BuildCountry(db, customs, nil)
	} else {
		country = FetchOrBuildCountry(db, []Customization{
			{
				Model: models.Country{
					Country:     "US",
					CountryName: "UNITED STATES",
				},
			},
		}, nil)
	}

	address.Country = &country
	address.CountryId = &country.ID

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&address, cAddress)

	// This helps assign counties when the factory is called for seed data or tests
	// Additionally, also only run if not 90210. 90210's county is by default populated
	if db != nil && address.PostalCode != "90210" {
		county, err := models.FindCountyByZipCode(db, address.PostalCode)
		if err != nil {
			// A zip code that is not being tracked has been entered
			address.County = "does not exist"
		} else {
			// The zip code successfully found a county
			address.County = county
		}
	} else if db == nil && address.PostalCode != "90210" {
		// If no db supplied, mark that
		address.County = "db nil when created"
	}

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &address)
	}

	return address
}

func BuildMinimalAddress(db *pop.Connection, customs []Customization, traits []Trait) models.Address {
	customs = setupCustomizations(customs, traits)

	// Find address customization and extract the custom address
	var cAddress models.Address
	if result := findValidCustomization(customs, Address); result != nil {
		cAddress = result.Model.(models.Address)
		if result.LinkOnly {
			return cAddress
		}
	}

	// Create default Address
	address := models.Address{
		StreetAddress1: "N/A",
		City:           "Fort Gorden",
		State:          "GA",
		PostalCode:     "30813",
		County:         "RICHMOND",
	}

	// Find/create the Country if customization is provided
	var country models.Country
	if result := findValidCustomization(customs, Country); result != nil {
		country = BuildCountry(db, customs, nil)
	} else {
		country = BuildCountry(db, []Customization{
			{
				Model: models.Country{
					Country:     "US",
					CountryName: "UNITED STATES",
				},
			},
		}, nil)
	}

	address.Country = &country
	address.CountryId = &country.ID

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
	return BuildAddress(db, nil, nil)
}

// GetTraitAddress2 is a sample GetTraitFunc
func GetTraitAddress2() []Customization {
	return []Customization{
		{
			Model: models.Address{
				StreetAddress1: "987 Any Avenue",
				StreetAddress2: models.StringPointer("P.O. Box 9876"),
				StreetAddress3: models.StringPointer("c/o Some Person"),
				City:           "Fairfield",
				State:          "CA",
				PostalCode:     "94535",
			},
		},
	}
}

// GetTraitAddress3 is a sample GetTraitFunc
func GetTraitAddress3() []Customization {

	return []Customization{
		{
			Model: models.Address{
				StreetAddress1: "987 Other Avenue",
				StreetAddress2: models.StringPointer("P.O. Box 1234"),
				StreetAddress3: models.StringPointer("c/o Another Person"),
				City:           "Des Moines",
				State:          "IA",
				PostalCode:     "50309",
			},
		},
	}
}

// GetTraitAddress4 is a sample GetTraitFunc
func GetTraitAddress4() []Customization {

	return []Customization{
		{
			Model: models.Address{
				StreetAddress1: "987 Over There Avenue",
				StreetAddress2: models.StringPointer("P.O. Box 1234"),
				StreetAddress3: models.StringPointer("c/o Another Person"),
				City:           "Houston",
				State:          "TX",
				PostalCode:     "77083",
			},
		},
	}
}
