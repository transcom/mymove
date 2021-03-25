package mtoagent

import (
	"github.com/gobuffalo/pop/v5"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// mtoAgentCreator sets up the service object
type mtoAgentCreator struct {
	db                     *pop.Connection
	mtoAvailabilityChecker services.MoveTaskOrderChecker
}

// NewMTOAgentCreator creates a new struct with the service dependencies
func NewMTOAgentCreator(db *pop.Connection, mtoAvailabilityChecker services.MoveTaskOrderChecker) services.MTOAgentCreator {
	return &mtoAgentCreator{
		db,
		mtoAvailabilityChecker,
	}
}

// #TODO: Will this only be used by Prime?
// I believe so, which means we only need to
// worry about validating for Prime
func (f *mtoAgentCreator) CreateMTOAgentPrime(mtoAgent *models.MTOAgent) (*models.MTOAgent, error) {
	return f.CreateMTOAgent(mtoAgent, "prime") //#TODO idk, remove this probably
}

func (f *mtoAgentCreator) CreateMTOAgent(mtoAgent *models.MTOAgent, validatorKey string) (*models.MTOAgent, error) {
	// Get existing shipment and agents information for validation
	mtoShipment := &models.MTOShipment{}
	err := f.db.Eager("MTOAgents").Find(mtoShipment, mtoAgent.MTOShipmentID)

	if err != nil {
		return nil, services.NewNotFoundError(mtoAgent.MTOShipmentID, "while looking for MTOShipment")
	}

	if validatorKey == "prime" { // #TODO update
		isAvailable, err := f.mtoAvailabilityChecker.MTOAvailableToPrime(mtoShipment.MoveTaskOrderID)

		if !isAvailable || err != nil {
			return nil, services.NewNotFoundError(mtoAgent.MTOShipmentID, "while looking for MTOShipment")
		}
	}

	// Confirm that MTOAgent does not already exist for the specified MTOAgentType
	for _, agent := range mtoShipment.MTOAgents {
		if agent.MTOAgentType == mtoAgent.MTOAgentType {
			return nil, services.NewConflictError(agent.ID, " MTOAgent already exists for this shipment. Please use updateMTOAgent endpoint.")
		}
	}

	verrs, err := f.db.ValidateAndCreate(mtoAgent) // what does this do, tho

	// If there were validation errors create an InvalidInputError type
	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(mtoAgent.ID, err, verrs, "Invalid input found while updating the agent.")
	} else if err != nil {
		// If the error is something else (this is unexpected), we create a QueryError
		return nil, services.NewQueryError("MTOAgent", err, "")
	}

	return mtoAgent, nil

}
