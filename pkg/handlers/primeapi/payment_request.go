package primeapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/services/audit"

	paymentrequestop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/payment_requests"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func payloadForPaymentRequestModel(pr models.PaymentRequest) *primemessages.PaymentRequest {
	return &primemessages.PaymentRequest{
		ID:              *handlers.FmtUUID(pr.ID),
		MoveTaskOrderID: *handlers.FmtUUID(pr.MoveTaskOrderID),
		IsFinal:         &pr.IsFinal,
		RejectionReason: pr.RejectionReason,
	}
}

type CreatePaymentRequestHandler struct {
	handlers.HandlerContext
	services.PaymentRequestCreator
}

func (h CreatePaymentRequestHandler) Handle(params paymentrequestop.CreatePaymentRequestParams) middleware.Responder {
	// TODO: authorization to create payment request

	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	payload := params.Body

	// Capture creation attempt in audit log
	_, err := audit.Capture(&payload, nil, logger, session, params.HTTPRequest)
	if err != nil {
		logger.Error("Auditing service error for payment request creation.", zap.Error(err))
		return paymentrequestop.NewCreatePaymentRequestInternalServerError()
	}

	if payload == nil {
		logger.Error("Invalid payment request: params Body is nil")
		return paymentrequestop.NewCreatePaymentRequestBadRequest()
	}

	moveTaskOrderIDString := payload.MoveTaskOrderID.String()
	mtoID, err := uuid.FromString(moveTaskOrderIDString)
	if err != nil {
		logger.Error("Invalid payment request: params MoveTaskOrderID cannot be converted to a UUID",
			zap.String("MoveTaskOrderID", moveTaskOrderIDString), zap.Error(err))
		return paymentrequestop.NewCreatePaymentRequestBadRequest()
	}

	isFinal := false
	if payload.IsFinal != nil {
		isFinal = *payload.IsFinal
	}

	paymentRequest := models.PaymentRequest{
		IsFinal:         isFinal,
		MoveTaskOrderID: mtoID,
	}

	/*
		paymentRequest.PaymentServiceItems, err = h.buildPaymentServiceItems(payload, createdPaymentRequest.ID, createdPaymentRequest.RequestedAt)
		if err != nil {
			createdPaymentRequest.Status = models.PaymentRequestStatusPending // TODO: why don't we have a DENIED or FAILED status???
			logger.Error("could not build service items for payment request", zap.Error(err))
			// TODO don't think we need to bail here, we want a record of what's happening -- return paymentrequestop.NewCreatePaymentRequestBadRequest()
		}
	*/

	// Build up the paymentRequest.PaymentServiceItems using the incoming payload to offload Swagger data coming
	// in from the API. These paymentRequest.PaymentServiceItems will be used as a temp holder to process the incoming API data
	paymentRequest.PaymentServiceItems, err = h.buildPaymentServiceItems(payload)
	if err != nil {
		logger.Error("could not build service items", zap.Error(err))
		// TODO: do not bail out before creating the payment request, we need the failed record
		//       we should create the failed record and store it as failed with a rejection
		return paymentrequestop.NewCreatePaymentRequestBadRequest()
	}

	createdPaymentRequest, err := h.PaymentRequestCreator.CreatePaymentRequest(&paymentRequest)
	if err != nil {
		logger.Error("Error creating payment request", zap.Error(err))
		return paymentrequestop.NewCreatePaymentRequestInternalServerError()
	}

	logger.Info("Payment Request params",
		zap.Any("payload", payload),
		// TODO add ProofOfService object to log
	)

	returnPayload := payloadForPaymentRequestModel(*createdPaymentRequest)
	return paymentrequestop.NewCreatePaymentRequestCreated().WithPayload(returnPayload)
}

/*
func (h CreatePaymentRequestHandler) translatePayload(payload *primemessages.CreatePaymentRequestPayload) (*CreatePaymentRequestIncomingPayload, error) {

	var tanslatedPayload CreatePaymentRequestIncomingPayload

	moveTaskOrderIDString := payload.MoveTaskOrderID.String()
	mtoID, err := uuid.FromString(moveTaskOrderIDString)
	if err != nil {
		return &tanslatedPayload, fmt.Errorf("invalid payment request: params MoveTaskOrderID cannot be converted to a UUID from MoveTaskOrderIDString %s with error %w", moveTaskOrderIDString, err)
	}
	tanslatedPayload.MoveTaskOrderID = mtoID

	isFinal := false
	if payload.IsFinal != nil {
		isFinal = *payload.IsFinal
	}
	tanslatedPayload.IsFinal = isFinal

	for _, payloadMTOServiceItem := range payload.ServiceItems {
		mtoServiceItemID, err := uuid.FromString(payloadMTOServiceItem.ID.String())
		if err != nil {
			return nil, fmt.Errorf("could not convert payload service item ID [%v] to UUID: %w", payloadMTOServiceItem.ID, err)
		}
		tanslatedPayload.IncomingPayloadMTOServiceItems.MTOServiceID = mtoServiceItemID

		for _, payloadMTOServiceItemParam := range payloadMTOServiceItem.Params {
			tanslatedPayload.IncomingPayloadMTOServiceItems.Params = append(tanslatedPayload.IncomingPayloadMTOServiceItems.Params,
				struct{
					Key string
					Value string
				}{
					Key: payloadMTOServiceItemParam.Key,
					Value: payloadMTOServiceItemParam.Value,
				})
		}
	}

	return &tanslatedPayload, nil
}

*/

