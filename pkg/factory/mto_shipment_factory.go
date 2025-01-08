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

// GHCTestYear is the default for GHC rate engine testing
var GHCTestYear = 2020

type mtoShipmentBuildType byte

const (
	mtoShipmentBuildBasic mtoShipmentBuildType = iota
	mtoShipmentBuild
	mtoShipmentNTS
	mtoShipmentPPM
	mtoShipmentNTSR
)

func buildMTOShipmentWithBuildType(db *pop.Connection, customs []Customization, traits []Trait, buildType mtoShipmentBuildType) models.MTOShipment {
	customs = setupCustomizations(customs, traits)

	// Find mtoShipment customization and convert to mtoShipment model
	var cMtoShipment models.MTOShipment
	if result := findValidCustomization(customs, MTOShipment); result != nil {
		cMtoShipment = result.Model.(models.MTOShipment)
		if result.LinkOnly {
			return cMtoShipment
		}
	}

	move := BuildMove(db, customs, traits)

	// defaults change depending on mtoshipment build type
	defaultShipmentType := models.MTOShipmentTypeHHG
	defaultStatus := models.MTOShipmentStatusSubmitted
	defaultMarketCode := models.MarketCodeDomestic
	setupPickupAndDelivery := true
	hasStorageFacilityCustom := findValidCustomization(customs, StorageFacility) != nil
	buildStorageFacility :=
		cMtoShipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom ||
			cMtoShipment.ShipmentType == models.MTOShipmentTypeHHGIntoNTS
	shipmentHasPickupDetails := cMtoShipment.ShipmentType != models.MTOShipmentTypeHHGOutOfNTSDom && cMtoShipment.ShipmentType != models.MTOShipmentTypePPM
	shipmentHasDeliveryDetails := cMtoShipment.ShipmentType != models.MTOShipmentTypeHHGIntoNTS && cMtoShipment.ShipmentType != models.MTOShipmentTypePPM
	addPrimeActualWeight := true
	switch buildType {
	case mtoShipmentNTS:
		defaultShipmentType = models.MTOShipmentTypeHHGIntoNTS
		defaultStatus = models.MTOShipmentStatusDraft
		buildStorageFacility = hasStorageFacilityCustom
		shipmentHasPickupDetails = true
		shipmentHasDeliveryDetails = false
	case mtoShipmentNTSR:
		defaultShipmentType = models.MTOShipmentTypeHHGOutOfNTSDom
		defaultStatus = models.MTOShipmentStatusDraft
		buildStorageFacility = hasStorageFacilityCustom
		addPrimeActualWeight = false
		shipmentHasPickupDetails = false
		shipmentHasDeliveryDetails = true
	case mtoShipmentBuildBasic:
		setupPickupAndDelivery = false
	case mtoShipmentPPM:
		defaultShipmentType = models.MTOShipmentTypePPM
		setupPickupAndDelivery = false
	default:
		defaultShipmentType = models.MTOShipmentTypeHHG
		setupPickupAndDelivery = true
	}

	newMTOShipment := models.MTOShipment{
		MoveTaskOrder:   move,
		MoveTaskOrderID: move.ID,
		ShipmentType:    defaultShipmentType,
		Status:          defaultStatus,
		MarketCode:      defaultMarketCode,
	}

	if cMtoShipment.Status == models.MTOShipmentStatusApproved {
		approvedDate := time.Date(GHCTestYear, time.March, 20, 0, 0, 0, 0, time.UTC)
		newMTOShipment.ApprovedDate = &approvedDate
	}

	if setupPickupAndDelivery {
		newMTOShipment.Status = models.MTOShipmentStatusDraft

		if buildStorageFacility {
			storageFacility := BuildStorageFacility(db, customs, traits)
			// only set storage facility pointers if building a
			// storage facility
			newMTOShipment.StorageFacility = &storageFacility
			newMTOShipment.StorageFacilityID = &storageFacility.ID
		}

		if addPrimeActualWeight {
			actualWeight := unit.Pound(980)
			newMTOShipment.PrimeActualWeight = &actualWeight
		}
		newMTOShipment.CustomerRemarks = models.StringPointer("Please treat gently")

		if shipmentHasPickupDetails {
			newMTOShipment.RequestedPickupDate = models.TimePointer(time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC))
			newMTOShipment.ScheduledPickupDate = models.TimePointer(time.Date(GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC))
			newMTOShipment.ActualPickupDate = models.TimePointer(time.Date(GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC))
		}

		if shipmentHasPickupDetails || findValidCustomization(customs, Addresses.PickupAddress) != nil {
			// Find/create the Pickup Address
			tempPickupAddressCustoms := customs
			result := findValidCustomization(customs, Addresses.PickupAddress)
			if result != nil {
				tempPickupAddressCustoms = convertCustomizationInList(tempPickupAddressCustoms, Addresses.PickupAddress, Address)
			}

			pickupAddress := BuildAddress(db, tempPickupAddressCustoms, traits)
			if db == nil {
				// fake an id for stubbed address, needed by the MTOShipmentCreator
				pickupAddress.ID = uuid.Must(uuid.NewV4())
			}
			newMTOShipment.PickupAddress = &pickupAddress
			newMTOShipment.PickupAddressID = &pickupAddress.ID

			// Check that a GBLOC exists for pickup address postal code, make one if not
			if db != nil {
				FetchOrBuildPostalCodeToGBLOC(db, pickupAddress.PostalCode, "KKFA")
			}

			// Find Secondary Pickup Address
			tempSecondaryPicAddressCustoms := customs
			result = findValidCustomization(customs, Addresses.SecondaryPickupAddress)
			if result != nil {
				tempSecondaryPicAddressCustoms = convertCustomizationInList(tempSecondaryPicAddressCustoms, Addresses.SecondaryPickupAddress, Address)
				secondaryPickupAddress := BuildAddress(db, tempSecondaryPicAddressCustoms, traits)

				newMTOShipment.SecondaryPickupAddress = &secondaryPickupAddress
				newMTOShipment.SecondaryPickupAddressID = &secondaryPickupAddress.ID
				newMTOShipment.HasSecondaryPickupAddress = models.BoolPointer(true)
			}

			// Find Tertiary Pickup Address
			tempTertiaryPickupAddressCustoms := customs
			result = findValidCustomization(customs, Addresses.TertiaryPickupAddress)
			if result != nil {
				tempTertiaryPickupAddressCustoms = convertCustomizationInList(tempTertiaryPickupAddressCustoms, Addresses.TertiaryPickupAddress, Address)
				tertiaryPickupAddress := BuildAddress(db, tempTertiaryPickupAddressCustoms, traits)

				newMTOShipment.TertiaryPickupAddress = &tertiaryPickupAddress
				newMTOShipment.TertiaryPickupAddressID = &tertiaryPickupAddress.ID
				newMTOShipment.HasTertiaryPickupAddress = models.BoolPointer(true)
			}
		}

		if shipmentHasDeliveryDetails {
			newMTOShipment.RequestedDeliveryDate = models.TimePointer(time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC))
			newMTOShipment.ScheduledDeliveryDate = models.TimePointer(time.Date(GHCTestYear, time.March, 17, 0, 0, 0, 0, time.UTC))

			// Find/create the Delivery Address
			tempDeliveryAddressCustoms := customs
			result := findValidCustomization(customs, Addresses.DeliveryAddress)
			if result != nil {
				tempDeliveryAddressCustoms = convertCustomizationInList(tempDeliveryAddressCustoms, Addresses.DeliveryAddress, Address)
			}

			traits = append(traits, GetTraitAddress2)
			deliveryAddress := BuildAddress(db, tempDeliveryAddressCustoms, traits)
			if db == nil {
				// fake an id for stubbed address, needed by the MTOShipmentCreator
				deliveryAddress.ID = uuid.Must(uuid.NewV4())
			}
			newMTOShipment.DestinationAddress = &deliveryAddress
			newMTOShipment.DestinationAddressID = &deliveryAddress.ID

			// Find Secondary Delivery Address
			tempSecondaryDeliveryAddressCustoms := customs
			result = findValidCustomization(customs, Addresses.SecondaryDeliveryAddress)
			if result != nil {
				tempSecondaryDeliveryAddressCustoms = convertCustomizationInList(tempSecondaryDeliveryAddressCustoms, Addresses.SecondaryDeliveryAddress, Address)
				secondaryDeliveryAddress := BuildAddress(db, tempSecondaryDeliveryAddressCustoms, traits)

				newMTOShipment.SecondaryDeliveryAddress = &secondaryDeliveryAddress
				newMTOShipment.SecondaryDeliveryAddressID = &secondaryDeliveryAddress.ID
				newMTOShipment.HasSecondaryDeliveryAddress = models.BoolPointer(true)

			}

			// Find Tertiary Delivery Address
			tempTertiaryDeliveryAddressCustoms := customs
			result = findValidCustomization(customs, Addresses.TertiaryDeliveryAddress)
			if result != nil {
				tempTertiaryDeliveryAddressCustoms = convertCustomizationInList(tempTertiaryDeliveryAddressCustoms, Addresses.TertiaryDeliveryAddress, Address)
				tertiaryDeliveryAddress := BuildAddress(db, tempTertiaryDeliveryAddressCustoms, traits)

				newMTOShipment.TertiaryDeliveryAddress = &tertiaryDeliveryAddress
				newMTOShipment.TertiaryDeliveryAddressID = &tertiaryDeliveryAddress.ID
				newMTOShipment.HasTertiaryDeliveryAddress = models.BoolPointer(true)
			}
		}

		if cMtoShipment.ScheduledPickupDate != nil {
			requiredDeliveryDate := time.Date(GHCTestYear, time.April, 15, 0, 0, 0, 0, time.UTC)
			newMTOShipment.RequiredDeliveryDate = &requiredDeliveryDate
		}
	}

	testdatagen.MergeModels(&newMTOShipment, cMtoShipment)

	if db != nil {
		mustCreate(db, &newMTOShipment)
	}

	return newMTOShipment
}

