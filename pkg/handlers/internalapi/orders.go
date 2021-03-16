package internalapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

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

	var dBAuthorizedWeight *int64
	dBAuthorizedWeight = nil
	if order.Entitlement != nil {
		dBAuthorizedWeight = swag.Int64(int64(*order.Entitlement.AuthorizedWeight()))
	}
	var originDutyStation models.DutyStation
	originDutyStation = models.DutyStation{}
	if order.OriginDutyStation != nil {
		originDutyStation = *order.OriginDutyStation
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
		OriginDutyStation:   payloadForDutyStationModel(originDutyStation),
		Grade:               order.Grade,
		NewDutyStation:      payloadForDutyStationModel(order.NewDutyStation),
		HasDependents:       handlers.FmtBool(order.HasDependents),
		SpouseHasProGear:    handlers.FmtBool(order.SpouseHasProGear),
		UploadedOrders:      documentPayload,
		OrdersNumber:        order.OrdersNumber,
		Moves:               moves,
		Tac:                 order.TAC,
		Sac:                 order.SAC,
		DepartmentIndicator: (*internalmessages.DeptIndicator)(order.DepartmentIndicator),
		Status:              internalmessages.OrdersStatus(order.Status),
		AuthorizedWeight:    dBAuthorizedWeight,
	}

	return payload, nil
}

// CreateOrdersHandler creates new orders via POST /orders
type CreateOrdersHandler struct {
	handlers.HandlerContext
}

// Handle ... creates new Orders from a request payload
func (h CreateOrdersHandler) Handle(params ordersop.CreateOrdersParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	payload := params.CreateOrders

	serviceMemberID, err := uuid.FromString(payload.ServiceMemberID.String())
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	serviceMember, err := models.FetchServiceMemberForUser(h.DB(), session, serviceMemberID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	stationID, err := uuid.FromString(payload.NewDutyStationID.String())
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	newDutyStation, err := models.FetchDutyStation(h.DB(), stationID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	originDutyStation := serviceMember.DutyStation
	grade := (*string)(serviceMember.Rank)

	weight, entitlementErr := models.GetEntitlement(*serviceMember.Rank, *payload.HasDependents)
	if entitlementErr != nil {
		return handlers.ResponseForError(logger, entitlementErr)
	}
	entitlement := models.Entitlement{
		DependentsAuthorized: payload.HasDependents,
		DBAuthorizedWeight:   models.IntPointer(weight),
	}

	if saveEntitlementErr := h.DB().Save(&entitlement); saveEntitlementErr != nil {
		return handlers.ResponseForError(logger, saveEntitlementErr)
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
		newDutyStation,
		payload.OrdersNumber,
		payload.Tac,
		payload.Sac,
		deptIndicator,
		&originDutyStation,
		grade,
		&entitlement,
	)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(logger, verrs, err)
	}

	moveOptions := models.MoveOptions{
		SelectedType: nil,
		Show:         swag.Bool(true),
	}
	newMove, verrs, err := newOrder.CreateNewMove(h.DB(), moveOptions)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(logger, verrs, err)
	}
	newOrder.Moves = append(newOrder.Moves, *newMove)

	orderPayload, err := payloadForOrdersModel(h.FileStorer(), newOrder)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	return ordersop.NewCreateOrdersCreated().WithPayload(orderPayload)
}

// ShowOrdersHandler returns orders for a user and order ID
type ShowOrdersHandler struct {
	handlers.HandlerContext
}

// Handle retrieves orders in the system belonging to the logged in user given order ID
func (h ShowOrdersHandler) Handle(params ordersop.ShowOrdersParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	// #nosec swagger verifies uuid format
	orderID, _ := uuid.FromString(params.OrdersID.String())
	order, err := models.FetchOrderForUser(h.DB(), session, orderID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	orderPayload, err := payloadForOrdersModel(h.FileStorer(), order)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	return ordersop.NewShowOrdersOK().WithPayload(orderPayload)
}

// UpdateOrdersHandler updates an order via PUT /orders/{orderId}
type UpdateOrdersHandler struct {
	handlers.HandlerContext
}

// Handle ... updates an order from a request payload
func (h UpdateOrdersHandler) Handle(params ordersop.UpdateOrdersParams) middleware.Responder {

	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	orderID, err := uuid.FromString(params.OrdersID.String())
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	order, err := models.FetchOrderForUser(h.DB(), session, orderID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	payload := params.UpdateOrders
	stationID, err := uuid.FromString(payload.NewDutyStationID.String())
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	dutyStation, err := models.FetchDutyStation(h.DB(), stationID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
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
	order.SAC = payload.Sac

	if payload.DepartmentIndicator != nil {
		order.DepartmentIndicator = handlers.FmtString(string(*payload.DepartmentIndicator))
	}

	verrs, err := models.SaveOrder(h.DB(), &order)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(logger, verrs, err)
	}

	orderPayload, err := payloadForOrdersModel(h.FileStorer(), order)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	return ordersop.NewUpdateOrdersOK().WithPayload(orderPayload)
}
