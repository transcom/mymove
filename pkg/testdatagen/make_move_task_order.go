package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMoveTaskOrder creates a single MoveTaskOrder and associated set relationships
func MakeMoveTaskOrder(db *pop.Connection, assertions Assertions) models.MoveTaskOrder {

	// Create new Orders if not provided
	// ID is required because it must be populated for Eager saving to work.
	var move models.Move
	if isZeroUUID(assertions.Move.ID) {
		move = MakeMove(db, assertions)
	}
	sm := assertions.Order.ServiceMember
	if isZeroUUID(sm.ID) {
		sm = move.Orders.ServiceMember
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
		MoveID:                   move.ID,
		CustomerID:               sm.ID,
		Customer:                 sm,
		OriginDutyStationID:      sm.DutyStation.ID,
		OriginDutyStation:        sm.DutyStation,
		DestinationDutyStation:   move.Orders.NewDutyStation,
		DestinationDutyStationID: move.Orders.NewDutyStation.ID,
		PickupAddress:            pickupAddress,
		PickupAddressID:          pickupAddress.ID,
		DestinationAddress:       destinationAddress,
		DestinationAddressID:     destinationAddress.ID,
		RequestedPickupDates:     time.Date(TestYear, time.March, 15, 0, 0, 0, 0, time.UTC),
		CustomerRemarks:          "Park in the alley",
		WeightEntitlement:        1000,
		SitEntitlement:           200,
		POVEntitlement:           false,
		NTSEntitlement:           false,
	}

	// Overwrite values with those from assertions
	mergeModels(&moveTaskOrder, assertions.MoveTaskOrder)

	mustCreate(db, &moveTaskOrder)

	return moveTaskOrder
}
