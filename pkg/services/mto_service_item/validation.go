package mtoserviceitem

import (
	"fmt"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"
)

// Validator is the interface for the various validations we might want to
// define.
type validator interface {
	Validate(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem, eTag string) error
}

type validatorFunc func(appcontext.AppContext, *models.MTOServiceItem, string) error

func (fn validatorFunc) Validate(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem, eTag string) error {
	return fn(appCtx, serviceItem, eTag)
}

func validateServiceItem(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem, eTag string, checks ...validator) (result error) {
	verrs := validate.NewErrors()
	for _, checker := range checks {
		if err := checker.Validate(appCtx, serviceItem, eTag); err != nil {
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
		result = apperror.NewInvalidInputError(serviceItem.ID, nil, verrs, "")
	}
	return result
}

func checkMoveStatus() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem, _ string) error {
		verrs := validate.NewErrors()
		move := serviceItem.MoveTaskOrder

		if move.Status != models.MoveStatusAPPROVED && move.Status != models.MoveStatusAPPROVALSREQUESTED {
			message := fmt.Sprintf("Cannot approve or reject a service item if the move's status is neither Approved nor Approvals Requested. The current status for the move with ID %s is %s", move.ID, move.Status)
			verrs.Add("status", message)
		}

		return verrs
	})
}

func checkETag() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem, eTag string) error {
		existingETag := etag.GenerateEtag(serviceItem.UpdatedAt)
		if existingETag != eTag {
			return apperror.NewPreconditionFailedError(serviceItem.ID, query.StaleIdentifierError{StaleIdentifier: eTag})
		}
		return nil
	})
}
