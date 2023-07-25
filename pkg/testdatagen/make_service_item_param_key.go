package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakeServiceItemParamKey creates a single ServiceItemParamKey
func MakeServiceItemParamKey(db *pop.Connection, assertions Assertions) models.ServiceItemParamKey {
	serviceItemParamKey := models.ServiceItemParamKey{
		Key:         DefaultServiceItemParamKeyName,
		Description: "test key",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginPrime,
	}

	// Overwrite values with those from assertions
	mergeModels(&serviceItemParamKey, assertions.ServiceItemParamKey)

	mustCreate(db, &serviceItemParamKey, assertions.Stub)

	return serviceItemParamKey
}