// BuildBaseMTOShipment creates a single MTOShipment with the base set of data required for a shipment.
func BuildBaseMTOShipment(db *pop.Connection, customs []Customization, traits []Trait) models.MTOShipment {
	return buildMTOShipmentWithBuildType(db, customs, traits, mtoShipmentBuildBasic)
}

// BuildMTOShipment creates a single MTOShipment and associated set relationships
// It will make a move record, if one is not provided.
// It will make pickup addresses if the shipment type is not one of (HHGOutOfNTSDom, PPM)
// It will make delivery addresses if the shipment type is not one of (HHGIntoNTSDom, PPM)
// It will make a storage facility if the shipment type is HHGOutOfNTSDom
func BuildMTOShipment(db *pop.Connection, customs []Customization, traits []Trait) models.MTOShipment {
	return buildMTOShipmentWithBuildType(db, customs, traits, mtoShipmentBuild)
}

// BuildMTOShipmentMinimal creates a single MTOShipment with a minimal set of data as could be possible through the UI
// for any shipment that doesn't have a child table associated with the MTOShipment model. It does not create associated
// addresses.
func BuildMTOShipmentMinimal(db *pop.Connection, customs []Customization, traits []Trait) models.MTOShipment {
	mtoShipment := BuildBaseMTOShipment(db, customs, traits)

	// Get shipment_locator from DB that was generated from shipment INSERT.
	if db != nil {
		var dbMtoShipment models.MTOShipment
		err := db.Find(&dbMtoShipment, mtoShipment.ID)
		if err == nil {
			mtoShipment.ShipmentLocator = dbMtoShipment.ShipmentLocator
		}
	}

	customs = setupCustomizations(customs, traits)

	// Find pickup address in case it was added to customizations list
	tempPickupAddressCustoms := customs
	result := findValidCustomization(customs, Addresses.PickupAddress)
	if result != nil {
		tempPickupAddressCustoms = convertCustomizationInList(tempPickupAddressCustoms, Addresses.PickupAddress, Address)
		pickupAddress := BuildAddress(db, tempPickupAddressCustoms, traits)
		if db == nil {
			// fake an id for stubbed address, needed by the MTOShipmentCreator
			pickupAddress.ID = uuid.Must(uuid.NewV4())
		}
		mtoShipment.PickupAddress = &pickupAddress
		mtoShipment.PickupAddressID = &pickupAddress.ID

		if db != nil {
			mustSave(db, &mtoShipment)
		}
	}

	// Find secondary pickup address in case it was added to customizations list
	tempSecondaryPickupAddressCustoms := customs
	result = findValidCustomization(customs, Addresses.SecondaryPickupAddress)
	if result != nil {
		tempSecondaryPickupAddressCustoms = convertCustomizationInList(tempSecondaryPickupAddressCustoms, Addresses.SecondaryPickupAddress, Address)
		secondaryPickupAddress := BuildAddress(db, tempSecondaryPickupAddressCustoms, traits)
		if db == nil {
			// fake an id for stubbed address, needed by the MTOShipmentCreator
			secondaryPickupAddress.ID = uuid.Must(uuid.NewV4())
		}
		mtoShipment.SecondaryPickupAddress = &secondaryPickupAddress
		mtoShipment.SecondaryPickupAddressID = &secondaryPickupAddress.ID
		mtoShipment.HasSecondaryPickupAddress = models.BoolPointer(true)

		if db != nil {
			mustSave(db, &mtoShipment)
		}
	}

	// Find tertiary pickup address in case it was added to customizations list
	tempTertiaryPickupAddressCustoms := customs
	result = findValidCustomization(customs, Addresses.TertiaryPickupAddress)
	if result != nil {
		tempTertiaryPickupAddressCustoms = convertCustomizationInList(tempTertiaryPickupAddressCustoms, Addresses.TertiaryPickupAddress, Address)
		tertiaryPickupAddress := BuildAddress(db, tempTertiaryPickupAddressCustoms, traits)
		if db == nil {
			// fake an id for stubbed address, needed by the MTOShipmentCreator
			tertiaryPickupAddress.ID = uuid.Must(uuid.NewV4())
		}
		mtoShipment.TertiaryPickupAddress = &tertiaryPickupAddress
		mtoShipment.TertiaryPickupAddressID = &tertiaryPickupAddress.ID
		mtoShipment.HasTertiaryPickupAddress = models.BoolPointer(true)

		if db != nil {
			mustSave(db, &mtoShipment)
		}
	}

	// Find destination address in case it was added to customizations list
	tempDestinationAddressCustoms := customs
	result = findValidCustomization(customs, Addresses.DeliveryAddress)
	if result != nil {
		tempDestinationAddressCustoms = convertCustomizationInList(tempDestinationAddressCustoms, Addresses.DeliveryAddress, Address)
		deliveryAddress := BuildAddress(db, tempDestinationAddressCustoms, traits)
		if db == nil {
			// fake an id for stubbed address, needed by the MTOShipmentCreator
			deliveryAddress.ID = uuid.Must(uuid.NewV4())
		}
		mtoShipment.DestinationAddress = &deliveryAddress
		mtoShipment.DestinationAddressID = &deliveryAddress.ID

		if db != nil {
			mustSave(db, &mtoShipment)
		}
	}

	// Find secondary delivery address in case it was added to customizations list
	tempSecondaryDeliveryAddressCustoms := customs
	result = findValidCustomization(customs, Addresses.SecondaryDeliveryAddress)
	if result != nil {
		tempSecondaryDeliveryAddressCustoms = convertCustomizationInList(tempSecondaryDeliveryAddressCustoms, Addresses.SecondaryDeliveryAddress, Address)
		secondaryDeliveryAddress := BuildAddress(db, tempSecondaryDeliveryAddressCustoms, traits)
		if db == nil {
			// fake an id for stubbed address, needed by the MTOShipmentCreator
			secondaryDeliveryAddress.ID = uuid.Must(uuid.NewV4())
		}
		mtoShipment.SecondaryDeliveryAddress = &secondaryDeliveryAddress
		mtoShipment.SecondaryDeliveryAddressID = &secondaryDeliveryAddress.ID
		mtoShipment.HasSecondaryDeliveryAddress = models.BoolPointer(true)

		if db != nil {
			mustSave(db, &mtoShipment)
		}
	}

	// Find tertiary delivery address in case it was added to customizations list
	tempTertiaryDeliveryAddressCustoms := customs
	result = findValidCustomization(customs, Addresses.TertiaryDeliveryAddress)
	if result != nil {
		tempTertiaryDeliveryAddressCustoms = convertCustomizationInList(tempTertiaryDeliveryAddressCustoms, Addresses.TertiaryDeliveryAddress, Address)
		tertiaryDeliveryAddress := BuildAddress(db, tempTertiaryDeliveryAddressCustoms, traits)
		if db == nil {
			// fake an id for stubbed address, needed by the MTOShipmentCreator
			tertiaryDeliveryAddress.ID = uuid.Must(uuid.NewV4())
		}
		mtoShipment.TertiaryDeliveryAddress = &tertiaryDeliveryAddress
		mtoShipment.TertiaryDeliveryAddressID = &tertiaryDeliveryAddress.ID
		mtoShipment.HasTertiaryDeliveryAddress = models.BoolPointer(true)

		if db != nil {
			mustSave(db, &mtoShipment)
		}
	}

	// Find storage facility in case it was added to customizations list
	storageResult := findValidCustomization(customs, StorageFacility)
	if storageResult != nil {
		storageFacility := BuildStorageFacility(db, customs, traits)
		mtoShipment.StorageFacility = &storageFacility
		mtoShipment.StorageFacilityID = &storageFacility.ID

		if db != nil {
			mustSave(db, &mtoShipment)
		}
	}

	if mtoShipment.RequestedPickupDate == nil {
		requestedPickupDate := time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC)

		mtoShipment.RequestedPickupDate = &requestedPickupDate

		if db != nil {
			mustSave(db, &mtoShipment)
		}
	}

	mtoShipment.MarketCode = models.MarketCodeDomestic
	if db != nil {
		mustSave(db, &mtoShipment)
	}

	return mtoShipment
}

