package testdatagen

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeServiceParam creates a single ServiceParam
func MakeServiceParam(db *pop.Connection, assertions Assertions) models.ServiceParam {

	serviceParam := models.ServiceParam{}

	setServiceParamIDs(db, &serviceParam, assertions)

	serviceParam.IsOptional = false
	switch assertions.ServiceParam.ServiceItemParamKey.Key {
	case models.ServiceItemParamNameWeightEstimated,
		models.ServiceItemParamNameWeightReweigh,
		models.ServiceItemParamNameWeightAdjusted,
		models.ServiceItemParamNameRequestedPickupDate,
		models.ServiceItemParamNameActualPickupDate:
		serviceParam.IsOptional = true
		assertions.ServiceParam.IsOptional = true
	}

	switch serviceParam.ServiceItemParamKey.Key {
	case models.ServiceItemParamNameWeightEstimated,
		models.ServiceItemParamNameWeightReweigh,
		models.ServiceItemParamNameWeightAdjusted,
		models.ServiceItemParamNameRequestedPickupDate,
		models.ServiceItemParamNameActualPickupDate:
		serviceParam.IsOptional = true
		assertions.ServiceParam.IsOptional = true
	}

	switch assertions.ServiceItemParamKey.Key {
	case models.ServiceItemParamNameWeightEstimated,
		models.ServiceItemParamNameWeightReweigh,
		models.ServiceItemParamNameWeightAdjusted,
		models.ServiceItemParamNameRequestedPickupDate,
		models.ServiceItemParamNameActualPickupDate:
		serviceParam.IsOptional = true
		assertions.ServiceParam.IsOptional = true
	}

	// Overwrite values with those from assertions
	mergeModels(&serviceParam, assertions.ServiceParam)

	mustCreate(db, &serviceParam, assertions.Stub)

	return serviceParam
}

func setServiceParamIDs(db *pop.Connection, serviceParam *models.ServiceParam, assertions Assertions) {
	// Make sure we have a ServiceID
	var reServiceItem models.ReService
	if assertions.ServiceParam.ServiceID == uuid.Nil && assertions.ServiceParam.Service.ID == uuid.Nil && assertions.ReService.ID == uuid.Nil {
		reServiceItem = MakeDefaultReService(db)
		serviceParam.ServiceID = reServiceItem.ID
	} else if assertions.ServiceParam.ServiceID != uuid.Nil {
		serviceParam.ServiceID = assertions.ServiceParam.ServiceID
	} else if assertions.ServiceParam.Service.ID != uuid.Nil {
		serviceParam.ServiceID = assertions.ServiceParam.Service.ID
	} else if assertions.ReService.ID != uuid.Nil {
		serviceParam.ServiceID = assertions.ReService.ID
	}

	// Make sure we have a ServiceItemParamKeyID
	var serviceItemParamKey models.ServiceItemParamKey
	if assertions.ServiceParam.ServiceItemParamKeyID == uuid.Nil && assertions.ServiceParam.ServiceItemParamKey.ID == uuid.Nil && assertions.ServiceItemParamKey.ID == uuid.Nil {
		serviceItemParamKey = MakeDefaultServiceItemParamKey(db)
		serviceParam.ServiceItemParamKeyID = serviceItemParamKey.ID
	} else if assertions.ServiceParam.ServiceItemParamKeyID != uuid.Nil {
		serviceParam.ServiceItemParamKeyID = assertions.ServiceParam.ServiceItemParamKeyID
	} else if assertions.ServiceParam.ServiceItemParamKey.ID != uuid.Nil {
		serviceParam.ServiceItemParamKeyID = assertions.ServiceParam.ServiceItemParamKey.ID
	} else if assertions.ServiceItemParamKey.ID != uuid.Nil {
		serviceParam.ServiceItemParamKeyID = assertions.ServiceItemParamKey.ID
	}
}
