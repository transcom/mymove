package supportapi

import (
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/services/support"

	"github.com/transcom/mymove/pkg/handlers/supportapi/internal/payloads"

	"github.com/go-openapi/runtime/middleware"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

// ListMTOsHandler lists move task orders with the option to filter since a particular date
type ListMTOsHandler struct {
	handlers.HandlerContext
	services.MoveTaskOrderFetcher
}

// Handle fetches all move task orders with the option to filter since a particular date
func (h ListMTOsHandler) Handle(params movetaskorderops.ListMTOsParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	mtos, err := h.MoveTaskOrderFetcher.ListAllMoveTaskOrders(false, params.Since)

	if err != nil {
		logger.Error("Unable to fetch records:", zap.Error(err))
		return movetaskorderops.NewListMTOsInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
	}

	payload := payloads.MoveTaskOrders(&mtos)

	return movetaskorderops.NewListMTOsOK().WithPayload(payload)
}

// MakeMoveTaskOrderAvailableHandlerFunc updates the status of a Move Task Order
type MakeMoveTaskOrderAvailableHandlerFunc struct {
	handlers.HandlerContext
	moveTaskOrderAvailabilityUpdater services.MoveTaskOrderUpdater
}

// Handle updates the prime availability of a MoveTaskOrder
func (h MakeMoveTaskOrderAvailableHandlerFunc) Handle(params movetaskorderops.MakeMoveTaskOrderAvailableParams) middleware.Responder {
	_, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	eTag := params.IfMatch

	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)

	mto, err := h.moveTaskOrderAvailabilityUpdater.MakeAvailableToPrime(moveTaskOrderID, eTag, false, false)

	if err != nil {
		logger.Error("supportapi.MakeMoveTaskOrderAvailableHandlerFunc error", zap.Error(err))
		switch typedErr := err.(type) {
		case services.NotFoundError:
			return movetaskorderops.NewMakeMoveTaskOrderAvailableNotFound().WithPayload(
				payloads.ClientError(handlers.NotFoundMessage, *handlers.FmtString(err.Error()), h.GetTraceID()))
		case services.InvalidInputError:
			return movetaskorderops.NewMakeMoveTaskOrderAvailableUnprocessableEntity().WithPayload(
				payloads.ValidationError(err.Error(), h.GetTraceID(), typedErr.ValidationErrors))
		case services.PreconditionFailedError:
			return movetaskorderops.NewMakeMoveTaskOrderAvailablePreconditionFailed().WithPayload(
				payloads.ClientError(handlers.PreconditionErrMessage, *handlers.FmtString(err.Error()), h.GetTraceID()))
		default:
			return movetaskorderops.NewMakeMoveTaskOrderAvailableInternalServerError().WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
		}
	}

	moveTaskOrderPayload := payloads.MoveTaskOrder(mto)

	return movetaskorderops.NewMakeMoveTaskOrderAvailableOK().WithPayload(moveTaskOrderPayload)
}

// HideNonFakeMoveTaskOrdersHandlerFunc calls service to hide MTOs that are not using fake data
type HideNonFakeMoveTaskOrdersHandlerFunc struct {
	handlers.HandlerContext
	services.MoveTaskOrderHider
}

// Handle hides any mto that doesnt have valid fake data
func (h HideNonFakeMoveTaskOrdersHandlerFunc) Handle(params movetaskorderops.HideNonFakeMoveTaskOrdersParams) middleware.Responder {
	_, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	hiddenMTOIds, err := h.Hide()
	if err != nil {
		logger.Error("supportapi.HideNonFakeMoveTaskOrdersHandlerFunc error", zap.Error(err))
		return movetaskorderops.NewHideNonFakeMoveTaskOrdersInternalServerError().WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
	}

	payload := payloads.MoveTaskOrderIDs(&hiddenMTOIds)

	return movetaskorderops.NewHideNonFakeMoveTaskOrdersOK().WithPayload(payload)
}

// GetMoveTaskOrderHandlerFunc updates the status of a Move Task Order
type GetMoveTaskOrderHandlerFunc struct {
	handlers.HandlerContext
	moveTaskOrderFetcher services.MoveTaskOrderFetcher
}

// Handle updates the status of a MoveTaskOrder
func (h GetMoveTaskOrderHandlerFunc) Handle(params movetaskorderops.GetMoveTaskOrderParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	mto, err := h.moveTaskOrderFetcher.FetchMoveTaskOrder(moveTaskOrderID)
	if err != nil {
		logger.Error("primeapi.support.GetMoveTaskOrderHandler error", zap.Error(err))
		switch err.(type) {
		case services.NotFoundError:
			return movetaskorderops.NewGetMoveTaskOrderNotFound().WithPayload(
				payloads.ClientError(handlers.NotFoundMessage, *handlers.FmtString(err.Error()), h.GetTraceID()))
		default:
			return movetaskorderops.NewGetMoveTaskOrderInternalServerError().WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
		}
	}
	moveTaskOrderPayload := payloads.MoveTaskOrder(mto)
	return movetaskorderops.NewGetMoveTaskOrderOK().WithPayload(moveTaskOrderPayload)
}

// CreateMoveTaskOrderHandler creates a move task order
type CreateMoveTaskOrderHandler struct {
	handlers.HandlerContext
	moveTaskOrderCreator support.InternalMoveTaskOrderCreator
}

// Handle updates to move task order post-counseling
func (h CreateMoveTaskOrderHandler) Handle(params movetaskorderops.CreateMoveTaskOrderParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	moveTaskOrder, err := h.moveTaskOrderCreator.InternalCreateMoveTaskOrder(*params.Body, logger)

	if err != nil {
		logger.Error("primeapi.support.CreateMoveTaskOrderHandler error", zap.Error(err))
		switch typedErr := err.(type) {
		case services.NotFoundError:
			return movetaskorderops.NewCreateMoveTaskOrderNotFound().WithPayload(
				payloads.ClientError(handlers.NotFoundMessage, *handlers.FmtString(err.Error()), h.GetTraceID()))
		case services.InvalidInputError:
			errPayload := payloads.ValidationError(err.Error(), h.GetTraceID(), typedErr.ValidationErrors)
			return movetaskorderops.NewCreateMoveTaskOrderUnprocessableEntity().WithPayload(errPayload)
		case services.QueryError:
			// This error is generated when the validation passed but there was an error in creation
			// Usually this is due to a more complex dependency like a foreign key constraint
			return movetaskorderops.NewCreateMoveTaskOrderBadRequest().WithPayload(
				payloads.ClientError(handlers.SQLErrMessage, *handlers.FmtString(err.Error()), h.GetTraceID()))
		default:
			return movetaskorderops.NewCreateMoveTaskOrderInternalServerError().WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceID()))
		}
	}

	moveTaskOrderPayload := payloads.MoveTaskOrder(moveTaskOrder)
	return movetaskorderops.NewCreateMoveTaskOrderCreated().WithPayload(moveTaskOrderPayload)

}
