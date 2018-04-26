package handlers

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth/context"
	ordersop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/orders"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForOrdersModel(user models.User, order models.Order) *internalmessages.OrdersPayload {
	return &internalmessages.OrdersPayload{
		ID:              fmtUUID(order.ID),
		CreatedAt:       fmtDateTime(order.CreatedAt),
		UpdatedAt:       fmtDateTime(order.UpdatedAt),
		ServiceMemberID: fmtUUID(order.ServiceMember.ID),
		IssueDate:       fmtDate(order.IssueDate),
		ReportByDate:    fmtDate(order.ReportByDate),
		OrdersType:      swag.String(order.OrdersType),
		NewDutyStation:  payloadForDutyStationModel(order.NewDutyStation),
		HasDependents:   fmtBool(order.HasDependents),
	}
}

// CreateOrdersHandler creates a new service member via POST /serviceMember
type CreateOrdersHandler HandlerContext

// Handle ... creates a new ServiceMember from a request payload
func (h CreateOrdersHandler) Handle(params ordersop.CreateOrdersParams) middleware.Responder {
	// User should always be populated by middleware
	user, _ := context.GetUser(params.HTTPRequest.Context())

	payload := params.CreateOrdersPayload

	serviceMemberID, err := uuid.FromString(payload.ServiceMemberID.String())
	if err != nil {
		return responseForError(h.logger, err)
	}
	serviceMember, err := models.FetchServiceMember(h.db, user, serviceMemberID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	stationID, err := uuid.FromString(payload.NewDutyStation.ID.String())
	if err != nil {
		return responseForError(h.logger, err)
	}
	dutyStation, err := models.FetchDutyStation(h.db, stationID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	newOrder := models.Order{
		ServiceMemberID:  serviceMember.ID,
		ServiceMember:    serviceMember,
		IssueDate:        time.Time(*payload.IssueDate),
		ReportByDate:     time.Time(*payload.ReportByDate),
		OrdersType:       *payload.OrdersType,
		HasDependents:    *payload.HasDependents,
		NewDutyStationID: dutyStation.ID,
		NewDutyStation:   dutyStation,
	}

	verrs, err := models.SaveOrder(h.db, &newOrder)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	orderPayload := payloadForOrdersModel(user, newOrder)
	return ordersop.NewCreateOrdersCreated().WithPayload(orderPayload)
}

// UpdateOrdersHandler updates an order via PUT /orders/{orderId}
type UpdateOrdersHandler HandlerContext

// Handle ... updates an order from a request payload
func (h UpdateOrdersHandler) Handle(params ordersop.UpdateOrdersParams) middleware.Responder {
	// User should always be populated by middleware
	user, _ := context.GetUser(params.HTTPRequest.Context())

	orderID, err := uuid.FromString(params.OrderID.String())
	if err != nil {
		return responseForError(h.logger, err)
	}
	order, err := models.FetchOrder(h.db, user, orderID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	payload := params.UpdateOrdersPayload
	stationID, err := uuid.FromString(payload.NewDutyStation.ID.String())
	if err != nil {
		return responseForError(h.logger, err)
	}
	dutyStation, err := models.FetchDutyStation(h.db, stationID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	order.IssueDate = time.Time(*payload.IssueDate)
	order.ReportByDate = time.Time(*payload.ReportByDate)
	order.OrdersType = *payload.OrdersType
	order.HasDependents = *payload.HasDependents
	order.NewDutyStationID = dutyStation.ID
	order.NewDutyStation = dutyStation

	verrs, err := models.SaveOrder(h.db, &order)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	orderPayload := payloadForOrdersModel(user, order)
	return ordersop.NewUpdateOrdersOK().WithPayload(orderPayload)
}
