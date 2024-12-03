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

	if db != nil {
		var existingPort models.Port
		err := db.Where("port_code = ?", cPort.PortCode).First(&existingPort)
		if err == nil {
			return existingPort
		}
	}

	port := models.Port{
		PortCode: "PDX",
		PortName: "PORTLAND INTL",
		PortType: "A",
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&port, cPort)
	return port
}
