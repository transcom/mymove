package ghcapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	paymentServiceItemOp "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_service_item"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	modelToPayload "github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/audit"
	"github.com/transcom/mymove/pkg/services/event"
	"github.com/transcom/mymove/pkg/services/query"
)

// UpdatePaymentServiceItemStatusHandler updates payment service item status
type UpdatePaymentServiceItemStatusHandler struct {
	handlers.HandlerContext
	services.PaymentServiceItemStatusUpdater
}

// Handle handles the handling for UpdatePaymentServiceItemStatusHandler
func (h UpdatePaymentServiceItemStatusHandler) Handle(
	params paymentServiceItemOp.UpdatePaymentServiceItemStatusParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			paymentServiceItemID, err := uuid.FromString(params.PaymentServiceItemID)
			newStatus := models.PaymentServiceItemStatus(params.Body.Status)

			if err != nil {
				appCtx.Logger().
					Error(fmt.Sprintf("Error parsing payment service item id: %s", params.PaymentServiceItemID), zap.Error(err))
			}

			updatedPaymentServiceItem, verrs, err := h.PaymentServiceItemStatusUpdater.UpdatePaymentServiceItemStatus(
				appCtx,
				paymentServiceItemID,
				newStatus,
				params.Body.RejectionReason,
				params.IfMatch,
			)
			if err != nil {
				appCtx.Logger().Error("Error updating payment service item status", zap.Error(err))

				switch e := err.(type) {
				case query.StaleIdentifierError:
					return paymentServiceItemOp.NewUpdatePaymentServiceItemStatusPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.NotFoundError:
					return paymentServiceItemOp.NewUpdatePaymentServiceItemStatusNotFound().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Validation errors", "UpdatePaymentServiceItemStatus", h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
					return paymentServiceItemOp.NewUpdatePaymentServiceItemStatusUnprocessableEntity().WithPayload(payload), err
				default:
					return paymentServiceItemOp.NewUpdatePaymentServiceItemStatusInternalServerError().WithPayload(&ghcmessages.Error{Message: handlers.FmtString("Error updating payment service item status")}), err
				}
			}
			if verrs != nil {
				payload := payloadForValidationError(
					"Validation errors",
					"UpdatePaymentServiceItemStatus",
					h.GetTraceIDFromRequest(params.HTTPRequest),
					verrs,
				)
				return paymentServiceItemOp.NewUpdatePaymentServiceItemStatusUnprocessableEntity().
					WithPayload(payload), verrs
			}

			// Capture update attempt in audit log
			_, err = audit.Capture(appCtx, &updatedPaymentServiceItem, nil, params.HTTPRequest)
			if err != nil {
				appCtx.Logger().
					Error("Auditing service error for payment service item status change.", zap.Error(err))
				return paymentServiceItemOp.NewUpdatePaymentServiceItemStatusInternalServerError(), err
			}

			_, err = event.TriggerEvent(event.Event{
				EventKey:        event.PaymentRequestUpdateEventKey,
				MtoID:           updatedPaymentServiceItem.PaymentRequest.MoveTaskOrderID,
				UpdatedObjectID: updatedPaymentServiceItem.PaymentRequestID,
				EndpointKey:     event.GhcUpdatePaymentServiceItemStatusEndpointKey,
				AppContext:      appCtx,
				TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
			})
			if err != nil {
				appCtx.Logger().
					Error("ghcapi.UpdatePaymentServiceItemStatusHandler could not generate the event")
			}

			// Make the payload and return it with a 200
			payload := modelToPayload.PaymentServiceItem(&updatedPaymentServiceItem)
			return paymentServiceItemOp.NewUpdatePaymentServiceItemStatusOK().WithPayload(payload), nil
		})
}
