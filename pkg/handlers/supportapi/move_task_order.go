package supportapi

import (
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/services/support"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/handlers/supportapi/internal/payloads"

	"github.com/go-openapi/runtime/middleware"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

// UpdateMoveTaskOrderStatusHandlerFunc updates the status of a Move Task Order
type UpdateMoveTaskOrderStatusHandlerFunc struct {
	handlers.HandlerContext
	moveTaskOrderStatusUpdater services.MoveTaskOrderUpdater
}

// Handle updates the status of a MoveTaskOrder
func (h UpdateMoveTaskOrderStatusHandlerFunc) Handle(params movetaskorderops.UpdateMoveTaskOrderStatusParams) middleware.Responder {
	_, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	eTag := params.IfMatch

	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)

	mto, err := h.moveTaskOrderStatusUpdater.MakeAvailableToPrime(moveTaskOrderID, eTag)

	if err != nil {
		logger.Error("supportapi.MoveTaskOrderHandler error", zap.Error(err))
		switch typedErr := err.(type) {
		case services.NotFoundError:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusNotFound().WithPayload(
				payloads.ClientError(handlers.NotFoundMessage, *handlers.FmtString(err.Error()), h.GetTraceID()))
		case services.InvalidInputError:
			return movetaskorderops.NewCreateMoveTaskOrderUnprocessableEntity().WithPayload(
				payloads.ValidationError(err.Error(), h.GetTraceID(), typedErr.ValidationErrors))
		case services.PreconditionFailedError:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusPreconditionFailed().WithPayload(
				payloads.ClientError(handlers.PreconditionErrMessage, *handlers.FmtString(err.Error()), h.GetTraceID()))
		default:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusInternalServerError().WithPayload(&supportmessages.Error{Message: handlers.FmtString(err.Error())})
		}
	}

	moveTaskOrderPayload := payloads.MoveTaskOrder(mto)

	return movetaskorderops.NewUpdateMoveTaskOrderStatusOK().WithPayload(moveTaskOrderPayload)
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
		switch typedErr := err.(type) {
		case services.NotFoundError:
			return movetaskorderops.NewGetMoveTaskOrderNotFound().WithPayload(
				payloads.ClientError(handlers.NotFoundMessage, *handlers.FmtString(err.Error()), h.GetTraceID()))
		case services.InvalidInputError:
			return movetaskorderops.NewCreateMoveTaskOrderUnprocessableEntity().WithPayload(
				payloads.ValidationError(err.Error(), h.GetTraceID(), typedErr.ValidationErrors))
		default:
			return movetaskorderops.NewGetMoveTaskOrderInternalServerError().WithPayload(&supportmessages.Error{Message: handlers.FmtString(err.Error())})
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
			return movetaskorderops.NewCreateMoveTaskOrderInternalServerError().WithPayload(&supportmessages.Error{Message: handlers.FmtString(err.Error())})
		}
	}
	moveTaskOrderPayload := payloads.MoveTaskOrder(moveTaskOrder)
	return movetaskorderops.NewCreateMoveTaskOrderCreated().WithPayload(moveTaskOrderPayload)

}
