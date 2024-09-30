package factory

import (
	"log"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

type mtoServiceItemBuildType byte

const (
	mtoServiceItemBuildBasic mtoServiceItemBuildType = iota
	mtoServiceItemBuildExtended
)

// buildMTOServiceItemWithBuildType creates a single MTOServiceItem.
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func buildMTOServiceItemWithBuildType(db *pop.Connection, customs []Customization, traits []Trait, buildType mtoServiceItemBuildType) models.MTOServiceItem {
	customs = setupCustomizations(customs, traits)

	// Find address customization and extract the custom address
	var cMTOServiceItem models.MTOServiceItem
	if result := findValidCustomization(customs, MTOServiceItem); result != nil {
		cMTOServiceItem = result.Model.(models.MTOServiceItem)
		if result.LinkOnly {
			return cMTOServiceItem
		}
	}

	var mtoShipmentID *uuid.UUID
	var mtoShipment models.MTOShipment
	var move models.Move
	var isCustomerExpense = false
	if buildType == mtoServiceItemBuildExtended {
		// BuildMTOShipment creates a move as necessary
		mtoShipment = BuildMTOShipment(db, customs, traits)
		mtoShipmentID = &mtoShipment.ID
		move = mtoShipment.MoveTaskOrder
	} else {
		move = BuildMove(db, customs, traits)
	}

	var reService models.ReService
	if result := findValidCustomization(customs, ReService); result != nil {
		reService = FetchOrBuildReService(db, customs, nil)
	} else {
		reService = FetchOrBuildReServiceByCode(db, models.ReServiceCode("STEST"))
	}

	requestedApprovalsRequestedStatus := false

	var lockedPriceCents = unit.Cents(12303)

	// Create default MTOServiceItem
	mtoServiceItem := models.MTOServiceItem{
		MoveTaskOrder:                     move,
		MoveTaskOrderID:                   move.ID,
		MTOShipment:                       mtoShipment,
		MTOShipmentID:                     mtoShipmentID,
		ReService:                         reService,
		ReServiceID:                       reService.ID,
		Status:                            models.MTOServiceItemStatusSubmitted,
		RequestedApprovalsRequestedStatus: &requestedApprovalsRequestedStatus,
		CustomerExpense:                   isCustomerExpense,
		LockedPriceCents:                  &lockedPriceCents,
	}

	// only set SITOriginHHGOriginalAddress if a customization is provided
	if result := findValidCustomization(customs, Addresses.SITOriginHHGOriginalAddress); result != nil {
		addressCustoms := convertCustomizationInList(customs, Addresses.SITOriginHHGOriginalAddress, Address)
		address := BuildAddress(db, addressCustoms, traits)
		mtoServiceItem.SITOriginHHGOriginalAddress = &address
		mtoServiceItem.SITOriginHHGOriginalAddressID = &address.ID
	}

	// only set SITOriginHHGActualAddress if a customization is provided
	if result := findValidCustomization(customs, Addresses.SITOriginHHGActualAddress); result != nil {
		addressCustoms := convertCustomizationInList(customs, Addresses.SITOriginHHGActualAddress, Address)
		address := BuildAddress(db, addressCustoms, traits)
		mtoServiceItem.SITOriginHHGActualAddress = &address
		mtoServiceItem.SITOriginHHGActualAddressID = &address.ID
	}

	// only set SITDestinationFinalAddress if a customization is provided
	if result := findValidCustomization(customs, Addresses.SITDestinationFinalAddress); result != nil {
		addressCustoms := convertCustomizationInList(customs, Addresses.SITDestinationFinalAddress, Address)
		address := BuildAddress(db, addressCustoms, traits)
		mtoServiceItem.SITDestinationFinalAddress = &address
		mtoServiceItem.SITDestinationFinalAddressID = &address.ID
	}

	// only set SITDestinationOriginalAddress if a customization is provided
	if result := findValidCustomization(customs, Addresses.SITDestinationOriginalAddress); result != nil {
		addressCustoms := convertCustomizationInList(customs, Addresses.SITDestinationOriginalAddress, Address)
		address := BuildAddress(db, addressCustoms, traits)
		mtoServiceItem.SITDestinationOriginalAddress = &address
		mtoServiceItem.SITDestinationOriginalAddressID = &address.ID
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&mtoServiceItem, cMTOServiceItem)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &mtoServiceItem)
	}

	return mtoServiceItem
}

