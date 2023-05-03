package factory

import (
	"log"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildPostalCodeToGBLOC creates a single PostalCodeToGBLOC entry.
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildPostalCodeToGBLOC(db *pop.Connection, customs []Customization, traits []Trait) models.PostalCodeToGBLOC {
	customs = setupCustomizations(customs, traits)

	// Find PostalCodeToGBLOC assertion and convert to models.PostalCodeToGBLOC
	var cPostalCodeToGBLOC models.PostalCodeToGBLOC
	if result := findValidCustomization(customs, PostalCodeToGBLOC); result != nil {
		cPostalCodeToGBLOC = result.Model.(models.PostalCodeToGBLOC)
		if result.LinkOnly {
			return cPostalCodeToGBLOC
		}
	}

	// Create PostalCodeToGBLOC
	postalCodeToGBLOC := models.PostalCodeToGBLOC{
		PostalCode: "90210",
		GBLOC:      "KKFA",
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&postalCodeToGBLOC, cPostalCodeToGBLOC)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &postalCodeToGBLOC)
	}
	return postalCodeToGBLOC
}

func FetchOrBuildPostalCodeToGBLOC(db *pop.Connection, postalCode string, gbloc string) models.PostalCodeToGBLOC {
	if postalCode == "" {
		log.Panic("Cannot FetchOrBuildPostalCodeToGBLOC with empty postalCode")
	}
	gblocForPostalCode, err := models.FetchGBLOCForPostalCode(db, postalCode)
	if err != nil && err != models.ErrFetchNotFound {
		log.Panicf("Cannot fetch gbloc for postal code %s: %s", postalCode, err)
	}
	if gblocForPostalCode.GBLOC == "" || err != nil {
		gblocForPostalCode = BuildPostalCodeToGBLOC(db, []Customization{
			{
				Model: models.PostalCodeToGBLOC{
					PostalCode: postalCode,
					GBLOC:      gbloc,
				},
			},
		}, nil)
	}
	return gblocForPostalCode
}
