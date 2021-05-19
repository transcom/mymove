package mtoagentvalidate

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// NOTE: These validator keys are used the same way for both Update and Create actions.
// The values passed into the AgentValidationData struct are what differentiate the two.

// BasicAgentValidatorKey is the key for generic validation on the Agent
const BasicAgentValidatorKey string = "BasicAgentValidatorKey"

// PrimeAgentValidatorKey is the key for validating the Agent for the Prime contractor
const PrimeAgentValidatorKey string = "PrimeAgentValidatorKey"

// PrimeAvailableAgentValidatorKey is the key for validating an Agent that we already know is Prime-available
const PrimeAvailableAgentValidatorKey string = "PrimeAvailableAgentValidatorKey"

// agentValidators is the map connecting the constant keys to the correct validator
// NOTE: This and the following validate functions are not exportable so that devs will be forced to call them through
// the ValidateAgent function, which is more complete.
var agentValidators = map[string]agentValidator{
	BasicAgentValidatorKey:          new(basicAgentValidator),
	PrimeAgentValidatorKey:          new(primeAgentValidator),
	PrimeAvailableAgentValidatorKey: primeAgentValidator{isAvailableToPrime: true},
}

type agentValidator interface {
	validate(agentData *AgentValidationData) error
}

// basicAgentValidator is the type for validation that should happen no matter who uses this service object
type basicAgentValidator struct{}

func (v basicAgentValidator) validate(agentData *AgentValidationData) error {
	checks := []services.ValidationFunc{
		agentData.checkShipmentID,
		agentData.checkAgentID,
	}
	return services.CheckValidationData(checks)
}

// primeAgentValidator is the type for validation of agent data submitted by the Prime contractor
type primeAgentValidator struct {
	// isAvailableToPrime allows us to tell the validator if we already know the agent's availability to the Prime.
	// checkPrimeAvailability hits the DB, so it's convenient to skip that check if we already know its status.
	isAvailableToPrime bool
}

func (v primeAgentValidator) validate(agentData *AgentValidationData) error {
	checks := []services.ValidationFunc{
		agentData.checkShipmentID,
		agentData.checkAgentID,
		agentData.checkContactInfo,
		agentData.checkAgentType,
	}
	if !v.isAvailableToPrime {
		checks = append(checks, agentData.checkPrimeAvailability)
	}
	return services.CheckValidationData(checks)
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
	err := validator.validate(agentData)
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