// BuildMTOServiceItem creates a single extended MTOServiceItem
func BuildMTOServiceItem(db *pop.Connection, customs []Customization, traits []Trait) models.MTOServiceItem {
	return buildMTOServiceItemWithBuildType(db, customs, traits, mtoServiceItemBuildExtended)
}

// BuildMTOServiceItemBasic creates a single basic MTOServiceItem
func BuildMTOServiceItemBasic(db *pop.Connection, customs []Customization, traits []Trait) models.MTOServiceItem {
	return buildMTOServiceItemWithBuildType(db, customs, traits, mtoServiceItemBuildBasic)
}

// Needed by BuildRealMTOServiceItemWithAllDeps

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
	paramDistanceZip = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameDistanceZip,
		Description: "distance zip",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramDistanceZipSITOrigin = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameDistanceZipSITOrigin,
		Description: "distance zip SIT origin",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramDistanceZipSITDest = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameDistanceZipSITDest,
		Description: "distance zip SIT destination",
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
	paramServiceAreaDest = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameServiceAreaDest,
		Description: "Destination service area",
		Type:        models.ServiceItemParamTypeString,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramServiceAreaOrigin = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameServiceAreaOrigin,
		Description: "Origin service area",
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
	paramSITScheduleDest = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameSITScheduleDest,
		Description: "Origin SIT schedule",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramSITScheduleOrigin = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameSITScheduleOrigin,
		Description: "Dest SIT schedule",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramSITServiceAreaOrigin = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameSITServiceAreaOrigin,
		Description: "SIT Origin service area",
		Type:        models.ServiceItemParamTypeString,
		Origin:      models.ServiceItemParamOriginSystem,
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
	paramZipSITOriginHHGActualAddress = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameZipSITOriginHHGActualAddress,
		Description: "zip dest address SIT origin HHG actual address",
		Type:        models.ServiceItemParamTypeString,
		Origin:      models.ServiceItemParamOriginPrime,
	}
	paramZipSITDestHHGFinalAddress = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameZipSITDestHHGFinalAddress,
		Description: "zip dest address SIT destination HHG final address",
		Type:        models.ServiceItemParamTypeString,
		Origin:      models.ServiceItemParamOriginPrime,
	}
	paramZipPickupAddress = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameZipPickupAddress,
		Description: "zip pickup address",
		Type:        models.ServiceItemParamTypeString,
		Origin:      models.ServiceItemParamOriginPrime,
	}
	paramZipSITOriginHHGOriginalAddress = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameZipSITOriginHHGOriginalAddress,
		Description: "zip pickup address SIT origin HHG original address",
		Type:        models.ServiceItemParamTypeString,
		Origin:      models.ServiceItemParamOriginPrime,
	}
	paramZipSITDestHHGOriginalAddress = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameZipSITDestHHGOriginalAddress,
		Description: "zip pickup address SIT destination HHG original address",
		Type:        models.ServiceItemParamTypeString,
		Origin:      models.ServiceItemParamOriginPrime,
	}
	paramLockedPriceCents = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameLockedPriceCents,
		Description: "locked price cents",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	fixtureServiceItemParamsMap = map[models.ReServiceCode]models.ServiceItemParamKeys{
		models.ReServiceCodeCS: {
			paramMTOAvailableAToPrimeAt,
			paramContractCode,
			paramLockedPriceCents,
			paramPriceRateOrFactor,
		},
		models.ReServiceCodeMS: {
			paramMTOAvailableAToPrimeAt,
			paramContractCode,
			paramLockedPriceCents,
			paramPriceRateOrFactor,
		},
		models.ReServiceCodeDLH: {
			paramActualPickupDate,
			paramContractCode,
			paramContractYearName,
			paramDistanceZip,
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
		models.ReServiceCodeDSH: {
			paramActualPickupDate,
			paramContractCode,
			paramContractYearName,
			paramDistanceZip,
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
			paramDistanceZip,
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
			paramSITServiceAreaOrigin,
			paramSITPaymentRequestEnd,
			paramSITPaymentRequestStart,
			paramWeightAdjusted,
			paramWeightBilled,
			paramWeightEstimated,
			paramWeightOriginal,
			paramWeightReweigh,
			paramZipPickupAddress,
		},
		models.ReServiceCodeDOFSIT: {
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
		models.ReServiceCodeDOPSIT: {
			paramActualPickupDate,
			paramContractCode,
			paramContractYearName,
			paramDistanceZipSITOrigin,
			paramEscalationCompounded,
			paramIsPeak,
			paramPriceRateOrFactor,
			paramReferenceDate,
			paramRequestedPickupDate,
			paramServiceAreaOrigin,
			paramSITScheduleOrigin,
			paramWeightAdjusted,
			paramWeightBilled,
			paramWeightEstimated,
			paramWeightOriginal,
			paramWeightReweigh,
			paramZipSITOriginHHGActualAddress,
			paramZipSITOriginHHGOriginalAddress,
		},
		models.ReServiceCodeDOSFSC: {
			paramActualPickupDate,
			paramDistanceZipSITOrigin,
			paramFSCWeightBasedDistanceMultiplier,
			paramEIAFuelPrice,
			paramWeightBilled,
			paramWeightAdjusted,
			paramWeightOriginal,
			paramWeightEstimated,
			paramZipSITOriginHHGOriginalAddress,
			paramZipSITOriginHHGActualAddress,
			paramFSCPriceDifferenceInCents,
			paramContractCode,
			paramFSCMultiplier,
		},
		models.ReServiceCodeDDASIT: {
			paramActualPickupDate,
			paramContractCode,
			paramContractYearName,
			paramEscalationCompounded,
			paramIsPeak,
			paramNumberDaysSIT,
			paramPriceRateOrFactor,
			paramReferenceDate,
			paramRequestedPickupDate,
			paramServiceAreaDest,
			paramSITPaymentRequestEnd,
			paramSITPaymentRequestStart,
			paramWeightAdjusted,
			paramWeightBilled,
			paramWeightEstimated,
			paramWeightOriginal,
			paramWeightReweigh,
			paramZipDestAddress,
		},
		models.ReServiceCodeDDFSIT: {
			paramActualPickupDate,
			paramContractCode,
			paramContractYearName,
			paramEscalationCompounded,
			paramIsPeak,
			paramPriceRateOrFactor,
			paramReferenceDate,
			paramRequestedPickupDate,
			paramServiceAreaDest,
			paramWeightAdjusted,
			paramWeightBilled,
			paramWeightEstimated,
			paramWeightOriginal,
			paramWeightReweigh,
			paramZipDestAddress,
		},
		models.ReServiceCodeDDDSIT: {
			paramActualPickupDate,
			paramContractCode,
			paramContractYearName,
			paramDistanceZipSITDest,
			paramEscalationCompounded,
			paramIsPeak,
			paramPriceRateOrFactor,
			paramReferenceDate,
			paramRequestedPickupDate,
			paramSITScheduleDest,
			paramServiceAreaDest,
			paramWeightAdjusted,
			paramWeightBilled,
			paramWeightEstimated,
			paramWeightOriginal,
			paramWeightReweigh,
			paramZipSITDestHHGFinalAddress,
			paramZipSITDestHHGOriginalAddress,
		},
		models.ReServiceCodeDDSFSC: {
			paramActualPickupDate,
			paramDistanceZipSITDest,
			paramFSCWeightBasedDistanceMultiplier,
			paramEIAFuelPrice,
			paramWeightBilled,
			paramWeightAdjusted,
			paramWeightOriginal,
			paramWeightEstimated,
			paramZipSITDestHHGOriginalAddress,
			paramZipSITDestHHGFinalAddress,
			paramFSCPriceDifferenceInCents,
			paramContractCode,
			paramFSCMultiplier,
		},
	}
)

// BuildRealMTOServiceItemWithAllDeps builds an MTOServiceItem along with its service params if they don't already exist
// Customizations that link the MTOServiceItem to the move, shipment, and service are created by default
func BuildRealMTOServiceItemWithAllDeps(db *pop.Connection, serviceCode models.ReServiceCode, mto models.Move, mtoShipment models.MTOShipment, customs []Customization, traits []Trait) models.MTOServiceItem {
	// look up the service item param keys we need
	if serviceItemParamKeys, ok := fixtureServiceItemParamsMap[serviceCode]; ok {
		// get or create the ReService
		reService := FetchOrBuildReServiceByCode(db, serviceCode)

		// create all params defined for this particular service
		for _, serviceParamKeyToCreate := range serviceItemParamKeys {
			serviceItemParamKey := FetchOrBuildServiceItemParamKey(db, []Customization{
				{
					Model: serviceParamKeyToCreate,
				},
			}, nil)
			FetchOrBuildServiceParam(db, []Customization{
				{
					Model:    reService,
					LinkOnly: true,
				},
				{
					Model:    serviceItemParamKey,
					LinkOnly: true,
				},
			}, nil)
		}

		allCustoms := []Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
			{
				Model:    reService,
				LinkOnly: true,
			},
		}
		if len(customs) > 0 {
			allCustoms = append(allCustoms, customs...)
		}

		allTraits := []Trait{GetTraitServiceItemStatusApproved}
		if len(traits) > 0 {
			allTraits = append(allTraits, traits...)
		}

		// create a service item and return it
		mtoServiceItem := BuildMTOServiceItem(db, allCustoms, allTraits)
		return mtoServiceItem
	}

	log.Panicf("couldn't create service item service code %s not defined", serviceCode)
	return models.MTOServiceItem{}

}

