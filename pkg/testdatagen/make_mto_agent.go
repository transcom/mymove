package testdatagen

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMTOAgent creates a single MTOAgent and associated MTOShipment
func MakeMTOAgent(db *pop.Connection, assertions Assertions) models.MTOAgent {
	var mtoShipmentID uuid.UUID
	mtoShipment := assertions.MTOAgent.MTOShipment

	if isZeroUUID(assertions.MTOAgent.MTOShipmentID) {
		mtoShipment = MakeMTOShipment(db, assertions)
		mtoShipmentID = mtoShipment.ID
	}

	firstName := "Jason"
	lastName := "Ash"
	email := "jason.ash@gmail.com"
	phone := "2025559301"

	mtoAgentType := models.MTOAgentReleasing
	if string(assertions.MTOAgent.MTOAgentType) != "" {
		mtoAgentType = assertions.MTOAgent.MTOAgentType
	}

	MTOAgent := models.MTOAgent{
		MTOShipment:   mtoShipment,
		MTOShipmentID: mtoShipmentID,
		FirstName:     &firstName,
		LastName:      &lastName,
		Email:         &email,
		Phone:         &phone,
		MTOAgentType:  mtoAgentType,
	}

	mergeModels(&MTOAgent, assertions.MTOAgent)

	mustCreate(db, &MTOAgent)

	return MTOAgent
}

// MakeDefaultMTOAgent returns a MTOAgent with default values
func MakeDefaultMTOAgent(db *pop.Connection) models.MTOAgent {
	return MakeMTOAgent(db, Assertions{})
}