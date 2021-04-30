package mtoagent

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
)

// BasicAgentValidatorKey is the key for generic validation on the Agent
const BasicAgentValidatorKey string = "BasicAgentValidatorKey"

// PrimeAgentValidatorKey is the key for validating the Agent for the Prime contractor
const PrimeAgentValidatorKey string = "PrimeAgentValidatorKey"

// agentValidators is the map connecting the constant keys to the correct validator
// NOTE: This and the following Validate functions are non-importable so that devs will be forced to call them through
// the ValidateAgent function, which is more complete.
var agentValidators = map[string]func(agentData *AgentValidationData) error{
	BasicAgentValidatorKey: basicValidate,
	PrimeAgentValidatorKey: primeValidate,
}

// basicValidate performs the necessary checks for validation that should happen no matter who uses this service object
func basicValidate(agentData *AgentValidationData) error {
	err := agentData.checkShipmentID()
	if err != nil {
		return err
	}
	return nil
}

// primeValidate peforms the necessary functions to validate agent data submitted by the Prime contractor
func primeValidate(agentData *AgentValidationData) error {
	err := agentData.checkShipmentID()
	if err != nil {
		return err
	}
	err = agentData.checkPrimeAvailability()
	if err != nil {
		return err
	}
	err = agentData.checkContactInfo()
	if err != nil {
		return err
	}
	err = agentData.getVerrs()
	if err != nil {
		return err
	}
	return nil
}

// ValidateAgent checks the provided agentData struct against the validator indicated by validatorKey.
// Defaults to base validation if the empty string is entered as the key.
// Returns an MTOAgent that has been set up for update.
func ValidateAgent(agentData *AgentValidationData, validatorKey string) (*models.MTOAgent, error) {
	if validatorKey == "" {
		validatorKey = BasicAgentValidatorKey
	}
	validator, ok := agentValidators[validatorKey]
	if !ok {
		err := fmt.Errorf("could not find agent validator with key %s", validatorKey)
		return nil, err
	}
	err := validator(agentData)
	if err != nil {
		return nil, err
	}

	err = agentData.getVerrs()
	if err != nil {
		return nil, err
	}

	validatedAgent := agentData.setFullAgent()

	return validatedAgent, nil
}
