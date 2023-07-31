package paymentrequest

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
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

// prevent creating new payment requests for service items that already been paid or requested
func checkStatusOfExistingPaymentRequest() paymentRequestValidator {
	return paymentRequestValidatorFunc(func(appCtx appcontext.AppContext, paymentRequest models.PaymentRequest, oldPaymentRequest *models.PaymentRequest) error {
		// var paymentRequestServiceItems = paymentRequest.PaymentServiceItems
		shipmentID := paymentRequest.PaymentServiceItems[0].MTOServiceItem.MTOShipmentID
		currentServiceItemID := paymentRequest.PaymentServiceItems[0].MTOServiceItemID
		reserviceCode := paymentRequest.PaymentServiceItems[0].MTOServiceItem.ReService.Code
		status := paymentRequest.PaymentServiceItems[0].Status

		//get all the payment requests that exist for a shipment
		//append
		//mtoServiceItemID
		// loop through all requests to find if a payment service item already exists (same one as the one thats being created)
		// if exists --> check status if requested or paid --> conflict error

		shipment, err := mtoshipment.FindShipment(appCtx, *shipmentID, "MoveTaskOrder.PaymentRequests")
		if err != nil {
			return err
		}

		var existingPaymentRequests = shipment.MoveTaskOrder.PaymentRequests

		if len(existingPaymentRequests) > 0 {
			for _, paymentRequest := range existingPaymentRequests {
				if paymentRequest.PaymentServiceItems[0].MTOServiceItemID == currentServiceItemID {
					if (reserviceCode != models.ReServiceCodeDDASIT && reserviceCode != models.ReServiceCodeDOASIT) && (status == models.PaymentServiceItemStatusRequested || status == models.PaymentServiceItemStatusPaid) {
						return apperror.NewConflictError(paymentRequest.PaymentServiceItems[0].MTOServiceItemID, "Conflict Error: Payment Request for Service Item is already paid or requested")
					}
				}
			}
		}
		// for _, paymentRequestServiceItem := range paymentRequestServiceItems {

		// 	status := paymentRequestServiceItem.Status
		// 	reserviceCode := paymentRequestServiceItem.MTOServiceItem.ReService.Code

		// 	if (reserviceCode != models.ReServiceCodeDDASIT && reserviceCode != models.ReServiceCodeDOASIT) && (status == models.PaymentServiceItemStatusRequested || status == models.PaymentServiceItemStatusPaid) {
		// 		return apperror.NewConflictError(paymentRequestServiceItem.MTOServiceItemID, "Conflict Error: Payment Request for Service Item is already paid or requested")
		// 	}
		// }
		return nil
	})
}
