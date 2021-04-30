package mtoagentvalidate

import (
	"fmt"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// AgentValidationData represents the data needed to validate an Agent before a create/update action.
// It will be set in the service object that calls the validator.
type AgentValidationData struct {
	NewAgent            models.MTOAgent
	OldAgent            *models.MTOAgent // not required for create
	Shipment            *models.MTOShipment
	AvailabilityChecker services.MoveTaskOrderChecker
	Verrs               *validate.Errors
}

// checkShipmentID checks that the user didn't attempt to change the agent's Shipment ID
func (v *AgentValidationData) checkShipmentID() error {
	if v.OldAgent == nil {
		if v.NewAgent.MTOShipmentID == uuid.Nil {
			v.Verrs.Add("mtoShipmentID", "Shipment ID is required")
		}
	} else {
		if v.NewAgent.MTOShipmentID != uuid.Nil && v.NewAgent.MTOShipmentID != v.OldAgent.MTOShipmentID {
			v.Verrs.Add("mtoShipmentID", "cannot be updated")
		}
	}
	return nil
}

// checkAgentID checks that the new agent's ID matches the old agent's ID (or is nil)
func (v *AgentValidationData) checkAgentID() error {
	if v.OldAgent == nil {
		if v.NewAgent.ID != uuid.Nil {
			v.Verrs.Add("ID", "cannot manually set a new agent's UUID")
		}
	} else {
		if v.NewAgent.ID != v.OldAgent.ID {
			return services.NewImplementationError(
				fmt.Sprintf("In AgentValidationData, the NewAgent's ID (%s) must match OldAgent's ID (%s).", v.NewAgent.ID, v.OldAgent.ID),
			)
		}
	}
	return nil
}

// checkAgentType checks that there is, at most, one RELEASING and one RECEIVING agent (each) on a shipment.
// It also checks that we're not adding more than the max number of agents.
// NOTE: You need to make sure MTOShipment.MTOAgents is populated for the results of this check to be accurate.
func (v *AgentValidationData) checkAgentType() error {
	if v.NewAgent.MTOAgentType == "" {
		return nil // We don't need to check the MTOAgentType if it's not being updated
	}

	agents := v.Shipment.MTOAgents
	maxAgents := 2
	if len(agents) >= maxAgents && v.NewAgent.ID == uuid.Nil { // a nil UUID here means we're creating a new agent
		return services.NewConflictError(
			v.Shipment.ID, fmt.Sprintf("This shipment already has %d agents - no more can be added.", maxAgents))
	}

	for _, agent := range agents {
		if agent.ID == v.NewAgent.ID {
			continue // since we're looking at the same agent, there's no need to check anything else here
		}

		if agent.MTOAgentType == v.NewAgent.MTOAgentType {
			return services.NewConflictError(
				v.NewAgent.ID, fmt.Sprintf("There is already an agent with type %s on the shipment", v.NewAgent.MTOAgentType))
		}
	}
	return nil
}

// checkPrimeAvailability checks that agent is connected to a Prime-available Shipment
func (v *AgentValidationData) checkPrimeAvailability() error {
	if v.Shipment == nil {
		return services.NewNotFoundError(v.NewAgent.ID, "while looking for Prime-available Shipment")
	}
	isAvailable, err := v.AvailabilityChecker.MTOAvailableToPrime(v.Shipment.MoveTaskOrderID)

	if !isAvailable || err != nil {
		return services.NewNotFoundError(
			v.NewAgent.ID, fmt.Sprintf("while looking for Prime-available Shipment with id: %s", v.Shipment.ID))
	}
	return nil
}

// checkContactInfo checks that the new agent has the minimum required contact info: First Name and one of Email or Phone
func (v *AgentValidationData) checkContactInfo() error {
	var firstName *string
	var email *string
	var phone *string

	// Set any pre-existing values as the baseline:
	if v.OldAgent != nil {
		firstName = v.OldAgent.FirstName
		email = v.OldAgent.Email
		phone = v.OldAgent.Phone
	}

	// Override pre-existing values with anything sent in for the update/create:
	if v.NewAgent.FirstName != nil {
		firstName = v.NewAgent.FirstName
	}
	if v.NewAgent.Email != nil {
		email = v.NewAgent.Email
	}
	if v.NewAgent.Phone != nil {
		phone = v.NewAgent.Phone
	}

	// Check that we have something in the FirstName field:
	if firstName == nil || *firstName == "" {
		v.Verrs.Add("firstName", "cannot be blank")
	}

	// Check that we have one method of contacting the agent:
	if (email == nil || *email == "") && (phone == nil || *phone == "") {
		v.Verrs.Add("contactInfo", "agent must have at least one contact method provided")
	}
	return nil
}

// getVerrs looks for any validation errors and returns a formatted InvalidInputError if any are found.
// Should only be called after the other check methods have been called.
func (v *AgentValidationData) getVerrs() error {
	if v.Verrs.HasAny() {
		return services.NewInvalidInputError(v.NewAgent.ID, nil, v.Verrs, "Invalid input found while validating the agent.")
	}
	return nil
}

// setFullAgent compares NewAgent and OldAgent and updates a new MTOAgent instance with all data
// (changed and unchanged) filled in. Does not return an error, data must be checked for validation before this step.
func (v *AgentValidationData) setFullAgent() *models.MTOAgent {
	agent := v.NewAgent
	if v.OldAgent != nil {
		agent = *v.OldAgent
	}

	if v.NewAgent.MTOAgentType != "" {
		agent.MTOAgentType = v.NewAgent.MTOAgentType
	}
	agent.FirstName = services.SetOptionalStringField(v.NewAgent.FirstName, agent.FirstName)
	agent.LastName = services.SetOptionalStringField(v.NewAgent.LastName, agent.LastName)
	agent.Email = services.SetOptionalStringField(v.NewAgent.Email, agent.Email)
	agent.Phone = services.SetOptionalStringField(v.NewAgent.Phone, agent.Phone)

	return &agent
}
