package ghcapi

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/services/event"

	"github.com/go-openapi/swag"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	modelToPayload "github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/audit"
	"github.com/transcom/mymove/pkg/services/query"

	paymentServiceItemOp "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_service_item"
)

// UpdatePaymentServiceItemStatusHandler updates payment service item status
type UpdatePaymentServiceItemStatusHandler struct {
	handlers.HandlerContext
	services.Fetcher
	query.Builder
}

// Handle handles the handling for UpdatePaymentServiceItemStatusHandler
func (h UpdatePaymentServiceItemStatusHandler) Handle(params paymentServiceItemOp.UpdatePaymentServiceItemStatusParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	paymentServiceItemID, err := uuid.FromString(params.PaymentServiceItemID)
	// Create a zero paymentServiceRequest for us to use in FetchRecord
	var paymentServiceItem models.PaymentServiceItem

	if err != nil {
		logger.Error(fmt.Sprintf("Error parsing payment service item id: %s", params.PaymentServiceItemID), zap.Error(err))
	}

	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", paymentServiceItemID.String())}
	// Get the existing record
	err = h.Fetcher.FetchRecord(&paymentServiceItem, filters)

	if err != nil {
		logger.Error(fmt.Sprintf("Error finding payment service item for status update with ID: %s", params.PaymentServiceItemID), zap.Error(err))
		payload := payloadForClientError("Unknown UUID(s)", "Unknown UUID(s) used to update a payment service item ", h.GetTraceID())
		return paymentServiceItemOp.NewUpdatePaymentServiceItemStatusNotFound().WithPayload(payload)
	}
	// Create a model object to use for the update and set the status
	newStatus := models.PaymentServiceItemStatus(params.Body.Status)
	paymentServiceItem.Status = newStatus

	if params.Body.RejectionReason != nil {
		paymentServiceItem.RejectionReason = params.Body.RejectionReason
	}
	// If we're approving this thing then we don't want there to be a rejection reason
	// We also will want to update the ApprovedAt field and nil out the DeniedAt field.
	if paymentServiceItem.Status == models.PaymentServiceItemStatusApproved {
		paymentServiceItem.RejectionReason = nil
		paymentServiceItem.ApprovedAt = swag.Time(time.Now())
		paymentServiceItem.DeniedAt = nil
	}
	// If we're denying this thing we want to make sure to update the DeniedAt field and nil out ApprovedAt.
	if paymentServiceItem.Status == models.PaymentServiceItemStatusDenied {
		paymentServiceItem.DeniedAt = swag.Time(time.Now())
		paymentServiceItem.ApprovedAt = nil
	}

	// Capture update attempt in audit log
	_, err = audit.Capture(&paymentServiceItem, nil, logger, session, params.HTTPRequest)
	if err != nil {
		logger.Error("Auditing service error for payment service item status change.", zap.Error(err))
		return paymentServiceItemOp.NewUpdatePaymentServiceItemStatusInternalServerError()
	}
	// Do the update
	verrs, err := h.UpdateOne(&paymentServiceItem, &params.IfMatch)
	// Using a switch to match error causes to appropriate return type in gen code
	if err != nil {
		switch err.(type) {
		case query.StaleIdentifierError:
			return paymentServiceItemOp.NewUpdatePaymentServiceItemStatusPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case services.NotFoundError:
			payload := payloadForClientError("Unknown UUID(s)", "Unknown UUID(s) used to update a payment service item ", h.GetTraceID())
			return paymentServiceItemOp.NewUpdatePaymentServiceItemStatusNotFound().WithPayload(payload)
		}
	}
	if verrs != nil {
		return paymentServiceItemOp.NewUpdatePaymentServiceItemStatusInternalServerError().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(verrs.String())})
	}

	_, err = event.TriggerEvent(event.Event{
		EventKey:        event.PaymentRequestUpdateEventKey,
		MtoID:           paymentServiceItem.PaymentRequest.MoveTaskOrderID,
		UpdatedObjectID: paymentServiceItem.PaymentRequestID,
		Request:         params.HTTPRequest,
		EndpointKey:     event.GhcUpdatePaymentServiceItemStatusEndpointKey,
		DBConnection:    h.DB(),
		HandlerContext:  h,
	})
	if err != nil {
		logger.Error("ghcapi.UpdatePaymentServiceItemStatusHandler could not generate the event")
	}

	// Make the payload and return it with a 200
	payload := modelToPayload.PaymentServiceItem(&paymentServiceItem)
	return paymentServiceItemOp.NewUpdatePaymentServiceItemStatusOK().WithPayload(payload)
}
