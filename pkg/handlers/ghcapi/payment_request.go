package ghcapi

import (
	"fmt"
	"reflect"
	"time"

	"github.com/transcom/mymove/pkg/models/roles"

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

// GetPaymentRequestForMoveHandler gets payment requests associated with a move
type GetPaymentRequestForMoveHandler struct {
	handlers.HandlerContext
	services.PaymentRequestListFetcher
}

// Handle handles the HTTP handling for GetPaymentRequestForMoveHandler
func (h GetPaymentRequestForMoveHandler) Handle(params paymentrequestop.GetPaymentRequestsForMoveParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	if !session.IsOfficeUser() || !session.Roles.HasRole(roles.RoleTypeTIO) {
		logger.Error("user is not authenticated with TIO office role")
		return paymentrequestop.NewGetPaymentRequestsForMoveForbidden()
	}

	locator := params.Locator

	paymentRequests, err := h.FetchPaymentRequestListByMove(session.OfficeUserID, locator)
	if err != nil {
		logger.Error(fmt.Sprintf("Error fetching Payment Request for locator: %s", locator), zap.Error(err))
		return paymentrequestop.NewGetPaymentRequestNotFound()
	}

	returnPayload, err := payloads.PaymentRequests(paymentRequests, h.FileStorer())

	if err != nil {
		logger.Error(fmt.Sprintf("Error building payment requests payload for locator: %s", locator), zap.Error(err))
		return paymentrequestop.NewGetPaymentRequestsForMoveInternalServerError()
	}

	return paymentrequestop.NewGetPaymentRequestsForMoveOK().WithPayload(*returnPayload)
}

// GetPaymentRequestHandler gets payment requests
type GetPaymentRequestHandler struct {
	handlers.HandlerContext
	services.PaymentRequestFetcher
}

// Handle gets payment requests
func (h GetPaymentRequestHandler) Handle(params paymentrequestop.GetPaymentRequestParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	if !session.IsOfficeUser() || !session.Roles.HasRole(roles.RoleTypeTIO) {
		logger.Error("user is not authenticated with TIO office role")
		return paymentrequestop.NewGetPaymentRequestForbidden()
	}

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

	if !session.IsOfficeUser() || !session.Roles.HasRole(roles.RoleTypeTIO) {
		logger.Error("user is not authenticated with TIO office role")
		return paymentrequestop.NewUpdatePaymentRequestStatusForbidden()
	}

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

	now := time.Now()
	existingPaymentRequest.Status = models.PaymentRequestStatus(params.Body.Status)

	// Let's map the incoming status to our enumeration type
	switch existingPaymentRequest.Status {
	case models.PaymentRequestStatusReviewed, models.PaymentRequestStatusReviewedAllRejected:
		existingPaymentRequest.ReviewedAt = &now
	case models.PaymentRequestStatusSentToGex:
		existingPaymentRequest.SentToGexAt = &now
	case models.PaymentRequestStatusReceivedByGex:
		existingPaymentRequest.ReceivedByGexAt = &now
	case models.PaymentRequestStatusPaid:
		existingPaymentRequest.PaidAt = &now
	}

	// If we got a rejection reason let's use it
	if params.Body.RejectionReason != nil {
		existingPaymentRequest.RejectionReason = params.Body.RejectionReason
	}

	// Capture update attempt in audit log
	_, err = audit.Capture(&existingPaymentRequest, nil, logger, session, params.HTTPRequest)
	if err != nil {
		logger.Error("Auditing service error for payment request update.", zap.Error(err))
		return paymentrequestop.NewUpdatePaymentRequestStatusInternalServerError()
	}

	// And now let's save our updated model object using the PaymentRequestUpdater service object.
	updatedPaymentRequest, err := h.PaymentRequestStatusUpdater.UpdatePaymentRequestStatus(&existingPaymentRequest, params.IfMatch)

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
