package ghcapi

import (
	"database/sql"
	"time"

	"github.com/transcom/mymove/pkg/models/roles"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	moveorderop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_order"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
)

// GetMoveOrdersHandler fetches the information of a specific move order
type GetMoveOrdersHandler struct {
	handlers.HandlerContext
	services.MoveOrderFetcher
}

// Handle getting the information of a specific move order
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
	// get the session from http request
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	officeUserAuthorized := session.Roles.HasRole(roles.RoleTypeTOO)
	if !officeUserAuthorized {
		return moveorderop.NewListMoveOrdersForbidden()
	}

	// list move orders and pass in office user ID as argument to filter list
	moveOrders, err := h.MoveOrderFetcher.ListMoveOrders(session.OfficeUserID)
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

// UpdateMoveOrderHandler updates an order via PATCH /move-orders/{moveOrderId}
type UpdateMoveOrderHandler struct {
	handlers.HandlerContext
	moveOrderUpdater services.MoveOrderUpdater
}

// Handle ... updates an order from a request payload
func (h UpdateMoveOrderHandler) Handle(params moveorderop.UpdateMoveOrderParams) middleware.Responder {

	_, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	orderID, err := uuid.FromString(params.MoveOrderID.String())
	if err != nil {
		logger.Error("unable to parse move order id param to uuid", zap.Error(err))
		return moveorderop.NewUpdateMoveOrderBadRequest()
	}

	newOrder, err := MoveOrder(*params.Body)
	if err != nil {
		logger.Error("error converting payload to move order model", zap.Error(err))
		return moveorderop.NewUpdateMoveOrderBadRequest()
	}
	newOrder.ID = orderID

	updatedOrder, err := h.moveOrderUpdater.UpdateMoveOrder(orderID, params.IfMatch, newOrder)
	if err != nil {
		logger.Error("error updating move order", zap.Error(err))
		switch err.(type) {
		case services.NotFoundError:
			return moveorderop.NewUpdateMoveOrderNotFound()
		case services.InvalidInputError:
			return moveorderop.NewUpdateMoveOrderBadRequest()
		case services.PreconditionFailedError:
			return moveorderop.NewUpdateMoveOrderPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return moveorderop.NewUpdateMoveOrderInternalServerError()
		}
	}

	moveOrderPayload := payloads.MoveOrder(updatedOrder)

	return moveorderop.NewUpdateMoveOrderOK().WithPayload(moveOrderPayload)
}

// MoveOrder transforms UpdateMoveOrderPayload to Order model
func MoveOrder(payload ghcmessages.UpdateMoveOrderPayload) (models.Order, error) {

	ordersTypeDetail := internalmessages.OrdersTypeDetail(payload.OrdersTypeDetail)

	var originDutyStationID uuid.UUID
	if payload.OriginDutyStationID != nil {
		originDutyStationID = uuid.FromStringOrNil(payload.OriginDutyStationID.String())
	}

	newDutyStationID, err := uuid.FromString(payload.NewDutyStationID.String())
	if err != nil {
		return models.Order{}, err
	}

	departmentIndicator := string(payload.DepartmentIndicator)

	return models.Order{
		IssueDate:           time.Time(*payload.IssueDate),
		ReportByDate:        time.Time(*payload.ReportByDate),
		OrdersType:          internalmessages.OrdersType(payload.OrdersType),
		OrdersTypeDetail:    &ordersTypeDetail,
		NewDutyStationID:    newDutyStationID,
		OrdersNumber:        payload.OrdersNumber,
		TAC:                 payload.Tac,
		SAC:                 payload.Sac,
		DepartmentIndicator: &departmentIndicator,
		OriginDutyStationID: &originDutyStationID,
	}, nil
}
