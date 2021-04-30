package mtoagent

import (
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	mtoagentvalidate "github.com/transcom/mymove/pkg/services/mto_agent/validation"

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

// CreateMTOAgentBasic passes the Prime validator key to CreateMTOAgent
func (f *mtoAgentCreator) CreateMTOAgentBasic(mtoAgent *models.MTOAgent) (*models.MTOAgent, error) {
	return f.CreateMTOAgent(mtoAgent, mtoagentvalidate.BasicAgentValidatorKey)
}

// CreateMTOAgentPrime passes the Prime validator key to CreateMTOAgent
func (f *mtoAgentCreator) CreateMTOAgentPrime(mtoAgent *models.MTOAgent) (*models.MTOAgent, error) {
	return f.CreateMTOAgent(mtoAgent, mtoagentvalidate.PrimeAgentValidatorKey)
}

// CreateMTOAgent creates an MTO Agent
func (f *mtoAgentCreator) CreateMTOAgent(mtoAgent *models.MTOAgent, validatorKey string) (*models.MTOAgent, error) {
	// Get existing shipment and agents information for validation
	mtoShipment := &models.MTOShipment{}
	err := f.db.Eager("MTOAgents").Find(mtoShipment, mtoAgent.MTOShipmentID)

	if err != nil {
		return nil, services.NewNotFoundError(mtoAgent.MTOShipmentID, "while looking for MTOShipment")
	}

	agentData := mtoagentvalidate.AgentValidationData{
		NewAgent:            *mtoAgent,
		Shipment:            mtoShipment,
		AvailabilityChecker: f.mtoAvailabilityChecker,
		Verrs:               validate.NewErrors(),
	}

	mtoAgent, err = mtoagentvalidate.ValidateAgent(&agentData, validatorKey)
	if err != nil {
		return nil, err
	}

	verrs, err := f.db.ValidateAndCreate(mtoAgent)

	// If there were validation errors create an InvalidInputError type
	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(uuid.Nil, err, verrs, "Invalid input found while creating the agent.")
	} else if err != nil {
		// If the error is something else (this is unexpected), we create a QueryError
		return nil, services.NewQueryError("MTOAgent", err, "")
	}

	return mtoAgent, nil

}
