package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeServiceItemParamKey creates a single ServiceItemParamKey
func MakeServiceItemParamKey(db *pop.Connection, assertions Assertions) models.ServiceItemParamKey {
	serviceItemParamKey := models.ServiceItemParamKey{
		Key:         "testkey",
		Description: "test key",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginPrime,
	}

	// Overwrite values with those from assertions
	mergeModels(&serviceItemParamKey, assertions.ServiceItemParamKey)

	mustCreate(db, &serviceItemParamKey)

	return serviceItemParamKey
}

// MakeDefaultServiceItemParamKey makes a ServiceItemParamKey with default values
func MakeDefaultServiceItemParamKey(db *pop.Connection) models.ServiceItemParamKey {
	return MakeServiceItemParamKey(db, Assertions{})
}
