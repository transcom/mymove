package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// MakeMTOShipment creates a single MTOShipment and associated set relationships
func MakeMTOShipment(db *pop.Connection, assertions Assertions) models.MTOShipment {
	moveTaskOrder := assertions.Move
	if isZeroUUID(moveTaskOrder.ID) {
		moveTaskOrder = MakeMove(db, assertions)
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
	requestedDeliveryDate := time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC)

	MTOShipment := models.MTOShipment{
		MoveTaskOrder:            moveTaskOrder,
		MoveTaskOrderID:          moveTaskOrder.ID,
		RequestedPickupDate:      &requestedPickupDate,
		ScheduledPickupDate:      &scheduledPickupDate,
		ActualPickupDate:         &actualPickupDate,
		RequestedDeliveryDate:    &requestedDeliveryDate,
		CustomerRemarks:          &remarks,
		PickupAddress:            &pickupAddress,
		PickupAddressID:          &pickupAddress.ID,
		DestinationAddress:       &destinationAddress,
		DestinationAddressID:     &destinationAddress.ID,
		PrimeActualWeight:        &actualWeight,
		SecondaryPickupAddress:   &secondaryPickupAddress,
		SecondaryDeliveryAddress: &secondaryDeliveryAddress,
		ShipmentType:             shipmentType,
		Status:                   "DRAFT",
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

// MakeDefaultMTOShipment makes an MTOShipment with default values
func MakeDefaultMTOShipment(db *pop.Connection) models.MTOShipment {
	return MakeMTOShipment(db, Assertions{})
}

// MakeMTOShipmentMinimal creates a single MTOShipment with a minimal set of data as could be possible
// through milmove UI.
func MakeMTOShipmentMinimal(db *pop.Connection, assertions Assertions) models.MTOShipment {
	moveTaskOrder := assertions.Move
	if isZeroUUID(moveTaskOrder.ID) {
		moveTaskOrder = MakeMove(db, assertions)
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

// MakeDefaultMTOShipmentMinimal makes a minimal MTOShipment with default values
func MakeDefaultMTOShipmentMinimal(db *pop.Connection) models.MTOShipment {
	return MakeMTOShipmentMinimal(db, Assertions{})
}
