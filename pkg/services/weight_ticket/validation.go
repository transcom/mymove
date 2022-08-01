package weightticket

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

// weightTicketValidator defines the interface for checking business rules for a weightTicket
type weightTicketValidator interface {
	// The newWeightTicket is optional, as create requires no payload.
	// The originalWeightTicket is optional, so it's a pointer type.
	Validate(appCtx appcontext.AppContext, newWeightTicket *models.WeightTicket, originalWeightTicket *models.WeightTicket) error
}

// weightTicketValidatorFunc is an adapter that will convert a function into an implementation of weightTicketValidator
type weightTicketValidatorFunc func(appcontext.AppContext, *models.WeightTicket, *models.WeightTicket) error

// Validate fulfills the weightTicketValidator interface
func (fn weightTicketValidatorFunc) Validate(appCtx appcontext.AppContext, newWeightTicket *models.WeightTicket, originalWeightTicket *models.WeightTicket) error {
	return fn(appCtx, newWeightTicket, originalWeightTicket)
}

func validateWeightTicket(appCtx appcontext.AppContext, newWeightTicket *models.WeightTicket, originalWeightTicket *models.WeightTicket, checks ...weightTicketValidator) error {
	verrs := validate.NewErrors()

	for _, check := range checks {
		if err := check.Validate(appCtx, newWeightTicket, originalWeightTicket); err != nil {
			switch e := err.(type) {
			case *validate.Errors:
				// Accumulate all validation errors
				verrs.Append(e)
			default:
				// Non-validation errors have priority and short-circuit doing any further checks
				return err
			}
		}
	}

	if verrs.HasAny() {
		currentID := uuid.Nil
		if newWeightTicket != nil {
			currentID = newWeightTicket.ID
		}
		return apperror.NewInvalidInputError(currentID, nil, verrs, "Invalid input found while validating the weight ticket.")
	}

	return nil
}
