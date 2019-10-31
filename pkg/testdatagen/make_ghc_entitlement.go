package testdatagen

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

//  MakeEntitlement creates a single GHCEntitlement and associated set of relationships
func MakeEntitlement(db *pop.Connection, assertions Assertions) models.GHCEntitlement {

	var moveTaskOrderID uuid.UUID
	if assertions.GHCEntitlement.MoveTaskOrder != nil {
		moveTaskOrderID = assertions.GHCEntitlement.MoveTaskOrder.ID
	}
	var moveTaskOrder models.MoveTaskOrder
	if isZeroUUID(moveTaskOrderID) {
		moveTaskOrder = MakeDefaultMoveTaskOrder(db)
	}

	entitlement := models.GHCEntitlement{
		DependentsAuthorized:  true,
		TotalDependents:       1,
		NonTemporaryStorage:   true,
		PrivatelyOwnedVehicle: true,
		ProGearWeight:         100,
		ProGearWeightSpouse:   200,
		StorageInTransit:      2,
		MoveTaskOrderID:       moveTaskOrder.ID,
		MoveTaskOrder:         &moveTaskOrder,
	}
	// Overwrite values with those from assertions
	mergeModels(&entitlement, assertions.GHCEntitlement)

	mustCreate(db, &entitlement)

	return entitlement
}
