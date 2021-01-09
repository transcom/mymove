package testdatagen

import (
	"log"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// makeServiceItem creates a single service item and associated set relationships
func makeServiceItem(db *pop.Connection, assertions Assertions, isBasicServiceItem bool) models.MTOServiceItem {
	moveTaskOrder := assertions.Move
	if isZeroUUID(moveTaskOrder.ID) {
		moveTaskOrder = MakeMove(db, assertions)
	}

	var MTOShipmentID *uuid.UUID
	var MTOShipment models.MTOShipment
	if !isBasicServiceItem {
		if isZeroUUID(assertions.MTOShipment.ID) {
			MTOShipment = MakeMTOShipment(db, assertions)
			MTOShipmentID = &MTOShipment.ID
		} else {
			MTOShipment = assertions.MTOShipment
			MTOShipmentID = &assertions.MTOShipment.ID
		}
	}

	reService := assertions.ReService
	if isZeroUUID(reService.ID) {
		reService = FetchOrMakeReService(db, assertions)
	}

	status := assertions.MTOServiceItem.Status
	if status == "" {
		status = models.MTOServiceItemStatusSubmitted
	}

	MTOServiceItem := models.MTOServiceItem{
		MoveTaskOrder:   moveTaskOrder,
		MoveTaskOrderID: moveTaskOrder.ID,
		MTOShipment:     MTOShipment,
		MTOShipmentID:   MTOShipmentID,
		ReService:       reService,
		ReServiceID:     reService.ID,
		Status:          status,
	}

	// Overwrite values with those from assertions
	mergeModels(&MTOServiceItem, assertions.MTOServiceItem)

	mustCreate(db, &MTOServiceItem, assertions.Stub)

	return MTOServiceItem
}

// MakeMTOServiceItem creates a single MTOServiceItem and associated set relationships
func MakeMTOServiceItem(db *pop.Connection, assertions Assertions) models.MTOServiceItem {
	return makeServiceItem(db, assertions, false)
}

// MakeDefaultMTOServiceItem returns a MTOServiceItem with default values
func MakeDefaultMTOServiceItem(db *pop.Connection) models.MTOServiceItem {
	return MakeMTOServiceItem(db, Assertions{})
}

// MakeMTOServiceItemBasic creates a single MTOServiceItem that is a basic type, meaning no shipment id associated.
func MakeMTOServiceItemBasic(db *pop.Connection, assertions Assertions) models.MTOServiceItem {
	return makeServiceItem(db, assertions, true)
}

// MakeMTOServiceItems makes an array of MTOServiceItems
func MakeMTOServiceItems(db *pop.Connection) models.MTOServiceItems {
	var serviceItemList models.MTOServiceItems
	serviceItemList = append(serviceItemList, MakeDefaultMTOServiceItem(db))
	return serviceItemList
}

type realMTOServiceParamData struct {
	ServiceCode          models.ReServiceCode
	ServiceItemParamKeys []models.ServiceItemParamKey
}

func fixtureMapOfServiceItemParams() map[models.ReServiceCode]realMTOServiceParamData {
	serviceParams := make(map[models.ReServiceCode]realMTOServiceParamData)
	// CS
	serviceParams[models.ReServiceCodeCS] = realMTOServiceParamData{
		ServiceCode: models.ReServiceCodeCS,
		ServiceItemParamKeys: []models.ServiceItemParamKey{
			{
				Key:         models.ServiceItemParamNameMTOAvailableToPrimeAt,
				Description: "mto available to prime at",
				Type:        models.ServiceItemParamTypeTimestamp,
				Origin:      models.ServiceItemParamOriginSystem,
			},
			{
				Key:         models.ServiceItemParamNameContractCode,
				Description: "contract code",
				Type:        models.ServiceItemParamTypeString,
				Origin:      models.ServiceItemParamOriginSystem,
			},
		},
	}
	// MS
	serviceParams[models.ReServiceCodeMS] = realMTOServiceParamData{
		ServiceCode: models.ReServiceCodeMS,
		ServiceItemParamKeys: []models.ServiceItemParamKey{
			{
				Key:         models.ServiceItemParamNameMTOAvailableToPrimeAt,
				Description: "mto available to prime at",
				Type:        models.ServiceItemParamTypeTimestamp,
				Origin:      models.ServiceItemParamOriginSystem,
			},
			{
				Key:         models.ServiceItemParamNameContractCode,
				Description: "contract code",
				Type:        models.ServiceItemParamTypeString,
				Origin:      models.ServiceItemParamOriginSystem,
			},
		},
	}
	// DLH
	serviceParams[models.ReServiceCodeDLH] = realMTOServiceParamData{
		ServiceCode: models.ReServiceCodeDLH,
		ServiceItemParamKeys: []models.ServiceItemParamKey{
			{
				Key:         models.ServiceItemParamNameWeightEstimated,
				Description: "estimated weight",
				Type:        models.ServiceItemParamTypeInteger,
				Origin:      models.ServiceItemParamOriginPrime,
			},
			{
				Key:         models.ServiceItemParamNameRequestedPickupDate,
				Description: "requested pickup date",
				Type:        models.ServiceItemParamTypeDate,
				Origin:      models.ServiceItemParamOriginPrime,
			},
			{
				Key:         models.ServiceItemParamNameContractCode,
				Description: "contract code",
				Type:        models.ServiceItemParamTypeString,
				Origin:      models.ServiceItemParamOriginSystem,
			},
			{
				Key:         models.ServiceItemParamNameDistanceZip3,
				Description: "distance zip3",
				Type:        models.ServiceItemParamTypeInteger,
				Origin:      models.ServiceItemParamOriginSystem,
			},
			{
				Key:         models.ServiceItemParamNameZipPickupAddress,
				Description: "zip pickup address",
				Type:        models.ServiceItemParamTypeString,
				Origin:      models.ServiceItemParamOriginPrime,
			},
			{
				Key:         models.ServiceItemParamNameZipDestAddress,
				Description: "zip destination address",
				Type:        models.ServiceItemParamTypeString,
				Origin:      models.ServiceItemParamOriginPrime,
			},
			{
				Key:         models.ServiceItemParamNameWeightBilledActual,
				Description: "weight billed actual",
				Type:        models.ServiceItemParamTypeInteger,
				Origin:      models.ServiceItemParamOriginSystem,
			},
			{
				Key:         models.ServiceItemParamNameWeightActual,
				Description: "weight actual",
				Type:        models.ServiceItemParamTypeInteger,
				Origin:      models.ServiceItemParamOriginPrime,
			},
			{
				Key:         models.ServiceItemParamNameServiceAreaOrigin,
				Description: "service area actual",
				Type:        models.ServiceItemParamTypeString,
				Origin:      models.ServiceItemParamOriginPrime,
			},
		},
	}

	return serviceParams
}

// MakeRealMTOServiceItemWithAllDeps Takes a service code, move, shipment
// and creates or finds all the needed data to create a service item all its params ready for pricing
func MakeRealMTOServiceItemWithAllDeps(db *pop.Connection, serviceCode models.ReServiceCode, mto models.Move, mtoShipment models.MTOShipment) models.MTOServiceItem {
	serviceParams := fixtureMapOfServiceItemParams()

	// look up the data we need
	data := serviceParams[serviceCode]
	if data.ServiceCode == serviceCode {
		// get or create the ReService
		reService := FetchOrMakeReService(db, Assertions{
			ReService: models.ReService{
				Code: serviceCode,
			},
		})

		// create all params defined for this particular service
		for _, serviceParamKeyToCreate := range data.ServiceItemParamKeys {
			serviceItemParamKey := FetchOrMakeServiceItemParamKey(db, Assertions{
				ServiceItemParamKey: serviceParamKeyToCreate,
			})
			_ = MakeServiceParam(db, Assertions{
				ServiceParam: models.ServiceParam{
					ServiceID:             reService.ID,
					ServiceItemParamKeyID: serviceItemParamKey.ID,
					ServiceItemParamKey:   serviceItemParamKey,
				},
			})
		}

		// create a service item and return it
		mtoServiceItem := MakeMTOServiceItem(db, Assertions{
			Move:        mto,
			MTOShipment: mtoShipment,
		})

		return mtoServiceItem
	}

	log.Panicf("couldn't create service item service code %s not defiled", serviceCode)
	return models.MTOServiceItem{}
}
