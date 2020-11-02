package ghcapi

import (
	"fmt"
	"reflect"
	"time"

	"github.com/transcom/mymove/pkg/services/event"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services/audit"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	paymentrequestop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_requests"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// GetPaymentRequestHandler gets payment requests
type GetPaymentRequestHandler struct {
	handlers.HandlerContext
	services.PaymentRequestFetcher
}

// Handle gets payment requests
func (h GetPaymentRequestHandler) Handle(params paymentrequestop.GetPaymentRequestParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	paymentRequestID, err := uuid.FromString(params.PaymentRequestID.String())

	if err != nil {
		logger.Error(fmt.Sprintf("Error parsing payment request id: %s", params.PaymentRequestID.String()), zap.Error(err))
		return paymentrequestop.NewGetPaymentRequestInternalServerError()
	}

	paymentRequest, err := h.FetchPaymentRequest(paymentRequestID)

	if err != nil {
		logger.Error(fmt.Sprintf("Error fetching Payment Request with ID: %s", params.PaymentRequestID.String()), zap.Error(err))
		return paymentrequestop.NewGetPaymentRequestNotFound()
	}

	if reflect.DeepEqual(paymentRequest, models.PaymentRequest{}) {
		logger.Info(fmt.Sprintf("Could not find a Payment Request with ID: %s", params.PaymentRequestID.String()))
		return paymentrequestop.NewGetPaymentRequestNotFound()
	}

	returnPayload, err := payloads.PaymentRequest(&paymentRequest, h.FileStorer())
	if err != nil {
		return paymentrequestop.NewGetPaymentRequestInternalServerError()
	}

	response := paymentrequestop.NewGetPaymentRequestOK().WithPayload(returnPayload)

	return response
}

// UpdatePaymentRequestStatusHandler updates payment requests status
type UpdatePaymentRequestStatusHandler struct {
	handlers.HandlerContext
	services.PaymentRequestStatusUpdater
	services.PaymentRequestFetcher
}

// Handle updates payment requests status
func (h UpdatePaymentRequestStatusHandler) Handle(params paymentrequestop.UpdatePaymentRequestStatusParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	paymentRequestID, err := uuid.FromString(params.PaymentRequestID.String())

	if err != nil {
		logger.Error(fmt.Sprintf("Error parsing payment request id: %s", params.PaymentRequestID.String()), zap.Error(err))
		return paymentrequestop.NewGetPaymentRequestInternalServerError()
	}

	// Let's fetch the existing payment request using the PaymentRequestFetcher service object
	existingPaymentRequest, err := h.PaymentRequestFetcher.FetchPaymentRequest(paymentRequestID)

	if err != nil {
		logger.Error(fmt.Sprintf("Error finding Payment Request for status update with ID: %s", params.PaymentRequestID.String()), zap.Error(err))
		return paymentrequestop.NewGetPaymentRequestNotFound()
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

	// Capture update attempt in audit log
	_, err = audit.Capture(&paymentRequestForUpdate, nil, logger, session, params.HTTPRequest)
	if err != nil {
		logger.Error("Auditing service error for payment request update.", zap.Error(err))
		return paymentrequestop.NewUpdatePaymentRequestStatusInternalServerError()
	}

	// And now let's save our updated model object using the PaymentRequestUpdater service object.
	updatedPaymentRequest, err := h.PaymentRequestStatusUpdater.UpdatePaymentRequestStatus(&paymentRequestForUpdate, params.IfMatch)

	if err != nil {
		switch err.(type) {
		case services.NotFoundError:
			payload := payloadForClientError("Unknown UUID(s)", "Unknown UUID(s) used to update a payment request ", h.GetTraceID())
			return paymentrequestop.NewUpdatePaymentRequestStatusNotFound().WithPayload(payload)
		case services.PreconditionFailedError:
			return paymentrequestop.NewUpdatePaymentRequestStatusPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case services.InvalidInputError:
			payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceID(), validate.NewErrors())
			return paymentrequestop.NewUpdatePaymentRequestStatusUnprocessableEntity().WithPayload(payload)
		default:
			logger.Error(fmt.Sprintf("Error saving payment request status for ID: %s: %s", paymentRequestID, err))
			return paymentrequestop.NewUpdatePaymentRequestStatusInternalServerError()
		}
	}

	_, err = event.TriggerEvent(event.Event{
		EventKey:        event.PaymentRequestUpdateEventKey,
		MtoID:           updatedPaymentRequest.MoveTaskOrderID,
		UpdatedObjectID: updatedPaymentRequest.ID,
		Request:         params.HTTPRequest,
		EndpointKey:     event.GhcUpdatePaymentRequestStatusEndpointKey,
		DBConnection:    h.DB(),
		HandlerContext:  h,
	})
	if err != nil {
		logger.Error("ghcapi.UpdatePaymentRequestStatusHandler could not generate the event")
	}

	returnPayload, err := payloads.PaymentRequest(updatedPaymentRequest, h.FileStorer())
	if err != nil {
		return paymentrequestop.NewGetPaymentRequestInternalServerError()
	}

	return paymentrequestop.NewUpdatePaymentRequestStatusOK().WithPayload(returnPayload)
}
