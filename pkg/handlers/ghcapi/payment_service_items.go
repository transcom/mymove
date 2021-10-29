package ghcapi

import (
	"fmt"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/services/event"

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
	services.PaymentServiceItemStatusUpdater
}

// Handle handles the handling for UpdatePaymentServiceItemStatusHandler
func (h UpdatePaymentServiceItemStatusHandler) Handle(params paymentServiceItemOp.UpdatePaymentServiceItemStatusParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)
	paymentServiceItemID, err := uuid.FromString(params.PaymentServiceItemID)
	newStatus := models.PaymentServiceItemStatus(params.Body.Status)

	if err != nil {
		appCtx.Logger().Error(fmt.Sprintf("Error parsing payment service item id: %s", params.PaymentServiceItemID), zap.Error(err))
	}

	updatedPaymentServiceItem, verrs, err := h.PaymentServiceItemStatusUpdater.UpdatePaymentServiceItemStatus(appCtx,
		paymentServiceItemID, newStatus, params.Body.RejectionReason, params.IfMatch)

	if err != nil {
		appCtx.Logger().Error("Error updating payment service item status", zap.Error(err))

		switch e := err.(type) {
		case query.StaleIdentifierError:
			return paymentServiceItemOp.NewUpdatePaymentServiceItemStatusPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case apperror.NotFoundError:
			return paymentServiceItemOp.NewUpdatePaymentServiceItemStatusNotFound().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case apperror.InvalidInputError:
			payload := payloadForValidationError("Validation errors", "UpdatePaymentServiceItemStatus", h.GetTraceID(), e.ValidationErrors)
			return paymentServiceItemOp.NewUpdatePaymentServiceItemStatusUnprocessableEntity().WithPayload(payload)
		default:
			return paymentServiceItemOp.NewUpdatePaymentServiceItemStatusInternalServerError().WithPayload(&ghcmessages.Error{Message: handlers.FmtString("Error updating payment service item status")})
		}
	}
	if verrs != nil {
		payload := payloadForValidationError("Validation errors", "UpdatePaymentServiceItemStatus", h.GetTraceID(), verrs)
		return paymentServiceItemOp.NewUpdatePaymentServiceItemStatusUnprocessableEntity().WithPayload(payload)
	}

	// Capture update attempt in audit log
	_, err = audit.Capture(&updatedPaymentServiceItem, nil, appCtx.Logger(), appCtx.Session(), params.HTTPRequest)
	if err != nil {
		appCtx.Logger().Error("Auditing service error for payment service item status change.", zap.Error(err))
		return paymentServiceItemOp.NewUpdatePaymentServiceItemStatusInternalServerError()
	}

	_, err = event.TriggerEvent(event.Event{
		EventKey:        event.PaymentRequestUpdateEventKey,
		MtoID:           updatedPaymentServiceItem.PaymentRequest.MoveTaskOrderID,
		UpdatedObjectID: updatedPaymentServiceItem.PaymentRequestID,
		Request:         params.HTTPRequest,
		EndpointKey:     event.GhcUpdatePaymentServiceItemStatusEndpointKey,
		HandlerContext:  h,
	})
	if err != nil {
		appCtx.Logger().Error("ghcapi.UpdatePaymentServiceItemStatusHandler could not generate the event")
	}

	// Make the payload and return it with a 200
	payload := modelToPayload.PaymentServiceItem(&updatedPaymentServiceItem)
	return paymentServiceItemOp.NewUpdatePaymentServiceItemStatusOK().WithPayload(payload)
}
