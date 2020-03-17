package paymentrequest

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	serviceparamlookups "github.com/transcom/mymove/pkg/payment_request/service_param_value_lookups"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type paymentRequestCreator struct {
	db *pop.Connection
}

// NewPaymentRequestCreator returns a new payment request creator
func NewPaymentRequestCreator(db *pop.Connection) services.PaymentRequestCreator {
	return &paymentRequestCreator{db: db}
}

func (p *paymentRequestCreator) createPaymentRequestSaveToDB(tx *pop.Connection, paymentRequest *models.PaymentRequest, requestedAt time.Time) (*models.PaymentRequest, error) {
	// Verify that the MTO ID exists
	//
	// Lock on the parent row to keep multiple transactions from getting this count at the same time
	// for the same move_task_order_id.  This should block if another payment request comes in for the
	// same move_task_order_id.  Payment requests for other move_task_order_ids should run concurrently.
	// Also note that we use "FOR NO KEY UPDATE" to allow concurrent mods to other tables that have a
	// FK to move_task_orders.
	var moveTaskOrder models.MoveTaskOrder
	sqlString, sqlArgs := tx.Where("id = $1", paymentRequest.MoveTaskOrderID).ToSQL(&pop.Model{Value: &moveTaskOrder})
	sqlString += " FOR NO KEY UPDATE"
	err := tx.RawQuery(sqlString, sqlArgs...).First(&moveTaskOrder)
	if err != nil {
		return nil, fmt.Errorf("could not find MoveTaskOrderID [%s]: %w", paymentRequest.MoveTaskOrderID, err)
	}

	// Update PaymentRequest
	paymentRequest.MoveTaskOrder = moveTaskOrder
	paymentRequest.Status = models.PaymentRequestStatusPending
	paymentRequest.RequestedAt = requestedAt

	uniqueIdentifier, sequenceNumber, err := p.makeUniqueIdentifier(tx, moveTaskOrder)
	if err != nil {
		return nil, fmt.Errorf("issue creating payment request unique identifier: %w", err)
	}
	paymentRequest.PaymentRequestNumber = uniqueIdentifier
	paymentRequest.SequenceNumber = sequenceNumber

	// Create the payment request for the database
	verrs, err := tx.ValidateAndCreate(paymentRequest)
	if verrs.HasAny() {
		return nil, fmt.Errorf("validation error creating payment request: %w for %s", verrs, paymentRequest.ID.String())
	}
	if err != nil {
		return nil, fmt.Errorf("failure creating payment request: %w for %s", err, paymentRequest.ID.String())
	}

	return paymentRequest, nil
}

func (p *paymentRequestCreator) createPaymentServiceItem(tx *pop.Connection, paymentServiceItem *models.PaymentServiceItem, paymentRequest *models.PaymentRequest, requestedAt time.Time) (models.PaymentServiceItem, models.MTOServiceItem, error) {
	// Verify that the MTO service item ID exists
	var mtoServiceItem models.MTOServiceItem
	err := tx.Eager("ReService").Find(&mtoServiceItem, paymentServiceItem.MTOServiceItemID)
	if err != nil {
		return *paymentServiceItem, models.MTOServiceItem{}, fmt.Errorf("could not find MTO MTOServiceItemID [%s]: %w", paymentServiceItem.MTOServiceItemID.String(), err)
	}

	paymentServiceItem.MTOServiceItemID = mtoServiceItem.ID
	paymentServiceItem.MTOServiceItem = mtoServiceItem
	paymentServiceItem.PaymentRequestID = paymentRequest.ID
	paymentServiceItem.PaymentRequest = *paymentRequest
	paymentServiceItem.Status = models.PaymentServiceItemStatusRequested
	// TODO: should PriceCents be a pointer? "0 cents " might be a valid value
	paymentServiceItem.PriceCents = unit.Cents(0) // TODO: Placeholder until we have pricing ready.
	paymentServiceItem.RequestedAt = requestedAt

	verrs, err := tx.ValidateAndCreate(paymentServiceItem)
	if err != nil {
		return *paymentServiceItem, mtoServiceItem, fmt.Errorf("failure creating payment service item: %w for MTO Service Item ID <%s>", err, paymentServiceItem.MTOServiceItemID.String())
	}
	if verrs.HasAny() {
		return *paymentServiceItem, mtoServiceItem, fmt.Errorf("validation error creating payment service item: %w for MTO Service Item ID <%s>", verrs, paymentServiceItem.MTOServiceItemID.String())
	}

	return *paymentServiceItem, mtoServiceItem, nil
}

