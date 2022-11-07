package progearweightticket

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

// progearWeightTicketValidator defines the interface for checking business rules for a progearWeightTicket
type progearWeightTicketValidator interface {
	// The newProgearWeightTicket is optional, as create requires no payload.
	// The originalProgearWeightTicket is optional, so it's a pointer type.
	Validate(appCtx appcontext.AppContext, newProgearWeightTicket *models.ProgearWeightTicket, originalProgearWeightTicket *models.ProgearWeightTicket) error
}

// progearWeightTicketValidatorFunc is an adapter that will convert a function into an implementation of progearWeightTicketValidator
type progearWeightTicketValidatorFunc func(appcontext.AppContext, *models.ProgearWeightTicket, *models.ProgearWeightTicket) error

// Validate fulfills the progearWeightTicketValidator interface
func (fn progearWeightTicketValidatorFunc) Validate(appCtx appcontext.AppContext, newProgearWeightTicket *models.ProgearWeightTicket, originalProgearWeightTicket *models.ProgearWeightTicket) error {
	return fn(appCtx, newProgearWeightTicket, originalProgearWeightTicket)
}

func validateProgearWeightTicket(appCtx appcontext.AppContext, newProgearWeightTicket *models.ProgearWeightTicket, originalProgearWeightTicket *models.ProgearWeightTicket, checks ...progearWeightTicketValidator) error {
	verrs := validate.NewErrors()

	for _, check := range checks {
		if err := check.Validate(appCtx, newProgearWeightTicket, originalProgearWeightTicket); err != nil {
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
		if newProgearWeightTicket != nil {
			currentID = newProgearWeightTicket.ID
		}
		return apperror.NewInvalidInputError(currentID, nil, verrs, "Invalid input found while validating the weight ticket.")
	}

	return nil
}
