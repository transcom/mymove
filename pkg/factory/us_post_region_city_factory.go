package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// Creates a UsPostRegionCity for MacDill AFB, Florida
func BuildUsPostRegionCity(db *pop.Connection, customs []Customization, traits []Trait) models.UsPostRegionCity {
	customs = setupCustomizations(customs, traits)

	var cUsPostRegionCity models.UsPostRegionCity
	if result := findValidCustomization(customs, UsPostRegionCity); result != nil {
		cUsPostRegionCity = result.Model.(models.UsPostRegionCity)
		if result.LinkOnly {
			return cUsPostRegionCity
		}
	}

	usPostRegionCity := models.UsPostRegionCity{
		UsprZipID:               "33608",
		USPostRegionCityNm:      "MacDill AFB",
		UsprcPrfdLstLineCtystNm: "MacDill",
		UsprcCountyNm:           "Hillsborough",
		CtryGencDgphCd:          "US",
		State:                   "FL",
	}

	testdatagen.MergeModels(&usPostRegionCity, cUsPostRegionCity)

	if db != nil {
		mustCreate(db, &usPostRegionCity)
	}

	return usPostRegionCity
}

// Creates a default UsPostRegionCity for MacDill AFB, Florida
func BuildDefaultUsPostRegionCity(db *pop.Connection) models.UsPostRegionCity {
	return BuildUsPostRegionCity(db, nil, nil)
}
