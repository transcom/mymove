package mtoshipment

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

type validator interface {
	Validate(appCtx appcontext.AppContext, newer *models.MTOShipment, older *models.MTOShipment) error
}

type validatorFunc func(appcontext.AppContext, *models.MTOShipment, *models.MTOShipment) error

func (fn validatorFunc) Validate(appCtx appcontext.AppContext, newer *models.MTOShipment, older *models.MTOShipment) error {
	return fn(appCtx, newer, older)
}

func validateShipment(appCtx appcontext.AppContext, newer *models.MTOShipment, older *models.MTOShipment, checks ...validator) (result error) {
	verrs := validate.NewErrors()
	for _, checker := range checks {
		if err := checker.Validate(appCtx, newer, older); err != nil {
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
		result = apperror.NewInvalidInputError(newer.ID, nil, verrs, "Invalid input found while updating the shipment.")
	}
	return result
}
