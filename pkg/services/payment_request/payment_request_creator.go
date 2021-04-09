package paymentrequest

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	serviceparamlookups "github.com/transcom/mymove/pkg/payment_request/service_param_value_lookups"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
)

type paymentRequestCreator struct {
	db      *pop.Connection
	planner route.Planner
	pricer  services.ServiceItemPricer
}

// NewPaymentRequestCreator returns a new payment request creator
func NewPaymentRequestCreator(db *pop.Connection, planner route.Planner, pricer services.ServiceItemPricer) services.PaymentRequestCreator {
	return &paymentRequestCreator{
		db:      db,
		planner: planner,
		pricer:  pricer,
	}
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
			var badDataError *services.BadDataError

			if _, ok := err.(services.InvalidCreateInputError); ok {
				return err
			}
			if _, ok := err.(services.NotFoundError); ok {
				return err
			}
			if _, ok := err.(services.ConflictError); ok {
				return err
			}
			if _, ok := err.(services.InvalidInputError); ok {
				return err
			}
			if _, ok := err.(services.QueryError); ok {
				return err
			}
			if errors.As(err, &badDataError) {
				return err
			}
			return fmt.Errorf("failure creating payment request: %w for %s", err, mtoMessageString+prMessageString)
		}
		if paymentRequestArg == nil {
			return fmt.Errorf("failure creating payment request <nil> for %s", mtoMessageString+prMessageString)
		}

		// Service Item Param Cache
		serviceParamCache := serviceparamlookups.ServiceParamsCache{}
		serviceParamCache.Initialize(p.db)

		// Run the pricer within this transactional context
		txPricer := p.pricer.UsingConnection(tx)

		// Create a payment service item for each incoming payment service item in the payment request
		// These incoming payment service items have not been created in the database yet
		var newPaymentServiceItems models.PaymentServiceItems
		for _, paymentServiceItem := range paymentRequestArg.PaymentServiceItems {
			var mtoServiceItem models.MTOServiceItem

			// Gather message information for logging
			mtoServiceItemString := " MTO Service item ID <" + paymentServiceItem.MTOServiceItemID.String() + ">"
			reServiceItem := paymentServiceItem.MTOServiceItem.ReService
			serviceItemMessageString := " RE Service Item Code: <" + string(reServiceItem.Code) + "> Name: <" + reServiceItem.Name + ">"
			errMessageString := mtoMessageString + prMessageString + mtoServiceItemString + serviceItemMessageString
			// Create the payment service item
			paymentServiceItem, mtoServiceItem, err = p.createPaymentServiceItem(tx, paymentServiceItem, paymentRequestArg, now)
			if err != nil {
				if _, ok := err.(services.InvalidCreateInputError); ok {
					return err
				}
				if _, ok := err.(services.NotFoundError); ok {
					return err
				}

				return fmt.Errorf("failure creating payment service item: %w for %s", err, errMessageString)
			}

			// store param Key:Value pairs coming from the create payment request payload, sent by the user when requesting payment
			incomingMTOServiceItemParams := make(map[string]string)

			// Create a payment service item parameter for each of the incoming payment service item params
			var newPaymentServiceItemParams models.PaymentServiceItemParams
			for _, paymentServiceItemParam := range paymentServiceItem.PaymentServiceItemParams {
				var param models.PaymentServiceItemParam
				var key, value *string
				param, key, value, err = p.createPaymentServiceItemParam(tx, paymentServiceItemParam, paymentServiceItem)
				if err != nil {
					if _, ok := err.(*services.BadDataError); ok {
						return err
					}
					if _, ok := err.(services.NotFoundError); ok {
						return err
					}
					return fmt.Errorf("failed to create payment service item param [%s]: %w for %s", paymentServiceItemParam.ServiceItemParamKeyID, err, errMessageString)
				}

				if param.ID != uuid.Nil && key != nil && value != nil {
					incomingMTOServiceItemParams[*key] = *value
					newPaymentServiceItemParams = append(newPaymentServiceItemParams, param)
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
			paramLookup, err := serviceparamlookups.ServiceParamLookupInitialize(tx, p.planner, paymentServiceItem.MTOServiceItemID, paymentServiceItem.ID, paymentRequestArg.MoveTaskOrderID, &serviceParamCache)
			if err != nil {
				return err
			}
			for _, reServiceParam := range reServiceParams {
				if _, found := incomingMTOServiceItemParams[reServiceParam.ServiceItemParamKey.Key.String()]; !found {
					// create the missing service item param
					var param *models.PaymentServiceItemParam
					param, err = p.createServiceItemParamFromLookup(tx, paramLookup, reServiceParam, paymentServiceItem)
					if err != nil {
						errMessage := fmt.Sprintf("Failed to create service item param for param key <%s> %s", reServiceParam.ServiceItemParamKey.Key, errMessageString)
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

			//
			// Validate that all params are available to prices the service item that is in
			// the payment request
			//
			validParamList, validateMessage := paymentHelper.ValidServiceParamList(mtoServiceItem, reServiceParams, paymentServiceItem.PaymentServiceItemParams)
			if validParamList == false {
				errMessage := "service item param list is not valid (will not be able to price the item) " + validateMessage + " for " + errMessageString
				return fmt.Errorf("%s err: %w", errMessage, err)
			}

			// Price the payment service item
			var psItem models.PaymentServiceItem
			var displayParams models.PaymentServiceItemParams
			psItem, displayParams, err = p.pricePaymentServiceItem(tx, txPricer, paymentServiceItem)
			if err != nil {
				return fmt.Errorf("failure pricing service %s for MTO service item ID %s: %w",
					paymentServiceItem.MTOServiceItem.ReService.Code, paymentServiceItem.MTOServiceItemID, err)
			}
			if len(displayParams) > 0 {
				psItem.PaymentServiceItemParams = append(paymentServiceItem.PaymentServiceItemParams, displayParams...)
			}
			newPaymentServiceItems = append(newPaymentServiceItems, psItem)
		}

		paymentRequestArg.PaymentServiceItems = newPaymentServiceItems

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return paymentRequestArg, nil
}

func (p *paymentRequestCreator) createPaymentRequestSaveToDB(tx *pop.Connection, paymentRequest *models.PaymentRequest, requestedAt time.Time) (*models.PaymentRequest, error) {
	// Verify that the MTO ID exists
	if paymentRequest.MoveTaskOrderID == uuid.Nil {
		return nil, services.NewInvalidCreateInputError(nil, "Invalid Create Input Error: MoveTaskOrderID is required on PaymentRequest create")
	}

	// Lock on the parent row to keep multiple transactions from getting this count at the same time
	// for the same move_id.  This should block if another payment request comes in for the
	// same move_id.  Payment requests for other move_ids should run concurrently.
	// Also note that we use "FOR NO KEY UPDATE" to allow concurrent mods to other tables that have a
	// FK to move_task_orders.
	var moveTaskOrder models.Move
	sqlString, sqlArgs := tx.Where("id = $1", paymentRequest.MoveTaskOrderID).ToSQL(&pop.Model{Value: &moveTaskOrder})
	sqlString += " FOR NO KEY UPDATE"
	err := tx.RawQuery(sqlString, sqlArgs...).First(&moveTaskOrder)

	if err != nil {
		if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			msg := fmt.Sprint("for Move")
			return nil, services.NewNotFoundError(paymentRequest.MoveTaskOrderID, msg)
		}
		return nil, fmt.Errorf("could not retrieve Move with ID [%s]: %w", paymentRequest.MoveTaskOrderID, err)
	}

	// Verify the Orders on the MTO
	err = tx.Load(&moveTaskOrder, "Orders")

	if err != nil {
		return nil, services.NewNotFoundError(moveTaskOrder.OrdersID, fmt.Sprintf("Orders on MoveTaskOrder (ID: %s) missing", moveTaskOrder.ID))
	}

	// Verify that the Orders has LOA
	if moveTaskOrder.Orders.TAC == nil || *moveTaskOrder.Orders.TAC == "" {
		return nil, services.NewConflictError(moveTaskOrder.OrdersID, fmt.Sprintf("Orders on MoveTaskOrder (ID: %s) missing Lines of Accounting TAC", moveTaskOrder.ID))
	}
	// Verify that the Orders have OriginDutyStation
	if moveTaskOrder.Orders.OriginDutyStationID == nil {
		return nil, services.NewConflictError(moveTaskOrder.OrdersID, fmt.Sprintf("Orders on MoveTaskOrder (ID: %s) missing OriginDutyStation", moveTaskOrder.ID))
	}
	// Verify that ServiceMember is Valid
	err = tx.Load(&moveTaskOrder.Orders, "ServiceMember")
	if err != nil {
		return nil, services.NewNotFoundError(moveTaskOrder.Orders.ServiceMemberID, fmt.Sprintf("ServiceMember on MoveTaskOrder (ID: %s) not valid", moveTaskOrder.ID))
	}

	serviceMember := moveTaskOrder.Orders.ServiceMember
	// Verify First Name
	if serviceMember.FirstName == nil || *serviceMember.FirstName == "" {
		return nil, services.NewConflictError(moveTaskOrder.Orders.ServiceMemberID, fmt.Sprintf("ServiceMember on MoveTaskOrder (ID: %s) missing First Name", moveTaskOrder.ID))
	}
	// Verify Last Name
	if serviceMember.LastName == nil || *serviceMember.LastName == "" {
		return nil, services.NewConflictError(moveTaskOrder.Orders.ServiceMemberID, fmt.Sprintf("ServiceMember on MoveTaskOrder (ID: %s) missing Last Name", moveTaskOrder.ID))
	}
	// Verify Rank
	if serviceMember.Rank == nil || *serviceMember.Rank == "" {
		return nil, services.NewConflictError(moveTaskOrder.Orders.ServiceMemberID, fmt.Sprintf("ServiceMember on MoveTaskOrder (ID: %s) missing Rank", moveTaskOrder.ID))
	}
	// Verify Affiliation
	if serviceMember.Affiliation == nil || *serviceMember.Affiliation == "" {
		return nil, services.NewConflictError(moveTaskOrder.Orders.ServiceMemberID, fmt.Sprintf("ServiceMember on MoveTaskOrder (ID: %s) missing Affiliation", moveTaskOrder.ID))
	}

	// Verify that there were no previous requests that were marked as final
	var finalPaymentRequests models.PaymentRequests
	count, err := tx.Q().Where("move_id = $1 AND is_final = TRUE AND status <> $2", paymentRequest.MoveTaskOrderID, models.PaymentRequestStatusReviewedAllRejected).Count(&finalPaymentRequests)

	if err != nil {
		return nil, services.NewQueryError("PaymentRequests", err, fmt.Sprintf("Error while querying final payment request for MTO %s: %s", paymentRequest.MoveTaskOrderID, err.Error()))
	}

	if count != 0 {
		return nil, services.NewInvalidInputError(moveTaskOrder.ID, nil, nil, fmt.Sprintf("Cannot create PaymentRequest because a final PaymentRequest has already been submitted for MoveTaskOrder (ID: %s)", moveTaskOrder.ID))
	}

	// Update PaymentRequest
	paymentRequest.MoveTaskOrder = moveTaskOrder
	paymentRequest.Status = models.PaymentRequestStatusPending
	paymentRequest.RequestedAt = requestedAt

	uniqueIdentifier, sequenceNumber, err := p.makeUniqueIdentifier(tx, moveTaskOrder)
	if err != nil {
		errMsg := fmt.Sprintf("issue creating payment request unique identifier: %s", err.Error())
		return nil, services.NewInvalidCreateInputError(nil, errMsg)
	}
	paymentRequest.PaymentRequestNumber = uniqueIdentifier
	paymentRequest.SequenceNumber = sequenceNumber

	// Create the payment request for the database
	verrs, err := tx.ValidateAndCreate(paymentRequest)
	if verrs.HasAny() {
		msg := fmt.Sprint("validation error creating payment request")
		return nil, services.NewInvalidCreateInputError(verrs, msg)
	}
	if err != nil {
		return nil, fmt.Errorf("failure creating payment request: %w for %s", err, paymentRequest.ID.String())
	}

	return paymentRequest, nil
}

func (p *paymentRequestCreator) createPaymentServiceItem(tx *pop.Connection, paymentServiceItem models.PaymentServiceItem, paymentRequest *models.PaymentRequest, requestedAt time.Time) (models.PaymentServiceItem, models.MTOServiceItem, error) {
	// Verify that the MTO service item ID exists
	var mtoServiceItem models.MTOServiceItem
	err := tx.Eager("ReService").Find(&mtoServiceItem, paymentServiceItem.MTOServiceItemID)
	if err != nil {
		if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			msg := fmt.Sprint("for MTO Service Item")
			return models.PaymentServiceItem{}, models.MTOServiceItem{}, services.NewNotFoundError(paymentServiceItem.MTOServiceItemID, msg)
		}
		return paymentServiceItem, models.MTOServiceItem{}, fmt.Errorf("could not fetch MTOServiceItem with ID [%s]: %w", paymentServiceItem.MTOServiceItemID.String(), err)
	}

	paymentServiceItem.MTOServiceItemID = mtoServiceItem.ID
	paymentServiceItem.MTOServiceItem = mtoServiceItem
	paymentServiceItem.PaymentRequestID = paymentRequest.ID
	paymentServiceItem.PaymentRequest = *paymentRequest
	paymentServiceItem.Status = models.PaymentServiceItemStatusRequested
	// No pricing at this point, so skipping the PriceCents field.
	paymentServiceItem.RequestedAt = requestedAt

	verrs, err := tx.ValidateAndCreate(&paymentServiceItem)
	if verrs.HasAny() {
		msg := fmt.Sprint("validation error creating payment request service item in payment request creation")
		return paymentServiceItem, mtoServiceItem, services.NewInvalidCreateInputError(verrs, msg)
	}
	if err != nil {
		return paymentServiceItem, mtoServiceItem, fmt.Errorf("failure creating payment service item: %w for MTO Service Item ID <%s>", err, paymentServiceItem.MTOServiceItemID.String())
	}

	return paymentServiceItem, mtoServiceItem, nil
}

func (p *paymentRequestCreator) pricePaymentServiceItem(tx *pop.Connection, pricer services.ServiceItemPricer, paymentServiceItem models.PaymentServiceItem) (models.PaymentServiceItem, models.PaymentServiceItemParams, error) {
	price, displayParams, err := pricer.PriceServiceItem(paymentServiceItem)
	if err != nil {
		// If a pricer isn't implemented yet, just skip saving any pricing for now.
		// TODO: Once all pricers are implemented, this should be removed.
		if _, ok := err.(services.NotImplementedError); ok {
			return paymentServiceItem, displayParams, nil
		}

		return models.PaymentServiceItem{}, displayParams, err
	}

	paymentServiceItem.PriceCents = &price

	verrs, err := tx.ValidateAndUpdate(&paymentServiceItem)
	if verrs.HasAny() {
		return models.PaymentServiceItem{}, displayParams, services.NewInvalidInputError(paymentServiceItem.ID, err, verrs, "")
	}
	if err != nil {
		return models.PaymentServiceItem{}, displayParams, fmt.Errorf("could not update payment service item for MTO service item ID %s: %w",
			paymentServiceItem.ID, err)
	}

	return paymentServiceItem, displayParams, nil
}

func (p *paymentRequestCreator) createPaymentServiceItemParam(tx *pop.Connection, paymentServiceItemParam models.PaymentServiceItemParam, paymentServiceItem models.PaymentServiceItem) (models.PaymentServiceItemParam, *string, *string, error) {
	// If the ServiceItemParamKeyID is provided, verify it exists; otherwise, lookup
	// via the IncomingKey field
	var key, value string
	createParam := false
	var serviceItemParamKey models.ServiceItemParamKey
	if paymentServiceItemParam.ServiceItemParamKeyID != uuid.Nil {
		err := tx.Find(&serviceItemParamKey, paymentServiceItemParam.ServiceItemParamKeyID)
		if err != nil {
			if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
				msg := fmt.Sprintf("Service Item Param Key ID")
				return models.PaymentServiceItemParam{}, nil, nil, services.NewNotFoundError(paymentServiceItemParam.ServiceItemParamKeyID, msg)
			}
			return models.PaymentServiceItemParam{}, nil, nil, fmt.Errorf("could not fetch ServiceItemParamKey with ID [%s]: %w", paymentServiceItemParam.ServiceItemParamKeyID, err)
		}
		key = serviceItemParamKey.Key.String()
		value = paymentServiceItemParam.Value
		createParam = true
	} else if paymentServiceItemParam.IncomingKey != "" {
		err := tx.Where("key = ?", paymentServiceItemParam.IncomingKey).First(&serviceItemParamKey)
		if err != nil {
			if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
				errorString := fmt.Sprintf("Service Item Param Key %s: %s", paymentServiceItemParam.IncomingKey, models.ErrFetchNotFound)
				return models.PaymentServiceItemParam{}, nil, nil, services.NewNotFoundError(uuid.Nil, errorString)
			}
			return models.PaymentServiceItemParam{}, nil, nil, fmt.Errorf("could not retrieve param key [%s]: %w", paymentServiceItemParam.IncomingKey, err)
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

		var err error
		verrs, err := tx.ValidateAndCreate(&paymentServiceItemParam)
		if verrs.HasAny() {
			msg := fmt.Sprint("validation error creating payment service item param in payment request creation")
			return models.PaymentServiceItemParam{}, nil, nil, services.NewInvalidCreateInputError(verrs, msg)
		}
		if err != nil {
			return models.PaymentServiceItemParam{}, nil, nil, fmt.Errorf("failure creating payment service item param: %w for Payment Service Item ID <%s> Service Item Param Key <%s>", err, paymentServiceItem.ID.String(), serviceItemParamKey.Key)
		}

		return paymentServiceItemParam, &key, &value, nil
	}

	// incoming provided paymentServiceItemParam is empty
	return models.PaymentServiceItemParam{}, nil, nil, nil
}

func (p *paymentRequestCreator) createServiceItemParamFromLookup(tx *pop.Connection, paramLookup *serviceparamlookups.ServiceItemParamKeyData, serviceParam models.ServiceParam, paymentServiceItem models.PaymentServiceItem) (*models.PaymentServiceItemParam, error) {
	// Pricing/pricer functions will create the params originating from pricers. Nothing to do here.
	if serviceParam.ServiceItemParamKey.Origin == models.ServiceItemParamOriginPricer {
		return &models.PaymentServiceItemParam{}, nil
	}

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
		IncomingKey:           serviceParam.ServiceItemParamKey.Key.String(),
		Value:                 value,
	}

	var verrs *validate.Errors
	verrs, err = tx.ValidateAndCreate(&paymentServiceItemParam)
	if verrs.HasAny() {
		msg := fmt.Sprintf("validation error creating payment service item param: for payment service item ID <%s> and service item key <%s>", paymentServiceItem.ID.String(), serviceParam.ServiceItemParamKey.Key)
		return nil, services.NewInvalidCreateInputError(verrs, msg)
	}

	if err != nil {
		return nil, fmt.Errorf("failure creating payment service item param: %w for payment service item ID <%s> and service item key <%s>", err, paymentServiceItem.ID.String(), serviceParam.ServiceItemParamKey.Key)
	}

	return &paymentServiceItemParam, nil
}

func (p *paymentRequestCreator) makeUniqueIdentifier(tx *pop.Connection, mto models.Move) (string, int, error) {
	if mto.ReferenceID == nil || *mto.ReferenceID == "" {
		errMsg := fmt.Sprintf("MTO %s has missing ReferenceID", mto.ID.String())
		return "", 0, errors.New(errMsg)
	}
	// Get the max sequence number that exists for the payment requests associated with the given MTO.
	// Since we have a lock to prevent concurrent payment requests for this MTO, this should be safe.
	var max int
	err := tx.RawQuery("SELECT COALESCE(MAX(sequence_number),0) FROM payment_requests WHERE move_id = $1", mto.ID).First(&max)
	if err != nil {
		return "", 0, fmt.Errorf("max sequence_number for MoveTaskOrderID [%s] failed: %w", mto.ID, err)
	}

	nextSequence := max + 1
	paymentRequestNumber := fmt.Sprintf("%s-%d", *mto.ReferenceID, nextSequence)

	return paymentRequestNumber, nextSequence, nil
}
