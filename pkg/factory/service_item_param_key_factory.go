package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildServiceItemParamKey creates a single ServiceItemParamKey
//
// Params:
//   - customs is a slice that will be modified by the factory
//   - db can be set to nil to create a stubbed model that is not stored in DB.
//
//  serviceItemParamKey := BuildServiceItemParamKey(suite.DB(), []Customization{
//  	{Model: customServiceItemParamKey},
//  }, nil)
func BuildServiceItemParamKey(db *pop.Connection, customs []Customization, traits []Trait) models.ServiceItemParamKey {
	customs = setupCustomizations(customs, traits)

	// Find serviceItemParamKey customization and extract the custom serviceItemParamKey
	var cServiceItemParamKey models.ServiceItemParamKey
	if result := findValidCustomization(customs, ServiceItemParamKey); result != nil {
		cServiceItemParamKey = result.Model.(models.ServiceItemParamKey)
		if result.LinkOnly {
			return cServiceItemParamKey
		}
	}

	serviceItemParamKey := models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameWeightEstimated,
		Description: "test name weight estimated description",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginPrime,
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&serviceItemParamKey, cServiceItemParamKey)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &serviceItemParamKey)
	}

	return serviceItemParamKey

}
