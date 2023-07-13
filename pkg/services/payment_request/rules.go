package paymentrequest

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

// verify that the MoveTaskOrderID on the payment request is not a nil uuid
func checkMTOIDField() paymentRequestValidator {
	return paymentRequestValidatorFunc(func(_ appcontext.AppContext, paymentRequest models.PaymentRequest, oldPaymentRequest *models.PaymentRequest) error {
		// Verify that the MTO ID exists
		if paymentRequest.MoveTaskOrderID == uuid.Nil {
			return apperror.NewInvalidCreateInputError(nil, "Invalid Create Input Error: MoveTaskOrderID is required on PaymentRequest create")
		}

		return nil
	})
}

func checkMTOIDMatchesServiceItemMTOID() paymentRequestValidator {
	return paymentRequestValidatorFunc(func(_ appcontext.AppContext, paymentRequest models.PaymentRequest, oldPaymentRequest *models.PaymentRequest) error {
		var paymentRequestServiceItems = paymentRequest.PaymentServiceItems
		for _, paymentRequestServiceItem := range paymentRequestServiceItems {
			if paymentRequest.MoveTaskOrderID != paymentRequestServiceItem.MTOServiceItem.MoveTaskOrderID && paymentRequestServiceItem.MTOServiceItemID != uuid.Nil {
				return apperror.NewConflictError(paymentRequestServiceItem.MTOServiceItem.MoveTaskOrderID, "Conflict Error: Payment Request MoveTaskOrderID does not match Service Item MoveTaskOrderID")
			}
		}
		return nil
	})
}

func checkStatusOfPaymentRequest() paymentRequestValidator {
	return paymentRequestValidatorFunc(func(_ appcontext.AppContext, paymentRequest models.PaymentRequest, oldPaymentRequest *models.PaymentRequest) error {
		var paymentRequestServiceItems = paymentRequest.PaymentServiceItems
		for _, paymentRequestServiceItem := range paymentRequestServiceItems {
			// fmt.Println(paymentRequest.PaymentServiceItems[0].MTOServiceItemID)
			// fmt.Println(paymentRequestServiceItem.MTOServiceItemID)
			// if paymentRequest.PaymentServiceItems[0].MTOServiceItemID == paymentRequestServiceItem.MTOServiceItemID && (paymentRequestServiceItem.PaymentRequest.Status == models.PaymentRequestStatusPending || paymentRequestServiceItem.PaymentRequest.Status == models.PaymentRequestStatusPaid) {
			if paymentRequest.PaymentServiceItems[0].MTOServiceItemID == paymentRequestServiceItem.MTOServiceItemID && paymentRequestServiceItem.MTOServiceItemID != uuid.Nil {

				// if paymentRequestServiceItem.Status == models.PaymentServiceItemStatusRequested || paymentRequestServiceItem.Status == models.PaymentServiceItemStatusPaid {

				// 	return apperror.NewConflictError(paymentRequestServiceItem.MTOServiceItemID, "Conflict Error: Payment Request for service item is already paid or requested")
				// }
				// if paymentRequestServiceItem.PaymentRequest.Status == models.PaymentRequestStatusPending || paymentRequestServiceItem.PaymentRequest.Status == models.PaymentRequestStatusPaid {
				// }
				// if paymentRequestServiceItem.PaymentRequest.Status == models.PaymentRequestStatusPending || paymentRequestServiceItem.PaymentRequest.Status == models.PaymentRequestStatusPaid {
				return apperror.NewConflictError(paymentRequestServiceItem.MTOServiceItemID, "Conflict Error: Payment Request for service item is already paid or requested")
				// }
				// if paymentRequestServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeDDASIT || paymentRequestServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeDOASIT {
				//     //   newStart, newEnd, newErr := getStartAndEndParams(paymentRequestServiceItem.PaymentServiceItemParams)
				// 	newStart, newEnd, _ := getStartAndEndParams(paymentRequestServiceItem.PaymentServiceItemParams)
				// 	if paymentRequestServiceItem.Status == models.PaymentServiceItemStatusRequested || paymentRequestServiceItem.Status == models.PaymentServiceItemStatusPaid {
				// 		// start, end, err := getStartAndEndParams(paymentRequestServiceItem.PaymentServiceItemParams)
				// 		start, end, _ := getStartAndEndParams(paymentRequestServiceItem.PaymentServiceItemParams)
				// 		fmt.Println(start, end, newStart, newEnd)

				// 		// inTimeFrame := time.After(start) && time.Before(end)
				// 	}
			}
		}
		return nil
	})
}
