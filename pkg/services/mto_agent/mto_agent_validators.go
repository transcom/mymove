package mtoagent

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// BasicAgentValidatorKey is the key for generic validation on the MTO Agent
const BasicAgentValidatorKey string = "BasicAgentValidatorKey"

// PrimeAgentValidatorKey is the key for validating the MTO Agent for the Prime contractor
const PrimeAgentValidatorKey string = "PrimeAgentValidatorKey"

// CreateMTOAgentPrimeValidator is the key for validating the MTO Agent for the Prime
const CreateMTOAgentPrimeValidator = "CreateMTOAgentPrimeValidator"

// agentValidators is the map connecting the constant keys to the correct validator
var agentValidators = map[string]AgentValidator{
	BasicAgentValidatorKey: new(BasicAgentValidator),
	PrimeAgentValidatorKey: new(PrimeAgentValidator),
}

// AgentValidator is the base interface for all MTO Agent validator types
type AgentValidator interface {
	Validate(agentData *AgentValidationData) error
}

// BasicAgentValidator is the type for validation that should happen no matter who uses this service object
type BasicAgentValidator struct{}

// Validate performs the necessary functions for basic validation
func (v *BasicAgentValidator) Validate(agentData *AgentValidationData) error {
	err := agentData.checkShipmentID()
	if err != nil {
		return err
	}

	err = agentData.getVerrs()
	if err != nil {
		return err
	}

	return nil
}

// PrimeAgentValidator is the type for validation that is just for updates from the Prime contractor
type PrimeAgentValidator struct{}

// Validate peforms the necessary functions to validate agent data from a Prime user
func (v *PrimeAgentValidator) Validate(agentData *AgentValidationData) error {
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

// AgentValidationData represents the data needed to validate an Agent before a create/update action
type AgentValidationData struct {
	newAgent            models.MTOAgent
	oldAgent            *models.MTOAgent // not required for create
	moveID              uuid.UUID
	availabilityChecker services.MoveTaskOrderChecker
	verrs               *validate.Errors
}

// checkShipmentID checks that the user didn't attempt to change the agent's shipment ID
func (v *AgentValidationData) checkShipmentID() error {
	if v.oldAgent == nil {
		if v.newAgent.MTOShipmentID == uuid.Nil {
			v.verrs.Add("mtoShipmentID", "shipment ID is required")
		}
	} else {
		if v.newAgent.MTOShipmentID != uuid.Nil && v.newAgent.MTOShipmentID != v.oldAgent.MTOShipmentID {
			v.verrs.Add("mtoShipmentID", "cannot be updated")
		}
	}

	return nil
}

// checkPrimeAvailability checks that agent is connected to a Prime-available shipment
func (v *AgentValidationData) checkPrimeAvailability() error {
	isAvailable, err := v.availabilityChecker.MTOAvailableToPrime(v.moveID)

	if !isAvailable || err != nil {
		return services.NewNotFoundError(v.newAgent.ID, "while looking for Prime-available MTOAgent")
	}

	return nil
}

// checkContactInfo checks that the new agent has the minimum required contact info: First Name and one of Email or Phone
func (v *AgentValidationData) checkContactInfo() error {
	var firstName *string
	var email *string
	var phone *string

	// Set any pre-existing values as the baseline:
	if v.oldAgent != nil {
		firstName = v.oldAgent.FirstName
		email = v.oldAgent.Email
		phone = v.oldAgent.Phone
	}

	// Override pre-existing values with anything sent in for the update/create:
	if v.newAgent.FirstName != nil {
		firstName = v.newAgent.FirstName
	}
	if v.newAgent.Email != nil {
		email = v.newAgent.Email
	}
	if v.newAgent.Phone != nil {
		phone = v.newAgent.Phone
	}

	// Check that we have something in the FirstName field:
	if firstName == nil || *firstName == "" {
		v.verrs.Add("firstName", "cannot be blank")
	}

	// Check that we have one method of contacting the agent:
	if (email == nil || *email == "") && (phone == nil || *phone == "") {
		v.verrs.Add("contactInfo", "agent must have at least one contact method provided")
	}

	return nil
}

// getVerrs looks for any validation errors and returns a formatted InvalidInputError if any are found.
// Should only be called after the other check methods have been called.
func (v *AgentValidationData) getVerrs() error {
	if v.verrs.HasAny() {
		return services.NewInvalidInputError(v.newAgent.ID, nil, v.verrs, "Invalid input found while validating the agent.")
	}

	return nil
}

// setNewMTOAgent compares newAgent and oldAgent and updates a new MTOAgent instance with all data
// (changed and unchanged) filled in. Does not return an error, data must be checked for validation before this step.
func (v *AgentValidationData) setNewMTOAgent() *models.MTOAgent {
	agent := v.newAgent
	if v.oldAgent != nil {
		agent = *v.oldAgent
	}

	if v.newAgent.MTOAgentType != "" {
		agent.MTOAgentType = v.newAgent.MTOAgentType
	}
	if v.newAgent.FirstName != nil {
		agent.FirstName = v.newAgent.FirstName

		if *v.newAgent.FirstName == "" {
			agent.FirstName = nil
		}
	}
	if v.newAgent.LastName != nil {
		agent.LastName = v.newAgent.LastName

		if *v.newAgent.LastName == "" {
			agent.LastName = nil
		}
	}
	if v.newAgent.Email != nil {
		agent.Email = v.newAgent.Email

		if *v.newAgent.Email == "" {
			agent.Email = nil
		}
	}
	if v.newAgent.Phone != nil {
		agent.Phone = v.newAgent.Phone

		if *v.newAgent.Phone == "" {
			agent.Phone = nil
		}
	}

	return &agent
}