/*
	TODO: This function is up in the air. The important thing is to the get the lookups satisfied, because there
   		  are two ways of looking params for a service item.
          1.) I think the first approach is to use what is provided in the payload which is already handled by CreatePaymentRequest
          2.) then, look for what is missing and fill in the rest of the params by using the service_params table
          3.) let any params that come from the payload remain and if any are missing see #2 above
		  4.) then go and find the values for all params that are NEEDED for pricing -- not necessarily all that are there
*/
/*
func (h CreatePaymentRequestHandler) buildPaymentServiceItems(payload *primemessages.CreatePaymentRequestPayload, paymentRequestID uuid.UUID, requestedAt time.Time) (models.PaymentServiceItems, error) {
	var paymentServiceItems models.PaymentServiceItems
	var serviceParamErrorsString *string

	for _, payloadMTOServiceItem := range payload.ServiceItems {
		mtoServiceItemID, err := uuid.FromString(payloadMTOServiceItem.ID.String())
		if err != nil {
			return nil, fmt.Errorf("could not convert service item ID [%v] to UUID: %w", payloadMTOServiceItem.ID, err)
		}

		paymentServiceItem := models.PaymentServiceItem{
			PaymentRequestID: paymentRequestID,
			MTOServiceItemID: mtoServiceItemID,
			RequestedAt:      requestedAt,
			Status:           models.PaymentServiceItemStatusRequested,
		}

		moveTaskOrderID, err := uuid.FromString(payload.MoveTaskOrderID.String())
		var errMessage *string
		paymentServiceItem.PaymentServiceItemParams, errMessage = h.buildPaymentServiceItemParams(payloadMTOServiceItem, paymentRequestID, moveTaskOrderID)

		if errMessage != nil {
			paymentServiceItem.RejectionReason = errMessage
			denied := string(models.PaymentServiceItemStatusDenied)
			paymentServiceItem.RejectionReason = &denied
			// we have an error to accumulate
			if serviceParamErrorsString != nil {
				*serviceParamErrorsString += *errMessage
			} else {
				serviceParamErrorsString = errMessage //TODO does this need to be string copy or is this OK?
			}
		}

		paymentServiceItems = append(paymentServiceItems, paymentServiceItem)
	}

	// Return from function
	if serviceParamErrorsString != nil {
		return paymentServiceItems, fmt.Errorf("error(s) found  while processing params for service items: %s",*serviceParamErrorsString)
	}

	return paymentServiceItems, nil
}
*/

func (h CreatePaymentRequestHandler) buildPaymentServiceItems(payload *primemessages.CreatePaymentRequestPayload) (models.PaymentServiceItems, error) {
	var paymentServiceItems models.PaymentServiceItems

	for _, payloadServiceItem := range payload.ServiceItems {
		mtoServiceItemID, err := uuid.FromString(payloadServiceItem.ID.String())
		if err != nil {
			return nil, fmt.Errorf("could not convert service item ID [%v] to UUID: %w", payloadServiceItem.ID, err)
		}

		paymentServiceItem := models.PaymentServiceItem{
			// The rest of the model will be filled in when the payment request is created
			MTOServiceItemID: mtoServiceItemID,
		}

		paymentServiceItem.PaymentServiceItemParams = h.buildPaymentServiceItemParams(payloadServiceItem)

		paymentServiceItems = append(paymentServiceItems, paymentServiceItem)
	}

	return paymentServiceItems, nil
}

func (h CreatePaymentRequestHandler) buildPaymentServiceItemParams(payloadMTOServiceItem *primemessages.ServiceItem) models.PaymentServiceItemParams {
	var paymentServiceItemParams models.PaymentServiceItemParams

	for _, payloadServiceItemParam := range payloadMTOServiceItem.Params {
		paymentServiceItemParam := models.PaymentServiceItemParam{
			// ID and PaymentServiceItemID to be filled in when payment request is created
			IncomingKey: payloadServiceItemParam.Key,
			Value:       payloadServiceItemParam.Value,
		}

		paymentServiceItemParams = append(paymentServiceItemParams, paymentServiceItemParam)
	}

	return paymentServiceItemParams
}

