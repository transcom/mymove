package factory

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// Creates a UsPostRegionCity for Beverly Hills, California 90210
func BuildUsPostRegionCity(db *pop.Connection, customs []Customization, traits []Trait) models.UsPostRegionCity {
	customs = setupCustomizations(customs, traits)

	var cUsPostRegionCity models.UsPostRegionCity
	if result := findValidCustomization(customs, UsPostRegionCity); result != nil {
		cUsPostRegionCity = result.Model.(models.UsPostRegionCity)
		if result.LinkOnly {
			return cUsPostRegionCity
		}
	}

	city := BuildCity(db, customs, nil)
	usPostRegion := BuildUsPostRegion(db, customs, nil)

	usPostRegionCity := models.UsPostRegionCity{
		ID:                 uuid.Must(uuid.NewV4()),
		UsprZipID:          "90210",
		USPostRegionCityNm: "Beverly Hills",
		UsprcCountyNm:      "LOS ANGELES",
		CtryGencDgphCd:     "US",
		City:               city,
		UsPostRegion:       usPostRegion,
		UsPostRegionId:     usPostRegion.ID,
	}

	testdatagen.MergeModels(&usPostRegionCity, cUsPostRegionCity)

	if db != nil {
		mustCreate(db, &usPostRegionCity)
	}

	return usPostRegionCity
}

// Creates a default UsPostRegionCity for Beverly Hills, California 90210
func BuildDefaultUsPostRegionCity(db *pop.Connection) models.UsPostRegionCity {
	return BuildUsPostRegionCity(db, nil, nil)
}
