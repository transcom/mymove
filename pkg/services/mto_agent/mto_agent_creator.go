package mtoagent

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// mtoAgentCreator sets up the service object
type mtoAgentCreator struct {
	basicChecks []mtoAgentValidator
	primeChecks []mtoAgentValidator
}

// NewMTOAgentCreator creates a new struct with the service dependencies
func NewMTOAgentCreator(mtoAvailabilityChecker services.MoveTaskOrderChecker) services.MTOAgentCreator {
	return &mtoAgentCreator{
		basicChecks: []mtoAgentValidator{
			checkShipmentID(),
			checkAgentID(),
		},
		primeChecks: []mtoAgentValidator{
			checkShipmentID(),
			checkAgentID(),
			checkContactInfo(),
			checkAgentType(),
			checkPrimeAvailability(mtoAvailabilityChecker),
		},
	}
}

// CreateMTOAgentBasic passes the Prime validator key to CreateMTOAgent
func (f *mtoAgentCreator) CreateMTOAgentBasic(appCtx appcontext.AppContext, mtoAgent *models.MTOAgent) (*models.MTOAgent, error) {
	return f.createMTOAgent(appCtx, mtoAgent, f.basicChecks...)
}

// CreateMTOAgentPrime passes the Prime validator key to CreateMTOAgent
func (f *mtoAgentCreator) CreateMTOAgentPrime(appCtx appcontext.AppContext, mtoAgent *models.MTOAgent) (*models.MTOAgent, error) {
	return f.createMTOAgent(appCtx, mtoAgent, f.primeChecks...)
}

// CreateMTOAgent creates an MTO Agent
func (f *mtoAgentCreator) createMTOAgent(appCtx appcontext.AppContext, mtoAgent *models.MTOAgent, checks ...mtoAgentValidator) (*models.MTOAgent, error) {
	// Get existing shipment and agents information for validation
	mtoShipment := &models.MTOShipment{}
	err := appCtx.DB().Eager("MTOAgents").Find(mtoShipment, mtoAgent.MTOShipmentID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(mtoAgent.MTOShipmentID, "while looking for MTOShipment")
		default:
			return nil, apperror.NewQueryError("MTOShipment", err, "")
		}
	}

	err = validateMTOAgent(appCtx, *mtoAgent, nil, mtoShipment, checks...)
	if err != nil {
		return nil, err
	}

	// TODO: this basically operates as a way to make a (shallow?) copy of the
	// mtoAgent (because pass-by-value parameter); haven't dug into
	// why this is necessary
	mtoAgent = mergeAgent(*mtoAgent, nil)

	verrs, err := appCtx.DB().ValidateAndCreate(mtoAgent)

	// If there were validation errors create an InvalidInputError type
	if verrs != nil && verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Invalid input found while creating the agent.")
	} else if err != nil {
		// If the error is something else (this is unexpected), we create a QueryError
		return nil, apperror.NewQueryError("MTOAgent", err, "")
	}

	return mtoAgent, nil

}