func (p *paymentRequestCreator) createPaymentServiceItemParam(tx *pop.Connection, paymentServiceItemParam *models.PaymentServiceItemParam, paymentServiceItem models.PaymentServiceItem) (*models.PaymentServiceItemParam, *string, *string, error) {
	// If the ServiceItemParamKeyID is provided, verify it exists; otherwise, lookup
	// via the IncomingKey field
	var key, value string
	createParam := false
	var serviceItemParamKey models.ServiceItemParamKey
	if paymentServiceItemParam.ServiceItemParamKeyID != uuid.Nil {
		err := tx.Find(&serviceItemParamKey, paymentServiceItemParam.ServiceItemParamKeyID)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("could not find ServiceItemParamKeyID [%s]: %w", paymentServiceItemParam.ServiceItemParamKeyID, err)
		}
		if serviceItemParamKey.ID == uuid.Nil || serviceItemParamKey.Key == "" {
			return nil, nil, nil, fmt.Errorf("ServiceItemParamKeyID [%s]: has invalid Key <%s> or UUID <%s> ", paymentServiceItemParam.ServiceItemParamKeyID, serviceItemParamKey.Key, serviceItemParamKey.ID.String())
		}
		key = serviceItemParamKey.Key
		value = paymentServiceItemParam.Value
		createParam = true
	} else if paymentServiceItemParam.IncomingKey != "" {
		err := tx.Where("key = ?", paymentServiceItemParam.IncomingKey).First(&serviceItemParamKey)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("could not find param key [%s]: %w", paymentServiceItemParam.IncomingKey, err)
		}

		key = paymentServiceItemParam.IncomingKey
		value = paymentServiceItemParam.Value
		createParam = true
	}
	if createParam {
		paymentServiceItemParam.ServiceItemParamKeyID = serviceItemParamKey.ID
		paymentServiceItemParam.ServiceItemParamKey = serviceItemParamKey

		paymentServiceItemParam.PaymentServiceItemID = paymentServiceItem.ID
		paymentServiceItemParam.PaymentServiceItem = paymentServiceItem

		var verrs *validate.Errors
		var err error
		verrs, err = tx.ValidateAndCreate(paymentServiceItemParam)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failure creating payment service item param: %w for Payment Service Item ID <%s> Service Item Param Key <%s>", err, paymentServiceItem.ID.String(), serviceItemParamKey.Key)
		}
		if verrs.HasAny() {
			return nil, nil, nil, fmt.Errorf("validation error creating payment service item param: %w for Payment Service Item ID <%s> Service Item Param Key <%s>", verrs, paymentServiceItem.ID.String(), serviceItemParamKey.Key)

		}

		return paymentServiceItemParam, &key, &value, nil
	}

	// incoming provided paymentServiceItemParam is empty
	return nil, nil, nil, nil
}

