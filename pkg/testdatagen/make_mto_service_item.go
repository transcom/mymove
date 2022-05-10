package testdatagen

import (
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

var (
	paramActualPickupDate = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameActualPickupDate,
		Description: "actual pickup date",
		Type:        models.ServiceItemParamTypeDate,
		Origin:      models.ServiceItemParamOriginPrime,
	}
	paramContractCode = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameContractCode,
		Description: "contract code",
		Type:        models.ServiceItemParamTypeString,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramContractYearName = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameContractYearName,
		Description: "contract year name",
		Type:        models.ServiceItemParamTypeString,
		Origin:      models.ServiceItemParamOriginPricer,
	}
	paramDistanceZip3 = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameDistanceZip3,
		Description: "distance zip3",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramDistanceZip5 = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameDistanceZip5,
		Description: "distance zip5",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramEIAFuelPrice = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameEIAFuelPrice,
		Description: "eia fuel price",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramEscalationCompounded = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameEscalationCompounded,
		Description: "escalation compounded",
		Type:        models.ServiceItemParamTypeDecimal,
		Origin:      models.ServiceItemParamOriginPricer,
	}
	paramFSCMultiplier = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameFSCMultiplier,
		Description: "fsc multiplier",
		Type:        models.ServiceItemParamTypeDecimal,
		Origin:      models.ServiceItemParamOriginPricer,
	}
	paramFSCPriceDifferenceInCents = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameFSCPriceDifferenceInCents,
		Description: "fsc price difference in cents",
		Type:        models.ServiceItemParamTypeDecimal,
		Origin:      models.ServiceItemParamOriginPricer,
	}
	paramFSCWeightBasedDistanceMultiplier = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
		Description: "fsc weight based multiplier",
		Type:        models.ServiceItemParamTypeDecimal,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramIsPeak = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameIsPeak,
		Description: "is peak",
		Type:        models.ServiceItemParamTypeBoolean,
		Origin:      models.ServiceItemParamOriginPricer,
	}
	paramMTOAvailableAToPrimeAt = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameMTOAvailableToPrimeAt,
		Description: "mto available to prime at",
		Type:        models.ServiceItemParamTypeTimestamp,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramNumberDaysSIT = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameNumberDaysSIT,
		Description: "number days SIT",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginPrime,
	}
	paramPriceRateOrFactor = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNamePriceRateOrFactor,
		Description: "price, rate, or factor",
		Type:        models.ServiceItemParamTypeDecimal,
		Origin:      models.ServiceItemParamOriginPricer,
	}
	paramReferenceDate = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameReferenceDate,
		Description: "reference date",
		Type:        models.ServiceItemParamTypeDate,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramRequestedPickupDate = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameRequestedPickupDate,
		Description: "requested pickup date",
		Type:        models.ServiceItemParamTypeDate,
		Origin:      models.ServiceItemParamOriginPrime,
	}
	paramServiceAreaOrigin = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameServiceAreaOrigin,
		Description: "service area origin",
		Type:        models.ServiceItemParamTypeString,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramServicesScheduleOrigin = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameServicesScheduleOrigin,
		Description: "services schedule origin",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramSITPaymentRequestEnd = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameSITPaymentRequestEnd,
		Description: "SIT payment request end",
		Type:        models.ServiceItemParamTypeDate,
		Origin:      models.ServiceItemParamOriginPaymentRequest,
	}
	paramSITPaymentRequestStart = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameSITPaymentRequestStart,
		Description: "SIT payment request start",
		Type:        models.ServiceItemParamTypeDate,
		Origin:      models.ServiceItemParamOriginPaymentRequest,
	}
	paramWeightAdjusted = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameWeightAdjusted,
		Description: "weight adjusted",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramWeightBilled = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameWeightBilled,
		Description: "weight billed",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramWeightEstimated = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameWeightEstimated,
		Description: "weight estimated",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginPrime,
	}
	paramWeightOriginal = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameWeightOriginal,
		Description: "weight original",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginPrime,
	}
	paramWeightReweigh = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameWeightReweigh,
		Description: "weight reweigh",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginPrime,
	}
	paramZipDestAddress = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameZipDestAddress,
		Description: "zip dest address",
		Type:        models.ServiceItemParamTypeString,
		Origin:      models.ServiceItemParamOriginPrime,
	}
	paramZipPickupAddress = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameZipPickupAddress,
		Description: "zip pickup address",
		Type:        models.ServiceItemParamTypeString,
		Origin:      models.ServiceItemParamOriginPrime,
	}
)

