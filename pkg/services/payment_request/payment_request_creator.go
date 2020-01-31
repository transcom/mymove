package paymentrequest

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gobuffalo/pop"
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

func (p *paymentRequestCreator) CreatePaymentRequest(paymentRequest *models.PaymentRequest) (*models.PaymentRequest, error) {
	transactionError := p.db.Transaction(func(tx *pop.Connection) error {
		now := time.Now()

		// Verify that the MTO ID exists
		var moveTaskOrder models.MoveTaskOrder
		err := tx.Find(&moveTaskOrder, paymentRequest.MoveTaskOrderID)
		if err != nil {
			return fmt.Errorf("could not find MoveTaskOrderID [%s]: %w", paymentRequest.MoveTaskOrderID, err)
		}
		paymentRequest.MoveTaskOrder = moveTaskOrder

		paymentRequest.Status = models.PaymentRequestStatusPending
		paymentRequest.RequestedAt = now

		uniqueIdentifier, err := p.makeUniqueIdentifier(paymentRequest.MoveTaskOrderID)
		if err != nil {
			return fmt.Errorf("issue creating payment request unique identifier: %w", err)
		}
		paymentRequest.PaymentRequestNumber = uniqueIdentifier

		// Create the payment request first
		verrs, err := tx.ValidateAndCreate(paymentRequest)
		if err != nil {
			return fmt.Errorf("failure creating payment request: %w", err)
		}
		if verrs.HasAny() {
			return fmt.Errorf("validation error creating payment request: %w", verrs)
		}

		// Create each payment service item for the payment request
		var newPaymentServiceItems models.PaymentServiceItems
		for _, paymentServiceItem := range paymentRequest.PaymentServiceItems {
			// Verify that the service item ID exists
			var mtoServiceItem models.MTOServiceItem
			err := tx.Find(&mtoServiceItem, paymentServiceItem.ServiceItemID)
			if err != nil {
				return fmt.Errorf("could not find ServiceItemID [%s]: %w", paymentServiceItem.ServiceItemID, err)
			}
			paymentServiceItem.ServiceItem = mtoServiceItem

			paymentServiceItem.PaymentRequestID = paymentRequest.ID
			paymentServiceItem.PaymentRequest = *paymentRequest
			paymentServiceItem.Status = models.PaymentServiceItemStatusRequested
			paymentServiceItem.PriceCents = unit.Cents(0) // TODO: Placeholder until we have pricing ready.
			paymentServiceItem.RequestedAt = now

			verrs, err := tx.ValidateAndCreate(&paymentServiceItem)
			if err != nil {
				return fmt.Errorf("failure creating payment service item: %w", err)
			}
			if verrs.HasAny() {
				return fmt.Errorf("validation error creating payment service item: %w", verrs)
			}

			// Create each payment service item parameter for the payment service item
			var newPaymentServiceItemParams models.PaymentServiceItemParams
			for _, paymentServiceItemParam := range paymentServiceItem.PaymentServiceItemParams {
				// If the ServiceItemParamKeyID is provided, verify it exists; otherwise, lookup
				// via the IncomingKey field
				var serviceItemParamKey models.ServiceItemParamKey
				if paymentServiceItemParam.ServiceItemParamKeyID != uuid.Nil {
					err := tx.Find(&serviceItemParamKey, paymentServiceItemParam.ServiceItemParamKeyID)
					if err != nil {
						return fmt.Errorf("could not find ServiceItemParamKeyID [%s]: %w", paymentServiceItemParam.ServiceItemParamKeyID, err)
					}
				} else {
					err := tx.Where("key = ?", paymentServiceItemParam.IncomingKey).First(&serviceItemParamKey)
					if err != nil {
						return fmt.Errorf("could not find param key [%s]: %w", paymentServiceItemParam.IncomingKey, err)
					}
				}
				paymentServiceItemParam.ServiceItemParamKeyID = serviceItemParamKey.ID
				paymentServiceItemParam.ServiceItemParamKey = serviceItemParamKey

				paymentServiceItemParam.PaymentServiceItemID = paymentServiceItem.ID
				paymentServiceItemParam.PaymentServiceItem = paymentServiceItem

				verrs, err := tx.ValidateAndCreate(&paymentServiceItemParam)
				if err != nil {
					return fmt.Errorf("failure creating payment service item param: %w", err)
				}
				if verrs.HasAny() {
					return fmt.Errorf("validation error creating payment service item param: %w", verrs)
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
		return nil, transactionError
	}

	return paymentRequest, nil
}

func (p *paymentRequestCreator) makeUniqueIdentifier(mtoID uuid.UUID) (string, error) {
	paymentRequestNumber, err := p.determinePaymentRequestNumberSuffix(mtoID)
	if err != nil {
		return "", fmt.Errorf("error determining Payment RequestNumber: %w", err)
	}
	var moveTaskOrder models.MoveTaskOrder
	err = p.db.Where("id = $1", mtoID).First(&moveTaskOrder)
	if err != nil {
		return "", fmt.Errorf("failed fetch of move task order: %w", err)
	}

	uniqueIdentifier := fmt.Sprintf("%s-%s", *moveTaskOrder.ReferenceID, strconv.Itoa(paymentRequestNumber))

	return uniqueIdentifier, err
}

func (p *paymentRequestCreator) determinePaymentRequestNumberSuffix(mtoID uuid.UUID) (int, error) {
	query := p.db.Where("move_task_order_id = $1", mtoID)
	count, err := query.Count(models.PaymentRequest{})
	if err != nil {
		return 0, fmt.Errorf("count of payment requests on move task order failed: %w", err)
	}

	paymentRequestNumber := count + 1

	return paymentRequestNumber, err
}
