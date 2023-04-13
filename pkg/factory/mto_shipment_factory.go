package factory

import (
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
	shipmentType := models.MTOShipmentTypeHHG

	newMTOShipment := models.MTOShipment{
		MoveTaskOrder:   move,
		MoveTaskOrderID: move.ID,
		ShipmentType:    shipmentType,
		Status:          models.MTOShipmentStatusSubmitted,
	}

	if buildType == mtoShipmentBuild {
		newMTOShipment.Status = models.MTOShipmentStatusDraft

		if cMtoShipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom || cMtoShipment.ShipmentType == models.MTOShipmentTypeHHGIntoNTSDom {
			storageFacility := BuildStorageFacility(db, customs, traits)
			// only set storage facility pointers if building a
			// storage facility
			newMTOShipment.StorageFacility = &storageFacility
			newMTOShipment.StorageFacilityID = &storageFacility.ID
		}

		actualWeight := unit.Pound(980)
		newMTOShipment.PrimeActualWeight = &actualWeight
		newMTOShipment.CustomerRemarks = models.StringPointer("Please treat gently")

		shipmentHasPickupDetails := cMtoShipment.ShipmentType != models.MTOShipmentTypeHHGOutOfNTSDom && cMtoShipment.ShipmentType != models.MTOShipmentTypePPM
		shipmentHasDeliveryDetails := cMtoShipment.ShipmentType != models.MTOShipmentTypeHHGIntoNTSDom && cMtoShipment.ShipmentType != models.MTOShipmentTypePPM

		if shipmentHasPickupDetails {
			newMTOShipment.RequestedPickupDate = models.TimePointer(time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC))
			newMTOShipment.ScheduledPickupDate = models.TimePointer(time.Date(GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC))
			// if cMtoShipment.Status != "" && cMtoShipment.Status != models.MTOShipmentStatusDraft {
			newMTOShipment.ActualPickupDate = models.TimePointer(time.Date(GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC))
			// }
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
			if db == nil {
				BuildPostalCodeToGBLOC(nil, []Customization{
					{
						Model: models.PostalCodeToGBLOC{
							PostalCode: pickupAddress.PostalCode,
							GBLOC:      "KKFA",
						},
					},
				}, nil)
			} else {
				gbloc, err := models.FetchGBLOCForPostalCode(db, pickupAddress.PostalCode)
				if gbloc.GBLOC == "" || err != nil {
					FetchOrBuildPostalCodeToGBLOC(db, pickupAddress.PostalCode, "KKFA")
				}
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
		}

		if shipmentHasDeliveryDetails {
			newMTOShipment.RequestedDeliveryDate = models.TimePointer(time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC))

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
		}

		if cMtoShipment.Status == models.MTOShipmentStatusApproved {
			approvedDate := time.Date(GHCTestYear, time.March, 20, 0, 0, 0, 0, time.UTC)
			newMTOShipment.ApprovedDate = &approvedDate
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

	if mtoShipment.RequestedPickupDate == nil {
		requestedPickupDate := time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC)

		mtoShipment.RequestedPickupDate = &requestedPickupDate

		if db != nil {
			mustSave(db, &mtoShipment)
		}
	}

	return mtoShipment
}
