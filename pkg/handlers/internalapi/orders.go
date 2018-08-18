package internalapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/auth"
	ordersop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/orders"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

func payloadForOrdersModel(storer storage.FileStorer, order models.Order) (*internalmessages.Orders, error) {
	documentPayload, err := payloadForDocumentModel(storer, order.UploadedOrders)
	if err != nil {
		return nil, err
	}

	var moves internalmessages.IndexMovesPayload
	for _, move := range order.Moves {
		payload, err := payloadForMoveModel(storer, order, move)
		if err != nil {
			return nil, err
		}
		moves = append(moves, payload)
	}

	payload := &internalmessages.Orders{
		ID:                  utils.FmtUUID(order.ID),
		CreatedAt:           utils.FmtDateTime(order.CreatedAt),
		UpdatedAt:           utils.FmtDateTime(order.UpdatedAt),
		ServiceMemberID:     utils.FmtUUID(order.ServiceMemberID),
		IssueDate:           utils.FmtDate(order.IssueDate),
		ReportByDate:        utils.FmtDate(order.ReportByDate),
		OrdersType:          order.OrdersType,
		OrdersTypeDetail:    order.OrdersTypeDetail,
		NewDutyStation:      payloadForDutyStationModel(order.NewDutyStation),
		HasDependents:       utils.FmtBool(order.HasDependents),
		SpouseHasProGear:    utils.FmtBool(order.SpouseHasProGear),
		UploadedOrders:      documentPayload,
		OrdersNumber:        order.OrdersNumber,
		Moves:               moves,
		Tac:                 order.TAC,
		DepartmentIndicator: (*internalmessages.DeptIndicator)(order.DepartmentIndicator),
		Status:              internalmessages.OrdersStatus(order.Status),
	}

	return payload, nil
}

// CreateOrdersHandler creates new orders via POST /orders
type CreateOrdersHandler HandlerContext

// Handle ... creates new Orders from a request payload
func (h CreateOrdersHandler) Handle(params ordersop.CreateOrdersParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	payload := params.CreateOrders

	serviceMemberID, err := uuid.FromString(payload.ServiceMemberID.String())
	if err != nil {
		return utils.ResponseForError(h.logger, err)
	}
	serviceMember, err := models.FetchServiceMemberForUser(h.db, session, serviceMemberID)
	if err != nil {
		return utils.ResponseForError(h.logger, err)
	}

	stationID, err := uuid.FromString(payload.NewDutyStationID.String())
	if err != nil {
		return utils.ResponseForError(h.logger, err)
	}
	dutyStation, err := models.FetchDutyStation(h.db, stationID)
	if err != nil {
		return utils.ResponseForError(h.logger, err)
	}

	newOrder, verrs, err := serviceMember.CreateOrder(
		h.db,
		time.Time(*payload.IssueDate),
		time.Time(*payload.ReportByDate),
		payload.OrdersType,
		*payload.HasDependents,
		*payload.SpouseHasProGear,
		dutyStation)
	if err != nil || verrs.HasAny() {
		return utils.ResponseForVErrors(h.logger, verrs, err)
	}

	// TODO: Don't default to PPM when we start supporting HHG
	newMoveType := internalmessages.SelectedMoveTypePPM
	newMove, verrs, err := newOrder.CreateNewMove(h.db, &newMoveType)
	if err != nil || verrs.HasAny() {
		return utils.ResponseForVErrors(h.logger, verrs, err)
	}
	newOrder.Moves = append(newOrder.Moves, *newMove)

	orderPayload, err := payloadForOrdersModel(h.storage, newOrder)
	if err != nil {
		return utils.ResponseForError(h.logger, err)
	}
	return ordersop.NewCreateOrdersCreated().WithPayload(orderPayload)
}

// ShowOrdersHandler returns orders for a user and order ID
type ShowOrdersHandler HandlerContext

// Handle retrieves orders in the system belonging to the logged in user given order ID
func (h ShowOrdersHandler) Handle(params ordersop.ShowOrdersParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// #nosec swagger verifies uuid format
	orderID, _ := uuid.FromString(params.OrdersID.String())
	order, err := models.FetchOrderForUser(h.db, session, orderID)
	if err != nil {
		return utils.ResponseForError(h.logger, err)
	}

	orderPayload, err := payloadForOrdersModel(h.storage, order)
	if err != nil {
		return utils.ResponseForError(h.logger, err)
	}
	return ordersop.NewShowOrdersOK().WithPayload(orderPayload)
}

// UpdateOrdersHandler updates an order via PUT /orders/{orderId}
type UpdateOrdersHandler HandlerContext

// Handle ... updates an order from a request payload
func (h UpdateOrdersHandler) Handle(params ordersop.UpdateOrdersParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	orderID, err := uuid.FromString(params.OrdersID.String())
	if err != nil {
		return utils.ResponseForError(h.logger, err)
	}
	order, err := models.FetchOrderForUser(h.db, session, orderID)
	if err != nil {
		return utils.ResponseForError(h.logger, err)
	}

	payload := params.UpdateOrders
	stationID, err := uuid.FromString(payload.NewDutyStationID.String())
	if err != nil {
		return utils.ResponseForError(h.logger, err)
	}
	dutyStation, err := models.FetchDutyStation(h.db, stationID)
	if err != nil {
		return utils.ResponseForError(h.logger, err)
	}

	order.OrdersNumber = payload.OrdersNumber
	order.IssueDate = time.Time(*payload.IssueDate)
	order.ReportByDate = time.Time(*payload.ReportByDate)
	order.OrdersType = payload.OrdersType
	order.OrdersTypeDetail = payload.OrdersTypeDetail
	order.HasDependents = *payload.HasDependents
	order.SpouseHasProGear = *payload.SpouseHasProGear
	order.NewDutyStationID = dutyStation.ID
	order.NewDutyStation = dutyStation
	order.TAC = payload.Tac

	if payload.DepartmentIndicator != nil {
		order.DepartmentIndicator = utils.FmtString(string(*payload.DepartmentIndicator))
	}

	verrs, err := models.SaveOrder(h.db, &order)
	if err != nil || verrs.HasAny() {
		return utils.ResponseForVErrors(h.logger, verrs, err)
	}

	orderPayload, err := payloadForOrdersModel(h.storage, order)
	if err != nil {
		return utils.ResponseForError(h.logger, err)
	}
	return ordersop.NewUpdateOrdersOK().WithPayload(orderPayload)
}
