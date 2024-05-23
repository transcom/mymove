package paymentrequest

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

// verify that the MoveTaskOrderID on the payment request is not a nil uuid
func checkMTOIDField() paymentRequestValidator {
	return paymentRequestValidatorFunc(func(_ appcontext.AppContext, paymentRequest models.PaymentRequest, _ *models.PaymentRequest) error {
		// Verify that the MTO ID exists
		if paymentRequest.MoveTaskOrderID == uuid.Nil {
			return apperror.NewInvalidCreateInputError(nil, "Invalid Create Input Error: MoveTaskOrderID is required on PaymentRequest create")
		}

		return nil
	})
}

func checkMTOIDMatchesServiceItemMTOID() paymentRequestValidator {
	return paymentRequestValidatorFunc(func(_ appcontext.AppContext, paymentRequest models.PaymentRequest, _ *models.PaymentRequest) error {
		var paymentRequestServiceItems = paymentRequest.PaymentServiceItems
		for _, paymentRequestServiceItem := range paymentRequestServiceItems {
			if paymentRequest.MoveTaskOrderID != paymentRequestServiceItem.MTOServiceItem.MoveTaskOrderID && paymentRequestServiceItem.MTOServiceItemID != uuid.Nil {
				return apperror.NewConflictError(paymentRequestServiceItem.MTOServiceItem.MoveTaskOrderID, "Conflict Error: Payment Request MoveTaskOrderID does not match Service Item MoveTaskOrderID")
			}
		}
		return nil
	})
}

// This rule enforces valid date inputs for payment request service items for additional days of SIT
func checkValidSitAddlDates() paymentRequestValidator {
	return paymentRequestValidatorFunc(func(_ appcontext.AppContext, pr models.PaymentRequest, _ *models.PaymentRequest) error {
		const format = "2006-01-02"

		for _, psi := range pr.PaymentServiceItems {
			var sitStartDate, sitEndDate *time.Time

			for _, param := range psi.PaymentServiceItemParams {
				// Look for SITPaymentRequestStart and then parse it
				if param.IncomingKey == models.ServiceItemParamNameSITPaymentRequestStart.String() && sitStartDate == nil {
					paramValue, err := time.Parse(format, param.Value)
					if err != nil {
						return apperror.NewInvalidCreateInputError(nil, "Invalid Create Input Error: SITPaymentRequestStart must be a valid date value of YYYY-MM-DD")
					}
					sitStartDate = &paramValue
					continue
				}
				// Look for SITPaymentRequestEnd and then parse it
				if param.IncomingKey == models.ServiceItemParamNameSITPaymentRequestEnd.String() && sitEndDate == nil {
					paramValue, err := time.Parse(format, param.Value)
					if err != nil {
						return apperror.NewInvalidCreateInputError(nil, "Invalid Create Input Error: SITPaymentRequestEnd must be a valid date value of YYYY-MM-DD")
					}
					sitEndDate = &paramValue
					continue
				}

			}
			// Check that both SITPaymentRequestStart and SITPaymentRequestEnd exist
			// If exist, compare dates to enforce rule
			if sitStartDate != nil && sitEndDate != nil {
				if sitStartDate.After(*sitEndDate) {
					return apperror.NewInvalidCreateInputError(nil, "Invalid Create Input Error: SITPaymentRequestStart must be a date that comes before SITPaymentRequestEnd")
				}
			}
		}
		return nil
	})
}
