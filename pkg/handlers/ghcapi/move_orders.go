package ghcapi

import (
	"database/sql"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/services"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	moveorderop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_order"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
)

// GetMoveOrdersHandler fetches the information of a specific customer
type GetMoveOrdersHandler struct {
	handlers.HandlerContext
	services.MoveOrderFetcher
}

// Handle getting the information of a specific customer
func (h GetMoveOrdersHandler) Handle(params moveorderop.GetMoveOrderParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	moveOrderID, _ := uuid.FromString(params.MoveOrderID.String())
	moveOrder, err := h.FetchMoveOrder(moveOrderID)
	if err != nil {
		logger.Error("fetching move order", zap.Error(err))
		switch err {
		case sql.ErrNoRows:
			return moveorderop.NewGetMoveOrderNotFound()
		default:
			return moveorderop.NewGetMoveOrderInternalServerError()
		}
	}
	moveOrderPayload := payloads.MoveOrder(moveOrder)
	return moveorderop.NewGetMoveOrderOK().WithPayload(moveOrderPayload)
}

// ListMoveOrdersHandler fetches all the move orders
type ListMoveOrdersHandler struct {
	handlers.HandlerContext
	services.MoveOrderFetcher
}

// Handle getting the all move orders
func (h ListMoveOrdersHandler) Handle(params moveorderop.ListMoveOrdersParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	moveOrders, err := h.ListMoveOrders()
	if err != nil {
		logger.Error("fetching all move orders", zap.Error(err))
		switch err {
		case sql.ErrNoRows:
			return moveorderop.NewListMoveOrdersNotFound()
		default:
			return moveorderop.NewListMoveOrdersInternalServerError()
		}
	}
	moveOrdersPayload := make(ghcmessages.MoveOrders, len(moveOrders))
	for i, moveOrder := range moveOrders {
		moveOrdersPayload[i] = payloads.MoveOrder(&moveOrder)
	}
	return moveorderop.NewListMoveOrdersOK().WithPayload(moveOrdersPayload)
}

// ListMoveTaskOrdersHandler fetches all the move orders
type ListMoveTaskOrdersHandler struct {
	handlers.HandlerContext
	services.MoveTaskOrderFetcher
}

// Handle getting the all move orders
func (h ListMoveTaskOrdersHandler) Handle(params moveorderop.ListMoveTaskOrdersParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	moveOrderID, _ := uuid.FromString(params.MoveOrderID.String())
	moveTaskOrders, err := h.ListMoveTaskOrders(moveOrderID)
	if err != nil {
		logger.Error("fetching all move orders", zap.Error(err))
		switch err {
		case sql.ErrNoRows:
			return moveorderop.NewListMoveTaskOrdersNotFound()
		default:
			return moveorderop.NewListMoveTaskOrdersInternalServerError()
		}
	}
	moveTaskOrdersPayload := make(ghcmessages.MoveTaskOrders, len(moveTaskOrders))
	for i, moveTaskOrder := range moveTaskOrders {
		moveTaskOrdersPayload[i] = payloads.MoveTaskOrder(&moveTaskOrder)
	}
	return moveorderop.NewListMoveTaskOrdersOK().WithPayload(moveTaskOrdersPayload)
}
