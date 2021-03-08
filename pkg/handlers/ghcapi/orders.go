package ghcapi

import (
	"database/sql"
	"time"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/event"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	orderop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/order"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
)

// GetOrdersHandler fetches the information of a specific move order
type GetOrdersHandler struct {
	handlers.HandlerContext
	services.OrderFetcher
}

// Handle getting the information of a specific move order
func (h GetOrdersHandler) Handle(params orderop.GetOrderParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	orderID, _ := uuid.FromString(params.OrderID.String())
	order, err := h.FetchOrder(orderID)
	if err != nil {
		logger.Error("fetching move order", zap.Error(err))
		switch err {
		case sql.ErrNoRows:
			return orderop.NewGetOrderNotFound()
		default:
			return orderop.NewGetOrderInternalServerError()
		}
	}
	orderPayload := payloads.Order(order)
	return orderop.NewGetOrderOK().WithPayload(orderPayload)
}

// ListMoveTaskOrdersHandler fetches all the move orders
type ListMoveTaskOrdersHandler struct {
	handlers.HandlerContext
	services.MoveTaskOrderFetcher
}

// Handle getting the all move orders
func (h ListMoveTaskOrdersHandler) Handle(params orderop.ListMoveTaskOrdersParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	orderID, _ := uuid.FromString(params.OrderID.String())
	moveTaskOrders, err := h.ListMoveTaskOrders(orderID, nil) // nil searchParams exclude disabled MTOs by default
	if err != nil {
		logger.Error("fetching all move orders", zap.Error(err))
		switch err {
		case sql.ErrNoRows:
			return orderop.NewListMoveTaskOrdersNotFound()
		default:
			return orderop.NewListMoveTaskOrdersInternalServerError()
		}
	}
	moveTaskOrdersPayload := make(ghcmessages.MoveTaskOrders, len(moveTaskOrders))
	for i, moveTaskOrder := range moveTaskOrders {
		copyOfMto := moveTaskOrder // Make copy to avoid implicit memory aliasing of items from a range statement.
		moveTaskOrdersPayload[i] = payloads.MoveTaskOrder(&copyOfMto)
	}
	return orderop.NewListMoveTaskOrdersOK().WithPayload(moveTaskOrdersPayload)
}

// UpdateMoveOrderHandler updates an order via PATCH /move-orders/{orderId}
type UpdateMoveOrderHandler struct {
	handlers.HandlerContext
	orderUpdater services.OrderUpdater
}

// Handle ... updates an order from a request payload
func (h UpdateMoveOrderHandler) Handle(params orderop.UpdateMoveOrderParams) middleware.Responder {
	_, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	orderID, err := uuid.FromString(params.OrderID.String())
	if err != nil {
		logger.Error("unable to parse move order id param to uuid", zap.Error(err))
		return orderop.NewUpdateMoveOrderBadRequest()
	}

	newOrder, err := Order(*params.Body)
	if err != nil {
		logger.Error("error converting payload to move order model", zap.Error(err))
		return orderop.NewUpdateMoveOrderBadRequest()
	}
	newOrder.ID = orderID

	updatedOrder, err := h.orderUpdater.UpdateOrder(params.IfMatch, newOrder)

	if err != nil {
		logger.Error("error updating move order", zap.Error(err))
		switch err.(type) {
		case services.NotFoundError:
			return orderop.NewUpdateMoveOrderNotFound()
		case services.InvalidInputError:
			return orderop.NewUpdateMoveOrderBadRequest()
		case services.PreconditionFailedError:
			return orderop.NewUpdateMoveOrderPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return orderop.NewUpdateMoveOrderInternalServerError()
		}
	}

	// Find the record where orderID matches order.ID
	var move models.Move
	query := h.DB().Where("orders_id = ?", updatedOrder.ID)
	err = query.First(&move)

	var moveID = move.ID

	if err != nil {
		logger.Error("ghcapi.UpdateMoveOrderHandler could not find move")
		moveID = uuid.Nil
	}

	// UpdateMoveOrder event Trigger for the first updated move:
	_, err = event.TriggerEvent(event.Event{
		EndpointKey: event.GhcUpdateMoveOrderEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.MoveOrderUpdateEventKey, // Event that you want to trigger
		UpdatedObjectID: updatedOrder.ID,               // ID of the updated logical object (look at what the payload returns)
		MtoID:           moveID,                        // ID of the associated Move
		Request:         params.HTTPRequest,            // Pass on the http.Request
		DBConnection:    h.DB(),                        // Pass on the pop.Connection
		HandlerContext:  h,                             // Pass on the handlerContext
	})
	// If the event trigger fails, just log the error.
	if err != nil {
		logger.Error("ghcapi.UpdateMoveOrderHandler could not generate the event")
	}

	orderPayload := payloads.Order(updatedOrder)

	return orderop.NewUpdateMoveOrderOK().WithPayload(orderPayload)
}

// Order transforms UpdateMoveOrderPayload to Order model
func Order(payload ghcmessages.UpdateMoveOrderPayload) (models.Order, error) {

	var originDutyStationID uuid.UUID
	if payload.OriginDutyStationID != nil {
		originDutyStationID = uuid.FromStringOrNil(payload.OriginDutyStationID.String())
	}

	newDutyStationID, err := uuid.FromString(payload.NewDutyStationID.String())
	if err != nil {
		return models.Order{}, err
	}

	var departmentIndicator *string
	if payload.DepartmentIndicator != nil {
		departmentIndicator = (*string)(payload.DepartmentIndicator)
	}

	var grade *string
	if payload.Grade != nil {
		grade = (*string)(payload.Grade)
	}

	var entitlement models.Entitlement
	if payload.AuthorizedWeight != nil {
		entitlement.DBAuthorizedWeight = swag.Int(int(*payload.AuthorizedWeight))
	}

	if payload.DependentsAuthorized != nil {
		entitlement.DependentsAuthorized = payload.DependentsAuthorized
	}

	var ordersTypeDetail *internalmessages.OrdersTypeDetail
	if payload.OrdersTypeDetail != nil {
		orderTypeDetail := internalmessages.OrdersTypeDetail(*payload.OrdersTypeDetail)
		ordersTypeDetail = &orderTypeDetail
	}

	var serviceMember models.ServiceMember
	if payload.Agency != "" {
		serviceMemberAffiliation := models.ServiceMemberAffiliation(payload.Agency)
		serviceMember.Affiliation = &serviceMemberAffiliation
	}

	return models.Order{
		ServiceMember:       serviceMember,
		DepartmentIndicator: departmentIndicator,
		Entitlement:         &entitlement,
		Grade:               grade,
		IssueDate:           time.Time(*payload.IssueDate),
		NewDutyStationID:    newDutyStationID,
		OrdersNumber:        payload.OrdersNumber,
		OrdersType:          internalmessages.OrdersType(payload.OrdersType),
		OrdersTypeDetail:    ordersTypeDetail,
		OriginDutyStationID: &originDutyStationID,
		ReportByDate:        time.Time(*payload.ReportByDate),
		SAC:                 payload.Sac,
		TAC:                 payload.Tac,
	}, nil

}