func BuildMTOShipmentWithMove(move *models.Move, db *pop.Connection, customs []Customization, traits []Trait) models.MTOShipment {
	customs = setupCustomizations(customs, traits)

	// Cannot provide move customization to this Build
	if result := findValidCustomization(customs, Move); result != nil {
		log.Panicf("Cannot provide Move customization to BuildMTOShipmentWithMove")
	}

	// provide linkonly customization for the provided move
	customs = append(customs, Customization{
		Model:    *move,
		LinkOnly: true,
	})

	shipment := BuildMTOShipment(db, customs, traits)

	move.MTOShipments = append(move.MTOShipments, shipment)

	return shipment

}

func BuildNTSShipment(db *pop.Connection, customs []Customization, traits []Trait) models.MTOShipment {
	// add secondary if not already customized
	secondaryAddressResult := findValidCustomization(customs, Addresses.SecondaryPickupAddress)
	if secondaryAddressResult == nil {
		// we already know customs do not apply
		secondaryAddress := BuildAddress(db, nil, traits)
		customs = append(customs, Customization{
			Model:    secondaryAddress,
			LinkOnly: true,
			Type:     &Addresses.SecondaryPickupAddress,
		})
	}
	tertiaryAddressResult := findValidCustomization(customs, Addresses.TertiaryPickupAddress)
	if tertiaryAddressResult == nil {
		// we already know customs do not apply
		tertiaryAddress := BuildAddress(db, nil, traits)
		customs = append(customs, Customization{
			Model:    tertiaryAddress,
			LinkOnly: true,
			Type:     &Addresses.TertiaryPickupAddress,
		})
	}

	return buildMTOShipmentWithBuildType(db, customs, traits, mtoShipmentNTS)
}

