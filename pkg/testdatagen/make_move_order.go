package testdatagen

import (
	"math/rand"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

func MakeGrade() string {
	grades := [28]string{"E_1",
		"E_2",
		"E_3",
		"E_4",
		"E_5",
		"E_6",
		"E_7",
		"E_8",
		"E_9",
		"O_1_ACADEMY_GRADUATE",
		"O_2",
		"O_3",
		"O_4",
		"O_5",
		"O_6",
		"O_7",
		"O_8",
		"O_9",
		"O_10",
		"W_1",
		"W_2",
		"W_3",
		"W_4",
		"W_5",
		"AVIATION_CADET",
		"CIVILIAN_EMPLOYEE",
		"ACADEMY_CADET",
		"MIDSHIPMAN"}
	return grades[rand.Intn(len(grades))]
}

// MakeMoveOrder creates a single MoveOrder and associated set relationships
func MakeMoveOrder(db *pop.Connection, assertions Assertions) models.MoveOrder {
	grade := assertions.MoveOrder.Grade
	if grade == "" {
		grade = MakeGrade()
	}
	customer := assertions.Customer
	if isZeroUUID(customer.ID) {
		customer = MakeCustomer(db, assertions)
	}
	entitlement := assertions.Entitlement
	if isZeroUUID(entitlement.ID) {
		assertions.MoveOrder.Grade = grade
		entitlement = MakeEntitlement(db, assertions)
	}
	originDutyStation := assertions.OriginDutyStation
	if isZeroUUID(originDutyStation.ID) {
		originDutyStation = MakeDutyStation(db, assertions)
	}
	destinationDutyStation := assertions.DestinationDutyStation
	if isZeroUUID(destinationDutyStation.ID) {
		destinationDutyStation = MakeDutyStation(db, assertions)
	}

	moveOrder := models.MoveOrder{
		Customer:                 customer,
		CustomerID:               customer.ID,
		ConfirmationNumber:       models.GenerateLocator(),
		Entitlement:              entitlement,
		EntitlementID:            entitlement.ID,
		DestinationDutyStation:   destinationDutyStation,
		DestinationDutyStationID: destinationDutyStation.ID,
		Grade:                    grade,
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
