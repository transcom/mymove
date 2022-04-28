package testdatagen

import (
	"database/sql"
	"log"

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

// FetchOrMakeServiceItemParamKey returns the ServiceItemParamKey for a given param key, or creates one if
// the param key does not exist yet.
func FetchOrMakeServiceItemParamKey(db *pop.Connection, assertions Assertions) models.ServiceItemParamKey {
	var existingServiceItemParamKeys models.ServiceItemParamKeys
	key := DefaultServiceItemParamKeyName
	if assertions.ServiceItemParamKey.Key != "" {
		key = assertions.ServiceItemParamKey.Key.String()
	}
	err := db.Where("key = ?", key).All(&existingServiceItemParamKeys)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	}

	if len(existingServiceItemParamKeys) == 0 {
		return MakeServiceItemParamKey(db, assertions)
	}

	return existingServiceItemParamKeys[0]
}

// MakeDefaultServiceItemParamKey makes a ServiceItemParamKey with default values
func MakeDefaultServiceItemParamKey(db *pop.Connection) models.ServiceItemParamKey {
	return MakeServiceItemParamKey(db, Assertions{})
}
