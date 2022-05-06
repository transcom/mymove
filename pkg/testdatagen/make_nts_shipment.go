package testdatagen

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakeNTSShipment creates a single MTOShipment of type NTS and associated set relationships
func MakeNTSShipment(db *pop.Connection, assertions Assertions) models.MTOShipment {

	// Make move if it was not provided
	moveTaskOrder := assertions.Move
	if isZeroUUID(moveTaskOrder.ID) {
		moveTaskOrder = MakeMove(db, assertions)
	}

	// Make pickup address if it was not provided
	pickupAddress := assertions.PickupAddress
	if isZeroUUID(pickupAddress.ID) {
		pickupAddress = MakeAddress(db, Assertions{
			Address: assertions.PickupAddress,
		})
	}

	// Make secondary pickup address if it was not provided
	secondaryPickupAddress := assertions.SecondaryPickupAddress
	if isZeroUUID(secondaryPickupAddress.ID) {
		secondaryPickupAddress = MakeAddress(db, Assertions{
			Address: assertions.SecondaryPickupAddress,
		})
	}

	// mock dates
	requestedPickupDate := time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC)
	// TODO: add releasing agent

	MTOShipment := models.MTOShipment{
		MoveTaskOrder:          moveTaskOrder,
		MoveTaskOrderID:        moveTaskOrder.ID,
		RequestedPickupDate:    &requestedPickupDate,
		CustomerRemarks:        swag.String("Please treat gently"),
		PickupAddress:          &pickupAddress,
		PickupAddressID:        &pickupAddress.ID,
		SecondaryPickupAddress: &secondaryPickupAddress,
		ShipmentType:           models.MTOShipmentTypeHHGIntoNTSDom,
		Status:                 "DRAFT",
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

// MakeNTSRShipment creates a single MTOShipment of type NTS-release and associated set relationships
func MakeNTSRShipment(db *pop.Connection, assertions Assertions) models.MTOShipment {

	// Make move if it was not provided
	moveTaskOrder := assertions.Move
	if isZeroUUID(moveTaskOrder.ID) {
		moveTaskOrder = MakeMove(db, assertions)
	}

	// Make destination address if it was not provided
	destinationAddress := assertions.DestinationAddress
	if isZeroUUID(destinationAddress.ID) {
		destinationAddress = MakeAddress2(db, Assertions{
			Address: assertions.DestinationAddress,
		})
	}

	// Make secondary delivery address if it was not provided
	secondaryDeliveryAddress := assertions.SecondaryDeliveryAddress
	if isZeroUUID(secondaryDeliveryAddress.ID) {
		secondaryDeliveryAddress = MakeAddress(db, Assertions{
			Address: assertions.SecondaryDeliveryAddress,
		})
	}

	// mock dates
	requestedDeliveryDate := time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC)
	// TODO: add receiving agent
	MTOShipment := models.MTOShipment{
		MoveTaskOrder:            moveTaskOrder,
		MoveTaskOrderID:          moveTaskOrder.ID,
		RequestedDeliveryDate:    &requestedDeliveryDate,
		CustomerRemarks:          swag.String("Please treat gently"),
		DestinationAddress:       &destinationAddress,
		DestinationAddressID:     &destinationAddress.ID,
		SecondaryDeliveryAddress: &secondaryDeliveryAddress,
		ShipmentType:             models.MTOShipmentTypeHHGOutOfNTSDom,
		Status:                   "DRAFT",
	}

	if assertions.MTOShipment.Status == models.MTOShipmentStatusApproved {
		approvedDate := time.Date(GHCTestYear, time.March, 20, 0, 0, 0, 0, time.UTC)
		MTOShipment.ApprovedDate = &approvedDate
	}

	// Overwrite values with those from assertions
	mergeModels(&MTOShipment, assertions.MTOShipment)

	mustCreate(db, &MTOShipment, assertions.Stub)

	return MTOShipment
}

// MakeDefaultNTSShipment makes an MTOShipment of type NTS with default values
func MakeDefaultNTSShipment(db *pop.Connection) models.MTOShipment {
	return MakeNTSShipment(db, Assertions{})
}

// MakeDefaultNTSRShipment makes an MTOShipment of type NTS with default values
func MakeDefaultNTSRShipment(db *pop.Connection) models.MTOShipment {
	return MakeNTSRShipment(db, Assertions{})
}