// BuildFullDLHMTOServiceItems makes a DLH type service item along with
// all its expected parameters returns the created move and all
// service items
//
// NOTE: the original did an override of the MTOShipment.Status to
// ensure it was Approved, but that is now the responsibility of the
// caller
func BuildFullDLHMTOServiceItems(db *pop.Connection, customs []Customization, traits []Trait) (models.Move, models.MTOServiceItems) {

	mtoShipment := BuildMTOShipment(db, customs, traits)

	move := mtoShipment.MoveTaskOrder
	move.MTOShipments = models.MTOShipments{mtoShipment}

	var mtoServiceItems models.MTOServiceItems
	// Service Item MS
	mtoServiceItemMS := BuildRealMTOServiceItemWithAllDeps(db,
		models.ReServiceCodeMS, move, mtoShipment, nil, nil)
	mtoServiceItems = append(mtoServiceItems, mtoServiceItemMS)
	// Service Item CS
	mtoServiceItemCS := BuildRealMTOServiceItemWithAllDeps(db,
		models.ReServiceCodeCS, move, mtoShipment, nil, nil)
	mtoServiceItems = append(mtoServiceItems, mtoServiceItemCS)
	// Service Item DLH
	mtoServiceItemDLH := BuildRealMTOServiceItemWithAllDeps(db,
		models.ReServiceCodeDLH, move, mtoShipment, nil, nil)
	mtoServiceItems = append(mtoServiceItems, mtoServiceItemDLH)
	// Service Item FSC
	mtoServiceItemFSC := BuildRealMTOServiceItemWithAllDeps(db,
		models.ReServiceCodeFSC, move, mtoShipment, nil, nil)
	mtoServiceItems = append(mtoServiceItems, mtoServiceItemFSC)

	return move, mtoServiceItems
}

