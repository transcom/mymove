package paymentrequest

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	serviceparamlookups "github.com/transcom/mymove/pkg/payment_request/service_param_value_lookups"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
)

type paymentRequestCreator struct {
	planner route.Planner
	pricer  services.ServiceItemPricer
}

// NewPaymentRequestCreator returns a new payment request creator
func NewPaymentRequestCreator(planner route.Planner, pricer services.ServiceItemPricer) services.PaymentRequestCreator {
	return &paymentRequestCreator{
		planner: planner,
		pricer:  pricer,
	}
}

func (p *paymentRequestCreator) CreatePaymentRequest(appCtx appcontext.AppContext, paymentRequestArg *models.PaymentRequest) (*models.PaymentRequest, error) {
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		var err error
		now := time.Now()

		// Gather information for logging
		mtoMessageString := " MTO ID <" + paymentRequestArg.MoveTaskOrderID.String() + ">"
		prMessageString := " paymentRequestID <" + paymentRequestArg.ID.String() + ">"

		// Create the payment request
		paymentRequestArg, err = p.createPaymentRequestSaveToDB(txnAppCtx, paymentRequestArg, now)

		if err != nil {
			var badDataError *apperror.BadDataError

			if _, ok := err.(apperror.InvalidCreateInputError); ok {
				return err
			}
			if _, ok := err.(apperror.NotFoundError); ok {
				return err
			}
			if _, ok := err.(apperror.ConflictError); ok {
				return err
			}
			if _, ok := err.(apperror.InvalidInputError); ok {
				return err
			}
			if _, ok := err.(apperror.QueryError); ok {
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
		serviceParamCache := serviceparamlookups.NewServiceParamsCache()

		// Track which shipments have been verified already
		shipmentIDs := make(map[uuid.UUID]bool)

		// Create a payment service item for each incoming payment service item in the payment request
		// These incoming payment service items have not been created in the database yet
		var newPaymentServiceItems models.PaymentServiceItems
		for _, paymentServiceItem := range paymentRequestArg.PaymentServiceItems {

			// check if shipment is valid for creating a payment request
			validShipmentError := p.validShipment(appCtx, paymentServiceItem.MTOServiceItem.MTOShipmentID, shipmentIDs)
			if validShipmentError != nil {
				return validShipmentError
			}

			// Gather message information for logging
			errMessageString := p.serviceItemErrorMessage(paymentServiceItem.MTOServiceItemID, paymentServiceItem.MTOServiceItem, mtoMessageString, prMessageString)

			// Create the payment service item
			var mtoServiceItem models.MTOServiceItem
			paymentServiceItem, mtoServiceItem, err = p.createPaymentServiceItem(txnAppCtx, paymentServiceItem, paymentRequestArg, now)
			if err != nil {
				if _, ok := err.(apperror.InvalidCreateInputError); ok {
					return err
				}
				if _, ok := err.(apperror.NotFoundError); ok {
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
				param, key, value, err = p.createPaymentServiceItemParam(txnAppCtx, paymentServiceItemParam, paymentServiceItem)
				if err != nil {
					if _, ok := err.(*apperror.BadDataError); ok {
						return err
					}
					if _, ok := err.(apperror.NotFoundError); ok {
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
			paymentHelper := paymentrequesthelper.RequestPaymentHelper{}

			reServiceParams, err := paymentHelper.FetchServiceParamList(txnAppCtx, paymentServiceItem.MTOServiceItem)
			if err != nil {
				errMessage := "Failed to retrieve service item param list for " + errMessageString
				return fmt.Errorf("%s err: %w", errMessage, err)
			}

			// Get values for needed service item params (do lookups)
			paramLookup, err := serviceparamlookups.ServiceParamLookupInitialize(txnAppCtx, p.planner, paymentServiceItem.MTOServiceItem, paymentRequestArg.ID, paymentRequestArg.MoveTaskOrderID, &serviceParamCache)
			if err != nil {
				return err
			}
			for _, reServiceParam := range reServiceParams {
				if _, found := incomingMTOServiceItemParams[reServiceParam.ServiceItemParamKey.Key.String()]; !found {
					// create the missing service item param
					var param *models.PaymentServiceItemParam
					param, err = p.createServiceItemParamFromLookup(txnAppCtx, paramLookup, reServiceParam, paymentServiceItem)
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
			if !validParamList {
				errMessage := "service item param list is not valid (will not be able to price the item) " + validateMessage + " for " + errMessageString
				return fmt.Errorf("%s err: %w", errMessage, err)
			}

			// Price the payment service item
			var psItem models.PaymentServiceItem
			var displayParams models.PaymentServiceItemParams
			psItem, displayParams, err = p.pricePaymentServiceItem(txnAppCtx, paymentServiceItem)
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

func (p *paymentRequestCreator) serviceItemErrorMessage(
	mtoServiceItemID uuid.UUID,
	mtoServiceItem models.MTOServiceItem,
	mtoMessageString string,
	prMessageString string,
) string {
	mtoServiceItemString := " MTO Service item ID <" + mtoServiceItemID.String() + ">"
	reServiceItem := mtoServiceItem.ReService
	serviceItemMessageString := " RE Service Item Code: <" + string(reServiceItem.Code) + "> Name: <" + reServiceItem.Name + ">"
	return mtoMessageString + prMessageString + mtoServiceItemString + serviceItemMessageString
}

func (p *paymentRequestCreator) validShipment(appCtx appcontext.AppContext, shipmentID *uuid.UUID, shipmentIDs map[uuid.UUID]bool) error {
	if shipmentID != nil {
		if _, found := shipmentIDs[*shipmentID]; !found {
			shipmentIDs[*shipmentID] = true
			var mtoShipment models.MTOShipment
			err := appCtx.DB().Find(&mtoShipment, *shipmentID)
			if err != nil {
				switch err {
				case sql.ErrNoRows:
					appCtx.Logger().Error(fmt.Sprintf("paymentRequestCreator.validShipment:MTOShipmentID %s not found", shipmentID.String()), zap.Error(err))
					return apperror.NewNotFoundError(*shipmentID, "for MTOShipment")
				default:
					appCtx.Logger().Error(fmt.Sprintf("paymentRequestCreator.validShipment: query error MTOShipmentID %s", shipmentID.String()), zap.Error(err))
					return apperror.NewQueryError("MTOShipment", err, fmt.Sprintf("paymentRequestCreator.validShipment:MTOShipmentID %s", shipmentID.String()))
				}
			}
			if mtoShipment.UsesExternalVendor {
				appCtx.Logger().Error("paymentRequestCreator.validShipment",
					zap.Any("mtoShipment.UsesExternalVendor", mtoShipment.UsesExternalVendor),
					zap.String("MTOShipmentID", shipmentID.String()))
				return apperror.NewConflictError(*shipmentID, fmt.Sprintf("paymentRequestCreator.validShipment: Shipment uses external vendor for MTOShipmentID %s", shipmentID.String()))
			}
		}
	}
	return nil
}

func (p *paymentRequestCreator) createPaymentRequestSaveToDB(appCtx appcontext.AppContext, paymentRequest *models.PaymentRequest, requestedAt time.Time) (*models.PaymentRequest, error) {
	// Verify that the MTO ID exists
	if paymentRequest.MoveTaskOrderID == uuid.Nil {
		return nil, apperror.NewInvalidCreateInputError(nil, "Invalid Create Input Error: MoveTaskOrderID is required on PaymentRequest create")
	}

	// Lock on the parent row to keep multiple transactions from getting this count at the same time
	// for the same move_id.  This should block if another payment request comes in for the
	// same move_id.  Payment requests for other move_ids should run concurrently.
	// Also note that we use "FOR NO KEY UPDATE" to allow concurrent mods to other tables that have a
	// FK to move_task_orders.
	var moveTaskOrder models.Move
	sqlString, sqlArgs := appCtx.DB().Where("id = $1", paymentRequest.MoveTaskOrderID).ToSQL(&pop.Model{Value: &moveTaskOrder})
	sqlString += " FOR NO KEY UPDATE"
	err := appCtx.DB().RawQuery(sqlString, sqlArgs...).First(&moveTaskOrder)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(paymentRequest.MoveTaskOrderID, "for Move")
		default:
			return nil, apperror.NewQueryError("Move", err, fmt.Sprintf("could not retrieve Move with ID [%s]", paymentRequest.MoveTaskOrderID))
		}
	}

	// Verify the Orders on the MTO
	err = appCtx.DB().Load(&moveTaskOrder, "Orders")

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(moveTaskOrder.OrdersID, fmt.Sprintf("Orders on MoveTaskOrder (ID: %s) missing", moveTaskOrder.ID))
		default:
			return nil, apperror.NewQueryError("Orders", err, "")
		}
	}

	// Verify that the Orders has LOA
	if moveTaskOrder.Orders.TAC == nil || *moveTaskOrder.Orders.TAC == "" {
		return nil, apperror.NewConflictError(moveTaskOrder.OrdersID, fmt.Sprintf("Orders on MoveTaskOrder (ID: %s) missing Lines of Accounting TAC", moveTaskOrder.ID))
	}
	// Verify that the Orders have OriginDutyLocation
	if moveTaskOrder.Orders.OriginDutyLocationID == nil {
		return nil, apperror.NewConflictError(moveTaskOrder.OrdersID, fmt.Sprintf("Orders on MoveTaskOrder (ID: %s) missing OriginDutyLocation", moveTaskOrder.ID))
	}
	// Verify that ServiceMember is Valid
	err = appCtx.DB().Load(&moveTaskOrder.Orders, "ServiceMember")
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(moveTaskOrder.Orders.ServiceMemberID, fmt.Sprintf("ServiceMember on MoveTaskOrder (ID: %s) not valid", moveTaskOrder.ID))
		default:
			return nil, apperror.NewQueryError("ServiceMember", err, "")
		}
	}

	serviceMember := moveTaskOrder.Orders.ServiceMember
	// Verify First Name
	if serviceMember.FirstName == nil || *serviceMember.FirstName == "" {
		return nil, apperror.NewConflictError(moveTaskOrder.Orders.ServiceMemberID, fmt.Sprintf("ServiceMember on MoveTaskOrder (ID: %s) missing First Name", moveTaskOrder.ID))
	}
	// Verify Last Name
	if serviceMember.LastName == nil || *serviceMember.LastName == "" {
		return nil, apperror.NewConflictError(moveTaskOrder.Orders.ServiceMemberID, fmt.Sprintf("ServiceMember on MoveTaskOrder (ID: %s) missing Last Name", moveTaskOrder.ID))
	}
	// Verify Rank
	if serviceMember.Rank == nil || *serviceMember.Rank == "" {
		return nil, apperror.NewConflictError(moveTaskOrder.Orders.ServiceMemberID, fmt.Sprintf("ServiceMember on MoveTaskOrder (ID: %s) missing Rank", moveTaskOrder.ID))
	}
	// Verify Affiliation
	if serviceMember.Affiliation == nil || *serviceMember.Affiliation == "" {
		return nil, apperror.NewConflictError(moveTaskOrder.Orders.ServiceMemberID, fmt.Sprintf("ServiceMember on MoveTaskOrder (ID: %s) missing Affiliation", moveTaskOrder.ID))
	}

	// Verify that there were no previous requests that were marked as final
	var finalPaymentRequests models.PaymentRequests
	count, err := appCtx.DB().Q().Where("move_id = $1 AND is_final = TRUE AND status <> $2", paymentRequest.MoveTaskOrderID, models.PaymentRequestStatusReviewedAllRejected).Count(&finalPaymentRequests)

	if err != nil {
		return nil, apperror.NewQueryError("PaymentRequests", err, fmt.Sprintf("Error while querying final payment request for MTO %s: %s", paymentRequest.MoveTaskOrderID, err.Error()))
	}

	if count != 0 {
		return nil, apperror.NewInvalidInputError(moveTaskOrder.ID, nil, nil, fmt.Sprintf("Cannot create PaymentRequest because a final PaymentRequest has already been submitted for MoveTaskOrder (ID: %s)", moveTaskOrder.ID))
	}

	// Update PaymentRequest
	paymentRequest.MoveTaskOrder = moveTaskOrder
	paymentRequest.Status = models.PaymentRequestStatusPending
	paymentRequest.RequestedAt = requestedAt

	uniqueIdentifier, sequenceNumber, err := p.makeUniqueIdentifier(appCtx, moveTaskOrder)
	if err != nil {
		errMsg := fmt.Sprintf("issue creating payment request unique identifier: %s", err.Error())
		return nil, apperror.NewInvalidCreateInputError(nil, errMsg)
	}
	paymentRequest.PaymentRequestNumber = uniqueIdentifier
	paymentRequest.SequenceNumber = sequenceNumber

	// Create the payment request for the database
	verrs, err := appCtx.DB().ValidateAndCreate(paymentRequest)
	if verrs.HasAny() {
		msg := "validation error creating payment request"
		return nil, apperror.NewInvalidCreateInputError(verrs, msg)
	}
	if err != nil {
		return nil, fmt.Errorf("failure creating payment request: %w for %s", err, paymentRequest.ID.String())
	}

	return paymentRequest, nil
}

func (p *paymentRequestCreator) createPaymentServiceItem(appCtx appcontext.AppContext, paymentServiceItem models.PaymentServiceItem, paymentRequest *models.PaymentRequest, requestedAt time.Time) (models.PaymentServiceItem, models.MTOServiceItem, error) {
	// Verify that the MTO service item ID exists
	var mtoServiceItem models.MTOServiceItem
	err := appCtx.DB().Eager("ReService").Find(&mtoServiceItem, paymentServiceItem.MTOServiceItemID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return models.PaymentServiceItem{}, models.MTOServiceItem{}, apperror.NewNotFoundError(paymentServiceItem.MTOServiceItemID, "for MTO Service Item")
		default:
			return models.PaymentServiceItem{}, models.MTOServiceItem{}, apperror.NewQueryError("MTOServiceItem", err, fmt.Sprintf("could not fetch MTOServiceItem with ID [%s]", paymentServiceItem.MTOServiceItemID.String()))
		}
	}

	paymentServiceItem.MTOServiceItemID = mtoServiceItem.ID
	paymentServiceItem.MTOServiceItem = mtoServiceItem
	paymentServiceItem.PaymentRequestID = paymentRequest.ID
	paymentServiceItem.PaymentRequest = *paymentRequest
	paymentServiceItem.Status = models.PaymentServiceItemStatusRequested
	// No pricing at this point, so skipping the PriceCents field.
	paymentServiceItem.RequestedAt = requestedAt

	verrs, err := appCtx.DB().ValidateAndCreate(&paymentServiceItem)
	if verrs.HasAny() {
		msg := "validation error creating payment request service item in payment request creation"
		return paymentServiceItem, mtoServiceItem, apperror.NewInvalidCreateInputError(verrs, msg)
	}
	if err != nil {
		return paymentServiceItem, mtoServiceItem, fmt.Errorf("failure creating payment service item: %w for MTO Service Item ID <%s>", err, paymentServiceItem.MTOServiceItemID.String())
	}

	return paymentServiceItem, mtoServiceItem, nil
}

func (p *paymentRequestCreator) pricePaymentServiceItem(appCtx appcontext.AppContext, paymentServiceItem models.PaymentServiceItem) (models.PaymentServiceItem, models.PaymentServiceItemParams, error) {
	price, displayParams, err := p.pricer.PriceServiceItem(appCtx, paymentServiceItem)
	if err != nil {
		// If a pricer isn't implemented yet, just skip saving any pricing for now.
		// TODO: Once all pricers are implemented, this should be removed.
		if _, ok := err.(apperror.NotImplementedError); ok {
			return paymentServiceItem, displayParams, nil
		}

		return models.PaymentServiceItem{}, displayParams, err
	}

	paymentServiceItem.PriceCents = &price

	verrs, err := appCtx.DB().ValidateAndUpdate(&paymentServiceItem)
	if verrs.HasAny() {
		return models.PaymentServiceItem{}, displayParams, apperror.NewInvalidInputError(paymentServiceItem.ID, err, verrs, "")
	}
	if err != nil {
		return models.PaymentServiceItem{}, displayParams, fmt.Errorf("could not update payment service item for MTO service item ID %s: %w",
			paymentServiceItem.ID, err)
	}

	return paymentServiceItem, displayParams, nil
}

func (p *paymentRequestCreator) createPaymentServiceItemParam(appCtx appcontext.AppContext, paymentServiceItemParam models.PaymentServiceItemParam, paymentServiceItem models.PaymentServiceItem) (models.PaymentServiceItemParam, *string, *string, error) {
	/* Note that we are not validating the param key type here.
	 * For now, invalid params will be caught when they are parsed in lookups.
	 * In the future we may want to add more validation here to catch things earlier.
	 */

	// If the ServiceItemParamKeyID is provided, verify it exists; otherwise, lookup
	// via the IncomingKey field
	var key, value string
	createParam := false
	var serviceItemParamKey models.ServiceItemParamKey
	if paymentServiceItemParam.ServiceItemParamKeyID != uuid.Nil {
		err := appCtx.DB().Find(&serviceItemParamKey, paymentServiceItemParam.ServiceItemParamKeyID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return models.PaymentServiceItemParam{}, nil, nil, apperror.NewNotFoundError(paymentServiceItemParam.ServiceItemParamKeyID, "Service Item Param Key ID")
			default:
				return models.PaymentServiceItemParam{}, nil, nil, apperror.NewQueryError("ServiceItemParamKey", err, fmt.Sprintf("could not fetch ServiceItemParamKey with ID [%s]", paymentServiceItemParam.ServiceItemParamKeyID))
			}
		}
		key = serviceItemParamKey.Key.String()
		value = paymentServiceItemParam.Value
		createParam = true
	} else if paymentServiceItemParam.IncomingKey != "" {
		err := appCtx.DB().Where("key = ?", paymentServiceItemParam.IncomingKey).First(&serviceItemParamKey)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				errorString := fmt.Sprintf("Service Item Param Key %s: %s", paymentServiceItemParam.IncomingKey, models.ErrFetchNotFound)
				return models.PaymentServiceItemParam{}, nil, nil, apperror.NewNotFoundError(uuid.Nil, errorString)
			default:
				return models.PaymentServiceItemParam{}, nil, nil, apperror.NewQueryError("ServiceItemParamKey", err, fmt.Sprintf("could not retrieve param key [%s]", paymentServiceItemParam.IncomingKey))
			}
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
		verrs, err := appCtx.DB().ValidateAndCreate(&paymentServiceItemParam)
		if verrs.HasAny() {
			msg := "validation error creating payment service item param in payment request creation"
			return models.PaymentServiceItemParam{}, nil, nil, apperror.NewInvalidCreateInputError(verrs, msg)
		}
		if err != nil {
			return models.PaymentServiceItemParam{}, nil, nil, fmt.Errorf("failure creating payment service item param: %w for Payment Service Item ID <%s> Service Item Param Key <%s>", err, paymentServiceItem.ID.String(), serviceItemParamKey.Key)
		}

		return paymentServiceItemParam, &key, &value, nil
	}

	// incoming provided paymentServiceItemParam is empty
	return models.PaymentServiceItemParam{}, nil, nil, nil
}

func (p *paymentRequestCreator) createServiceItemParamFromLookup(appCtx appcontext.AppContext, paramLookup *serviceparamlookups.ServiceItemParamKeyData, serviceParam models.ServiceParam, paymentServiceItem models.PaymentServiceItem) (*models.PaymentServiceItemParam, error) {
	// Pricing/pricer functions will create the params originating from pricers. Nothing to do here.
	if serviceParam.ServiceItemParamKey.Origin == models.ServiceItemParamOriginPricer {
		return nil, nil
	}

	// key not found in map
	// Did not find service item param needed for pricing, add it to the list
	value, err := paramLookup.ServiceParamValue(appCtx, serviceParam.ServiceItemParamKey.Key)
	if err != nil {
		errMessage := "Failed to lookup ServiceParamValue for param key <" + serviceParam.ServiceItemParamKey.Key + "> "
		return nil, fmt.Errorf("%s err: %w", errMessage, err)
	}

	// Some params are considered optional.  If this is an optional param and the value is an empty string,
	// do not try to save to the database.
	if value == "" && serviceParam.IsOptional {
		return nil, nil
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
	verrs, err = appCtx.DB().ValidateAndCreate(&paymentServiceItemParam)
	if verrs.HasAny() {
		msg := fmt.Sprintf("validation error creating payment service item param: for payment service item ID <%s> and service item key <%s>", paymentServiceItem.ID.String(), serviceParam.ServiceItemParamKey.Key)
		return nil, apperror.NewInvalidCreateInputError(verrs, msg)
	}

	if err != nil {
		return nil, fmt.Errorf("failure creating payment service item param: %w for payment service item ID <%s> and service item key <%s>", err, paymentServiceItem.ID.String(), serviceParam.ServiceItemParamKey.Key)
	}

	return &paymentServiceItemParam, nil
}

func (p *paymentRequestCreator) makeUniqueIdentifier(appCtx appcontext.AppContext, mto models.Move) (string, int, error) {
	if mto.ReferenceID == nil || *mto.ReferenceID == "" {
		errMsg := fmt.Sprintf("MTO %s has missing ReferenceID", mto.ID.String())
		return "", 0, errors.New(errMsg)
	}
	// Get the max sequence number that exists for the payment requests associated with the given MTO.
	// Since we have a lock to prevent concurrent payment requests for this MTO, this should be safe.
	var max int
	err := appCtx.DB().RawQuery("SELECT COALESCE(MAX(sequence_number),0) FROM payment_requests WHERE move_id = $1", mto.ID).First(&max)
	if err != nil {
		return "", 0, fmt.Errorf("max sequence_number for MoveTaskOrderID [%s] failed: %w", mto.ID, err)
	}

	nextSequence := max + 1
	paymentRequestNumber := fmt.Sprintf("%s-%d", *mto.ReferenceID, nextSequence)

	return paymentRequestNumber, nextSequence, nil
}
