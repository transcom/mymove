package internal

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
type CreateOrdersHandler utils.HandlerContext

// Handle ... creates new Orders from a request payload
func (h CreateOrdersHandler) Handle(params ordersop.CreateOrdersParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	payload := params.CreateOrders

	serviceMemberID, err := uuid.FromString(payload.ServiceMemberID.String())
	if err != nil {
		return responseForError(h.Logger, err)
	}
	serviceMember, err := models.FetchServiceMemberForUser(h.Db, session, serviceMemberID)
	if err != nil {
		return responseForError(h.Logger, err)
	}

	stationID, err := uuid.FromString(payload.NewDutyStationID.String())
	if err != nil {
		return responseForError(h.Logger, err)
	}
	dutyStation, err := models.FetchDutyStation(h.Db, stationID)
	if err != nil {
		return responseForError(h.Logger, err)
	}

	newOrder, verrs, err := serviceMember.CreateOrder(
		h.Db,
		time.Time(*payload.IssueDate),
		time.Time(*payload.ReportByDate),
		payload.OrdersType,
		*payload.HasDependents,
		*payload.SpouseHasProGear,
		dutyStation)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.Logger, verrs, err)
	}

	// TODO: Don't default to PPM when we start supporting HHG
	newMoveType := internalmessages.SelectedMoveTypePPM
	newMove, verrs, err := newOrder.CreateNewMove(h.Db, &newMoveType)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.Logger, verrs, err)
	}
	newOrder.Moves = append(newOrder.Moves, *newMove)

	orderPayload, err := payloadForOrdersModel(h.Storage, newOrder)
	if err != nil {
		return responseForError(h.Logger, err)
	}
	return ordersop.NewCreateOrdersCreated().WithPayload(orderPayload)
}

// ShowOrdersHandler returns orders for a user and order ID
type ShowOrdersHandler utils.HandlerContext

// Handle retrieves orders in the system belonging to the logged in user given order ID
func (h ShowOrdersHandler) Handle(params ordersop.ShowOrdersParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// #nosec swagger verifies uuid format
	orderID, _ := uuid.FromString(params.OrdersID.String())
	order, err := models.FetchOrderForUser(h.Db, session, orderID)
	if err != nil {
		return responseForError(h.Logger, err)
	}

	orderPayload, err := payloadForOrdersModel(h.Storage, order)
	if err != nil {
		return responseForError(h.Logger, err)
	}
	return ordersop.NewShowOrdersOK().WithPayload(orderPayload)
}

// UpdateOrdersHandler updates an order via PUT /orders/{orderId}
type UpdateOrdersHandler utils.HandlerContext

// Handle ... updates an order from a request payload
func (h UpdateOrdersHandler) Handle(params ordersop.UpdateOrdersParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	orderID, err := uuid.FromString(params.OrdersID.String())
	if err != nil {
		return responseForError(h.Logger, err)
	}
	order, err := models.FetchOrderForUser(h.Db, session, orderID)
	if err != nil {
		return responseForError(h.Logger, err)
	}

	payload := params.UpdateOrders
	stationID, err := uuid.FromString(payload.NewDutyStationID.String())
	if err != nil {
		return responseForError(h.Logger, err)
	}
	dutyStation, err := models.FetchDutyStation(h.Db, stationID)
	if err != nil {
		return responseForError(h.Logger, err)
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

	verrs, err := models.SaveOrder(h.Db, &order)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.Logger, verrs, err)
	}

	orderPayload, err := payloadForOrdersModel(h.Storage, order)
	if err != nil {
		return responseForError(h.Logger, err)
	}
	return ordersop.NewUpdateOrdersOK().WithPayload(orderPayload)
}
