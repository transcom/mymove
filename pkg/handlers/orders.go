package handlers

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/auth"
	ordersop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/orders"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForOrdersModel(storage FileStorer, order models.Order) (*internalmessages.Orders, error) {
	documentPayload, err := payloadForDocumentModel(storage, order.UploadedOrders)
	if err != nil {
		return nil, err
	}

	var moves internalmessages.IndexMovesPayload
	for _, move := range order.Moves {
		payload, err := payloadForMoveModel(storage, order, move)
		if err != nil {
			return nil, err
		}
		moves = append(moves, payload)
	}

	payload := &internalmessages.Orders{
		ID:                  fmtUUID(order.ID),
		CreatedAt:           fmtDateTime(order.CreatedAt),
		UpdatedAt:           fmtDateTime(order.UpdatedAt),
		ServiceMemberID:     fmtUUID(order.ServiceMemberID),
		IssueDate:           fmtDate(order.IssueDate),
		ReportByDate:        fmtDate(order.ReportByDate),
		OrdersType:          order.OrdersType,
		OrdersTypeDetail:    order.OrdersTypeDetail,
		NewDutyStation:      payloadForDutyStationModel(order.NewDutyStation),
		HasDependents:       fmtBool(order.HasDependents),
		SpouseHasProGear:    fmtBool(order.SpouseHasProGear),
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
		return responseForError(h.logger, err)
	}
	serviceMember, err := models.FetchServiceMember(h.db, session, serviceMemberID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	stationID, err := uuid.FromString(payload.NewDutyStationID.String())
	if err != nil {
		return responseForError(h.logger, err)
	}
	dutyStation, err := models.FetchDutyStation(h.db, stationID)
	if err != nil {
		return responseForError(h.logger, err)
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
		return responseForVErrors(h.logger, verrs, err)
	}

	// TODO: Don't default to PPM when we start supporting HHG
	newMoveType := internalmessages.SelectedMoveTypePPM
	newMove, verrs, err := newOrder.CreateNewMove(h.db, &newMoveType)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}
	newOrder.Moves = append(newOrder.Moves, *newMove)

	orderPayload, err := payloadForOrdersModel(h.storage, newOrder)
	if err != nil {
		return responseForError(h.logger, err)
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
	order, err := models.FetchOrder(h.db, session, orderID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	orderPayload, err := payloadForOrdersModel(h.storage, order)
	if err != nil {
		return responseForError(h.logger, err)
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
		return responseForError(h.logger, err)
	}
	order, err := models.FetchOrder(h.db, session, orderID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	payload := params.UpdateOrders
	stationID, err := uuid.FromString(payload.NewDutyStationID.String())
	if err != nil {
		return responseForError(h.logger, err)
	}
	dutyStation, err := models.FetchDutyStation(h.db, stationID)
	if err != nil {
		return responseForError(h.logger, err)
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
		order.DepartmentIndicator = fmtString(string(*payload.DepartmentIndicator))
	}

	verrs, err := models.SaveOrder(h.db, &order)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	orderPayload, err := payloadForOrdersModel(h.storage, order)
	if err != nil {
		return responseForError(h.logger, err)
	}
	return ordersop.NewUpdateOrdersOK().WithPayload(orderPayload)
}
