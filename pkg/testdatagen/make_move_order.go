package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMoveOrder creates a single MoveOrder and associated set relationships
func MakeMoveOrder(db *pop.Connection, assertions Assertions) models.MoveOrder {
	customer := assertions.Customer
	if isZeroUUID(customer.ID) {
		customer = MakeCustomer(db, assertions)
	}
	entitlement := assertions.Entitlement
	if isZeroUUID(entitlement.ID) {
		entitlement = MakeEntitlement(db, assertions)
	}
	originDutyStation := assertions.OriginDutyStation
	if isZeroUUID(originDutyStation.ID) {
		originDutyStation = MakeDutyStation(db, Assertions{
			DutyStation: models.DutyStation{
				Name: "Alamo",
			},
		})
	}
	destinationDutyStation := assertions.DestinationDutyStation
	if isZeroUUID(destinationDutyStation.ID) {
		destinationDutyStation = MakeDutyStation(db, Assertions{
			DutyStation: models.DutyStation{
				Name: "Versailles",
			},
		})
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
