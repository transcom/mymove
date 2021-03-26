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

// CreateMTOAgentPrime passes the Prime validator key to CreateMTOAgent
func (f *mtoAgentCreator) CreateMTOAgentPrime(mtoAgent *models.MTOAgent) (*models.MTOAgent, error) {
	return f.CreateMTOAgent(mtoAgent, CreateMTOAgentPrimeValidator) //#TODO Should I leave this to establish pattern in case we need to add more? Or remove since I'm not using all the validating functions

}

func (f *mtoAgentCreator) CreateMTOAgent(mtoAgent *models.MTOAgent, validatorKey string) (*models.MTOAgent, error) {
	// Get existing shipment and agents information for validation
	mtoShipment := &models.MTOShipment{}
	var err error
	err = f.db.Eager("MTOAgents").Find(mtoShipment, mtoAgent.MTOShipmentID)

	if err != nil {
		return nil, services.NewNotFoundError(mtoAgent.MTOShipmentID, "while looking for MTOShipment")
	}

	if validatorKey == CreateMTOAgentPrimeValidator {
		var isAvailable bool

		isAvailable, err = f.mtoAvailabilityChecker.MTOAvailableToPrime(mtoShipment.MoveTaskOrderID)

		if !isAvailable || err != nil {
			return nil, services.NewNotFoundError(mtoAgent.MTOShipmentID, "while looking for MTOShipment")
		}
	}

	// Confirm that MTOAgent does not already exist for the specified MTOAgentType
	for _, agent := range mtoShipment.MTOAgents {
		if agent.MTOAgentType == mtoAgent.MTOAgentType {
			return nil, services.NewConflictError(agent.ID, " MTOAgent already exists for this agent type and shipment. Please use updateMTOAgent endpoint.")
		}
	}

	verrs, err := f.db.ValidateAndCreate(mtoAgent)

	// If there were validation errors create an InvalidInputError type
	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(mtoAgent.ID, err, verrs, "Invalid input found while updating the agent.")
	} else if err != nil {
		// If the error is something else (this is unexpected), we create a QueryError
		return nil, services.NewQueryError("MTOAgent", err, "")
	}

	return mtoAgent, nil

}
