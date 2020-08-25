package testdatagen

import (
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMTOAgent creates a single MTOAgent and its associations.
func MakeMTOAgent(db *pop.Connection, assertions Assertions) models.MTOAgent {
	// Will make an MTO if one is not provided via assertion.
	mtoShipment := assertions.MTOShipment
	if isZeroUUID(mtoShipment.ID) {
		mtoShipment = MakeMTOShipment(db, assertions)
	}

	mtoAgent := models.MTOAgent{
		MTOShipment:   mtoShipment,
		MTOShipmentID: mtoShipment.ID,
		FirstName:     swag.String("Test"),
		LastName:      swag.String("Agent"),
		MTOAgentType:  models.MTOAgentReleasing,
		Email:         swag.String("test@test.email.com"),
	}

	// Overwrite default values with those from assertions.
	mergeModels(&mtoAgent, assertions.MTOAgent)
	mustCreate(db, &mtoAgent)

	return mtoAgent
}

// MakeDefaultMTOAgent makes an MTOAgent with default values
func MakeDefaultMTOAgent(db *pop.Connection) models.MTOAgent {
	return MakeMTOAgent(db, Assertions{})
}
