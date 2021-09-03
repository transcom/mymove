package paymentrequest

import (
	"fmt"
	"sort"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type paymentRequestRepricer struct {
	paymentRequestCreator services.PaymentRequestCreator
}

// NewPaymentRequestRepricer returns a new payment request repricer
func NewPaymentRequestRepricer(paymentRequestCreator services.PaymentRequestCreator) services.PaymentRequestRepricer {
	return &paymentRequestRepricer{
		paymentRequestCreator: paymentRequestCreator,
	}
}

func (p *paymentRequestRepricer) RepricePaymentRequest(appCtx appcontext.AppContext, paymentRequestID uuid.UUID) (*models.PaymentRequest, error) {
	var newPaymentRequest *models.PaymentRequest

	// Make sure we do this whole process in a transaction so partial changes do not get made committed
	// in the event of an error.
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// Try to fetch the existing payment request.
		var existingPaymentRequest models.PaymentRequest

		// Fetch the payment request and payment service items from the existing request.
		err := txnAppCtx.DB().
			Eager("PaymentServiceItems.MTOServiceItem.ReService").
			Find(&existingPaymentRequest, paymentRequestID)
		if err != nil {
			return err
		}

		// Re-create the payment request which will cause repricing to occur.
		inputPaymentRequest := buildPaymentRequestForRepricing(existingPaymentRequest)
		newPaymentRequest, err = p.paymentRequestCreator.CreatePaymentRequest(txnAppCtx, &inputPaymentRequest)
		if err != nil {
			return err
		}

		// Set the (now) old payment request's status.
		// TODO: We need a better status for this -- something like "REPRICED".
		newStatus := models.PaymentRequestStatusReviewedAllRejected
		existingPaymentRequest.Status = newStatus
		verrs, err := txnAppCtx.DB().ValidateAndUpdate(&existingPaymentRequest)
		if err != nil {
			return fmt.Errorf("failed to set existing payment request status to %v: %w", newStatus, err)
		}
		if verrs.HasAny() {
			return fmt.Errorf("failed to validate existing payment request when setting status to %v: %w", newStatus, verrs)
		}

		// TODO:
		//   - Link repriced payment request to old payment request
		//   - Need to re-associate proof of service docs

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return newPaymentRequest, nil
}

// buildPaymentRequestForRepricing builds up the expected payment request data based upon an existing payment request.
func buildPaymentRequestForRepricing(existingPaymentRequest models.PaymentRequest) models.PaymentRequest {
	newPaymentRequest := models.PaymentRequest{
		IsFinal:         existingPaymentRequest.IsFinal,
		MoveTaskOrderID: existingPaymentRequest.MoveTaskOrderID,
	}

	var newPaymentServiceItems models.PaymentServiceItems
	for _, existingPaymentServiceItem := range existingPaymentRequest.PaymentServiceItems {
		newPaymentServiceItem := models.PaymentServiceItem{
			MTOServiceItemID: existingPaymentServiceItem.MTOServiceItemID,
			MTOServiceItem:   existingPaymentServiceItem.MTOServiceItem,
		}

		newPaymentServiceItems = append(newPaymentServiceItems, newPaymentServiceItem)
	}

	sort.SliceStable(newPaymentServiceItems, func(i, j int) bool {
		return newPaymentServiceItems[i].MTOServiceItem.ReService.Priority < newPaymentServiceItems[j].MTOServiceItem.ReService.Priority
	})

	newPaymentRequest.PaymentServiceItems = newPaymentServiceItems

	return newPaymentRequest
}
