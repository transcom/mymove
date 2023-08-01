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

		// currentServiceItemID := paymentRequest.PaymentServiceItems[0].MTOServiceItemID
		// reserviceCode := paymentRequest.PaymentServiceItems[0].MTOServiceItem.ReService.Code
		// status := paymentRequest.PaymentServiceItems[0].Status

		//get all the payment requests that exist for a shipment
		//append
		//mtoServiceItemID
		// loop through all requests to find if a payment service item already exists (same one as the one thats being created)
		// if exists --> check status if requested or paid --> conflict error

		shipment, err := mtoshipment.FindShipment(appCtx, *shipmentID, "MoveTaskOrder", "MoveTaskOrder.PaymentRequests", "MoveTaskOrder.PaymentRequests.PaymentServiceItems", "MoveTaskOrder.PaymentRequests.PaymentServiceItems.MTOServiceItem", "MoveTaskOrder.PaymentRequests.PaymentServiceItems.MTOServiceItem.ReService.Code")
		if err != nil {
			return err
		}

		var existingPaymentRequests = shipment.MoveTaskOrder.PaymentRequests
		newPaymentServiceItems := paymentRequest.PaymentServiceItems

		// var existingShipments = paymentRequest.MoveTaskOrder.MTOShipments
		//paymentRequest.MoveTaskOrder.MTOShipment.ID
		// return apperror.NewConflictError(paymentRequest.ID, "Conflict Error: Payment Request for Service Item is already paid or requested")

		// return apperror.NewConflictError(paymentRequest.ID, "Conflict Error: Payment Request for Service Item is already paid or requested")
		if len(existingPaymentRequests) > 0 {
			for _, pr := range existingPaymentRequests {
				for _, existingPaymentServiceItem := range pr.PaymentServiceItems {

					// this vv check is needed for moves that have multiple shipments, we don't want to exclude
					// a paid/requested payment service item if it's associated with a different shipment
					// this check needs more work because currently it's coming back as always false
					// I believe it's an issue with how MTOService item is being eager loaded in
					if existingPaymentServiceItem.MTOServiceItem.MTOShipmentID.String() == shipmentID.String() {
						//   for _, existingShipment := range existingShipments {
						// if existingShipment.ID.String() == shipmentID.String() {
						for _, newPaymentServiceItem := range newPaymentServiceItems {
							if newPaymentServiceItem.MTOServiceItemID == existingPaymentServiceItem.MTOServiceItemID {
								// if newPaymentServiceItem.MTOServiceItem.ReService.Code != models.ReServiceCodeDDASIT && newPaymentServiceItem.MTOServiceItem.ReService.Code != models.ReServiceCodeDOASIT {
								if (newPaymentServiceItem.MTOServiceItem.ReService.Code != models.ReServiceCodeDDASIT && newPaymentServiceItem.MTOServiceItem.ReService.Code != models.ReServiceCodeDOASIT) && (existingPaymentServiceItem.Status == models.PaymentServiceItemStatusRequested || existingPaymentServiceItem.Status == models.PaymentServiceItemStatusPaid) {
									// need to add back the exception for DDA and DOASIT
									return apperror.NewConflictError(pr.ID, "Conflict Error: Payment Request for Service Item is already paid or requested")
								}
								// }
							}
						}
					}
				}
			}
		}

		return nil
	})
}
