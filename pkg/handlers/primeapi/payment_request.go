package primeapi

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/services/query"

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
		logger.Error("Invalid payment request: params Body is nil")
		return paymentrequestop.NewCreatePaymentRequestBadRequest()
	}

	logger.Info("primeapi.CreatePaymentRequestHandler info", zap.String("pointOfContact", params.Body.PointOfContact))

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

// UpdatePaymentRequestStatusHandler updates payment requests status
type UpdatePaymentRequestStatusHandler struct {
	handlers.HandlerContext
	services.PaymentRequestStatusUpdater
	services.PaymentRequestFetcher
}

// Handle updates payment requests status
func (h UpdatePaymentRequestStatusHandler) Handle(params paymentrequestop.UpdatePaymentRequestStatusParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	paymentRequestID, err := uuid.FromString(params.PaymentRequestID.String())

	if err != nil {
		logger.Error(fmt.Sprintf("Error parsing payment request id: %s", params.PaymentRequestID.String()), zap.Error(err))
		return paymentrequestop.NewUpdatePaymentRequestStatusInternalServerError()
	}

	// Let's fetch the existing payment request using the PaymentRequestFetcher service object
	filter := []services.QueryFilter{query.NewQueryFilter("id", "=", paymentRequestID.String())}
	existingPaymentRequest, err := h.PaymentRequestFetcher.FetchPaymentRequest(filter)

	if err != nil {
		logger.Error(fmt.Sprintf("Error finding Payment Request for status update with ID: %s", params.PaymentRequestID.String()), zap.Error(err))
		return paymentrequestop.NewUpdatePaymentRequestStatusInternalServerError()
	}

	status := existingPaymentRequest.Status
	var reviewedDate time.Time
	var recGexDate time.Time
	var sentGexDate time.Time
	var paidAtDate time.Time

	if existingPaymentRequest.ReviewedAt != nil {
		reviewedDate = *existingPaymentRequest.ReviewedAt
	}
	if existingPaymentRequest.ReceivedByGexAt != nil {
		recGexDate = *existingPaymentRequest.ReceivedByGexAt
	}
	if existingPaymentRequest.SentToGexAt != nil {
		sentGexDate = *existingPaymentRequest.SentToGexAt
	}
	if existingPaymentRequest.PaidAt != nil {
		paidAtDate = *existingPaymentRequest.PaidAt
	}

	// Let's map the incoming status to our enumeration type
	switch params.Body.Status {
	case "PENDING":
		status = models.PaymentRequestStatusPending
	case "REVIEWED":
		status = models.PaymentRequestStatusReviewed
		reviewedDate = time.Now()
	case "SENT_TO_GEX":
		status = models.PaymentRequestStatusSentToGex
		sentGexDate = time.Now()
	case "RECEIVED_BY_GEX":
		status = models.PaymentRequestStatusReceivedByGex
		recGexDate = time.Now()
	case "PAID":
		status = models.PaymentRequestStatusPaid
		paidAtDate = time.Now()
	}

	// If we got a rejection reason let's use it
	rejectionReason := existingPaymentRequest.RejectionReason
	if params.Body.RejectionReason != nil {
		rejectionReason = params.Body.RejectionReason
	}

	paymentRequestForUpdate := models.PaymentRequest{
		ID:                   existingPaymentRequest.ID,
		MoveTaskOrder:        existingPaymentRequest.MoveTaskOrder,
		MoveTaskOrderID:      existingPaymentRequest.MoveTaskOrderID,
		IsFinal:              existingPaymentRequest.IsFinal,
		Status:               status,
		RejectionReason:      rejectionReason,
		RequestedAt:          existingPaymentRequest.RequestedAt,
		ReviewedAt:           &reviewedDate,
		SentToGexAt:          &sentGexDate,
		ReceivedByGexAt:      &recGexDate,
		PaidAt:               &paidAtDate,
		PaymentRequestNumber: existingPaymentRequest.PaymentRequestNumber,
		SequenceNumber:       existingPaymentRequest.SequenceNumber,
	}

	// And now let's save our updated model object using the PaymentRequestUpdater service object.
	updatedPaymentRequest, err := h.PaymentRequestStatusUpdater.UpdatePaymentRequestStatus(&paymentRequestForUpdate, params.IfMatch)

	if err != nil {
		switch err.(type) {
		case services.NotFoundError:
			return paymentrequestop.NewUpdatePaymentRequestStatusNotFound()
		case services.PreconditionFailedError:
			return paymentrequestop.NewUpdatePaymentRequestStatusPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			logger.Error(fmt.Sprintf("Error saving payment request status for ID: %s: %s", paymentRequestID, err))
			return paymentrequestop.NewUpdatePaymentRequestStatusInternalServerError()
		}
	}

	returnPayload := payloads.PaymentRequest(updatedPaymentRequest)
	return paymentrequestop.NewUpdatePaymentRequestStatusOK().WithPayload(returnPayload)
}
