package mtoagent

import (
	"context"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// mtoAgentValidator describes a method for checking business requirements
type mtoAgentValidator interface {
	// Validate checks the newAgent for adherence to business rules. The
	// oldAgent parameter is expected to be nil in creation use cases.
	// It is safe to return a *validate.Errors with zero added errors as
	// a success case.
	Validate(c context.Context, newAgent models.MTOAgent, oldAgent *models.MTOAgent, shipment *models.MTOShipment) error
}

// mtoAgentValidatorFunc is an adapter type for converting a function into an implementation of mtoAgentValidator
type mtoAgentValidatorFunc func(context.Context, models.MTOAgent, *models.MTOAgent, *models.MTOShipment) error

// Validate fulfills the mtoAgentValidator interface
func (fn mtoAgentValidatorFunc) Validate(ctx context.Context, newer models.MTOAgent, older *models.MTOAgent, ship *models.MTOShipment) error {
	return fn(ctx, newer, older, ship)
}

// validateMTOAgent checks an MTOAgent against a passed-in set of business rule checks
// Validation errors are aggregated and returned to be reported on en masse. Other errors,
// of types not specifically associated with validations, are treated with higher priority
// and returned immediately, ignoring any accumulated validation errors and short circuiting
// the execution of any further mtoAgentValidator instances.
func validateMTOAgent(
	ctx context.Context,
	newAgent models.MTOAgent,
	oldAgent *models.MTOAgent,
	shipment *models.MTOShipment,
	checks ...mtoAgentValidator,
) (result error) {
	verrs := validate.NewErrors()
	for _, checker := range checks {
		if err := checker.Validate(ctx, newAgent, oldAgent, shipment); err != nil {
			switch e := err.(type) {
			case *validate.Errors:
				// accumulate validation errors
				verrs.Append(e)
			default:
				// non-validation errors have priority,
				// and short-circuit doing any further checks
				return err
			}
		}
	}
	if verrs.HasAny() {
		result = services.NewInvalidInputError(newAgent.ID, nil, verrs, "Invalid input found while validating the agent.")
	}
	return result
}

// mergeAgent compares NewAgent and OldAgent and updates a new MTOAgent instance with all data
// (changed and unchanged) filled in. Does not return an error, data must be checked for validation before this step.
func mergeAgent(newAgent models.MTOAgent, oldAgent *models.MTOAgent) *models.MTOAgent {
	if oldAgent == nil {
		return &newAgent
	}

	agent := *oldAgent

	if newAgent.MTOAgentType != "" {
		agent.MTOAgentType = newAgent.MTOAgentType
	}
	agent.FirstName = services.SetOptionalStringField(newAgent.FirstName, agent.FirstName)
	agent.LastName = services.SetOptionalStringField(newAgent.LastName, agent.LastName)
	agent.Email = services.SetOptionalStringField(newAgent.Email, agent.Email)
	agent.Phone = services.SetOptionalStringField(newAgent.Phone, agent.Phone)

	return &agent
}
