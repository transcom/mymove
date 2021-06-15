package mtoagent

import (
	"context"
	"fmt"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

var basicChecks = []mtoAgentValidator{
	checkShipmentID(),
	checkAgentID(),
}

var primeChecks = append(
	basicChecks,
	checkContactInfo(),
	checkAgentType(),
)

// checkShipmentID checks that the user didn't attempt to change the agent's Shipment ID
func checkShipmentID() mtoAgentValidator {
	return mtoAgentValidatorFunc(func(_ context.Context, newAgent models.MTOAgent, oldAgent *models.MTOAgent, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if oldAgent == nil {
			if newAgent.MTOShipmentID == uuid.Nil {
				verrs.Add("mtoShipmentID", "Shipment ID is required")
			}
		} else {
			if newAgent.MTOShipmentID != uuid.Nil && newAgent.MTOShipmentID != oldAgent.MTOShipmentID {
				verrs.Add("mtoShipmentID", "cannot be updated")
			}
		}
		return verrs
	})
}

// checkAgentID checks that the new agent's ID matches the old agent's ID (or is nil)
func checkAgentID() mtoAgentValidator {
	return mtoAgentValidatorFunc(func(_ context.Context, newAgent models.MTOAgent, oldAgent *models.MTOAgent, shipment *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if oldAgent == nil {
			if newAgent.ID != uuid.Nil {
				verrs.Add("ID", "cannot manually set a new agent's UUID")
			}
		} else {
			if newAgent.ID != oldAgent.ID {
				return services.NewImplementationError(
					fmt.Sprintf("the newAgent ID (%s) must match oldAgent ID (%s).", newAgent.ID, oldAgent.ID),
				)
			}
		}
		return verrs
	})
}

const maxAgents = 2

// checkAgentType checks that there is, at most, one RELEASING and one RECEIVING agent (each) on a shipment.
// It also checks that we're not adding more than the max number of agents.
// NOTE: You need to make sure MTOShipment.MTOAgents is populated for the results of this check to be accurate.
func checkAgentType() mtoAgentValidator {
	return mtoAgentValidatorFunc(func(_ context.Context, newAgent models.MTOAgent, oldAgent *models.MTOAgent, shipment *models.MTOShipment) error {
		if shipment == nil {
			return services.NewImplementationError(
				fmt.Sprintf("mtoAgent validation needs the shipment data in order to validate the AgentType for newAgent: %s", newAgent.ID),
			)
		}

		// We don't need to check the MTOAgentType if it's not being updated, or if there are no other agents:
		if newAgent.MTOAgentType == "" || shipment.MTOAgents == nil {
			return nil
		}

		agents := shipment.MTOAgents
		if len(agents) >= maxAgents && newAgent.ID == uuid.Nil { // a nil UUID here means we're creating a new agent
			return services.NewConflictError(
				shipment.ID, fmt.Sprintf("This shipment already has %d agents - no more can be added.", maxAgents))
		}

		for _, agent := range agents {
			if agent.ID != uuid.Nil && agent.ID == newAgent.ID {
				// Since we're looking at the same agent, there's no need to check anything else here.
				// Note that we might also have other agents with nil UUIDs if we're dealing with a bulk create -
				// we DON'T want to skip validation in that case.
				continue
			}

			if agent.MTOAgentType == newAgent.MTOAgentType {
				return services.NewConflictError(
					newAgent.ID, fmt.Sprintf("There is already an agent with type %s on the shipment", newAgent.MTOAgentType))
			}
		}
		return nil
	})
}

// checkContactInfo checks that the new agent has the minimum required contact info: First Name and one of Email or Phone
func checkContactInfo() mtoAgentValidator {
	return mtoAgentValidatorFunc(func(_ context.Context, newAgent models.MTOAgent, oldAgent *models.MTOAgent, shipment *models.MTOShipment) error {
		verrs := validate.NewErrors()

		var firstName *string
		var email *string
		var phone *string

		// Set any pre-existing values as the baseline:
		if oldAgent != nil {
			firstName = oldAgent.FirstName
			email = oldAgent.Email
			phone = oldAgent.Phone
		}

		// Override pre-existing values with anything sent in for the update/create:
		if newAgent.FirstName != nil {
			firstName = newAgent.FirstName
		}
		if newAgent.Email != nil {
			email = newAgent.Email
		}
		if newAgent.Phone != nil {
			phone = newAgent.Phone
		}

		// Check that we have something in the FirstName field:
		if firstName == nil || *firstName == "" {
			verrs.Add("firstName", "cannot be blank")
		}

		// Check that we have one method of contacting the agent:
		if (email == nil || *email == "") && (phone == nil || *phone == "") {
			verrs.Add("contactInfo", "agent must have at least one contact method provided")
		}
		return verrs
	})
}

// checkPrimeAvailability returns a type that checks that agent is connected to a Prime-available Shipment
func checkPrimeAvailability(checker services.MoveTaskOrderChecker) mtoAgentValidator {
	return mtoAgentValidatorFunc(func(ctx context.Context, newAgent models.MTOAgent, oldAgent *models.MTOAgent, shipment *models.MTOShipment) error {
		if shipment == nil {
			return services.NewNotFoundError(newAgent.ID, "while looking for Prime-available Shipment")
		}

		isAvailable, err := checker.MTOAvailableToPrime(shipment.MoveTaskOrderID)
		if !isAvailable || err != nil {
			return services.NewNotFoundError(
				newAgent.ID, fmt.Sprintf("while looking for Prime-available Shipment with id: %s", shipment.ID))
		}
		return nil
	})
}
