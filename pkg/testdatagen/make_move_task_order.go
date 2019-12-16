package testdatagen

import (
	"github.com/gobuffalo/pop"
	"time"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMoveTaskOrder creates a single MoveTaskOrder and associated set relationships
func MakeMoveTaskOrder(db *pop.Connection, assertions Assertions) models.MoveTaskOrder {
	moveOrder := assertions.MoveOrder
	customer := assertions.Customer
	if isZeroUUID(moveOrder.ID) {
		moveOrder = MakeMoveOrder(db, assertions)
	}
	var referenceID *string

	if isZeroUUID(customer.ID) {
		customer = MakeCustomer(db, assertions)
	}

	pickupAddress := assertions.MoveTaskOrder.PickupAddress
	if isZeroUUID(pickupAddress.ID) {
		pickupAddress = MakeAddress(db, assertions)
	}

	destinationAddress := assertions.MoveTaskOrder.DestinationAddress
	if isZeroUUID(destinationAddress.ID) {
		destinationAddress = MakeAddress2(db, assertions)
	}

	moveTaskOrder := models.MoveTaskOrder{
		MoveOrder:          moveOrder,
		MoveOrderID:        moveOrder.ID,
		ReferenceID:        referenceID,
		IsAvailableToPrime: false,
		IsCanceled:         false,
		CustomerID:         customer.ID,
		Customer:                 customer,
		OriginDutyStationID:      moveOrder.OriginDutyStationID,
		OriginDutyStation:        moveOrder.OriginDutyStation,
		DestinationDutyStation:   moveOrder.DestinationDutyStation,
		DestinationDutyStationID: moveOrder.DestinationDutyStation.ID,
		PickupAddress:            pickupAddress,
		PickupAddressID:          pickupAddress.ID,
		DestinationAddress:       destinationAddress,
		DestinationAddressID:     destinationAddress.ID,
		RequestedPickupDate:      time.Date(TestYear, time.March, 15, 0, 0, 0, 0, time.UTC),
	}

	// Overwrite values with those from assertions
	mergeModels(&moveTaskOrder, assertions.MoveTaskOrder)

	mustCreate(db, &moveTaskOrder)

	return moveTaskOrder
}

// MakeDefaultMoveTaskOrder makes an MoveTaskOrder with default values
func MakeDefaultMoveTaskOrder(db *pop.Connection) models.MoveTaskOrder {
	return MakeMoveTaskOrder(db, Assertions{})
}