// BuildFullOriginMTOServiceItems (follow-on to
// BuildFullDLHMTOServiceItem) makes a DLH type service item along
// with all its expected parameters returns the created move and all
// service items
//
// NOTE: the original did an override of the MTOShipment.Status to
// ensure it was Approved, but that is now the responsibility of the
// caller
func BuildFullOriginMTOServiceItems(db *pop.Connection, customs []Customization, traits []Trait) (models.Move, models.MTOServiceItems) {
	mtoShipment := BuildMTOShipment(db, customs, traits)

	move := mtoShipment.MoveTaskOrder
	move.MTOShipments = models.MTOShipments{mtoShipment}

	var mtoServiceItems models.MTOServiceItems
	// Service Item DPK
	mtoServiceItemDPK := BuildRealMTOServiceItemWithAllDeps(db,
		models.ReServiceCodeDPK, move, mtoShipment, nil, nil)
	mtoServiceItems = append(mtoServiceItems, mtoServiceItemDPK)
	// Service Item DOP
	mtoServiceItemDOP := BuildRealMTOServiceItemWithAllDeps(db,
		models.ReServiceCodeDOP, move, mtoShipment, nil, nil)
	mtoServiceItems = append(mtoServiceItems, mtoServiceItemDOP)

	return move, mtoServiceItems
}

