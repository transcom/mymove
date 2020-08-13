package converthelper

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
)

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

	return orders.ID, nil
}
