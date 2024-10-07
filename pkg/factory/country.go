package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildCountry creates a single Country.
// Also creates, if not provided
// - Country
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildUSCountry(db *pop.Connection, customs []Customization, traits []Trait) models.Country {
	customs = setupCustomizations(customs, traits)

	var cCountry models.Country
	if result := findValidCustomization(customs, Country); result != nil {
		cCountry = result.Model.(models.Country)
		if result.LinkOnly {
			return cCountry
		}
	}

	// Check if the "US" country already exists in the database
	if db != nil {
		var existingCountry models.Country
		err := db.Where("country = ?", "US").First(&existingCountry)
		if err == nil {
			return existingCountry
		}
	}

	country := models.Country{
		Country:     "US",
		CountryName: "UNITED STATES",
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&country, cCountry)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &country)
	}
	return country
}
