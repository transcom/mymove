package testdatagen

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// MakeMTOShipment creates a single MTOShipment and associated set relationships
func MakeMTOShipment(db *pop.Connection, assertions Assertions) models.MTOShipment {
	shipmentType := models.MTOShipmentTypeHHG
	shipmentStatus := models.MTOShipmentStatusDraft
	mtoShipment := assertions.MTOShipment
	counselorRemarks := mtoShipment.CounselorRemarks
	// Make move if it was not provided
	moveTaskOrder := assertions.Move
	if isZeroUUID(moveTaskOrder.ID) {
		moveTaskOrder = MakeMove(db, assertions)
	}

	if mtoShipment.ShipmentType != "" {
		shipmentType = mtoShipment.ShipmentType
	}

	if mtoShipment.Status != "" {
		shipmentStatus = mtoShipment.Status
	}

	if counselorRemarks != nil {
		counselorRemarks = mtoShipment.CounselorRemarks
	}

	shipmentHasPickupDetails := mtoShipment.ShipmentType != models.MTOShipmentTypeHHGOutOfNTSDom
	shipmentHasDeliveryDetails := mtoShipment.ShipmentType != models.MTOShipmentTypeHHGIntoNTSDom

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
	}

	var destinationAddress, secondaryDeliveryAddress models.Address
	if shipmentHasDeliveryDetails {
		// Make destination address if it was not provided
		destinationAddress = assertions.DestinationAddress
		if isZeroUUID(destinationAddress.ID) {
			destinationAddress = MakeAddress2(db, Assertions{
				Address: assertions.DestinationAddress,
			})
		}

		secondaryDeliveryAddress = assertions.SecondaryDeliveryAddress
	}

	// mock weights
	var estimatedWeight *unit.Pound
	if assertions.MTOShipment.PrimeEstimatedWeight != nil {
		estimatedWeight = assertions.MTOShipment.PrimeEstimatedWeight
	}
	actualWeight := unit.Pound(980)

	// mock dates
	var requestedPickupDate, scheduledPickupDate, actualPickupDate, requestedDeliveryDate time.Time

	if shipmentHasPickupDetails {
		requestedPickupDate = time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC)
		scheduledPickupDate = time.Date(GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)
		actualPickupDate = time.Date(GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)
	}
	if shipmentHasDeliveryDetails {
		requestedDeliveryDate = time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC)
	}

	MTOShipment := models.MTOShipment{
		MoveTaskOrder:         moveTaskOrder,
		MoveTaskOrderID:       moveTaskOrder.ID,
		RequestedPickupDate:   &requestedPickupDate,
		ScheduledPickupDate:   &scheduledPickupDate,
		ActualPickupDate:      &actualPickupDate,
		RequestedDeliveryDate: &requestedDeliveryDate,
		CustomerRemarks:       swag.String("Please treat gently"),
		CounselorRemarks:      counselorRemarks,
		PrimeEstimatedWeight:  estimatedWeight,
		PrimeActualWeight:     &actualWeight,
		ShipmentType:          shipmentType,
		Status:                shipmentStatus,
		RejectionReason:       swag.String("Not enough information"),
	}

	if shipmentHasDeliveryDetails {
		MTOShipment.DestinationAddress = &destinationAddress
		MTOShipment.DestinationAddressID = &destinationAddress.ID

		if !isZeroUUID(secondaryDeliveryAddress.ID) {
			MTOShipment.SecondaryDeliveryAddress = &secondaryDeliveryAddress
			MTOShipment.SecondaryDeliveryAddressID = &secondaryDeliveryAddress.ID
		}
	}

	if shipmentHasPickupDetails {
		MTOShipment.PickupAddress = &pickupAddress
		MTOShipment.PickupAddressID = &pickupAddress.ID

		if !isZeroUUID(secondaryPickupAddress.ID) {
			MTOShipment.SecondaryPickupAddress = &secondaryPickupAddress
			MTOShipment.SecondaryPickupAddressID = &secondaryPickupAddress.ID
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

// MakeDefaultMTOShipment makes an MTOShipment with default values
func MakeDefaultMTOShipment(db *pop.Connection) models.MTOShipment {
	return MakeMTOShipment(db, Assertions{})
}

// MakeMTOShipmentMinimal creates a single MTOShipment with a minimal set of data as could be possible
// through milmove UI.
// It does not create associated addresses.
func MakeMTOShipmentMinimal(db *pop.Connection, assertions Assertions) models.MTOShipment {
	moveTaskOrder := assertions.Move
	if isZeroUUID(moveTaskOrder.ID) {
		moveTaskOrder = MakeMove(db, assertions)
	}

	// mock dates
	requestedPickupDate := time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC)

	MTOShipment := models.MTOShipment{
		MoveTaskOrder:       moveTaskOrder,
		MoveTaskOrderID:     moveTaskOrder.ID,
		RequestedPickupDate: &requestedPickupDate,
		ShipmentType:        models.MTOShipmentTypeHHG,
		Status:              "SUBMITTED",
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

// MakeDefaultMTOShipmentMinimal makes a minimal MTOShipment with default values
func MakeDefaultMTOShipmentMinimal(db *pop.Connection) models.MTOShipment {
	return MakeMTOShipmentMinimal(db, Assertions{})
}

// MakeMTOShipmentWithMove makes a shipment connected to a given move and updates the move's MTOShipments array
func MakeMTOShipmentWithMove(db *pop.Connection, move *models.Move, assertions Assertions) models.MTOShipment {
	if move != nil {
		assertions.Move = *move
		assertions.MTOShipment.MoveTaskOrder = *move
		assertions.MTOShipment.MoveTaskOrderID = move.ID
	}
	shipment := MakeMTOShipment(db, assertions)
	if move != nil {
		// This will allow someone to easily create multiple test shipments for one move
		move.MTOShipments = append(move.MTOShipments, shipment)
	}
	return shipment
}

// MakeSubmittedMTOShipmentWithMove makes a shipment with the "SUBMITTED" status and a specific move
func MakeSubmittedMTOShipmentWithMove(db *pop.Connection, move *models.Move, assertions Assertions) models.MTOShipment {
	assertions.MTOShipment.Status = models.MTOShipmentStatusSubmitted
	return MakeMTOShipmentWithMove(db, move, assertions)
}

// MakeStubbedShipment makes a stubbed shipment
func MakeStubbedShipment(db *pop.Connection) models.MTOShipment {
	return MakeMTOShipmentMinimal(db, Assertions{
		MTOShipment: models.MTOShipment{
			ID: uuid.Must(uuid.NewV4()),
		},
		Stub: true,
	})
}