func (p *paymentRequestCreator) createServiceItemParamFromLookup(tx *pop.Connection, paramLookup *serviceparamlookups.ServiceItemParamKeyData, serviceParam models.ServiceParam, paymentServiceItem models.PaymentServiceItem) (*models.PaymentServiceItemParam, error) {
	// key not found in map
	// Did not find service item param needed for pricing, add it to the list
	value, err := paramLookup.ServiceParamValue(serviceParam.ServiceItemParamKey.Key)
	if err != nil {
		errMessage := "Failed to lookup ServiceParamValue for param key <" + serviceParam.ServiceItemParamKey.Key + "> "
		return nil, fmt.Errorf("%s err: %w", errMessage, err)
	}

	paymentServiceItemParam := models.PaymentServiceItemParam{
		// ID and PaymentServiceItemID to be filled in when payment request is created
		PaymentServiceItemID:  paymentServiceItem.ID,
		PaymentServiceItem:    paymentServiceItem,
		ServiceItemParamKeyID: serviceParam.ServiceItemParamKey.ID,
		ServiceItemParamKey:   serviceParam.ServiceItemParamKey,
		IncomingKey:           serviceParam.ServiceItemParamKey.Key,
		Value:                 value,
	}

	var verrs *validate.Errors
	verrs, err = tx.ValidateAndCreate(&paymentServiceItemParam)
	if err != nil {
		return nil, fmt.Errorf("failure creating payment service item param: %w for payment service item ID <%s> and service item key <%s>", err, paymentServiceItem.ID.String(), serviceParam.ServiceItemParamKey.Key)
	}
	if verrs.HasAny() {
		return nil, fmt.Errorf("validation error creating payment service item param: %w for payment service item ID <%s> and service item key <%s>", verrs, paymentServiceItem.ID.String(), serviceParam.ServiceItemParamKey.Key)
	}

	return &paymentServiceItemParam, nil
}

