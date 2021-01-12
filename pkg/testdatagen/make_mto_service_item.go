package testdatagen

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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
	// FSC
	serviceParams[models.ReServiceCodeFSC] = realMTOServiceParamData{
		ServiceCode: models.ReServiceCodeFSC,
		ServiceItemParamKeys: []models.ServiceItemParamKey{
			{
				Key:         models.ServiceItemParamNameEIAFuelPrice,
				Description: "eia fuel price",
				Type:        models.ServiceItemParamTypeInteger,
				Origin:      models.ServiceItemParamOriginSystem,
			},
			{
				Key:         models.ServiceItemParamNameContractCode,
				Description: "contract code",
				Type:        models.ServiceItemParamTypeString,
				Origin:      models.ServiceItemParamOriginSystem,
			},
			{
				Key:         models.ServiceItemParamNameActualPickupDate,
				Description: "actual pickup date",
				Type:        models.ServiceItemParamTypeDate,
				Origin:      models.ServiceItemParamOriginPrime,
			},
			{
				Key:         models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
				Description: "fsc weight based multiplier",
				Type:        models.ServiceItemParamTypeDecimal,
				Origin:      models.ServiceItemParamOriginSystem,
			},
			{
				Key:         models.ServiceItemParamNameDistanceZip3,
				Description: "distance zip 3",
				Type:        models.ServiceItemParamTypeInteger,
				Origin:      models.ServiceItemParamOriginSystem,
			},
			{
				Key:         models.ServiceItemParamNameDistanceZip5,
				Description: "distance zip 5",
				Type:        models.ServiceItemParamTypeInteger,
				Origin:      models.ServiceItemParamOriginSystem,
			},
			{
				Key:         models.ServiceItemParamNameWeightBilledActual,
				Description: "weight billed actual",
				Type:        models.ServiceItemParamTypeInteger,
				Origin:      models.ServiceItemParamOriginSystem,
			},
		},
	}

	return serviceParams
}

func loadFixtureDataServiceItemParams(db *pop.Connection) {
	// Loads fixture data. This data is pulled from the development migrations using the following pg_dump command
	// pg_dump --host localhost --username postgres --column-inserts --data-only --table service_params --table re_services --table service_item_param_keys dev_db > pkg/testdatagen/testdata/fixture_service_item_params.sql

	// we don't know where in the pkg dir the tests will be so we have to find pgk

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("==============%s\n\n", filepath.Clean(path))
	fixtureData := filepath.Join("..", "..", "testdatagen", "testdata", "fixture_service_item_params.sql") // relative path
	bytes, err := ioutil.ReadFile(filepath.Clean(fixtureData))
	if err != nil {
		log.Fatalf("loading fixture data <%s>: %+v", fixtureData, err)
	}

	err = db.RawQuery(string(bytes)).Exec()
	if err != nil {
		log.Fatalf("importing fixture data <%s>: %+v", bytes, err)
	}
}

// MakeRealMTOServiceItemWithAllDeps Takes a service code, move, shipment
// and creates or finds all the needed data to create a service item all its params ready for pricing
func MakeRealMTOServiceItemWithAllDeps(db *pop.Connection, serviceCode models.ReServiceCode, mto models.Move, mtoShipment models.MTOShipment) models.MTOServiceItem {
	loadFixtureDataServiceItemParams(db)
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

	log.Panicf("couldn't create service item service code %s not defined", serviceCode)
	return models.MTOServiceItem{}
}
