package mtoagent

import (
	"fmt"

	mtoagentvalidate "github.com/transcom/mymove/pkg/services/mto_agent/validate"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
)

// mtoAgentUpdater handles the db connection
type mtoAgentUpdater struct {
	db *pop.Connection
}

// NewMTOAgentUpdater creates a new struct with the service dependencies
func NewMTOAgentUpdater(db *pop.Connection) services.MTOAgentUpdater {
	return &mtoAgentUpdater{
		db: db,
	}
}

// UpdateMTOAgentBasic updates the MTO Agent using base validators
func (f *mtoAgentUpdater) UpdateMTOAgentBasic(mtoAgent *models.MTOAgent, eTag string) (*models.MTOAgent, error) {
	return f.UpdateMTOAgent(mtoAgent, eTag, mtoagentvalidate.BasicAgentValidatorKey)
}

// UpdateMTOAgentPrime updates the MTO Agent using Prime API validators
func (f *mtoAgentUpdater) UpdateMTOAgentPrime(mtoAgent *models.MTOAgent, eTag string) (*models.MTOAgent, error) {
	return f.UpdateMTOAgent(mtoAgent, eTag, mtoagentvalidate.PrimeAgentValidatorKey)
}

// UpdateMTOAgent updates the MTO Agent
func (f *mtoAgentUpdater) UpdateMTOAgent(mtoAgent *models.MTOAgent, eTag string, validatorKey string) (*models.MTOAgent, error) {
	oldAgent := models.MTOAgent{}

	// Find the agent, return error if not found
	err := f.db.Eager("MTOShipment.MTOAgents").Find(&oldAgent, mtoAgent.ID)
	if err != nil {
		return nil, services.NewNotFoundError(mtoAgent.ID, "while looking for MTOAgent")
	}

	agentData := mtoagentvalidate.AgentValidationData{
		NewAgent:            *mtoAgent,
		OldAgent:            &oldAgent,
		Shipment:            &oldAgent.MTOShipment,
		AvailabilityChecker: movetaskorder.NewMoveTaskOrderChecker(f.db),
		Verrs:               validate.NewErrors(),
	}

	newAgent, err := mtoagentvalidate.ValidateAgent(&agentData, validatorKey)
	if err != nil {
		return nil, err
	}

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
