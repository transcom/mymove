package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMoveOrder creates a single MoveOrder and associated set relationships
func MakeMoveOrder(db *pop.Connection, assertions Assertions) models.MoveOrder {
	var customer models.Customer
	if isZeroUUID(assertions.Customer.ID) {
		customer = MakeCustomer(db, assertions)
	}
	var entitlement models.Entitlement
	if isZeroUUID(assertions.Entitlement.ID) {
		entitlement = MakeEntitlement(db, assertions)
	}
	var originDutyStation models.DutyStation
	if isZeroUUID(assertions.OriginDutyStation.ID) {
		originDutyStation = MakeDutyStation(db, assertions)
	}
	var destinationDutyStation models.DutyStation
	if isZeroUUID(assertions.DestinationDutyStation.ID) {
		destinationDutyStation = MakeDutyStation(db, assertions)
	}
	moveOrder := models.MoveOrder{
		Customer:                 customer,
		CustomerID:               customer.ID,
		Entitlement:              entitlement,
		EntitlementID:            entitlement.ID,
		DestinationDutyStation:   destinationDutyStation,
		DestinationDutyStationID: destinationDutyStation.ID,
		OriginDutyStation:        originDutyStation,
		OriginDutyStationID:      originDutyStation.ID,
	}

	// Overwrite values with those from assertions
	mergeModels(&moveOrder, assertions.MoveOrder)

	mustCreate(db, &moveOrder)

	return moveOrder
}

// MakeDefaultMoveOrder makes a MoveOrder with default values
func MakeDefaultMoveOrder(db *pop.Connection) models.MoveOrder {
	return MakeMoveOrder(db, Assertions{})
}