/* TODO Remove this function move parts to the CreatePaymentRequest payment function as needed */
/*
func (h CreatePaymentRequestHandler) buildPaymentServiceItemParams(
	payloadMTOServiceItem *primemessages.ServiceItem,
	paymentRequestID uuid.UUID,
	moveTaskOrderID uuid.UUID) (models.PaymentServiceItemParams, *string) {
	var paymentServiceItemParams models.PaymentServiceItemParams

	mtoServiceItemID, err := uuid.FromString(payloadMTOServiceItem.ID.String())
	if err != nil {
		errMesage := "Invalid payloadServiceItem.ID " + payloadMTOServiceItem.ID.String() + ", mtoID " + moveTaskOrderID.String() + ", paymentRequestID " + paymentRequestID.String()
		return paymentServiceItemParams, &errMesage;
	}

	// Find input params needed for `service_id` to be priced
	paymentRequestHelper := paymentrequesthelper.PaymentRequestHelper{DB: h.DB()}
	// TODO this fetch is wrong, need to lookup based on re_service.id NOT mto_service.id
	reServiceParams, err := paymentRequestHelper.FetchServiceParamList(mtoServiceItemID)
	if err != nil {
		errMesage := "Failed to retrieve service item param list for serviceID " + mtoServiceItemID.String() + ", mtoID " + moveTaskOrderID.String() + ", paymentRequestID " + paymentRequestID.String()
		return paymentServiceItemParams, &errMesage;
	}

	// Param Key:Value pairs coming from the create payment request payload, sent by the user when requesting payment
	// store these values with the payment request service item.
	var payloadIncomingMTOServiceItemParams map[string]string
	payloadIncomingMTOServiceItemParams = make(map[string]string)
	for _, payloadMTOServiceItemParam := range payloadMTOServiceItem.Params {
		payloadIncomingMTOServiceItemParams[payloadMTOServiceItemParam.Key] = payloadMTOServiceItemParam.Value
	}

	// Get values for needed service item params (do lookups)
	paramLookup := service_param_value_lookups.ServiceParamLookupInitialize(*payloadMTOServiceItem, mtoServiceItemID, paymentRequestID, moveTaskOrderID )

	for _, reServiceParam := range reServiceParams {
		var value string
		if _, found := payloadIncomingMTOServiceItemParams[reServiceParam.ServiceItemParamKey.Key]; found {
			value = payloadIncomingMTOServiceItemParams[reServiceParam.ServiceItemParamKey.Key]
		} else {
			// did not get param input from request payment payload
			value, err = paramLookup.ServiceParamValue(reServiceParam.ServiceItemParamKey.Key)
			if err != nil {
				errMesage := "Failed to lookup ServiceParamValue item param list for param key " + reServiceParam.ServiceItemParamKey.Key +  ", mtoServiceID " + mtoServiceItemID.String() + ", mtoID " + moveTaskOrderID.String() + ", paymentRequestID " + paymentRequestID.String()
				return paymentServiceItemParams, &errMesage;
			}
		}

		paymentServiceItemParam := models.PaymentServiceItemParam{
			// ID and PaymentServiceItemID to be filled in when payment request is created
			PaymentServiceItemID: paymentRequestID,
			ServiceItemParamKeyID: reServiceParam.ID,
			IncomingKey: reServiceParam.ServiceItemParamKey.Key,
			Value:  value,
		}

		paymentServiceItemParams = append(paymentServiceItemParams, paymentServiceItemParam)
	}

	// Now check that all of the incoming keys were paired to a payment request service item, if not, then
	// create new service item params
	for key, element := range payloadIncomingMTOServiceItemParams {
		found := false
		for _, reServiceParam := range reServiceParams {
			if reServiceParam.ServiceItemParamKey.Key == key {
				found = true
				break
			}
		}
		if !found {
			paymentServiceItemParam := models.PaymentServiceItemParam{
				// ID and PaymentServiceItemID to be filled in when payment request is created
				PaymentServiceItemID: paymentRequestID,
				ServiceItemParamKeyID: reServiceParam.ID, //TODO db lookup to get param key
				IncomingKey: reServiceParam.ServiceItemParamKey.Key, //TODO db lookup to get param key
				Value:  element,
			}
			paymentServiceItemParams = append(paymentServiceItemParams, paymentServiceItemParam)
		}
	}

	// Check that all values have been saved for service item params
	_, errMessage := paymentRequestHelper.ValidServiceParamList(mtoServiceItemID, reServiceParams, paymentServiceItemParams)

	return paymentServiceItemParams, errMessage
}




*/