var fixtureServiceItemParamsMap = map[models.ReServiceCode]models.ServiceItemParamKeys{
	models.ReServiceCodeCS: {
		paramContractCode,
		paramMTOAvailableAToPrimeAt,
		paramPriceRateOrFactor,
	},
	models.ReServiceCodeMS: {
		paramContractCode,
		paramMTOAvailableAToPrimeAt,
		paramPriceRateOrFactor,
	},
	models.ReServiceCodeDLH: {
		paramActualPickupDate,
		paramContractCode,
		paramContractYearName,
		paramDistanceZip3,
		paramEscalationCompounded,
		paramIsPeak,
		paramPriceRateOrFactor,
		paramReferenceDate,
		paramRequestedPickupDate,
		paramServiceAreaOrigin,
		paramWeightAdjusted,
		paramWeightBilled,
		paramWeightEstimated,
		paramWeightOriginal,
		paramWeightReweigh,
		paramZipDestAddress,
		paramZipPickupAddress,
	},
	models.ReServiceCodeFSC: {
		paramActualPickupDate,
		paramContractCode,
		paramDistanceZip3,
		paramDistanceZip5,
		paramEIAFuelPrice,
		paramFSCMultiplier,
		paramFSCPriceDifferenceInCents,
		paramFSCWeightBasedDistanceMultiplier,
		paramWeightAdjusted,
		paramWeightBilled,
		paramWeightEstimated,
		paramWeightOriginal,
		paramWeightReweigh,
		paramZipDestAddress,
		paramZipPickupAddress,
	},
	models.ReServiceCodeDPK: {
		paramActualPickupDate,
		paramContractCode,
		paramContractYearName,
		paramEscalationCompounded,
		paramIsPeak,
		paramPriceRateOrFactor,
		paramReferenceDate,
		paramRequestedPickupDate,
		paramServiceAreaOrigin,
		paramServicesScheduleOrigin,
		paramWeightAdjusted,
		paramWeightBilled,
		paramWeightEstimated,
		paramWeightOriginal,
		paramWeightReweigh,
		paramZipPickupAddress,
	},
	models.ReServiceCodeDOP: {
		paramActualPickupDate,
		paramContractCode,
		paramContractYearName,
		paramEscalationCompounded,
		paramIsPeak,
		paramPriceRateOrFactor,
		paramReferenceDate,
		paramRequestedPickupDate,
		paramServiceAreaOrigin,
		paramWeightAdjusted,
		paramWeightBilled,
		paramWeightEstimated,
		paramWeightOriginal,
		paramWeightReweigh,
		paramZipPickupAddress,
	},
	models.ReServiceCodeDOASIT: {
		paramActualPickupDate,
		paramContractCode,
		paramContractYearName,
		paramEscalationCompounded,
		paramIsPeak,
		paramNumberDaysSIT,
		paramPriceRateOrFactor,
		paramReferenceDate,
		paramRequestedPickupDate,
		paramServiceAreaOrigin,
		paramSITPaymentRequestEnd,
		paramSITPaymentRequestStart,
		paramWeightAdjusted,
		paramWeightBilled,
		paramWeightEstimated,
		paramWeightOriginal,
		paramWeightReweigh,
		paramZipPickupAddress,
	},
}

// makeServiceItem creates a single service item and associated set relationships
func makeServiceItem(db *pop.Connection, assertions Assertions, isBasicServiceItem bool) models.MTOServiceItem {
	moveTaskOrder := assertions.Move
	if isZeroUUID(moveTaskOrder.ID) {
		moveTaskOrder = MakeMove(db, assertions)
	}

	var mtoShipmentID *uuid.UUID
	var mtoShipment models.MTOShipment
	if !isBasicServiceItem {
		if isZeroUUID(assertions.MTOShipment.ID) {
			mtoShipment = MakeMTOShipment(db, assertions)
			mtoShipmentID = &mtoShipment.ID
		} else {
			mtoShipment = assertions.MTOShipment
			mtoShipmentID = &assertions.MTOShipment.ID
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
		MTOShipment:     mtoShipment,
		MTOShipmentID:   mtoShipmentID,
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

// MakeMTOServiceItem creates a single MTOServiceItem and associated set relationships
func MakeStubbedMTOServiceItem(db *pop.Connection) models.MTOServiceItem {
	return makeServiceItem(db, Assertions{
		Stub: true,
	}, false)
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

// MakeRealMTOServiceItemWithAllDeps Takes a service code, move, shipment
// and creates or finds all the needed data to create a service item all its params ready for pricing
func MakeRealMTOServiceItemWithAllDeps(db *pop.Connection, serviceCode models.ReServiceCode, mto models.Move, mtoShipment models.MTOShipment) models.MTOServiceItem {
	// look up the service item param keys we need
	if serviceItemParamKeys, ok := fixtureServiceItemParamsMap[serviceCode]; ok {
		// get or create the ReService
		reService := FetchOrMakeReService(db, Assertions{
			ReService: models.ReService{
				Code: serviceCode,
			},
		})

		// create all params defined for this particular service
		for _, serviceParamKeyToCreate := range serviceItemParamKeys {
			serviceItemParamKey := FetchOrMakeServiceItemParamKey(db, Assertions{
				ServiceItemParamKey: serviceParamKeyToCreate,
			})
			_ = MakeServiceParam(db, Assertions{
				ServiceParam: models.ServiceParam{
					ServiceID:             reService.ID,
					Service:               reService,
					ServiceItemParamKeyID: serviceItemParamKey.ID,
					ServiceItemParamKey:   serviceItemParamKey,
				},
			})
		}

		// create a service item and return it
		mtoServiceItem := MakeMTOServiceItem(db, Assertions{
			Move:        mto,
			MTOShipment: mtoShipment,
			ReService:   reService,
		})

		return mtoServiceItem
	}

	log.Panicf("couldn't create service item service code %s not defined", serviceCode)
	return models.MTOServiceItem{}
}

// MakeMTOServiceItemDomesticCrating makes a domestic crating service item and its associated item and crate
func MakeMTOServiceItemDomesticCrating(db *pop.Connection, assertions Assertions) models.MTOServiceItem {
	mtoServiceItem := MakeMTOServiceItem(db, assertions)

	// Create item
	dimensionItem := MakeMTOServiceItemDimension(db, Assertions{
		MTOServiceItemDimension: assertions.MTOServiceItemDimension,
		MTOServiceItem:          mtoServiceItem,
	})

	// Create crate
	assertions.MTOServiceItemDimensionCrate.Type = models.DimensionTypeCrate
	crateItem := MakeMTOServiceItemDimension(db, Assertions{
		MTOServiceItemDimension: assertions.MTOServiceItemDimensionCrate,
		MTOServiceItem:          mtoServiceItem,
	})

	mtoServiceItem.Dimensions = append(mtoServiceItem.Dimensions, dimensionItem, crateItem)

	return mtoServiceItem
}
