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
	pickupAddress := assertions.MTOShipment.PickupAddress
	if isZeroUUID(pickupAddress.ID) {
		pickupAddress = MakeAddress(db, assertions)
	}
	destinationAddress := assertions.MTOShipment.DestinationAddress
	if isZeroUUID(destinationAddress.ID) {
		destinationAddress = MakeAddress2(db, assertions)
	}

	shipmentType := models.MTOShipmentTypeHHG
	if assertions.MTOShipment.ShipmentType != "" {
		shipmentType = assertions.MTOShipment.ShipmentType
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

	MTOShipment := models.MTOShipment{
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
		ShipmentType:                     shipmentType,
		Status:                           "SUBMITTED",
	}
	// Overwrite values with those from assertions
	mergeModels(&MTOShipment, assertions.MTOShipment)

	mustCreate(db, &MTOShipment)

	return MTOShipment
}