// BuildOriginSITServiceItems makes all of the service items that are
// associated with Origin SIT. A move and shipment that exist in the db
// are required params, and entryDate and departureDate can be specificed
// optionally.
func BuildOriginSITServiceItems(db *pop.Connection, move models.Move, shipment models.MTOShipment, entryDate *time.Time, departureDate *time.Time) models.MTOServiceItems {
	postalCode := "90210"
	reason := "peak season all trucks in use"
	defaultEntryDate := time.Now().AddDate(0, 0, -45)
	defaultApprovedAtDate := time.Now()
	if entryDate != nil {
		defaultEntryDate = *entryDate
	}
	var defaultDepartureDate *time.Time
	if departureDate != nil {
		defaultDepartureDate = departureDate
	}

	dofsit := BuildRealMTOServiceItemWithAllDeps(db, models.ReServiceCodeDOFSIT, move, shipment, []Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				ApprovedAt:    &defaultApprovedAtDate,
				SITEntryDate:  &defaultEntryDate,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
		{
			Model: models.Address{},
			Type:  &Addresses.SITOriginHHGActualAddress,
		},
		{
			Model: models.Address{},
			Type:  &Addresses.SITOriginHHGOriginalAddress,
		},
	}, nil)

	doasit := BuildRealMTOServiceItemWithAllDeps(db, models.ReServiceCodeDOASIT, move, shipment, []Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				ApprovedAt:    &defaultApprovedAtDate,
				SITEntryDate:  &defaultEntryDate,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
		{
			Model: models.Address{},
			Type:  &Addresses.SITOriginHHGActualAddress,
		},
		{
			Model: models.Address{},
			Type:  &Addresses.SITOriginHHGOriginalAddress,
		},
	}, nil)

	dopsit := BuildRealMTOServiceItemWithAllDeps(db, models.ReServiceCodeDOPSIT, move, shipment, []Customization{
		{
			Model: models.MTOServiceItem{
				Status:           models.MTOServiceItemStatusApproved,
				ApprovedAt:       &defaultApprovedAtDate,
				SITEntryDate:     &defaultEntryDate,
				SITDepartureDate: defaultDepartureDate,
				SITPostalCode:    &postalCode,
				Reason:           &reason,
			},
		},
		{
			Model: models.Address{},
			Type:  &Addresses.SITOriginHHGActualAddress,
		},
		{
			Model: models.Address{},
			Type:  &Addresses.SITOriginHHGOriginalAddress,
		},
	}, nil)

	dosfsc := BuildRealMTOServiceItemWithAllDeps(db, models.ReServiceCodeDOSFSC, move, shipment, []Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				ApprovedAt:    &defaultApprovedAtDate,
				SITEntryDate:  &defaultEntryDate,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
		{
			Model: models.Address{},
			Type:  &Addresses.SITOriginHHGActualAddress,
		},
		{
			Model: models.Address{},
			Type:  &Addresses.SITOriginHHGOriginalAddress,
		},
	}, nil)
	return []models.MTOServiceItem{dofsit, doasit, dopsit, dosfsc}
}

