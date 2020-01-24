package paymentrequest

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type paymentRequestCreator struct {
	db *pop.Connection
}

func NewPaymentRequestCreator(db *pop.Connection) services.PaymentRequestCreator {
	return &paymentRequestCreator{db: db}
}

func (p *paymentRequestCreator) CreatePaymentRequest(paymentRequestArg *models.PaymentRequest) (*models.PaymentRequest, error) {
	transactionError := p.db.Transaction(func(tx *pop.Connection) error {
		now := time.Now()

		// Verify that the MTO ID exists
		var moveTaskOrder models.MoveTaskOrder
		err := tx.Find(&moveTaskOrder, paymentRequestArg.MoveTaskOrderID)
		if err != nil {
			return fmt.Errorf("could not find MoveTaskOrderID [%s]: %w", paymentRequestArg.MoveTaskOrderID, err)
		}
		paymentRequestArg.MoveTaskOrder = moveTaskOrder

		paymentRequestArg.Status = models.PaymentRequestStatusPending
		paymentRequestArg.RequestedAt = now

		// Create the payment request first
		verrs, err := tx.ValidateAndCreate(paymentRequestArg)
		if verrs.HasAny() {
			return fmt.Errorf("validation error creating payment request: %w", verrs)
		}
		if err != nil {
			return fmt.Errorf("failure creating payment request: %w", err)
		}

		// Create each payment service item for the payment request
		var newPaymentServiceItems models.PaymentServiceItems
		for _, paymentServiceItem := range paymentRequestArg.PaymentServiceItems {
			fmt.Printf("===== MTO Service Item <%s>\n", paymentServiceItem.MTOServiceItem.ID.String())

			// Verify that the service item ID exists
			var mtoServiceItem models.MTOServiceItem
			err := tx.Find(&mtoServiceItem, paymentServiceItem.MTOServiceItemID)
			if err != nil {
				return fmt.Errorf("could not find MTO MTOServiceItemID [%s]: %w", paymentServiceItem.MTOServiceItemID, err)
			}
			paymentServiceItem.MTOServiceItemID = mtoServiceItem.ID
			paymentServiceItem.MTOServiceItem = mtoServiceItem
			paymentServiceItem.PaymentRequestID = paymentRequestArg.ID
			paymentServiceItem.PaymentRequest = *paymentRequestArg
			paymentServiceItem.Status = models.PaymentServiceItemStatusRequested
			// TODO: should PriceCents be a pointer? "0 cents " might be a valid value
			paymentServiceItem.PriceCents = unit.Cents(0) // TODO: Placeholder until we have pricing ready.
			paymentServiceItem.RequestedAt = now

			verrs, err := tx.ValidateAndCreate(&paymentServiceItem)
			if err != nil {
				return fmt.Errorf("failure creating payment service item: %w", err)
			}
			if verrs.HasAny() {
				return fmt.Errorf("validation error creating payment service item: %w", verrs)
			}

			// Param Key:Value pairs coming from the create payment request payload, sent by the user when requesting payment
			var incomingMTOServiceItemParams map[string]string
			incomingMTOServiceItemParams = make(map[string]string)

			// Create each payment service item parameter for the payment service item
			var newPaymentServiceItemParams models.PaymentServiceItemParams
			for _, paymentServiceItemParam := range paymentServiceItem.PaymentServiceItemParams {
				// If the ServiceItemParamKeyID is provided, verify it exists; otherwise, lookup
				// via the IncomingKey field
				foundParamKey := false
				var serviceItemParamKey models.ServiceItemParamKey
				if paymentServiceItemParam.ServiceItemParamKeyID != uuid.Nil {
					fmt.Printf("Key <%s> via ServiceItemParamKey.Key added to param list \n\n\n", paymentServiceItemParam.ServiceItemParamKey.Key)
					err = tx.Find(&serviceItemParamKey, paymentServiceItemParam.ServiceItemParamKeyID)
					if err != nil {
						return fmt.Errorf("could not find ServiceItemParamKeyID [%s]: %w", paymentServiceItemParam.ServiceItemParamKeyID, err)
					}
					if serviceItemParamKey.ID == uuid.Nil || serviceItemParamKey.Key == "" {
						return fmt.Errorf("ServiceItemParamKeyID [%s]: has invalid Key <%s> or UUID <%s>", paymentServiceItemParam.ServiceItemParamKeyID, serviceItemParamKey.Key, serviceItemParamKey.ID.String())
					}
					incomingMTOServiceItemParams[serviceItemParamKey.Key] = paymentServiceItemParam.Value
					foundParamKey = true
				} else if paymentServiceItemParam.IncomingKey != "" {
					fmt.Printf("Key <%s> via IncomingKey added to param list \n\n\n", paymentServiceItemParam.IncomingKey)

					err = tx.Where("key = ?", paymentServiceItemParam.IncomingKey).First(&serviceItemParamKey)
					if err != nil {
						return fmt.Errorf("could not find param key [%s]: %w", paymentServiceItemParam.IncomingKey, err)
					}
					incomingMTOServiceItemParams[paymentServiceItemParam.IncomingKey] = paymentServiceItemParam.Value
					foundParamKey = true
				}
				if foundParamKey {
					paymentServiceItemParam.ServiceItemParamKeyID = serviceItemParamKey.ID
					paymentServiceItemParam.ServiceItemParamKey = serviceItemParamKey

					paymentServiceItemParam.PaymentServiceItemID = paymentServiceItem.ID
					paymentServiceItemParam.PaymentServiceItem = paymentServiceItem

					var verrs *validate.Errors
					verrs, err = tx.ValidateAndCreate(&paymentServiceItemParam)
					if err != nil {
						return fmt.Errorf("failure creating payment service item param: %w", err)
					}
					if verrs.HasAny() {
						return fmt.Errorf("validation error creating payment service item param: %w", verrs)
					}

					newPaymentServiceItemParams = append(newPaymentServiceItemParams, paymentServiceItemParam)
				}
			}

			//
			// For the existing params for the current service item, find any missing params needed to price
			// this service item
			//

			// Retrieve all of the params needed to price this service item
			paymentHelper := paymentrequesthelper.RequestPaymentHelper{DB: p.db}
			reServiceParams, err := paymentHelper.FetchServiceParamList(paymentServiceItem.MTOServiceItemID)
			if err != nil {
				errMessage := "Failed to retrieve service item param list for MTO Service ID " + paymentServiceItem.MTOServiceItemID.String() + ", mtoServiceID " + paymentServiceItem.MTOServiceItemID.String() + ", paymentRequestID " + paymentRequestArg.ID.String()
				return fmt.Errorf("%s err: %w", errMessage, err)
			}

			// Get values for needed service item params (do lookups)
			//paramLookup := serviceparamlookups.ServiceParamLookupInitialize(paymentServiceItem.MTOServiceItemID, paymentServiceItem.ID, paymentRequestArg.MoveTaskOrderID)
			for _, reServiceParam := range reServiceParams {
				var value string
				if _, found := incomingMTOServiceItemParams[reServiceParam.ServiceItemParamKey.Key]; !found {
					// key not found in map

					fmt.Printf("\n\nKey not found <%s>, add to param list\n\n", reServiceParam.ServiceItemParamKey.Key)
					// Did not find service item param needed for pricing, add it to the list
					value = "8.88"
					//value, err = paramLookup.ServiceParamValue(reServiceParam.ServiceItemParamKey.Key)
					if err != nil {
						errMessage := "Failed to lookup ServiceParamValue for param key " + reServiceParam.ServiceItemParamKey.Key + ", mtoServiceID " + paymentServiceItem.MTOServiceItemID.String() + ", mtoID " + paymentRequestArg.MoveTaskOrderID.String() + ", paymentRequestID " + paymentRequestArg.ID.String()
						return fmt.Errorf("%s err: %w", errMessage, err)
					}

					paymentServiceItemParam := models.PaymentServiceItemParam{
						// ID and PaymentServiceItemID to be filled in when payment request is created
						PaymentServiceItemID:  paymentServiceItem.ID,
						PaymentServiceItem:    paymentServiceItem,
						ServiceItemParamKeyID: reServiceParam.ServiceItemParamKey.ID,
						ServiceItemParamKey:   reServiceParam.ServiceItemParamKey,
						IncomingKey:           reServiceParam.ServiceItemParamKey.Key,
						Value:                 value,
					}

					var verrs *validate.Errors
					verrs, err = tx.ValidateAndCreate(&paymentServiceItemParam)
					if err != nil {
						return fmt.Errorf("failure creating payment service item param: %w", err)
					}
					if verrs.HasAny() {
						return fmt.Errorf("validation error creating payment service item param: %w", verrs)
					}

					newPaymentServiceItemParams = append(newPaymentServiceItemParams, paymentServiceItemParam)
				}
			}

			//
			// Save all params for current service item
			// Save the new payment service item to the list of service items to be returned
			//

			paymentServiceItem.PaymentServiceItemParams = newPaymentServiceItemParams

			newPaymentServiceItems = append(newPaymentServiceItems, paymentServiceItem)
		}
		paymentRequestArg.PaymentServiceItems = newPaymentServiceItems

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return paymentRequestArg, nil
}
