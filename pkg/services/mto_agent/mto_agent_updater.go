package mtoagent

import (
	"fmt"

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

// UpdateMTOAgentBase updates the MTO Agent using base validators
func (f *mtoAgentUpdater) UpdateMTOAgentBase(mtoAgent *models.MTOAgent, eTag string) (*models.MTOAgent, error) {
	return f.UpdateMTOAgent(mtoAgent, eTag, UpdateMTOAgentBaseValidator)
}

// UpdateMTOAgentPrime updates the MTO Agent using Prime API validators
func (f *mtoAgentUpdater) UpdateMTOAgentPrime(mtoAgent *models.MTOAgent, eTag string) (*models.MTOAgent, error) {
	return f.UpdateMTOAgent(mtoAgent, eTag, UpdateMTOAgentPrimeValidator)
}

// UpdateMTOAgent updates the MTO Agent
func (f *mtoAgentUpdater) UpdateMTOAgent(mtoAgent *models.MTOAgent, eTag string, validatorKey string) (*models.MTOAgent, error) {
	oldAgent := models.MTOAgent{}

	// Find the agent, return error if not found
	err := f.db.Eager("MTOShipment").Find(&oldAgent, mtoAgent.ID)
	if err != nil {
		return nil, services.NewNotFoundError(mtoAgent.ID, "while looking for MTOAgent")
	}

	checker := movetaskorder.NewMoveTaskOrderChecker(f.db)
	agentData := updateMTOAgentData{
		updatedAgent:        *mtoAgent,
		oldAgent:            oldAgent,
		availabilityChecker: checker,
		verrs:               validate.NewErrors(),
	}

	newAgent, err := ValidateUpdateMTOAgent(&agentData, validatorKey)
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

	// Get the updated address and return
	updatedAgent := models.MTOAgent{}
	err = f.db.Find(&updatedAgent, newAgent.ID)
	if err != nil {
		return nil, services.NewQueryError("MTOAgent", err, fmt.Sprintf("Unexpected error after saving: %v", err))
	}
	return &updatedAgent, nil
}

// ValidateUpdateMTOAgent checks the provided agentData struct against the validator indicated by validatorKey.
// Defaults to base validation if the empty string is entered as the key.
// Returns an MTOAgent that has been set up for update.
func ValidateUpdateMTOAgent(agentData *updateMTOAgentData, validatorKey string) (*models.MTOAgent, error) {
	var newAgent models.MTOAgent

	if validatorKey == "" {
		validatorKey = UpdateMTOAgentBaseValidator
	}
	validator, ok := UpdateMTOAgentValidators[validatorKey]
	if !ok {
		err := fmt.Errorf("validator key %s was not found in update MTO Agent validators", validatorKey)
		return nil, err
	}
	err := validator.validate(agentData)
	if err != nil {
		return nil, err
	}
	err = agentData.setNewMTOAgent(&newAgent)
	if err != nil {
		return nil, err
	}

	return &newAgent, nil
}
