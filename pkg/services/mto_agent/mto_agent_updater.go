package mtoagent

import (
	"context"
	"fmt"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// mtoAgentUpdater handles the db connection
type mtoAgentUpdater struct {
	db          *pop.Connection
	basicChecks []mtoAgentValidator
	primeChecks []mtoAgentValidator
}

// NewMTOAgentUpdater creates a new struct with the service dependencies
func NewMTOAgentUpdater(db *pop.Connection, mtoChecker services.MoveTaskOrderChecker) services.MTOAgentUpdater {
	return &mtoAgentUpdater{
		db:          db,
		basicChecks: basicChecks,
		primeChecks: append(primeChecks, checkPrimeAvailability(mtoChecker)),
	}
}

// UpdateMTOAgentBasic updates the MTO Agent using base validators
func (f *mtoAgentUpdater) UpdateMTOAgentBasic(mtoAgent *models.MTOAgent, eTag string) (*models.MTOAgent, error) {
	return f.updateMTOAgent(mtoAgent, eTag, f.basicChecks...)
}

// UpdateMTOAgentPrime updates the MTO Agent using Prime API validators
func (f *mtoAgentUpdater) UpdateMTOAgentPrime(mtoAgent *models.MTOAgent, eTag string) (*models.MTOAgent, error) {
	return f.updateMTOAgent(mtoAgent, eTag, f.primeChecks...)
}

// UpdateMTOAgent updates the MTO Agent
func (f *mtoAgentUpdater) updateMTOAgent(mtoAgent *models.MTOAgent, eTag string, checks ...mtoAgentValidator) (*models.MTOAgent, error) {
	oldAgent := models.MTOAgent{}

	// Find the agent, return error if not found
	err := f.db.Eager("MTOShipment.MTOAgents").Find(&oldAgent, mtoAgent.ID)
	if err != nil {
		return nil, services.NewNotFoundError(mtoAgent.ID, "while looking for MTOAgent")
	}

	err = validateMTOAgent(context.TODO(), *mtoAgent, &oldAgent, &oldAgent.MTOShipment, checks...)
	if err != nil {
		return nil, err
	}
	newAgent := mergeAgent(*mtoAgent, &oldAgent)

	// Check the If-Match header against existing eTag before updating
	encodedUpdatedAt := etag.GenerateEtag(oldAgent.UpdatedAt)
	if encodedUpdatedAt != eTag {
		return nil, services.NewPreconditionFailedError(newAgent.ID, nil)
	}

	// Make the update and create a InvalidInputError if there were validation issues
	verrs, err := f.db.ValidateAndSave(newAgent)

	// If there were validation errors create an InvalidInputError type
	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(newAgent.ID, err, verrs, "Invalid input found while updating the agent.")
	} else if err != nil {
		// If the error is something else (this is unexpected), we create a QueryError
		return nil, services.NewQueryError("MTOAgent", err, "")
	}

	// Get the updated agent and return
	updatedAgent := models.MTOAgent{}
	err = f.db.Find(&updatedAgent, newAgent.ID)
	if err != nil {
		return nil, services.NewQueryError("MTOAgent", err, fmt.Sprintf("Unexpected error after saving: %v", err))
	}
	return &updatedAgent, nil
}