// BuildDestSITServiceItems makes all of the service items that are
// associated with Destination SIT. A move and shipment that exist in the db
// are required params, and entryDate and departureDate can be specificed
// optionally.
func BuildDestSITServiceItems(db *pop.Connection, move models.Move, shipment models.MTOShipment, entryDate *time.Time, departureDate *time.Time) models.MTOServiceItems {
	postalCode := "90210"
	reason := "peak season all trucks in use"
	defaultEntryDate := time.Now().AddDate(0, 0, -45)
	defaultApprovedAtDate := time.Now()
	if entryDate != nil {
		defaultEntryDate = *entryDate
	}
	var defaultDepartureDate *time.Time
	if departureDate != nil {
		defaultDepartureDate = departureDate
	}

	ddfsit := BuildRealMTOServiceItemWithAllDeps(db, models.ReServiceCodeDDFSIT, move, shipment, []Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				ApprovedAt:    &defaultApprovedAtDate,
				SITEntryDate:  &defaultEntryDate,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
	}, nil)

	ddasit := BuildRealMTOServiceItemWithAllDeps(db, models.ReServiceCodeDDASIT, move, shipment, []Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				ApprovedAt:    &defaultApprovedAtDate,
				SITEntryDate:  &defaultEntryDate,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
	}, nil)

	dddsit := BuildRealMTOServiceItemWithAllDeps(db, models.ReServiceCodeDDDSIT, move, shipment, []Customization{
		{
			Model: models.MTOServiceItem{
				Status:           models.MTOServiceItemStatusApproved,
				ApprovedAt:       &defaultApprovedAtDate,
				SITEntryDate:     &defaultEntryDate,
				SITDepartureDate: defaultDepartureDate,
				SITPostalCode:    &postalCode,
				Reason:           &reason,
			},
		},
		{
			Model: models.Address{},
			Type:  &Addresses.SITDestinationFinalAddress,
		},
		{
			Model: models.Address{},
			Type:  &Addresses.SITDestinationOriginalAddress,
		},
	}, nil)

	ddsfsc := BuildRealMTOServiceItemWithAllDeps(db, models.ReServiceCodeDDSFSC, move, shipment, []Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				ApprovedAt:    &defaultApprovedAtDate,
				SITEntryDate:  &defaultEntryDate,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
		{
			Model: models.Address{},
			Type:  &Addresses.SITDestinationFinalAddress,
		},
		{
			Model: models.Address{},
			Type:  &Addresses.SITDestinationOriginalAddress,
		},
	}, nil)
	return []models.MTOServiceItem{ddfsit, ddasit, dddsit, ddsfsc}
}

func BuildDestSITServiceItemsNoSITDepartureDate(db *pop.Connection, move models.Move, shipment models.MTOShipment, entryDate *time.Time) models.MTOServiceItems {
	postalCode := "90210"
	reason := "peak season all trucks in use"
	defaultEntryDate := time.Now().AddDate(0, 0, -45)
	defaultApprovedAtDate := time.Now()
	if entryDate != nil {
		defaultEntryDate = *entryDate
	}

	ddfsit := BuildRealMTOServiceItemWithAllDeps(db, models.ReServiceCodeDDFSIT, move, shipment, []Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				ApprovedAt:    &defaultApprovedAtDate,
				SITEntryDate:  &defaultEntryDate,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
	}, nil)

	ddasit := BuildRealMTOServiceItemWithAllDeps(db, models.ReServiceCodeDDASIT, move, shipment, []Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				ApprovedAt:    &defaultApprovedAtDate,
				SITEntryDate:  &defaultEntryDate,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
	}, nil)

	dddsit := BuildRealMTOServiceItemWithAllDeps(db, models.ReServiceCodeDDDSIT, move, shipment, []Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				ApprovedAt:    &defaultApprovedAtDate,
				SITEntryDate:  &defaultEntryDate,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
		{
			Model: models.Address{},
			Type:  &Addresses.SITDestinationFinalAddress,
		},
		{
			Model: models.Address{},
			Type:  &Addresses.SITDestinationOriginalAddress,
		},
	}, nil)

	ddsfsc := BuildRealMTOServiceItemWithAllDeps(db, models.ReServiceCodeDDSFSC, move, shipment, []Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				ApprovedAt:    &defaultApprovedAtDate,
				SITEntryDate:  &defaultEntryDate,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
		{
			Model: models.Address{},
			Type:  &Addresses.SITDestinationFinalAddress,
		},
		{
			Model: models.Address{},
			Type:  &Addresses.SITDestinationOriginalAddress,
		},
	}, nil)
	return []models.MTOServiceItem{ddfsit, ddasit, dddsit, ddsfsc}
}

// ------------------------
//        TRAITS
// ------------------------

func GetTraitServiceItemStatusApproved() []Customization {
	return []Customization{
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
		},
	}
}
