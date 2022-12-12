package factory

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// DefaultZip3 is the default zip3 for testing
var DefaultZip3 = "902"

// BuildTariff400ngZip3 finds or makes a single Tariff400ngZip3 record
func BuildTariff400ngZip3(db *pop.Connection, customs []Customization, traits []Trait) models.Tariff400ngZip3 {
	customs = setupCustomizations(customs, traits)

	var cTariff models.Tariff400ngZip3
	if result := findValidCustomization(customs, Tariff400ngZip3); result != nil {
		cTariff = result.Model.(models.Tariff400ngZip3)
	}

	zip3 := models.Tariff400ngZip3{
		Zip3:          DefaultZip3,
		BasepointCity: "Beverly Hills",
		State:         "CA",
		ServiceArea:   "56",
		RateArea:      "US88",
		Region:        "2",
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&zip3, cTariff)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &zip3)
	}

	return zip3
}

// FetchOrBuildTariff400ngZip3 tries fetching an existing zip3 first, then falls back to creating one
func FetchOrBuildTariff400ngZip3(db *pop.Connection, customs []Customization, traits []Trait) models.Tariff400ngZip3 {
	var existingZip3s models.Tariff400ngZip3s
	zip3 := DefaultZip3

	err := db.Where("zip3 = ?", zip3).All(&existingZip3s)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	}

	if len(existingZip3s) == 0 {
		return BuildTariff400ngZip3(db, customs, traits)
	}

	return existingZip3s[0]
}
