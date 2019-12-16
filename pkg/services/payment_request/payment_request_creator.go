package paymentrequest

import (
	"errors"
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type paymentRequestCreator struct {
	db *pop.Connection
}

func NewPaymentRequestCreator(db *pop.Connection) services.PaymentRequestCreator {
	return &paymentRequestCreator{db}
}

func (p *paymentRequestCreator) CreatePaymentRequest(paymentRequest *models.PaymentRequest) (*models.PaymentRequest, *validate.Errors, error) {
	responseVErrors := validate.NewErrors()
	var responseError error

	transactionError := p.db.Transaction(func(tx *pop.Connection) error {
		transactionError := errors.New("rollback the transaction")

		now := time.Now()

		// Verify that the MTO ID exists
		var moveTaskOrder models.MoveTaskOrder
		err := p.db.Find(&moveTaskOrder, paymentRequest.MoveTaskOrderID)
		if err != nil {
			responseError = fmt.Errorf("could not find MoveTaskOrderID [%s]: %w", paymentRequest.MoveTaskOrderID, err)
			return transactionError
		}
		paymentRequest.MoveTaskOrder = moveTaskOrder

		paymentRequest.Status = models.PaymentRequestStatusPending
		paymentRequest.RequestedAt = now

		// Create the payment request first
		verrs, err := p.db.ValidateAndCreate(paymentRequest)
		if err != nil || verrs.HasAny() {
			responseVErrors.Append(verrs)
			responseError = err
			return transactionError
		}

		// Create each payment service item for the payment request
		var newPaymentServiceItems models.PaymentServiceItems
		for _, paymentServiceItem := range paymentRequest.PaymentServiceItems {
			// Verify that the service item ID exists
			var mtoServiceItem models.MTOServiceItem
			err := p.db.Find(&mtoServiceItem, paymentServiceItem.ServiceItemID)
			if err != nil {
				responseError = fmt.Errorf("could not find ServiceItemID [%s]: %w", paymentServiceItem.ServiceItemID, err)
				return transactionError
			}
			paymentServiceItem.ServiceItem = mtoServiceItem

			paymentServiceItem.PaymentRequestID = paymentRequest.ID
			paymentServiceItem.PaymentRequest = *paymentRequest
			paymentServiceItem.Status = models.PaymentServiceItemStatusRequested
			paymentServiceItem.PriceCents = unit.Cents(0) // TODO: Placeholder until we have pricing ready.
			paymentServiceItem.RequestedAt = now

			verrs, err := p.db.ValidateAndCreate(&paymentServiceItem)
			if err != nil || verrs.HasAny() {
				responseVErrors.Append(verrs)
				responseError = err
				return transactionError
			}

			// Create each payment service item parameter for the payment service item
			var newPaymentServiceItemParams models.PaymentServiceItemParams
			for _, paymentServiceItemParam := range paymentServiceItem.PaymentServiceItemParams {
				// If the ServiceItemParamKeyID is provided, verify it exists; otherwise, lookup
				// via the IncomingKey field
				var serviceItemParamKey models.ServiceItemParamKey
				if paymentServiceItemParam.ServiceItemParamKeyID != uuid.Nil {
					err := p.db.Find(&serviceItemParamKey, paymentServiceItemParam.ServiceItemParamKeyID)
					if err != nil {
						responseError = fmt.Errorf("could not find ServiceItemParamKeyID [%s]: %w", paymentServiceItemParam.ServiceItemParamKeyID, err)
						return transactionError
					}
				} else {
					err := p.db.Where("key = ?", paymentServiceItemParam.IncomingKey).First(&serviceItemParamKey)
					if err != nil {
						responseError = fmt.Errorf("could not find param key [%s]: %w", paymentServiceItemParam.IncomingKey, err)
						return transactionError
					}
				}
				paymentServiceItemParam.ServiceItemParamKeyID = serviceItemParamKey.ID
				paymentServiceItemParam.ServiceItemParamKey = serviceItemParamKey

				paymentServiceItemParam.PaymentServiceItemID = paymentServiceItem.ID
				paymentServiceItemParam.PaymentServiceItem = paymentServiceItem

				verrs, err := p.db.ValidateAndCreate(&paymentServiceItemParam)
				if err != nil || verrs.HasAny() {
					responseVErrors.Append(verrs)
					responseError = err
					return transactionError
				}

				newPaymentServiceItemParams = append(newPaymentServiceItemParams, paymentServiceItemParam)
			}
			paymentServiceItem.PaymentServiceItemParams = newPaymentServiceItemParams

			newPaymentServiceItems = append(newPaymentServiceItems, paymentServiceItem)
		}
		paymentRequest.PaymentServiceItems = newPaymentServiceItems

		return nil
	})

	if transactionError != nil {
		return nil, responseVErrors, responseError
	}

	return paymentRequest, responseVErrors, responseError
}
