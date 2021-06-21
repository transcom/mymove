package mtoshipment

import (
	"context"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type validator interface {
	Validate(c context.Context, newer *models.MTOShipment, older *models.MTOShipment) error
}

type validatorFunc func(context.Context, *models.MTOShipment, *models.MTOShipment) error

func (fn validatorFunc) Validate(c context.Context, newer *models.MTOShipment, older *models.MTOShipment) error {
	return fn(c, newer, older)
}

func validateShipment(ctx context.Context, newer *models.MTOShipment, older *models.MTOShipment, checks ...validator) (result error) {
	verrs := validate.NewErrors()
	for _, checker := range checks {
		if err := checker.Validate(ctx, newer, older); err != nil {
			switch e := err.(type) {
			case *validate.Errors:
				// accumulate validation errors
				verrs.Append(e)
				// if true {
				// 	return err
				// }
			default:
				// non-validation errors have priority,
				// and short-circuit doing any further checks
				return err
			}
		}
	}
	if verrs.HasAny() {
		result = services.NewInvalidInputError(newer.ID, nil, verrs, "Invalid input found while updating the shipment.")
	}
	return result
}

func checkStatus() validator {
	return validatorFunc(func(_ context.Context, newer *models.MTOShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if newer.Status != "" && newer.Status != models.MTOShipmentStatusDraft && newer.Status != models.MTOShipmentStatusSubmitted {
			verrs.Add("status", "can only update status to DRAFT or SUBMITTED. use UpdateMTOShipmentStatus for other status updates")
		}
		return verrs
	})
}
