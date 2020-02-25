package testdatagen

import (
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMTOAgent creates a single MTOAgent and its associations.
func MakeATOAgent(db *pop.Connection, assertions Assertions) models.MTOAgent {
	// Will make an MTO if one is not provided via assertion.
	moveTaskOrder := assertions.PaymentRequest.MoveTaskOrder
	if isZeroUUID(moveTaskOrder.ID) {
		moveTaskOrder = MakeMoveTaskOrder(db, assertions)
	}

	mtoAgent := models.MTOAgent{
		MoveTaskOrder:   moveTaskOrder,
		MoveTaskOrderID: moveTaskOrder.ID,
		FirstName:       swag.String("Test"),
		LastName:        swag.String("Agent"),
		MTOAgentType:    models.MTOAgentReleasing,
		Email:           swag.String("test@test.email.com"),
	}

	// Overwrite default values with those from assertions.
	mergeModels(&mtoAgent, assertions.MTOAgent)
	mustCreate(db, &mtoAgent)

	return mtoAgent
}