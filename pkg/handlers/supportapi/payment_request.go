package supportapi

import (
	"fmt"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/services/event"
	"github.com/transcom/mymove/pkg/services/query"

	paymentrequestop "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/payment_requests"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/supportapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

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
		return paymentrequestop.NewUpdatePaymentRequestStatusInternalServerError().WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
	}

	// Let's fetch the existing payment request using the PaymentRequestFetcher service object
	filter := []services.QueryFilter{query.NewQueryFilter("id", "=", paymentRequestID.String())}
	existingPaymentRequest, err := h.PaymentRequestFetcher.FetchPaymentRequest(filter)

	if err != nil {
		msg := fmt.Sprintf("Error finding Payment Request for status update with ID: %s", params.PaymentRequestID.String())
		logger.Error(msg, zap.Error(err))
		return paymentrequestop.NewUpdatePaymentRequestStatusNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, msg, h.GetTraceID()))
	}

	status := existingPaymentRequest.Status
	mtoID := existingPaymentRequest.MoveTaskOrderID

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
			return paymentrequestop.NewUpdatePaymentRequestStatusNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID()))
		case services.PreconditionFailedError:
			return paymentrequestop.NewUpdatePaymentRequestStatusPreconditionFailed().WithPayload(payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceID()))
		default:
			logger.Error(fmt.Sprintf("Error saving payment request status for ID: %s: %s", paymentRequestID, err))
			return paymentrequestop.NewUpdatePaymentRequestStatusInternalServerError().WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
		}
	}

	_, err = event.TriggerEvent(event.Event{
		EventKey:        event.PaymentRequestUpdateEventKey,
		MtoID:           mtoID,
		UpdatedObjectID: updatedPaymentRequest.ID,
		Request:         params.HTTPRequest,
		EndpointKey:     event.SupportUpdatePaymentRequestStatusEndpointKey,
		DBConnection:    h.DB(),
		HandlerContext:  h,
	})
	if err != nil {
		logger.Error("supportapi.UpdatePaymentRequestStatusHandler could not generate the event")
	}

	returnPayload := payloads.PaymentRequest(updatedPaymentRequest)
	return paymentrequestop.NewUpdatePaymentRequestStatusOK().WithPayload(returnPayload)
}

// ListMTOPaymentRequestsHandler gets all payment requests for a given MTO
type ListMTOPaymentRequestsHandler struct {
	handlers.HandlerContext
}

// Handle getting payment requests for a given MTO
func (h ListMTOPaymentRequestsHandler) Handle(params paymentrequestop.ListMTOPaymentRequestsParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	mtoID, err := uuid.FromString(params.MoveTaskOrderID.String())

	if err != nil {
		logger.Error(fmt.Sprintf("Error parsing move task order id: %s", params.MoveTaskOrderID.String()), zap.Error(err))
		return paymentrequestop.NewListMTOPaymentRequestsInternalServerError()
	}

	var paymentRequests models.PaymentRequests

	query := h.DB().Where("move_task_order_id = ?", mtoID)

	err = query.All(&paymentRequests)

	if err != nil {
		logger.Error("Unable to fetch records:", zap.Error(err))
		return paymentrequestop.NewListMTOPaymentRequestsInternalServerError()
	}

	payload := payloads.PaymentRequests(&paymentRequests)

	return paymentrequestop.NewListMTOPaymentRequestsOK().WithPayload(*payload)
}
