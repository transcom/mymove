package progearweightticket

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

type progearWeightTicketValidator interface {
	Validate(appCtx appcontext.AppContext, newProgearWeightTicket *models.ProgearWeightTicket, oldProgearWeightTicket *models.ProgearWeightTicket) error
}

type progearWeightTicketValidatorFunc func(appCtx appcontext.AppContext, newProgearWeightTicket *models.ProgearWeightTicket, oldProgearWeightTicket *models.ProgearWeightTicket) error

func (fn progearWeightTicketValidatorFunc) Validate(appCtx appcontext.AppContext, newProgearWeightTicket *models.ProgearWeightTicket, oldProgearWeightTicket *models.ProgearWeightTicket) error {
	return fn(appCtx, newProgearWeightTicket, oldProgearWeightTicket)
}

func validateProgearWeightTicket(appCtx appcontext.AppContext, newProgearWeightTicket *models.ProgearWeightTicket, oldProgearWeightTicket *models.ProgearWeightTicket, checks ...progearWeightTicketValidator) error {
	verrs := validate.NewErrors()

	for _, check := range checks {
		if err := check.Validate(appCtx, newProgearWeightTicket, oldProgearWeightTicket); err != nil {
			switch e := err.(type) {
			case *validate.Errors:
				verrs.Append(e)
			default:
				return err
			}
		}
	}

	if verrs.HasAny() {
		var currentID uuid.UUID
		if newProgearWeightTicket != nil {
			currentID = newProgearWeightTicket.ID
		}
		return apperror.NewInvalidInputError(currentID, nil, verrs, "")
	}

	return nil
}
