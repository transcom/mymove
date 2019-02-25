package internalapi

import (
	"reflect"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/honeycombio/beeline-go"

	"github.com/transcom/mymove/pkg/auth"
	ordersop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/orders"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
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
		ID:                  handlers.FmtUUID(order.ID),
		CreatedAt:           handlers.FmtDateTime(order.CreatedAt),
		UpdatedAt:           handlers.FmtDateTime(order.UpdatedAt),
		ServiceMemberID:     handlers.FmtUUID(order.ServiceMemberID),
		IssueDate:           handlers.FmtDate(order.IssueDate),
		ReportByDate:        handlers.FmtDate(order.ReportByDate),
		OrdersType:          order.OrdersType,
		OrdersTypeDetail:    order.OrdersTypeDetail,
		NewDutyStation:      payloadForDutyStationModel(order.NewDutyStation),
		HasDependents:       handlers.FmtBool(order.HasDependents),
		SpouseHasProGear:    handlers.FmtBool(order.SpouseHasProGear),
		UploadedOrders:      documentPayload,
		OrdersNumber:        order.OrdersNumber,
		ParagraphNumber:     order.ParagraphNumber,
		OrdersIssuingAgency: order.OrdersIssuingAgency,
		Moves:               moves,
		Tac:                 order.TAC,
		Sac:                 order.SAC,
		DepartmentIndicator: (*internalmessages.DeptIndicator)(order.DepartmentIndicator),
		Status:              internalmessages.OrdersStatus(order.Status),
	}

	return payload, nil
}

// CreateOrdersHandler creates new orders via POST /orders
type CreateOrdersHandler struct {
	handlers.HandlerContext
}

// Handle ... creates new Orders from a request payload
func (h CreateOrdersHandler) Handle(params ordersop.CreateOrdersParams) middleware.Responder {

	ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	payload := params.CreateOrders

	serviceMemberID, err := uuid.FromString(payload.ServiceMemberID.String())
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	serviceMember, err := models.FetchServiceMemberForUser(ctx, h.DB(), session, serviceMemberID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	stationID, err := uuid.FromString(payload.NewDutyStationID.String())
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	dutyStation, err := models.FetchDutyStation(ctx, h.DB(), stationID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	var deptIndicator *string
	if payload.DepartmentIndicator != nil {
		converted := string(*payload.DepartmentIndicator)
		deptIndicator = &converted
	}

	newOrder, verrs, err := serviceMember.CreateOrder(
		h.DB(),
		time.Time(*payload.IssueDate),
		time.Time(*payload.ReportByDate),
		payload.OrdersType,
		*payload.HasDependents,
		*payload.SpouseHasProGear,
		dutyStation,
		payload.OrdersNumber,
		payload.ParagraphNumber,
		payload.OrdersIssuingAgency,
		payload.Tac,
		payload.Sac,
		deptIndicator)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	newMove, verrs, err := newOrder.CreateNewMove(h.DB(), nil)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}
	newOrder.Moves = append(newOrder.Moves, *newMove)

	orderPayload, err := payloadForOrdersModel(h.FileStorer(), newOrder)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	return ordersop.NewCreateOrdersCreated().WithPayload(orderPayload)
}

// ShowOrdersHandler returns orders for a user and order ID
type ShowOrdersHandler struct {
	handlers.HandlerContext
}

// Handle retrieves orders in the system belonging to the logged in user given order ID
func (h ShowOrdersHandler) Handle(params ordersop.ShowOrdersParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// #nosec swagger verifies uuid format
	orderID, _ := uuid.FromString(params.OrdersID.String())
	order, err := models.FetchOrderForUser(h.DB(), session, orderID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	orderPayload, err := payloadForOrdersModel(h.FileStorer(), order)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	return ordersop.NewShowOrdersOK().WithPayload(orderPayload)
}

// UpdateOrdersHandler updates an order via PUT /orders/{orderId}
type UpdateOrdersHandler struct {
	handlers.HandlerContext
}

// Handle ... updates an order from a request payload
func (h UpdateOrdersHandler) Handle(params ordersop.UpdateOrdersParams) middleware.Responder {

	ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	orderID, err := uuid.FromString(params.OrdersID.String())
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	order, err := models.FetchOrderForUser(h.DB(), session, orderID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	payload := params.UpdateOrders
	stationID, err := uuid.FromString(payload.NewDutyStationID.String())
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	dutyStation, err := models.FetchDutyStation(ctx, h.DB(), stationID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	order.OrdersNumber = payload.OrdersNumber
	order.ParagraphNumber = payload.ParagraphNumber
	order.OrdersIssuingAgency = payload.OrdersIssuingAgency
	order.IssueDate = time.Time(*payload.IssueDate)
	order.ReportByDate = time.Time(*payload.ReportByDate)
	order.OrdersType = payload.OrdersType
	order.OrdersTypeDetail = payload.OrdersTypeDetail
	order.HasDependents = *payload.HasDependents
	order.SpouseHasProGear = *payload.SpouseHasProGear
	order.NewDutyStationID = dutyStation.ID
	order.NewDutyStation = dutyStation
	order.TAC = payload.Tac
	order.SAC = payload.Sac

	if payload.DepartmentIndicator != nil {
		order.DepartmentIndicator = handlers.FmtString(string(*payload.DepartmentIndicator))
	}

	verrs, err := models.SaveOrder(h.DB(), &order)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	orderPayload, err := payloadForOrdersModel(h.FileStorer(), order)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	return ordersop.NewUpdateOrdersOK().WithPayload(orderPayload)
}
