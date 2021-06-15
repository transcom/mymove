package mtoagent

import (
	"context"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// mtoAgentCreator sets up the service object
type mtoAgentCreator struct {
	db          *pop.Connection
	basicChecks []mtoAgentValidator
	primeChecks []mtoAgentValidator
}

// NewMTOAgentCreator creates a new struct with the service dependencies
func NewMTOAgentCreator(db *pop.Connection, mtoAvailabilityChecker services.MoveTaskOrderChecker) services.MTOAgentCreator {
	return &mtoAgentCreator{
		db:          db,
		basicChecks: basicChecks,
		primeChecks: append(primeChecks, checkPrimeAvailability(mtoAvailabilityChecker)),
	}
}

// CreateMTOAgentBasic passes the Prime validator key to CreateMTOAgent
func (f *mtoAgentCreator) CreateMTOAgentBasic(mtoAgent *models.MTOAgent) (*models.MTOAgent, error) {
	return f.createMTOAgent(mtoAgent, f.basicChecks...)
}

// CreateMTOAgentPrime passes the Prime validator key to CreateMTOAgent
func (f *mtoAgentCreator) CreateMTOAgentPrime(mtoAgent *models.MTOAgent) (*models.MTOAgent, error) {
	return f.createMTOAgent(mtoAgent, f.primeChecks...)
}

// CreateMTOAgent creates an MTO Agent
func (f *mtoAgentCreator) createMTOAgent(mtoAgent *models.MTOAgent, checks ...mtoAgentValidator) (*models.MTOAgent, error) {
	// Get existing shipment and agents information for validation
	mtoShipment := &models.MTOShipment{}
	err := f.db.Eager("MTOAgents").Find(mtoShipment, mtoAgent.MTOShipmentID)
	if err != nil {
		return nil, services.NewNotFoundError(mtoAgent.MTOShipmentID, "while looking for MTOShipment")
	}

	err = validateMTOAgent(context.TODO(), *mtoAgent, nil, mtoShipment, checks...)
	if err != nil {
		return nil, err
	}

	// TODO: this basically operates as a way to make a (shallow?) copy of the
	// mtoAgent (because pass-by-value parameter); haven't dug into
	// why this is necessary
	mtoAgent = mergeAgent(*mtoAgent, nil)

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
