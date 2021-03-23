package ghcapi

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/services/event"

	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"

	"github.com/go-openapi/runtime/middleware"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/audit"
)

// GetMoveTaskOrderHandler fetches a Move Task Order
type GetMoveTaskOrderHandler struct {
	handlers.HandlerContext
	moveTaskOrderFetcher services.MoveTaskOrderFetcher
}

// Handle fetches a single MoveTaskOrder
func (h GetMoveTaskOrderHandler) Handle(params movetaskorderops.GetMoveTaskOrderParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)

	searchParams := services.FetchMoveTaskOrderParams{
		IncludeHidden: false,
	}
	mto, err := h.moveTaskOrderFetcher.FetchMoveTaskOrder(moveTaskOrderID, &searchParams)
	if err != nil {
		logger.Error("ghcapi.GetMoveTaskOrderHandler error", zap.Error(err))
		switch err.(type) {
		case services.NotFoundError:
			return movetaskorderops.NewGetMoveTaskOrderNotFound()
		case services.InvalidInputError:
			return movetaskorderops.NewGetMoveTaskOrderBadRequest()
		default:
			return movetaskorderops.NewGetMoveTaskOrderInternalServerError()
		}
	}
	moveTaskOrderPayload := payloads.MoveTaskOrder(mto)
	return movetaskorderops.NewGetMoveTaskOrderOK().WithPayload(moveTaskOrderPayload)
}

// UpdateMoveTaskOrderStatusHandlerFunc updates the status of a Move Task Order
type UpdateMoveTaskOrderStatusHandlerFunc struct {
	handlers.HandlerContext
	moveTaskOrderStatusUpdater services.MoveTaskOrderUpdater
}

// Handle updates the status of a MoveTaskOrder
func (h UpdateMoveTaskOrderStatusHandlerFunc) Handle(params movetaskorderops.UpdateMoveTaskOrderStatusParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	eTag := params.IfMatch

	// TODO how are we going to handle auth in new api? Do we need some sort of placeholder to remind us to
	// TODO to revisit?
	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)

	serviceItemCodes := ghcmessages.MTOApprovalServiceItemCodes{}
	if params.ServiceItemCodes != nil {
		serviceItemCodes = *params.ServiceItemCodes
	}

	mto, err := h.moveTaskOrderStatusUpdater.MakeAvailableToPrime(moveTaskOrderID, eTag,
		serviceItemCodes.ServiceCodeMS, serviceItemCodes.ServiceCodeCS)

	if err != nil {
		logger.Error("ghcapi.UpdateMoveTaskOrderStatusHandlerFunc error", zap.Error(err))
		switch err.(type) {
		case services.NotFoundError:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusNotFound()
		case services.InvalidInputError:
			payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceID(), validate.NewErrors())
			return movetaskorderops.NewUpdateMoveTaskOrderStatusUnprocessableEntity().WithPayload(payload)
		case services.PreconditionFailedError:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case services.ConflictError:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusInternalServerError()
		}
	}

	moveTaskOrderPayload := payloads.Move(mto)

	// Audit attempt to make MTO available to prime
	_, err = audit.Capture(mto, moveTaskOrderPayload, logger, session, params.HTTPRequest)
	if err != nil {
		logger.Error("Auditing service error for making MTO available to Prime.", zap.Error(err))
		return movetaskorderops.NewUpdateMoveTaskOrderStatusInternalServerError()
	}

	_, err = event.TriggerEvent(event.Event{
		EventKey:        event.MoveTaskOrderUpdateEventKey,
		MtoID:           mto.ID,
		UpdatedObjectID: mto.ID,
		Request:         params.HTTPRequest,
		EndpointKey:     event.GhcUpdateMoveTaskOrderStatusEndpointKey,
		DBConnection:    h.DB(),
		HandlerContext:  h,
	})
	if err != nil {
		logger.Error("ghcapi.UpdateMoveTaskOrderStatusHandlerFunc could not generate the event")
	}

	return movetaskorderops.NewUpdateMoveTaskOrderStatusOK().WithPayload(moveTaskOrderPayload)
}
