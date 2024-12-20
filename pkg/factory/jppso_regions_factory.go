package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildJppsoRegions creates a single JppsoRegions entry.
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildJppsoRegions(db *pop.Connection, customs []Customization, traits []Trait) models.JppsoRegions {
	customs = setupCustomizations(customs, traits)

	// Find JppsoRegions assertion and convert to models.JppsoRegions
	var cJppsoRegions models.JppsoRegions
	if result := findValidCustomization(customs, JppsoRegions); result != nil {
		cJppsoRegions = result.Model.(models.JppsoRegions)
		if result.LinkOnly {
			return cJppsoRegions
		}
	}

	// Create JppsoRegions
	jppsoRegions := models.JppsoRegions{
		Code: "KKFA",
		Name: "JPPSO-North Central",
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&jppsoRegions, cJppsoRegions)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &jppsoRegions)
	}
	return jppsoRegions
}