func BuildNTSRShipment(db *pop.Connection, customs []Customization, traits []Trait) models.MTOShipment {
	// add secondary if not already customized
	secondaryAddressResult := findValidCustomization(customs, Addresses.SecondaryDeliveryAddress)
	if secondaryAddressResult == nil {
		// we already know customs do not apply
		secondaryAddress := BuildAddress(db, nil, traits)
		customs = append(customs, Customization{
			Model:    secondaryAddress,
			LinkOnly: true,
			Type:     &Addresses.SecondaryDeliveryAddress,
		})
	}
	tertiaryAddressResult := findValidCustomization(customs, Addresses.TertiaryDeliveryAddress)
	if tertiaryAddressResult == nil {
		// we already know customs do not apply
		tertiaryAddress := BuildAddress(db, nil, traits)
		customs = append(customs, Customization{
			Model:    tertiaryAddress,
			LinkOnly: true,
			Type:     &Addresses.TertiaryDeliveryAddress,
		})
	}
	return buildMTOShipmentWithBuildType(db, customs, traits, mtoShipmentNTSR)
}
func AddPPMShipmentToMTOShipment(db *pop.Connection, mtoShipment *models.MTOShipment, ppmShipment models.PPMShipment) {
	if mtoShipment.ShipmentType != models.MTOShipmentTypePPM {
		log.Panic("mtoShipmentType must be MTOShipmentTypePPM")
	}
	if db == nil && ppmShipment.ID.IsNil() {
		// need to create an ID so we can use the ppmShipment as
		// LinkOnly
		ppmShipment.ID = uuid.Must(uuid.NewV4())
	}
	mtoShipment.PPMShipment = &ppmShipment
}

// ------------------------
//
//	TRAITS
//
// ------------------------
func GetTraitSubmittedShipment() []Customization {
	return []Customization{
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		},
	}
}
