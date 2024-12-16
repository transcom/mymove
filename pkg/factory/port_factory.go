package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func FetchPort(db *pop.Connection, customs []Customization, traits []Trait) models.Port {
	customs = setupCustomizations(customs, traits)

	var cPort models.Port
	if result := findValidCustomization(customs, Port); result != nil {
		cPort = result.Model.(models.Port)
		if result.LinkOnly {
			return cPort
		}
	}

	var port models.Port
	if db != nil {

		var err error
		if cPort.PortCode != "" {
			// Find the port based off the port code
			err = db.Where("port_code = ?", cPort.PortCode).First(&port)
		} else {
			//No port code provided, so find the default port
			err = db.Where("port_code = 'PDX'").First(&port)
		}
		if err == nil {
			return port
		}

	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&port, cPort)
	return port
}
