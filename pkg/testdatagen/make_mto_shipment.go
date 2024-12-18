package testdatagen

import (
	"log"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// makeMTOShipment creates a single MTOShipment and associated set relationships
// It will make a move record, if one is not provided.
// It will make pickup addresses if the shipment type is not one of (HHGOutOfNTSDom, PPM)
// It will make delivery addresses if the shipment type is not one of (HHGOutOfNTSDom, PPM)
// It will make a storage facility if the shipment type is
// HHGOutOfNTSDom
//
// Deprecated: use factory.BuildMTOShipment
func makeMTOShipment(db *pop.Connection, assertions Assertions) models.MTOShipment {
	shipmentType := models.MTOShipmentTypeHHG
	shipmentStatus := models.MTOShipmentStatusDraft
	mtoShipment := assertions.MTOShipment

	// Make move if it was not provided
	moveTaskOrder := assertions.Move
	if isZeroUUID(moveTaskOrder.ID) {
		moveTaskOrder = makeMove(db, assertions)
	}

	if mtoShipment.ShipmentType != "" {
		shipmentType = mtoShipment.ShipmentType
	}

	if mtoShipment.Status != "" {
		shipmentStatus = mtoShipment.Status
	}

	shipmentHasPickupDetails := mtoShipment.ShipmentType != models.MTOShipmentTypeHHGOutOfNTSDom && mtoShipment.ShipmentType != models.MTOShipmentTypePPM
	shipmentHasDeliveryDetails := mtoShipment.ShipmentType != models.MTOShipmentTypeHHGIntoNTS && mtoShipment.ShipmentType != models.MTOShipmentTypePPM

	var pickupAddress, secondaryPickupAddress models.Address
	if shipmentHasPickupDetails {
		// Make pickup address if it was not provided
		pickupAddress = assertions.PickupAddress
		if isZeroUUID(pickupAddress.ID) {
			pickupAddress = MakeAddress(db, Assertions{
				Address: assertions.PickupAddress,
			})
		}

		secondaryPickupAddress = assertions.SecondaryPickupAddress

		// Check that a GBLOC exists for pickup address postal code, make one if not
		gbloc, err := models.FetchGBLOCForPostalCode(db, pickupAddress.PostalCode)
		if gbloc.GBLOC == "" || err != nil {
			makePostalCodeToGBLOC(db, pickupAddress.PostalCode, "KKFA")
		}
	}

	var destinationAddress, secondaryDeliveryAddress, tertiaryDeliveryAddress models.Address
	if shipmentHasDeliveryDetails {
		// Make destination address if it was not provided
		destinationAddress = assertions.DestinationAddress
		if isZeroUUID(destinationAddress.ID) {
			destinationAddress = MakeAddress2(db, Assertions{
				Address: assertions.DestinationAddress,
			})
		}

		secondaryDeliveryAddress = assertions.SecondaryDeliveryAddress
		tertiaryDeliveryAddress = assertions.TertiaryDeliveryAddress
	}

	// mock weights
	var estimatedWeight *unit.Pound
	if assertions.MTOShipment.PrimeEstimatedWeight != nil {
		estimatedWeight = assertions.MTOShipment.PrimeEstimatedWeight
	}
	actualWeight := unit.Pound(980)

	// mock dates -- set to nil initially since these are all nullable
	var requestedPickupDate, scheduledPickupDate, actualPickupDate, requestedDeliveryDate *time.Time

	if shipmentHasPickupDetails {
		requestedPickupDate = models.TimePointer(time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC))
		scheduledPickupDate = models.TimePointer(time.Date(GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC))
		actualPickupDate = models.TimePointer(time.Date(GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC))
	}
	if shipmentHasDeliveryDetails {
		requestedDeliveryDate = models.TimePointer(time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC))
	}

	var storageFacilityID *uuid.UUID
	var storageFacility models.StorageFacility
	if mtoShipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom ||
		mtoShipment.ShipmentType == models.MTOShipmentTypeHHGIntoNTS {
		if mtoShipment.StorageFacility != nil {
			if isZeroUUID(mtoShipment.StorageFacility.ID) {
				storageFacility = MakeStorageFacility(db, Assertions{
					StorageFacility: *mtoShipment.StorageFacility,
				})
				storageFacilityID = &storageFacility.ID
			} else {
				storageFacilityID = &mtoShipment.StorageFacility.ID
				err := db.Eager().Find(&storageFacility, storageFacilityID)
				if err != nil {
					log.Panic(err)
				}
			}
		} else if !isZeroUUID(assertions.StorageFacility.ID) {
			storageFacilityID = &assertions.StorageFacility.ID
			err := db.Eager().Find(&storageFacility, storageFacilityID)
			if err != nil {
				log.Panic(err)
			}
		} else {
			storageFacility = MakeDefaultStorageFacility(db)
			storageFacilityID = &storageFacility.ID
		}
	}

	if uuid.Nil != storageFacility.AddressID {
		err := db.Find(&storageFacility.Address, storageFacility.AddressID)
		if err != nil {
			log.Panic(err)
		}
	}

	var storageFacilityPtr *models.StorageFacility
	if storageFacilityID != nil {
		storageFacilityPtr = &storageFacility
	}

	MTOShipment := models.MTOShipment{
		MoveTaskOrder:         moveTaskOrder,
		MoveTaskOrderID:       moveTaskOrder.ID,
		RequestedPickupDate:   requestedPickupDate,
		ScheduledPickupDate:   scheduledPickupDate,
		ActualPickupDate:      actualPickupDate,
		RequestedDeliveryDate: requestedDeliveryDate,
		CustomerRemarks:       models.StringPointer("Please treat gently"),
		PrimeEstimatedWeight:  estimatedWeight,
		PrimeActualWeight:     &actualWeight,
		ShipmentType:          shipmentType,
		Status:                shipmentStatus,
		StorageFacilityID:     storageFacilityID,
		StorageFacility:       storageFacilityPtr,
		MarketCode:            models.MarketCodeDomestic,
	}

	if assertions.MTOShipment.DestinationType != nil {
		MTOShipment.DestinationType = assertions.MTOShipment.DestinationType
	}

	if shipmentHasDeliveryDetails {
		MTOShipment.DestinationAddress = &destinationAddress
		MTOShipment.DestinationAddressID = &destinationAddress.ID

		if !isZeroUUID(secondaryDeliveryAddress.ID) {
			MTOShipment.SecondaryDeliveryAddress = &secondaryDeliveryAddress
			MTOShipment.SecondaryDeliveryAddressID = &secondaryDeliveryAddress.ID
			MTOShipment.HasSecondaryDeliveryAddress = models.BoolPointer(true)
		}

		if !isZeroUUID(tertiaryDeliveryAddress.ID) {
			MTOShipment.TertiaryDeliveryAddress = &tertiaryDeliveryAddress
			MTOShipment.TertiaryDeliveryAddressID = &tertiaryDeliveryAddress.ID
			MTOShipment.HasTertiaryDeliveryAddress = models.BoolPointer(true)
		}
	}

	if shipmentHasPickupDetails {
		MTOShipment.PickupAddress = &pickupAddress
		MTOShipment.PickupAddressID = &pickupAddress.ID

		if !isZeroUUID(secondaryPickupAddress.ID) {
			MTOShipment.SecondaryPickupAddress = &secondaryPickupAddress
			MTOShipment.SecondaryPickupAddressID = &secondaryPickupAddress.ID
			MTOShipment.HasSecondaryPickupAddress = models.BoolPointer(true)
		}
	}

	if assertions.MTOShipment.Status == models.MTOShipmentStatusApproved {
		approvedDate := time.Date(GHCTestYear, time.March, 20, 0, 0, 0, 0, time.UTC)
		MTOShipment.ApprovedDate = &approvedDate
	}

	if assertions.MTOShipment.ScheduledPickupDate != nil {
		requiredDeliveryDate := time.Date(GHCTestYear, time.April, 15, 0, 0, 0, 0, time.UTC)
		MTOShipment.RequiredDeliveryDate = &requiredDeliveryDate
	}

	// Overwrite values with those from assertions
	mergeModels(&MTOShipment, assertions.MTOShipment)

	mustCreate(db, &MTOShipment, assertions.Stub)
	return MTOShipment
}
