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

	paymentRequest.PaymentServiceItems, err = h.buildPaymentServiceItems(payload)
	if err != nil {
		logger.Error("could not build service items", zap.Error(err))
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

func (h CreatePaymentRequestHandler) buildPaymentServiceItems(payload *primemessages.CreatePaymentRequestPayload) (models.PaymentServiceItems, error) {
	var paymentServiceItems models.PaymentServiceItems

	for _, payloadServiceItem := range payload.ServiceItems {
		serviceItemID, err := uuid.FromString(payloadServiceItem.ID.String())
		if err != nil {
			return nil, fmt.Errorf("could not convert service item ID [%v] to UUID: %w", payloadServiceItem.ID, err)
		}

		paymentServiceItem := models.PaymentServiceItem{
			// The rest of the model will be filled in when the payment request is created
			ServiceItemID: serviceItemID,
		}

		paymentServiceItem.PaymentServiceItemParams = h.buildPaymentServiceItemParams(payloadServiceItem)

		paymentServiceItems = append(paymentServiceItems, paymentServiceItem)
	}

	return paymentServiceItems, nil
}

func (h CreatePaymentRequestHandler) buildPaymentServiceItemParams(payloadServiceItem *primemessages.ServiceItem) models.PaymentServiceItemParams {
	var paymentServiceItemParams models.PaymentServiceItemParams

	for _, payloadServiceItemParam := range payloadServiceItem.Params {
		paymentServiceItemParam := models.PaymentServiceItemParam{
			// ID and PaymentServiceItemID to be filled in when payment request is created
			IncomingKey: payloadServiceItemParam.Key,
			Value:       payloadServiceItemParam.Value,
		}

		paymentServiceItemParams = append(paymentServiceItemParams, paymentServiceItemParam)
	}

	return paymentServiceItemParams
}