func (p *paymentRequestCreator) CreatePaymentRequest(paymentRequestArg *models.PaymentRequest) (*models.PaymentRequest, error) {
	transactionError := p.db.Transaction(func(tx *pop.Connection) error {
		var err error
		now := time.Now()

		// Gather information for logging
		mtoMessageString := " MTO ID <" + paymentRequestArg.MoveTaskOrderID.String() + ">"
		prMessageString := " paymentRequestID <" + paymentRequestArg.ID.String() + ">"

		// Create the payment request
		paymentRequestArg, err = p.createPaymentRequestSaveToDB(tx, paymentRequestArg, now)
		if err != nil {
			return fmt.Errorf("failure creating payment request: %w for %s", err, mtoMessageString+prMessageString)
		}
		if paymentRequestArg == nil {
			return fmt.Errorf("failure creating payment request <nil> for %s", mtoMessageString+prMessageString)
		}

		// Create a payment service item for each incoming payment service item in the payment request
		// These incoming payment service items have not been created in the database yet
		var newPaymentServiceItems models.PaymentServiceItems
		for _, paymentServiceItem := range paymentRequestArg.PaymentServiceItems {

			var mtoServiceItem models.MTOServiceItem
			// Create the payment service item
			paymentServiceItem, mtoServiceItem, err = p.createPaymentServiceItem(tx, &paymentServiceItem, paymentRequestArg, now)

			// Gather message information for logging
			mtoServiceItemString := " MTO Service item ID <" + paymentServiceItem.MTOServiceItemID.String() + ">"
			reServiceItem := paymentServiceItem.MTOServiceItem.ReService
			serviceItemMessageString := " RE Service Item Code: <" + string(reServiceItem.Code) + "> Name: <" + reServiceItem.Name + ">"
			errMessageString := mtoMessageString + prMessageString + mtoServiceItemString + serviceItemMessageString

			// Check for error from creating payment service item
			if err != nil {
				return fmt.Errorf("failure creating payment service item: %w for %s", err, errMessageString)
			}

			// store param Key:Value pairs coming from the create payment request payload, sent by the user when requesting payment
			incomingMTOServiceItemParams := make(map[string]string)

			// Create a payment service item parameter for each of the incoming payment service item params
			var newPaymentServiceItemParams models.PaymentServiceItemParams
			for _, paymentServiceItemParam := range paymentServiceItem.PaymentServiceItemParams {
				var param *models.PaymentServiceItemParam
				var key, value *string
				param, key, value, err = p.createPaymentServiceItemParam(tx, &paymentServiceItemParam, paymentServiceItem)
				if err != nil {
					return fmt.Errorf("failed to create payment service item param [%s]: %w for %s", paymentServiceItemParam.ServiceItemParamKeyID, err, errMessageString)
				}
				if param != nil && key != nil && value != nil {
					incomingMTOServiceItemParams[*key] = *value
					newPaymentServiceItemParams = append(newPaymentServiceItemParams, paymentServiceItemParam)
				}
			}

			//
			// For the current service item, find any missing params needed to price
			// this service item
			//

			// Retrieve all of the params needed to price this service item
			paymentHelper := paymentrequesthelper.RequestPaymentHelper{DB: tx}
			reServiceParams, err := paymentHelper.FetchServiceParamList(paymentServiceItem.MTOServiceItemID)
			if err != nil {
				errMessage := "Failed to retrieve service item param list for " + errMessageString
				return fmt.Errorf("%s err: %w", errMessage, err)
			}

			// Get values for needed service item params (do lookups)
			paramLookup := serviceparamlookups.ServiceParamLookupInitialize(paymentServiceItem.MTOServiceItemID, paymentServiceItem.ID, paymentRequestArg.MoveTaskOrderID)
			for _, reServiceParam := range reServiceParams {
				if _, found := incomingMTOServiceItemParams[reServiceParam.ServiceItemParamKey.Key]; !found {
					// create the missing service item param
					var param *models.PaymentServiceItemParam
					param, err = p.createServiceItemParamFromLookup(tx, paramLookup, reServiceParam, paymentServiceItem)
					if err != nil {
						errMessage := "Failed to create service item param for param key <" + reServiceParam.ServiceItemParamKey.Key + "> " + errMessageString
						return fmt.Errorf("%s err: %w", errMessage, err)
					}
					if param != nil {
						newPaymentServiceItemParams = append(newPaymentServiceItemParams, *param)
					}
				}
			}

			//
			// Save all params for current service item
			// Save the new payment service item to the list of service items to be returned
			//

			paymentServiceItem.PaymentServiceItemParams = newPaymentServiceItemParams

			newPaymentServiceItems = append(newPaymentServiceItems, paymentServiceItem)

			//
			// Validate that all params are available to prices the service item that is in
			// the payment request
			//
			validParamList, validateMessage := paymentHelper.ValidServiceParamList(mtoServiceItem, reServiceParams, paymentServiceItem.PaymentServiceItemParams)
			if validParamList == false {
				errMessage := "service item param list is not valid (will not be able to price the item) " + validateMessage + " for " + errMessageString
				return fmt.Errorf("%s err: %w", errMessage, err)
			}
		}
		paymentRequestArg.PaymentServiceItems = newPaymentServiceItems

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return paymentRequestArg, nil
}

func (p *paymentRequestCreator) makeUniqueIdentifier(tx *pop.Connection, mto models.MoveTaskOrder) (string, int, error) {
	// Get the max sequence number that exists for the payment requests associated with the given MTO.
	// Since we have a lock to prevent concurrent payment requests for this MTO, this should be safe.
	var max int
	err := tx.RawQuery("SELECT COALESCE(MAX(sequence_number),0) FROM payment_requests WHERE move_task_order_id = $1", mto.ID).First(&max)
	if err != nil {
		return "", 0, fmt.Errorf("max sequence_number for MoveTaskOrderID [%s] failed: %w", mto.ID, err)
	}

	if mto.ReferenceID == "" {
		return "", 0, fmt.Errorf("could not find reference ID for MoveTaskOrderID [%s]", mto.ID)
	}

	nextSequence := max + 1
	paymentRequestNumber := fmt.Sprintf("%s-%d", mto.ReferenceID, nextSequence)

	return paymentRequestNumber, nextSequence, nil
}
