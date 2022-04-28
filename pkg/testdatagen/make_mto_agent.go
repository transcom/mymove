package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMTOAgent creates a single MTOAgent and associated MTOShipment
func MakeMTOAgent(db *pop.Connection, assertions Assertions) models.MTOAgent {

	mtoShipment := assertions.MTOShipment
	if isZeroUUID(mtoShipment.ID) {
		mtoShipment = MakeMTOShipment(db, assertions)
	}

	firstName := "Jason"
	lastName := "Ash"
	email := "jason.ash@example.com"
	phone := "202-555-9301"

	MTOAgent := models.MTOAgent{
		MTOShipment:   mtoShipment,
		MTOShipmentID: mtoShipment.ID,
		FirstName:     &firstName,
		LastName:      &lastName,
		Email:         &email,
		Phone:         &phone,
		MTOAgentType:  models.MTOAgentReleasing,
	}

	mergeModels(&MTOAgent, assertions.MTOAgent)

	mustCreate(db, &MTOAgent, assertions.Stub)

	return MTOAgent
}

// MakeDefaultMTOAgent returns a MTOAgent with default values
func MakeDefaultMTOAgent(db *pop.Connection) models.MTOAgent {
	return MakeMTOAgent(db, Assertions{})
}
