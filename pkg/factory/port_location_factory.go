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

	var portLocation models.PortLocation
	if db != nil {

		// Find the Port if customization is provided
		var port models.Port
		if result := findValidCustomization(customs, Port); result != nil {
			port = FetchPort(db, customs, nil)
		}

		// Find the port location based on the port code
		err := db.EagerPreload("Port", "City", "Country", "UsPostRegionCity.UsPostRegion.State").Where("is_active = TRUE").InnerJoin("ports p", "port_id = p.id").Where("p.port_code = $1", port.PortCode).First(&portLocation)
		if err == nil {
			return portLocation
		}

		// Didn't find a port location based on the custom port code, so grab the default port location
		err = db.EagerPreload("Port", "City", "Country", "UsPostRegionCity.UsPostRegion.State").Where("is_active = TRUE").InnerJoin("ports p", "port_id = p.id").Where("p.port_code = 'PDX'").First(&portLocation)
		if err == nil {
			return portLocation
		}
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&portLocation, cPortLocation)

	return portLocation
}
