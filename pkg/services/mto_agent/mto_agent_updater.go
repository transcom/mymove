package mtoagent

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// mtoAgentUpdater handles the db connection
type mtoAgentUpdater struct {
	basicChecks []mtoAgentValidator
	primeChecks []mtoAgentValidator
}

// NewMTOAgentUpdater creates a new struct with the service dependencies
func NewMTOAgentUpdater(mtoChecker services.MoveTaskOrderChecker) services.MTOAgentUpdater {
	return &mtoAgentUpdater{
		basicChecks: []mtoAgentValidator{
			checkShipmentID(),
			checkAgentID(),
		},
		primeChecks: []mtoAgentValidator{
			checkShipmentID(),
			checkAgentID(),
			checkContactInfo(),
			checkAgentType(),
			checkPrimeAvailability(mtoChecker),
		},
	}
}

// UpdateMTOAgentBasic updates the MTO Agent using base validators
func (f *mtoAgentUpdater) UpdateMTOAgentBasic(appCfg appconfig.AppConfig, mtoAgent *models.MTOAgent, eTag string) (*models.MTOAgent, error) {
	return f.updateMTOAgent(appCfg, mtoAgent, eTag, f.basicChecks...)
}

// UpdateMTOAgentPrime updates the MTO Agent using Prime API validators
func (f *mtoAgentUpdater) UpdateMTOAgentPrime(appCfg appconfig.AppConfig, mtoAgent *models.MTOAgent, eTag string) (*models.MTOAgent, error) {
	return f.updateMTOAgent(appCfg, mtoAgent, eTag, f.primeChecks...)
}

// UpdateMTOAgent updates the MTO Agent
func (f *mtoAgentUpdater) updateMTOAgent(appCfg appconfig.AppConfig, mtoAgent *models.MTOAgent, eTag string, checks ...mtoAgentValidator) (*models.MTOAgent, error) {
	oldAgent := models.MTOAgent{}

	// Find the agent, return error if not found
	err := appCfg.DB().Eager("MTOShipment.MTOAgents").Find(&oldAgent, mtoAgent.ID)
	if err != nil {
		return nil, services.NewNotFoundError(mtoAgent.ID, "while looking for MTOAgent")
	}

	err = validateMTOAgent(appCfg, *mtoAgent, &oldAgent, &oldAgent.MTOShipment, checks...)
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
	verrs, err := appCfg.DB().ValidateAndSave(newAgent)

	// If there were validation errors create an InvalidInputError type
	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(newAgent.ID, err, verrs, "Invalid input found while updating the agent.")
	} else if err != nil {
		// If the error is something else (this is unexpected), we create a QueryError
		return nil, services.NewQueryError("MTOAgent", err, "")
	}

	// Get the updated agent and return
	updatedAgent := models.MTOAgent{}
	err = appCfg.DB().Find(&updatedAgent, newAgent.ID)
	if err != nil {
		return nil, services.NewQueryError("MTOAgent", err, fmt.Sprintf("Unexpected error after saving: %v", err))
	}
	return &updatedAgent, nil
}
