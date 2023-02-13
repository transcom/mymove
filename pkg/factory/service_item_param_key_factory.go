package factory

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

var defaultSericeItemParamKey = models.ServiceItemParamNameWeightEstimated

// BuildServiceItemParamKey creates a single ServiceItemParamKey
//
// Params:
//
//   - customs is a slice that will be modified by the factory
//
//   - db can be set to nil to create a stubbed model that is not stored in DB.
//
//     serviceItemParamKey := BuildServiceItemParamKey(suite.DB(), []Customization{
//     {Model: customServiceItemParamKey},
//     }, nil)
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
		Key:         defaultSericeItemParamKey,
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

// FetchOrBuildServiceItemParamKey returns the ServiceItemParamKey for a given param key, or creates one if
// the param key does not exist yet.
func FetchOrBuildServiceItemParamKey(db *pop.Connection, customs []Customization, traits []Trait) models.ServiceItemParamKey {
	if db == nil {
		return BuildServiceItemParamKey(db, customs, traits)
	}
	customs = setupCustomizations(customs, traits)

	var existingServiceItemParamKeys models.ServiceItemParamKeys
	key := defaultSericeItemParamKey

	var cServiceItemParamKey models.ServiceItemParamKey
	if result := findValidCustomization(customs, ServiceItemParamKey); result != nil {
		cServiceItemParamKey = result.Model.(models.ServiceItemParamKey)
		if cServiceItemParamKey.Key != "" {
			key = cServiceItemParamKey.Key
		}
	}

	err := db.Where("key = ?", key).All(&existingServiceItemParamKeys)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	}

	if len(existingServiceItemParamKeys) == 0 {
		return BuildServiceItemParamKey(db, customs, traits)
	}

	return existingServiceItemParamKeys[0]

}
