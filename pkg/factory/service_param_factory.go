package factory

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildServiceParam creates a single ServiceParam
//
// Params:
//
//   - customs is a slice that will be modified by the factory
//
//   - db can be set to nil to create a stubbed model that is not stored in DB.
//
//     serviceParam := BuildServiceParam(suite.DB(), []Customization{
//     {Model: customServiceParam},
//     }, nil)
func BuildServiceParam(db *pop.Connection, customs []Customization, traits []Trait) models.ServiceParam {
	customs = setupCustomizations(customs, traits)

	// Find serviceParam customization and extract the custom serviceParam
	var cServiceParam models.ServiceParam
	if result := findValidCustomization(customs, ServiceParam); result != nil {
		cServiceParam = result.Model.(models.ServiceParam)
		if result.LinkOnly {
			return cServiceParam
		}
	}

	reService := FetchReService(db, customs, traits)

	serviceItemParamKey := FetchOrBuildServiceItemParamKey(db, customs, traits)

	serviceParam := models.ServiceParam{
		ServiceID:             reService.ID,
		Service:               reService,
		ServiceItemParamKeyID: serviceItemParamKey.ID,
		ServiceItemParamKey:   serviceItemParamKey,
		IsOptional:            false,
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&serviceParam, cServiceParam)

	// always override optional setting for certain keys
	switch serviceParam.ServiceItemParamKey.Key {
	case models.ServiceItemParamNameWeightEstimated,
		models.ServiceItemParamNameWeightReweigh,
		models.ServiceItemParamNameWeightAdjusted,
		models.ServiceItemParamNameRequestedPickupDate,
		models.ServiceItemParamNameActualPickupDate:
		serviceParam.IsOptional = true
	}

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &serviceParam)
	}

	return serviceParam

}

// FetchOrBuildServiceParam returns the ServiceParam for a given param key, or creates one if
// the param key does not exist yet.
func FetchOrBuildServiceParam(db *pop.Connection, customs []Customization, traits []Trait) models.ServiceParam {
	if db == nil {
		return BuildServiceParam(db, customs, traits)
	}
	customs = setupCustomizations(customs, traits)

	var cServiceParam models.ServiceParam
	if result := findValidCustomization(customs, ServiceParam); result != nil {
		cServiceParam = result.Model.(models.ServiceParam)
		if result.LinkOnly {
			return cServiceParam
		}
	}

	reService := FetchReService(db, customs, traits)
	serviceItemParamKey := FetchOrBuildServiceItemParamKey(db, customs, traits)

	existingServiceParam := models.ServiceParam{}
	err := db.Where(
		"service_params.service_id = ? AND service_params.service_item_param_key_id = ?",
		reService.ID,
		serviceItemParamKey.ID).First(&existingServiceParam)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	} else if err == nil {
		return existingServiceParam
	}

	return BuildServiceParam(db, customs, traits)
}
