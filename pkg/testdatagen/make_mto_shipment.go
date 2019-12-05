package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// MakeMTOShipment creates a single MTOShipment and associated set relationships
func MakeMTOShipment(db *pop.Connection, assertions Assertions) models.MTOShipment {
	var moveTaskOrder models.MoveTaskOrder
	if isZeroUUID(assertions.MoveTaskOrder.ID) {
		moveTaskOrder = MakeMoveTaskOrder(db, assertions)
	}
	pickupAddress := assertions.MTOShipment.PickupAddress
	if isZeroUUID(pickupAddress.ID) {
		pickupAddress = MakeAddress(db, assertions)
	}
	destinationAddress := assertions.MTOShipment.DestinationAddress
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

	mtoShipment := models.MTOShipment{
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
	mergeModels(&mtoShipment, assertions.MTOShipment)

	mustCreate(db, &mtoShipment)

	return mtoShipment
}
