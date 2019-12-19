package ghcapi

import (
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"

	movetaskorderservice "github.com/transcom/mymove/pkg/services/move_task_order"

	"github.com/go-openapi/runtime/middleware"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

// GetMoveTaskOrderHandler updates the status of a Move Task Order
type GetMoveTaskOrderHandler struct {
	handlers.HandlerContext
	moveTaskOrderFetcher services.MoveTaskOrderFetcher
}

// GetMoveTaskOrderHandler updates the status of a MoveTaskOrder
func (h GetMoveTaskOrderHandler) Handle(params movetaskorderops.GetMoveTaskOrderParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	mto, err := h.moveTaskOrderFetcher.FetchMoveTaskOrder(moveTaskOrderID)
	if err != nil {
		logger.Error("ghciap.GetMoveTaskOrderHandler error", zap.Error(err))
		switch err.(type) {
		case movetaskorderservice.ErrNotFound:
			return movetaskorderops.NewGetMoveTaskOrderNotFound()
		case movetaskorderservice.ErrInvalidInput:
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
	moveTaskOrderStatusUpdater services.MoveTaskOrderStatusUpdater
}

// UpdateMoveTaskOrderStatusHandlerFunc updates the status of a MoveTaskOrder
func (h UpdateMoveTaskOrderStatusHandlerFunc) Handle(params movetaskorderops.UpdateMoveTaskOrderStatusParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	// TODO how are we going to handle auth in new api? Do we need some sort of placeholder to remind us to
	// TODO to revisit?
	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	mto, err := h.moveTaskOrderStatusUpdater.MakeAvailableToPrime(moveTaskOrderID)
	if err != nil {
		logger.Error("ghciap.MoveTaskOrderHandler error", zap.Error(err))
		switch err.(type) {
		case movetaskorderservice.ErrNotFound:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusNotFound()
		case movetaskorderservice.ErrInvalidInput:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusBadRequest()
		default:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusInternalServerError()
		}
	}
	moveTaskOrderPayload := payloads.MoveTaskOrder(mto)
	return movetaskorderops.NewUpdateMoveTaskOrderStatusOK().WithPayload(moveTaskOrderPayload)
}
