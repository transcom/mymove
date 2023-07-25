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

// func findPaymentRequestStatus(appCtx appcontext.AppContext, paymentRequestID uuid.UUID) (models.PaymentRequest, error) {
// 	var paymentRequest models.PaymentRequest

// 	// err := appCtx.DB().Eager().Find(&paymentRequest, paymentRequestID)
// 	err := appCtx.DB().Eager(
// 		"PaymentServiceItems",
// 	).Find(&paymentRequest, paymentRequestID)

// 	if err != nil {
// 		switch err {
// 		case sql.ErrNoRows:
// 			return models.PaymentRequest{}, apperror.NewNotFoundError(paymentRequestID, "looking for PaymentRequest")
// 		default:
// 			return models.PaymentRequest{}, apperror.NewQueryError("PaymentRequest", err, "")
// 		}
// 	}

// 	return paymentRequest, err
// }

// 1) get mtoserviceitemid from the new payment request getting created
// 2) look to find if there is already an existing payment request for that service item
// 3) if YES --> check status of existing payment request - pending or paid
// 4) pending/paid --> conflict error
// 5) already reviewed --> allow creation of payment request

// ANOTHER OPTION --> FIND PAYMENT REQUEST FROM DATABASE (supportapi/payment_request.go)
// 1) after finding matching mtoserviceitem, get paymentrequest id from that service Item
// 2) use function  findPaymentRequestStatus to find status of above payment request id , searches thru DB
// 3) if status of existing payment request = pending or paid --> conflict error

// prevent creating new payment requests for service items that already been paid or requested
func checkStatusOfExistingPaymentRequest() paymentRequestValidator {
	return paymentRequestValidatorFunc(func(appCtx appcontext.AppContext, paymentRequest models.PaymentRequest, oldPaymentRequest *models.PaymentRequest) error {
		var paymentRequestServiceItems = paymentRequest.PaymentServiceItems

		// fetcher := NewPaymentRequestFetcher()

		for _, paymentRequestServiceItem := range paymentRequestServiceItems {
			// if paymentRequest.PaymentServiceItems[0].MTOServiceItemID == paymentRequestServiceItem.MTOServiceItemID {
			// var paymentRequestIDFromServiceItem = paymentRequestServiceItem.PaymentRequestID
			// foundPaymentRequest, err := fetcher.FetchPaymentRequest(appCtx, paymentRequestIDFromServiceItem)

			// if err != nil {
			// 	msg := fmt.Sprintf("Error finding Payment Request for status update with ID: %s", paymentRequestIDFromServiceItem)
			// 	appCtx.Logger().Error(msg, zap.Error(err))
			// 	// return paymentrequestop.NewUpdatePaymentRequestStatusNotFound().WithPayload(
			// 	// 	payloads.ClientError(handlers.NotFoundMessage, msg, h.GetTraceIDFromRequest(params.HTTPRequest))), err
			// }

			// status := foundPaymentRequest.Status
			status := paymentRequestServiceItem.Status

			if status == models.PaymentServiceItemStatusRequested || status == models.PaymentServiceItemStatusPaid {
				return apperror.NewConflictError(paymentRequestServiceItem.MTOServiceItemID, "Conflict Error: Payment Request for Service Item is already paid or requested")
			}

			// }
		}

		return nil
	})
}
