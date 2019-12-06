package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// MakeMtoShipment creates a single MtoShipment and associated set relationships
func MakeMtoShipment(db *pop.Connection, assertions Assertions) models.MtoShipment {
	moveTaskOrder := assertions.MoveTaskOrder
	if isZeroUUID(moveTaskOrder.ID) {
		moveTaskOrder = MakeMoveTaskOrder(db, assertions)
	}
	pickupAddress := assertions.MtoShipment.PickupAddress
	if isZeroUUID(pickupAddress.ID) {
		pickupAddress = MakeAddress(db, assertions)
	}
	destinationAddress := assertions.MtoShipment.DestinationAddress
	if isZeroUUID(destinationAddress.ID) {
		destinationAddress = MakeAddress2(db, assertions)
	}

	// mock remarks
	remarks := "please treat gently"

	// mock weights
	estimatedWeight := unit.Pound(1000)
	actualWeight := unit.Pound(980)

	// mock dates
	scheduledPickupDate := time.Date(TestYear, time.March, 16, 0, 0, 0, 0, time.UTC)
	requestedPickupDate := time.Date(TestYear, time.March, 15, 0, 0, 0, 0, time.UTC)
	primeEstimatedWeightDate := time.Date(TestYear, time.March, 20, 0, 0, 0, 0, time.UTC)

	mtoShipment := models.MtoShipment{
		MoveTaskOrder:                    moveTaskOrder,
		MoveTaskOrderID:                  moveTaskOrder.ID,
		ScheduledPickupDate:              &scheduledPickupDate,
		RequestedPickupDate:              &requestedPickupDate,
		CustomerRemarks:                  &remarks,
		PickupAddress:                    pickupAddress,
		PickupAddressID:                  pickupAddress.ID,
		DestinationAddress:               destinationAddress,
		DestinationAddressID:             destinationAddress.ID,
		PrimeEstimatedWeight:             &estimatedWeight,
		PrimeEstimatedWeightRecordedDate: &primeEstimatedWeightDate,
		PrimeActualWeight:                &actualWeight,
	}
	// Overwrite values with those from assertions
	mergeModels(&mtoShipment, assertions.MtoShipment)

	mustCreate(db, &mtoShipment)

	return mtoShipment
}
