package factory

import (
	"database/sql"
	"log"

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
func BuildCountry(db *pop.Connection, customs []Customization, traits []Trait) models.Country {
	customs = setupCustomizations(customs, traits)

	var cCountry models.Country
	if result := findValidCustomization(customs, Country); result != nil {
		cCountry = result.Model.(models.Country)
		if result.LinkOnly {
			return cCountry
		}
	}

	// Check if the country provided already exists in the database
	if db != nil {
		var existingCountry models.Country
		err := db.Where("country = ?", cCountry.Country).First(&existingCountry)
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

// FetchOrBuildCountry tries fetching a Country using a provided customization, then falls back to creating a default "US" country
func FetchOrBuildCountry(db *pop.Connection, customs []Customization, traits []Trait) models.Country {
	if db == nil {
		return BuildCountry(db, customs, traits)
	}

	customs = setupCustomizations(customs, traits)

	var cCountry models.Country
	if result := findValidCustomization(customs, Country); result != nil {
		cCountry = result.Model.(models.Country)
		if result.LinkOnly {
			return cCountry
		}
	}

	if !cCountry.ID.IsNil() {
		err := db.Where("ID = $1", cCountry.ID).First(&cCountry)
		if err != nil && err != sql.ErrNoRows {
			log.Panic(err)
		} else if err == nil {
			return cCountry
		}
	}

	if cCountry.Country != "" {
		err := db.Where("Country = $1", cCountry.Country).First(&cCountry)
		if err != nil && err != sql.ErrNoRows {
			log.Panic(err)
		} else if err == nil {
			return cCountry
		}
	}

	// search for the default code if one is not provided
	defaultCountryCode := "US"
	err := db.Where("country = $1", defaultCountryCode).First(&cCountry)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	} else if err == nil {
		return cCountry
	}

	return BuildCountry(db, customs, traits)
}
