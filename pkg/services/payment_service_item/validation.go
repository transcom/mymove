package paymentserviceitem

import (
	"fmt"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"
)

type validator interface {
	Validate(appCtx appcontext.AppContext, paymentServiceItem *models.PaymentServiceItem,
		desiredStatus models.PaymentServiceItemStatus, rejectionReason *string, eTag string) error
}

type validatorFunc func(appcontext.AppContext, *models.PaymentServiceItem, models.PaymentServiceItemStatus, *string, string) error

func (fn validatorFunc) Validate(appCtx appcontext.AppContext, paymentServiceItem *models.PaymentServiceItem,
	desiredStatus models.PaymentServiceItemStatus, rejectionReason *string, eTag string) error {
	return fn(appCtx, paymentServiceItem, desiredStatus, rejectionReason, eTag)
}

func validatePaymentServiceItem(appCtx appcontext.AppContext, paymentServiceItem *models.PaymentServiceItem,
	desiredStatus models.PaymentServiceItemStatus, rejectionReason *string, eTag string, checks ...validator) (result error) {
	verrs := validate.NewErrors()
	for _, checker := range checks {
		if err := checker.Validate(appCtx, paymentServiceItem, desiredStatus, rejectionReason, eTag); err != nil {
			switch e := err.(type) {
			case *validate.Errors:
				// gather validation errors
				verrs.Append(e)
			default:
				return err
			}
		}
	}
	if verrs.HasAny() {
		result = apperror.NewInvalidInputError(paymentServiceItem.ID, nil, verrs, "")
	}
	return result
}

func checkETag() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, paymentServiceItem *models.PaymentServiceItem,
		_ models.PaymentServiceItemStatus, rejectionReason *string, eTag string) error {
		existingETag := etag.GenerateEtag(paymentServiceItem.UpdatedAt)
		if existingETag != eTag {
			return apperror.NewPreconditionFailedError(paymentServiceItem.ID,
				query.StaleIdentifierError{StaleIdentifier: eTag})
		}
		return nil
	})
}

// Checks to make sure that a rejection reason was provided if the status
func rejectionRequiresRejectionReason() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, paymentServiceItem *models.PaymentServiceItem,
		desiredStatus models.PaymentServiceItemStatus, rejectionReason *string, _ string) error {
		verrs := validate.NewErrors()

		if desiredStatus == models.PaymentServiceItemStatusDenied &&
			rejectionReason == nil {
			message := fmt.Sprintf("Cannot reject a payment service item without providing a rejection reason."+
				"The current payment service item ID is %s with status %s", paymentServiceItem.ID, paymentServiceItem.Status)
			verrs.Add("status", message)
		}
		return verrs
	})

}
