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
	if v.NewAgent.FirstName != nil {
		agent.FirstName = v.NewAgent.FirstName

		if *v.NewAgent.FirstName == "" {
			agent.FirstName = nil
		}
	}
	if v.NewAgent.LastName != nil {
		agent.LastName = v.NewAgent.LastName

		if *v.NewAgent.LastName == "" {
			agent.LastName = nil
		}
	}
	if v.NewAgent.Email != nil {
		agent.Email = v.NewAgent.Email

		if *v.NewAgent.Email == "" {
			agent.Email = nil
		}
	}
	if v.NewAgent.Phone != nil {
		agent.Phone = v.NewAgent.Phone

		if *v.NewAgent.Phone == "" {
			agent.Phone = nil
		}
	}

	return &agent
}
