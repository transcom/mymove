package ghcapi

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/services/query"

	paymentrequestop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_requests"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func payloadForPaymentRequestModel(pr models.PaymentRequest) *ghcmessages.PaymentRequest {

	return &ghcmessages.PaymentRequest{
		ID:              *handlers.FmtUUID(pr.ID),
		IsFinal:         &pr.IsFinal,
		RejectionReason: pr.RejectionReason,
	}
}

type ListPaymentRequestsHandler struct {
	handlers.HandlerContext
	services.PaymentRequestListFetcher
}

func (h ListPaymentRequestsHandler) Handle(params paymentrequestop.ListPaymentRequestsParams) middleware.Responder {
	// TODO: add authorizations
	logger := h.LoggerFromRequest(params.HTTPRequest)

	paymentRequests, err := h.FetchPaymentRequestList()
	if err != nil {
		logger.Error("Error listing payment requests err", zap.Error(err))
		return paymentrequestop.NewListPaymentRequestsInternalServerError()
	}

	paymentRequestsList := make(ghcmessages.PaymentRequests, len(*paymentRequests))
	for i, paymentRequest := range *paymentRequests {
		paymentRequestsList[i] = payloadForPaymentRequestModel(paymentRequest)
	}

	return paymentrequestop.NewListPaymentRequestsOK().WithPayload(paymentRequestsList)
}

type GetPaymentRequestHandler struct {
	handlers.HandlerContext
	services.PaymentRequestFetcher
}

func (h GetPaymentRequestHandler) Handle(params paymentrequestop.GetPaymentRequestParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	paymentRequestID, err := uuid.FromString(params.PaymentRequestID.String())

	if err != nil {
		logger.Error(fmt.Sprintf("Error parsing payment request id: %s", params.PaymentRequestID.String()), zap.Error(err))
		return paymentrequestop.NewGetPaymentRequestInternalServerError()
	}

	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", paymentRequestID.String())}

	paymentRequest, err := h.FetchPaymentRequest(filters)

	if err != nil {
		logger.Error(fmt.Sprintf("Error fetching Payment Request with ID: %s", params.PaymentRequestID.String()), zap.Error(err))
		return paymentrequestop.NewGetPaymentRequestInternalServerError()
	}

	if reflect.DeepEqual(paymentRequest, models.PaymentRequest{}) {
		logger.Info(fmt.Sprintf("Could not find a Payment Request with ID: %s", params.PaymentRequestID.String()))
		return paymentrequestop.NewGetPaymentRequestNotFound()
	}

	returnPayload := payloadForPaymentRequestModel(paymentRequest)
	response := paymentrequestop.NewGetPaymentRequestOK().WithPayload(returnPayload)

	return response
}

type UpdatePaymentRequestStatusHandler struct {
	handlers.HandlerContext
	services.PaymentRequestStatusUpdater
	services.PaymentRequestFetcher
}

func (h UpdatePaymentRequestStatusHandler) Handle(params paymentrequestop.UpdatePaymentRequestStatusParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	paymentRequestID, err := uuid.FromString(params.PaymentRequestID.String())

	if err != nil {
		logger.Error(fmt.Sprintf("Error parsing payment request id: %s", params.PaymentRequestID.String()), zap.Error(err))
		return paymentrequestop.NewGetPaymentRequestInternalServerError()
	}

	// Let's fetch the existing payment request using the PaymentRequestFetcher service object

	filter := []services.QueryFilter{query.NewQueryFilter("id", "=", paymentRequestID.String())}
	existingPaymentRequest, err := h.PaymentRequestFetcher.FetchPaymentRequest(filter)

	if err != nil {
		logger.Error(fmt.Sprintf("Error finding Payment Request for status update with ID: %s", params.PaymentRequestID.String()), zap.Error(err))
		return paymentrequestop.NewGetPaymentRequestInternalServerError()
	}

	// Let's map the incoming status to our enumeration type
	status := existingPaymentRequest.Status
	switch params.Body.Status {
	case "PENDING":
		status = models.PaymentRequestStatusPending
	case "REVIEWED":
		status = models.PaymentRequestStatusReviewed
	case "SENT_TO_GEX":
		status = models.PaymentRequestStatusSentToGex
	case "RECEIVED_BY_GEX":
		status = models.PaymentRequestStatusReceivedByGex
	case "PAID":
		status = models.PaymentRequestStatusPaid
	}

	// If we got a rejection reason let's use it
	rejectionReason := existingPaymentRequest.RejectionReason
	if params.Body.RejectionReason != nil {
		rejectionReason = params.Body.RejectionReason
	}

	paymentRequestForUpdate := models.PaymentRequest{
		ID:              existingPaymentRequest.ID,
		MoveTaskOrder:   existingPaymentRequest.MoveTaskOrder,
		MoveTaskOrderID: existingPaymentRequest.MoveTaskOrderID,
		IsFinal:         existingPaymentRequest.IsFinal,
		Status:          status,
		RejectionReason: rejectionReason,
		RequestedAt:     existingPaymentRequest.RequestedAt,
		ReviewedAt:      existingPaymentRequest.ReviewedAt,
		SentToGexAt:     existingPaymentRequest.SentToGexAt,
		ReceivedByGexAt: existingPaymentRequest.ReceivedByGexAt,
		PaidAt:          existingPaymentRequest.PaidAt,
	}

	// And now let's save our updated model object using the PaymentRequestUpdater service object.
	verrs, err := h.PaymentRequestStatusUpdater.UpdatePaymentRequestStatus(&paymentRequestForUpdate)
	if err != nil {
		if errors.Cause(err).Error() == "sql: no rows in result set" {
			payload := payloadForClientError("Unknown UUID(s)", "Unknown UUID(s) used to update a payment request ", h.GetTraceID())
			return paymentrequestop.NewUpdatePaymentRequestStatusNotFound().WithPayload(payload)
		}
		logger.Error(fmt.Sprintf("Error saving payment request status for ID: %s", paymentRequestID))
		return paymentrequestop.NewUpdatePaymentRequestStatusInternalServerError()
	}
	if verrs != nil {
		logger.Error(fmt.Sprintf("Validation error saving payment request status for ID: %s", paymentRequestID))
		return paymentrequestop.NewUpdatePaymentRequestStatusInternalServerError()
	}

	returnPayload := payloadForPaymentRequestModel(paymentRequestForUpdate)
	return paymentrequestop.NewUpdatePaymentRequestStatusOK().WithPayload(returnPayload)
}