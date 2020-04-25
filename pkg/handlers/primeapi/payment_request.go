package primeapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	paymentrequestop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/payment_requests"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// CreatePaymentRequestHandler is the handler for creating payment requests
type CreatePaymentRequestHandler struct {
	handlers.HandlerContext
	services.PaymentRequestCreator
}

// Handle creates the payment request
func (h CreatePaymentRequestHandler) Handle(params paymentrequestop.CreatePaymentRequestParams) middleware.Responder {
	// TODO: authorization to create payment request

	logger := h.LoggerFromRequest(params.HTTPRequest)

	payload := params.Body

	if payload == nil {
		errPayload := &primemessages.ClientError{
			Title:    handlers.FmtString(handlers.SQLErrMessage),
			Detail:   handlers.FmtString("Invalid payment request: params Body is nil"),
			Instance: handlers.FmtUUID(h.GetTraceID()),
		}
		logger.Info("Payment Request",
			zap.Any("payload", errPayload))
		logger.Error("Invalid payment request: params Body is nil")
		return paymentrequestop.NewCreatePaymentRequestBadRequest().WithPayload(errPayload)
	}

	logger.Info("primeapi.CreatePaymentRequestHandler info", zap.String("pointOfContact", params.Body.PointOfContact))

	moveTaskOrderIDString := payload.MoveTaskOrderID.String()
	mtoID, err := uuid.FromString(moveTaskOrderIDString)
	if err != nil {
		logger.Error("Invalid payment request: params MoveTaskOrderID cannot be converted to a UUID",
			zap.String("MoveTaskOrderID", moveTaskOrderIDString), zap.Error(err))
		errPayload := &primemessages.ClientError{
			Title:    handlers.FmtString(handlers.SQLErrMessage),
			Detail:   handlers.FmtString(err.Error()),
			Instance: handlers.FmtUUID(h.GetTraceID()),
		}
		logger.Info("Payment Request",
			zap.Any("payload", payload))
		return paymentrequestop.NewCreatePaymentRequestBadRequest().WithPayload(errPayload)
	}

	isFinal := false
	if payload.IsFinal != nil {
		isFinal = *payload.IsFinal
	}

	paymentRequest := models.PaymentRequest{
		IsFinal:         isFinal,
		MoveTaskOrderID: mtoID,
	}

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
		if typedErr, ok := err.(services.InvalidCreateInputError); ok {
			verrs := typedErr.ValidationErrors
			payload := &primemessages.ValidationError{
				InvalidFields: handlers.NewValidationErrorsResponse(verrs).Errors,
			}

			payload.Title = handlers.FmtString(handlers.ValidationErrMessage)
			payload.Detail = handlers.FmtString(err.Error())
			payload.Instance = handlers.FmtUUID(h.GetTraceID())
			logger.Info("Payment Request",
				zap.Any("payload", payload))
			return paymentrequestop.NewCreatePaymentRequestUnprocessableEntity().WithPayload(payload)
		}

		if _, ok := err.(services.NotFoundError); ok {
			payload := &primemessages.ClientError{
				Title:    handlers.FmtString(handlers.NotFoundMessage),
				Detail:   handlers.FmtString(err.Error()),
				Instance: handlers.FmtUUID(h.GetTraceID()),
			}
			logger.Info("Payment Request",
				zap.Any("payload", payload))
			return paymentrequestop.NewCreatePaymentRequestNotFound().WithPayload(payload)
		}
		if _, ok := err.(*services.BadDataError); ok {
			payload := &primemessages.ClientError{
				Title:    handlers.FmtString(handlers.SQLErrMessage),
				Detail:   handlers.FmtString(err.Error()),
				Instance: handlers.FmtUUID(h.GetTraceID()),
			}
			logger.Info("Payment Request",
				zap.Any("payload", payload))
			return paymentrequestop.NewCreatePaymentRequestBadRequest().WithPayload(payload)
		}
		logger.Info("Payment Request",
			zap.Any("payload", payload))
		return paymentrequestop.NewCreatePaymentRequestInternalServerError()
	}

	returnPayload := payloads.PaymentRequest(createdPaymentRequest)
	return paymentrequestop.NewCreatePaymentRequestCreated().WithPayload(returnPayload)
}

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
