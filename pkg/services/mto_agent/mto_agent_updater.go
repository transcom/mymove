package mtoagent

import (
	"fmt"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
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

// UpdateMTOAgent updates the MTO Agent
func (f *mtoAgentUpdater) UpdateMTOAgent(mtoAgent *models.MTOAgent, eTag string, validatorKey string) (*models.MTOAgent, error) {
	oldAgent := models.MTOAgent{}

	// Find the agent, return error if not found
	err := f.db.Eager("MTOShipment").Find(&oldAgent, mtoAgent.ID)
	if err != nil {
		return nil, services.NewNotFoundError(mtoAgent.ID, "while looking for MTOAgent")
	}

	newAgent := models.MTOAgent{}
	// TODO validation etc

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
