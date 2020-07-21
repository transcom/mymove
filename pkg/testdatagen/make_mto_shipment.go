package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// MakeMTOShipment creates a single MTOShipment and associated set relationships
func MakeMTOShipment(db *pop.Connection, assertions Assertions) models.MTOShipment {
	moveTaskOrder := assertions.MoveTaskOrder
	if isZeroUUID(moveTaskOrder.ID) {
		moveTaskOrder = MakeMoveTaskOrder(db, assertions)
	}

	pickupAddress := MakeAddress(db, assertions)
	destinationAddress := MakeAddress2(db, assertions)
	secondaryPickupAddress := MakeAddress(db, assertions)
	secondaryDeliveryAddress := MakeAddress(db, assertions)
	shipmentType := models.MTOShipmentTypeHHG

	if assertions.MTOShipment.ShipmentType != "" {
		shipmentType = assertions.MTOShipment.ShipmentType
	}

	// mock remarks
	remarks := "please treat gently"
	rejectionReason := "shipment not good enough"

	// mock weights
	actualWeight := unit.Pound(980)

	// mock dates
	requestedPickupDate := time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC)
	scheduledPickupDate := time.Date(GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)
	actualPickupDate := time.Date(GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)

	MTOShipment := models.MTOShipment{
		MoveTaskOrder:            moveTaskOrder,
		MoveTaskOrderID:          moveTaskOrder.ID,
		RequestedPickupDate:      &requestedPickupDate,
		ScheduledPickupDate:      &scheduledPickupDate,
		ActualPickupDate: 		  &actualPickupDate,
		CustomerRemarks:          &remarks,
		PickupAddress:            &pickupAddress,
		PickupAddressID:          &pickupAddress.ID,
		DestinationAddress:       &destinationAddress,
		DestinationAddressID:     &destinationAddress.ID,
		PrimeActualWeight:        &actualWeight,
		SecondaryPickupAddress:   &secondaryPickupAddress,
		SecondaryDeliveryAddress: &secondaryDeliveryAddress,
		ShipmentType:             shipmentType,
		Status:                   "SUBMITTED",
		RejectionReason:          &rejectionReason,
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

	mustCreate(db, &MTOShipment)

	return MTOShipment
}

// MakeMTOShipmentMinimal creates a single MTOShipment with a minimal set of data as could be possible
// through milmove UI.
func MakeMTOShipmentMinimal(db *pop.Connection, assertions Assertions) models.MTOShipment {
	moveTaskOrder := assertions.MoveTaskOrder
	if isZeroUUID(moveTaskOrder.ID) {
		moveTaskOrder = MakeMoveTaskOrder(db, assertions)
	}
	pickupAddress := MakeAddress(db, assertions)
	destinationAddress := MakeAddress2(db, assertions)
	shipmentType := models.MTOShipmentTypeHHG

	// mock dates
	requestedPickupDate := time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC)

	MTOShipment := models.MTOShipment{
		MoveTaskOrder:        moveTaskOrder,
		MoveTaskOrderID:      moveTaskOrder.ID,
		RequestedPickupDate:  &requestedPickupDate,
		PickupAddress:        &pickupAddress,
		PickupAddressID:      &pickupAddress.ID,
		DestinationAddress:   &destinationAddress,
		DestinationAddressID: &destinationAddress.ID,
		ShipmentType:         shipmentType,
		Status:               "SUBMITTED",
	}

	if assertions.MTOShipment.Status == models.MTOShipmentStatusApproved {
		approvedDate := time.Date(GHCTestYear, time.March, 20, 0, 0, 0, 0, time.UTC)
		MTOShipment.ApprovedDate = &approvedDate
	}

	// Overwrite values with those from assertions
	mergeModels(&MTOShipment, assertions.MTOShipment)

	mustCreate(db, &MTOShipment)

	return MTOShipment
}
