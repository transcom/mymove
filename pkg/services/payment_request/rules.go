package paymentrequest

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	mtoFetcher "github.com/transcom/mymove/pkg/services/move_task_order"
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

		//new logic:1) determine if the incoming payment service item is MS or CS
		// a) filter through new service items to find which is MS/CS and create new list
		//if service item is MS OR CS then
		// 2) get all the payment requests for that specific move
		//3 loop through all payment requests payment service items for that move and find if MS or CS exist
		//4 if it exists, check if status is requested or paid
		//5 if yes, return conflict error
		//6 if 0 move level items, then continue w original function searching by shipmentID

		newPaymentServiceItems := paymentRequest.PaymentServiceItems
		moveID := paymentRequest.MoveTaskOrderID

		//create a new list for the filter
		var moveLevelItems []models.PaymentServiceItem

		searchParams := services.MoveTaskOrderFetcherParams{
			MoveTaskOrderID: moveID,
		}

		//fetching a move to then grab all payment requests for that move
		move, err := mtoFetcher.NewMoveTaskOrderFetcher().FetchMoveTaskOrder(appCtx, &searchParams)
		if err != nil {
			return err
		}

		allMovePaymentRequests := move.PaymentRequests

		//checking to see if new payment service items are MS or CS, then adding it to a filtered list
		//loop through all Payment service items to look for MS or CS
		for _, newPaymentServiceItem := range newPaymentServiceItems {
			if newPaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeMS || newPaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeCS {
				moveLevelItems = append(moveLevelItems, newPaymentServiceItem)
				for _, movePR := range allMovePaymentRequests {
					for _, movePaymentServiceItem := range movePR.PaymentServiceItems {
						if movePaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeMS || movePaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeCS {
							if movePaymentServiceItem.Status == models.PaymentServiceItemStatusRequested || movePaymentServiceItem.Status == models.PaymentServiceItemStatusPaid {
								return apperror.NewConflictError(movePR.ID, "Conflict Error: Payment Request for Service Item is already paid or requested")
							}
						}
					}
				}
			}
		}

		// if len(moveLevelItems) > 0 {

		// }

		// if len(moveLevelItems) > 0 {
		// 	for _, movePR := range allMovePaymentRequests {
		// 		for _, movePaymentServiceItem := range movePR.PaymentServiceItems {
		// 			if movePaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeMS || movePaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeCS {
		// 				if (movePaymentServiceItem.Status == models.PaymentServiceItemStatusRequested || movePaymentServiceItem.Status == models.PaymentServiceItemStatusPaid){
		// 					return apperror.NewConflictError(movePR.ID, "Conflict Error: Payment Request for Service Item is already paid or requested")
		// 				}
		// 			}
		// 		}
		// 	}

		// }

		// if there are 0 move level items, then run the original function searching by shipmentID
		if len(moveLevelItems) == 0 {

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
		}
		return nil
	})
}
