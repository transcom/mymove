package models

import (
	"fmt"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
)

// ConvertFromPPMToGHC creates models in the new GHC data model from data in a PPM move
func ConvertFromPPMToGHC(db *pop.Connection, moveID uuid.UUID) (uuid.UUID, error) {
	var move Move
	if err := db.Eager("Orders.ServiceMember").Find(&move, moveID); err != nil {
		return uuid.Nil, fmt.Errorf("Could not fetch move with id %s, %w", moveID, err)
	}

	// service member -> customer
	sm := move.Orders.ServiceMember
	var customer Customer
	customer.CreatedAt = sm.CreatedAt
	customer.UpdatedAt = sm.UpdatedAt
	customer.DODID = *sm.Edipi
	customer.UserID = sm.UserID
	customer.FirstName = *sm.FirstName
	customer.LastName = *sm.LastName
	customer.Email = sm.PersonalEmail
	customer.PhoneNumber = sm.Telephone
	customer.Agency = string(*sm.Affiliation)
	customer.CurrentAddressID = sm.ResidentialAddressID

	if err := db.Save(&customer); err != nil {
		return uuid.Nil, fmt.Errorf("Could not save customer, %w", err)
	}

	// create entitlement (required by move order)
	weight, entitlementErr := GetEntitlement(*sm.Rank, move.Orders.HasDependents, move.Orders.SpouseHasProGear)
	if entitlementErr != nil {
		return uuid.Nil, entitlementErr
	}
	entitlement := Entitlement{
		DependentsAuthorized: &move.Orders.HasDependents,
		DBAuthorizedWeight:   IntPointer(weight),
	}

	if err := db.Save(&entitlement); err != nil {
		return uuid.Nil, fmt.Errorf("Could not save entitlement, %w", err)
	}

	// orders -> move order
	orders := move.Orders
	var mo MoveOrder
	mo.CreatedAt = orders.CreatedAt
	mo.UpdatedAt = orders.UpdatedAt
	mo.Customer = &customer
	mo.CustomerID = &customer.ID
	mo.DestinationDutyStation = &orders.NewDutyStation
	mo.DestinationDutyStationID = &orders.NewDutyStationID

	orderType := "GHC"
	mo.OrderType = &orderType
	orderTypeDetail := "TBD"
	mo.OrderTypeDetail = &orderTypeDetail
	mo.OriginDutyStation = &sm.DutyStation
	mo.OriginDutyStationID = sm.DutyStationID
	mo.Entitlement = &entitlement
	mo.EntitlementID = &entitlement.ID
	mo.Grade = (*string)(sm.Rank)
	mo.DateIssued = &orders.IssueDate
	mo.ReportByDate = &orders.ReportByDate

	if err := db.Save(&mo); err != nil {
		return uuid.Nil, fmt.Errorf("Could not save move order, %w", err)
	}

	return mo.ID, nil
}
