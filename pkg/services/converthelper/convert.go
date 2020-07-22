package converthelper

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/models"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/unit"

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

	// add/update fields originally from move_order
	orders := move.Orders
	orders.Grade = (*string)(sm.Rank)
	orders.EntitlementID = &entitlement.ID
	orders.Entitlement = &entitlement
	orders.OriginDutyStationID = sm.DutyStationID
	orders.OriginDutyStation = &sm.DutyStation

	if err := db.Save(&orders); err != nil {
		return uuid.Nil, fmt.Errorf("Could not save order, %w", err)
	}

	var moveTaskOrders []models.MoveTaskOrder
	err := db.Where("move_order_id = $1", orders.ID).All(&moveTaskOrders)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Could not fetch move task order, %w", err)
	}
	// create mto -> move task order
	var mto models.MoveTaskOrder
	if len(moveTaskOrders) == 0 {
		var contractor models.Contractor

		err := db.Where("contract_number = ?", "HTC111-11-1-1111").First(&contractor)
		if err != nil {
			return uuid.Nil, fmt.Errorf("Could not find contractor, %w", err)
		}

		mto = models.MoveTaskOrder{
			MoveOrderID:  orders.ID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			ContractorID: contractor.ID,
		}

		builder := query.NewQueryBuilder(db)
		mtoCreator := movetaskorder.NewMoveTaskOrderCreator(builder, db)

		if _, verrs, err := mtoCreator.CreateMoveTaskOrder(&mto); err != nil || (verrs != nil && verrs.HasAny()) {
			return uuid.Nil, fmt.Errorf("Could not save move task order, %w, %v", err, verrs)
		}
	} else {
		mto = moveTaskOrders[0]
	}

	// create HHG -> house hold goods
	// mto shipment of type HHG
	// RequestedPickupDate 10 days ago
	requestedPickupDate := time.Now().AddDate(0, 0, -10)
	// add 7 days from requested pickup to scheduled pickup
	scheduledPickupDate := requestedPickupDate.AddDate(0, 0, 7)
	// set prime estmated record date to 3 days before scheduledPickupDate
	primeEstimatedWeightRecordDate := scheduledPickupDate.AddDate(0, 0, -3)
	primeEstimatedWeight := unit.Pound(4096)
	hhg := models.MTOShipment{
		MoveTaskOrderID:                  mto.ID,
		RequestedPickupDate:              &requestedPickupDate,
		ScheduledPickupDate:              &scheduledPickupDate,
		PickupAddressID:                  &sm.DutyStation.AddressID,
		DestinationAddressID:             &move.Orders.NewDutyStation.AddressID,
		ShipmentType:                     models.MTOShipmentTypeHHGLongHaulDom,
		Status:                           models.MTOShipmentStatusSubmitted,
		PrimeEstimatedWeight:             &primeEstimatedWeight,
		PrimeEstimatedWeightRecordedDate: &primeEstimatedWeightRecordDate,
		CreatedAt:                        time.Now(),
		UpdatedAt:                        time.Now(),
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
		MoveTaskOrderID:                  mto.ID,
		RequestedPickupDate:              &requestedPickupDate,
		ScheduledPickupDate:              &scheduledPickupDate,
		PickupAddressID:                  &sanDiego.AddressID,
		DestinationAddressID:             &miramar.AddressID,
		ShipmentType:                     models.MTOShipmentTypeHHGShortHaulDom,
		Status:                           models.MTOShipmentStatusSubmitted,
		PrimeEstimatedWeight:             &primeEstimatedWeight,
		PrimeEstimatedWeightRecordedDate: &primeEstimatedWeightRecordDate,
		CreatedAt:                        time.Now(),
		UpdatedAt:                        time.Now(),
	}

	if err := db.Save(&hhgDomShortHaul); err != nil {
		return uuid.Nil, fmt.Errorf("Could not save hhg domestic shorthaul shipment, %w", err)
	}

	return orders.ID, nil
}

// ConvertProfileOrdersToGHC creates models in the new GHC data model for pre move-setup data (no shipments created)
func ConvertProfileOrdersToGHC(db *pop.Connection, moveID uuid.UUID) (uuid.UUID, error) {
	var move models.Move
	if err := db.Eager("Orders.ServiceMember").Find(&move, moveID); err != nil {
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

	// add fields originally from move_order
	orders := move.Orders
	orders.Grade = (*string)(sm.Rank)
	orders.EntitlementID = &entitlement.ID
	orders.Entitlement = &entitlement
	orders.OriginDutyStationID = sm.DutyStationID
	orders.OriginDutyStation = &sm.DutyStation

	if err := db.Save(&orders); err != nil {
		return uuid.Nil, fmt.Errorf("Could not save order, %w", err)
	}

	var contractor models.Contractor

	err := db.Where("contract_number = ?", "HTC111-11-1-1111").First(&contractor)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Could not find contractor, %w", err)
	}
	// create mto -> move task order
	var mto = models.MoveTaskOrder{
		MoveOrderID:  orders.ID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ContractorID: contractor.ID,
	}

	builder := query.NewQueryBuilder(db)
	mtoCreator := movetaskorder.NewMoveTaskOrderCreator(builder, db)

	if _, verrs, err := mtoCreator.CreateMoveTaskOrder(&mto); err != nil || (verrs != nil && verrs.HasAny()) {
		return uuid.Nil, fmt.Errorf("Could not save move task order, %w, %v", err, verrs)
	}

	return orders.ID, nil
}
