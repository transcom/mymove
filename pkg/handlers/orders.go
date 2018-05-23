package handlers

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/app"
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
		payload := payloadForMoveModel(order, move)
		moves = append(moves, &payload)
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
		UploadedOrders:      documentPayload,
		OrdersNumber:        order.OrdersNumber,
		Moves:               moves,
		Tac:                 order.TAC,
		DepartmentIndicator: (*internalmessages.DeptIndicator)(order.DepartmentIndicator),
	}

	return payload, nil
}

// CreateOrdersHandler creates new orders via POST /orders
type CreateOrdersHandler HandlerContext

// Handle ... creates new Orders from a request payload
func (h CreateOrdersHandler) Handle(params ordersop.CreateOrdersParams) middleware.Responder {
	// User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	reqApp := app.GetAppFromContext(params.HTTPRequest)

	payload := params.CreateOrders

	serviceMemberID, err := uuid.FromString(payload.ServiceMemberID.String())
	if err != nil {
		return responseForError(h.logger, err)
	}
	serviceMember, err := models.FetchServiceMember(h.db, user, reqApp, serviceMemberID)
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
		dutyStation)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

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
	// User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	reqApp := app.GetAppFromContext(params.HTTPRequest)

	orderID, _ := uuid.FromString(params.OrdersID.String())
	order, err := models.FetchOrder(h.db, user, reqApp, orderID)
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
	// User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	reqApp := app.GetAppFromContext(params.HTTPRequest)

	orderID, err := uuid.FromString(params.OrdersID.String())
	if err != nil {
		return responseForError(h.logger, err)
	}
	order, err := models.FetchOrder(h.db, user, reqApp, orderID)
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
