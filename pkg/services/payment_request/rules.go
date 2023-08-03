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
		shipmentID := paymentRequest.PaymentServiceItems[0].MTOServiceItem.MTOShipmentID

		shipment, err := mtoshipment.FindShipment(appCtx, *shipmentID,
			"MoveTaskOrder",
			"MoveTaskOrder.PaymentRequests",
			"MoveTaskOrder.PaymentRequests.PaymentServiceItems",
			"MoveTaskOrder.PaymentRequests.PaymentServiceItems.MTOServiceItem",
			"MoveTaskOrder.PaymentRequests.PaymentServiceItems.MTOServiceItem.ReService.Code",
		)
		if err != nil {
			return err
		}

		var existingPaymentRequests = shipment.MoveTaskOrder.PaymentRequests
		newPaymentServiceItems := paymentRequest.PaymentServiceItems

		for _, pr := range existingPaymentRequests {
			if pr.Status == models.PaymentRequestStatusReviewedAllRejected || pr.Status == models.PaymentRequestStatusDeprecated {
				continue
			}
			for _, existingPaymentServiceItem := range pr.PaymentServiceItems {

				if existingPaymentServiceItem.MTOServiceItem.MTOShipmentID.String() != shipmentID.String() {
					continue
				}
				for _, newPaymentServiceItem := range newPaymentServiceItems {
					if newPaymentServiceItem.MTOServiceItemID != existingPaymentServiceItem.MTOServiceItemID {
						continue
					}
					if newPaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeDDASIT || newPaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeDOASIT {
						continue
					}
					if existingPaymentServiceItem.Status == models.PaymentServiceItemStatusRequested || existingPaymentServiceItem.Status == models.PaymentServiceItemStatusPaid {
						return apperror.NewConflictError(pr.ID, "Conflict Error: Payment Request for Service Item is already paid or requested")
					}
				}
			}
		}
		return nil
	})
}
