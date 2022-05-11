package paymentrequest

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

type paymentRequestValidator interface {
	Validate(appCtx appcontext.AppContext, newPaymentRequest models.PaymentRequest, oldPaymentRequest *models.PaymentRequest) error
}

func validatePaymentRequest(appCtx appcontext.AppContext, newPaymentRequest models.PaymentRequest, oldPaymentRequest *models.PaymentRequest, checks ...paymentRequestValidator) (result error) {
	verrs := validate.NewErrors()
	for _, checker := range checks {
		if err := checker.Validate(appCtx, newPaymentRequest, oldPaymentRequest); err != nil {
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
		result = apperror.NewInvalidInputError(newPaymentRequest.ID, nil, verrs, "Invalid input found while validating the payment request.")
	}
	return result
}

type paymentRequestValidatorFunc func(appcontext.AppContext, models.PaymentRequest, *models.PaymentRequest) error

func (fn paymentRequestValidatorFunc) Validate(appCtx appcontext.AppContext, newPaymentRequest models.PaymentRequest, oldPaymentRequest *models.PaymentRequest) error {
	return fn(appCtx, newPaymentRequest, oldPaymentRequest)
}
