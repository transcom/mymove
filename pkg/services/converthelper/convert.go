package converthelper

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/models"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
)

// ConvertFromPPMToGHC creates models in the new GHC data model from data in a PPM move
func ConvertFromPPMToGHC(db *pop.Connection, moveID uuid.UUID) (uuid.UUID, error) {
	var move models.Move
	if err := db.Eager("Orders.ServiceMember.DutyStation.Address", "Orders.NewDutyStation.Address").Find(&move, moveID); err != nil {
		return uuid.Nil, fmt.Errorf("Could not fetch move with id %s, %w", moveID, err)
	}

	sm := move.Orders.ServiceMember

	// create entitlement (required by move order)
	weight, entitlementErr := models.GetEntitlement(*sm.Rank, move.Orders.HasDependents, move.Orders.SpouseHasProGear)
	if entitlementErr != nil {
		return uuid.Nil, entitlementErr
	}
	entitlement := models.Entitlement{
		DependentsAuthorized: &move.Orders.HasDependents,
		DBAuthorizedWeight:   models.IntPointer(weight),
	}

	if err := db.Save(&entitlement); err != nil {
		return uuid.Nil, fmt.Errorf("Could not save entitlement, %w", err)
	}

	// orders -> move order
	var mo models.MoveOrder
	var moveOrders []models.MoveOrder
	err := db.Where("customer_id = $1", sm.ID).All(&moveOrders)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Could not fetch move order, %w", err)
	}
	// if a sm
	if len(moveOrders) == 0 {
		orders := move.Orders
		mo.CreatedAt = orders.CreatedAt
		mo.UpdatedAt = orders.UpdatedAt
		mo.Customer = &sm
		mo.CustomerID = &sm.ID
		mo.DestinationDutyStation = &orders.NewDutyStation
		mo.DestinationDutyStationID = &orders.NewDutyStationID

		orderType := "GHC"
		mo.OrderNumber = orders.OrdersNumber
		mo.OrderType = &orderType
		orderTypeDetail := "Shipment of HHG permitted"
		mo.OrderTypeDetail = &orderTypeDetail
		mo.OriginDutyStation = &sm.DutyStation
		mo.OriginDutyStationID = sm.DutyStationID
		mo.Entitlement = &entitlement
		mo.EntitlementID = &entitlement.ID
		mo.Grade = (*string)(sm.Rank)
		mo.DateIssued = &orders.IssueDate
		mo.ReportByDate = &orders.ReportByDate
		mo.LinesOfAccounting = orders.TAC

		if err = db.Save(&mo); err != nil {
			return uuid.Nil, fmt.Errorf("Could not save move order, %w", err)
		}
	} else {
		mo = moveOrders[0]
	}

	var contractor models.Contractor

	err = db.Where("contract_number = ?", "HTC111-11-1-1111").First(&contractor)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Could not find contractor, %w", err)
	}
	// create mto -> move task order
	var mto = models.MoveTaskOrder{
		MoveOrderID:  mo.ID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ContractorID: contractor.ID,
	}

	builder := query.NewQueryBuilder(db)
	mtoCreator := movetaskorder.NewMoveTaskOrderCreator(builder, db)

	if _, verrs, err := mtoCreator.CreateMoveTaskOrder(&mto); err != nil || (verrs != nil && verrs.HasAny()) {
		return uuid.Nil, fmt.Errorf("Could not save move task order, %w, %v", err, verrs)
	}

	// create HHG -> house hold goods
	// mto shipment of type HHG
	requestedPickupDate := time.Now()
	scheduledPickupDate := time.Now()
	// add 7 days from requested pickup to scheduled pickup
	h, _ := time.ParseDuration("168h")
	scheduledPickupDate.Add(h)
	hhg := models.MTOShipment{
		MoveTaskOrderID:      mto.ID,
		RequestedPickupDate:  &requestedPickupDate,
		ScheduledPickupDate:  &scheduledPickupDate,
		PickupAddressID:      &sm.DutyStation.AddressID,
		DestinationAddressID: &move.Orders.NewDutyStation.AddressID,
		ShipmentType:         models.MTOShipmentTypeHHGLongHaulDom,
		Status:               models.MTOShipmentStatusSubmitted,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := db.Save(&hhg); err != nil {
		return uuid.Nil, fmt.Errorf("Could not save hhg shipment, %w", err)
	}

	// Domestic Shorthaul is less than 50 miles
	// Hard coding the shipment to meet that requirment
	var miramar models.DutyStation
	var sanDiego models.DutyStation

	if err := db.Where("name = ?", "USMC Miramar").First(&miramar); err != nil {
		return uuid.Nil, fmt.Errorf("Could not find miramar, %w", err)
	}

	if err := db.Where("name = ?", "USMC San Diego").First(&sanDiego); err != nil {
		return uuid.Nil, fmt.Errorf("Could not find san diego, %w", err)
	}

	hhgDomShortHaul := models.MTOShipment{
		MoveTaskOrderID:      mto.ID,
		RequestedPickupDate:  &requestedPickupDate,
		ScheduledPickupDate:  &scheduledPickupDate,
		PickupAddressID:      &sanDiego.AddressID,
		DestinationAddressID: &miramar.AddressID,
		ShipmentType:         models.MTOShipmentTypeHHGShortHaulDom,
		Status:               models.MTOShipmentStatusSubmitted,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := db.Save(&hhgDomShortHaul); err != nil {
		return uuid.Nil, fmt.Errorf("Could not save hhg domestic shorthaul shipment, %w", err)
	}

	return mo.ID, nil
}

// ConvertProfileOrdersToGHC creates models in the new GHC data model for pre move-setup data (no shipments created)
func ConvertProfileOrdersToGHC(db *pop.Connection, moveID uuid.UUID) (uuid.UUID, error) {
	var move models.Move
	if err := db.Eager("Orders.ServiceMember.DutyStation.Address", "Orders.NewDutyStation.Address").Find(&move, moveID); err != nil {
		return uuid.Nil, fmt.Errorf("Could not fetch move with id %s, %w", moveID, err)
	}

	sm := move.Orders.ServiceMember

	// create entitlement (required by move order)
	weight, entitlementErr := models.GetEntitlement(*sm.Rank, move.Orders.HasDependents, move.Orders.SpouseHasProGear)
	if entitlementErr != nil {
		return uuid.Nil, entitlementErr
	}
	entitlement := models.Entitlement{
		DependentsAuthorized: &move.Orders.HasDependents,
		DBAuthorizedWeight:   models.IntPointer(weight),
	}

	if err := db.Save(&entitlement); err != nil {
		return uuid.Nil, fmt.Errorf("Could not save entitlement, %w", err)
	}

	// orders -> move order
	orders := move.Orders
	var mo models.MoveOrder
	mo.CreatedAt = orders.CreatedAt
	mo.UpdatedAt = orders.UpdatedAt
	mo.Customer = &sm
	mo.CustomerID = &sm.ID
	mo.DestinationDutyStation = &orders.NewDutyStation
	mo.DestinationDutyStationID = &orders.NewDutyStationID

	orderType := "GHC"
	mo.OrderNumber = orders.OrdersNumber
	mo.OrderType = &orderType
	orderTypeDetail := "Shipment of HHG permitted"
	mo.OrderTypeDetail = &orderTypeDetail
	mo.OriginDutyStation = &sm.DutyStation
	mo.OriginDutyStationID = sm.DutyStationID
	mo.Entitlement = &entitlement
	mo.EntitlementID = &entitlement.ID
	mo.Grade = (*string)(sm.Rank)
	mo.DateIssued = &orders.IssueDate
	mo.ReportByDate = &orders.ReportByDate
	mo.LinesOfAccounting = orders.TAC

	if err := db.Save(&mo); err != nil {
		return uuid.Nil, fmt.Errorf("Could not save move order, %w", err)
	}

	var contractor models.Contractor

	err := db.Where("contract_number = ?", "HTC111-11-1-1111").First(&contractor)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Could not find contractor, %w", err)
	}
	// create mto -> move task order
	var mto = models.MoveTaskOrder{
		MoveOrderID:  mo.ID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ContractorID: contractor.ID,
	}

	builder := query.NewQueryBuilder(db)
	mtoCreator := movetaskorder.NewMoveTaskOrderCreator(builder, db)

	if _, verrs, err := mtoCreator.CreateMoveTaskOrder(&mto); err != nil || (verrs != nil && verrs.HasAny()) {
		return uuid.Nil, fmt.Errorf("Could not save move task order, %w, %v", err, verrs)
	}

	return mo.ID, nil
}
