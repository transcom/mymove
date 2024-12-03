package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func FetchPortLocation(db *pop.Connection, customs []Customization, traits []Trait) models.PortLocation {
	customs = setupCustomizations(customs, traits)

	var cPortLocation models.PortLocation
	if result := findValidCustomization(customs, PortLocation); result != nil {
		cPortLocation = result.Model.(models.PortLocation)
		if result.LinkOnly {
			return cPortLocation
		}
	}

	// Find the Port if customization is provided
	var port models.Port
	if result := findValidCustomization(customs, Port); result != nil {
		port = FetchPort(db, customs, nil)
	}

	if db != nil {
		var existingPortLocation models.PortLocation
		err := db.EagerPreload("Port", "City", "Country", "UsPostRegionCity.UsPostRegion.State").Where("is_active = TRUE").InnerJoin("ports p", "port_id = p.id").Where("p.port_code = $1", port.PortCode).First(&existingPortLocation)
		if err == nil {
			return existingPortLocation
		}
	}

	defaultPortLocation := models.PortLocation{
		Port: models.Port{
			PortType: models.PortTypeAir,
			PortCode: "PDX",
			PortName: "PORTLAND INTL",
		},
		City: models.City{
			CityName: "PORTLAND",
		},
		UsPostRegionCity: models.UsPostRegionCity{
			UsprcCountyNm: "MULTNOMAH",
			UsprZipID:     "97220",
			UsPostRegion: models.UsPostRegion{
				State: models.State{
					StateName: "OREGON",
				},
			},
		},
		Country: models.Country{
			CountryName: "UNITED STATES",
		},
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&defaultPortLocation, cPortLocation)

	return defaultPortLocation
}
