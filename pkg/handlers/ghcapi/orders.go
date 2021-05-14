package ghcapi

import (
	"database/sql"
	"time"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/validate/v3"

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

// GetOrdersHandler fetches the information of a specific order
type GetOrdersHandler struct {
	handlers.HandlerContext
	services.OrderFetcher
}

// Handle getting the information of a specific order
func (h GetOrdersHandler) Handle(params orderop.GetOrderParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	orderID, _ := uuid.FromString(params.OrderID.String())
	order, err := h.FetchOrder(orderID)
	if err != nil {
		logger.Error("fetching order", zap.Error(err))
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

// ListMoveTaskOrdersHandler fetches all the moves
type ListMoveTaskOrdersHandler struct {
	handlers.HandlerContext
	services.MoveTaskOrderFetcher
}

// Handle getting the all moves
func (h ListMoveTaskOrdersHandler) Handle(params orderop.ListMoveTaskOrdersParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	orderID, _ := uuid.FromString(params.OrderID.String())
	moveTaskOrders, err := h.ListMoveTaskOrders(orderID, nil) // nil searchParams exclude disabled MTOs by default
	if err != nil {
		logger.Error("fetching all moves", zap.Error(err))
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

// UpdateOrderHandler updates an order via PATCH /orders/{orderId}
type UpdateOrderHandler struct {
	handlers.HandlerContext
	orderUpdater services.OrderUpdater
	orderFetcher services.OrderFetcher
}

// Handle ... updates an order from a request payload
func (h UpdateOrderHandler) Handle(params orderop.UpdateOrderParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	handleError := func(err error) middleware.Responder {
		logger.Error("error updating order", zap.Error(err))
		switch err.(type) {
		case *services.BadDataError:
			return orderop.NewUpdateOrderBadRequest()
		case services.NotFoundError:
			return orderop.NewUpdateOrderNotFound()
		case services.InvalidInputError:
			payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceID(), validate.NewErrors())
			return orderop.NewUpdateOrderUnprocessableEntity().WithPayload(payload)
		case services.PreconditionFailedError:
			return orderop.NewUpdateOrderPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case services.ForbiddenError:
			return orderop.NewUpdateAllowanceForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return orderop.NewUpdateOrderInternalServerError()
		}
	}

	if !session.IsOfficeUser() || (!session.Roles.HasRole(roles.RoleTypePPMOfficeUsers) &&
		!session.Roles.HasRole(roles.RoleTypeTOO) && !session.Roles.HasRole(roles.RoleTypeTIO) &&
		!session.Roles.HasRole(roles.RoleTypeServicesCounselor)) {
		return handleError(services.NewForbiddenError("is not an user with roles ppm office, TOO, TIO or Service Counselor"))
	}

	// Parsing order id
	orderID, err := uuid.FromString(params.OrderID.String())
	if err != nil {
		return handleError(services.NewBadDataError("unable to parse order id param to uuid"))
	}

	// make sure order exists
	existingOrder, err := h.orderFetcher.FetchOrder(orderID)
	if err != nil {
		return handleError(services.NewNotFoundError(orderID, "while looking for order"))
	}

	// make sure duty station exists

	// make sure eTag matches
	existingETag := etag.GenerateEtag(existingOrder.UpdatedAt)
	if existingETag != params.IfMatch {
		return handleError(services.NewPreconditionFailedError(orderID, query.StaleIdentifierError{StaleIdentifier: params.IfMatch}))
	}

	// good to update
	updatingOrder := Order(*existingOrder, *params.Body, session)
	updatedOrder, err := h.orderUpdater.UpdateOrder(updatingOrder)
	if err != nil {
		return handleError(err)
	}

	// For capturing event information
	// Find the record where orderID matches order.ID
	var move models.Move
	query := h.DB().Where("orders_id = ?", updatedOrder.ID)
	err = query.First(&move)

	var moveID = move.ID
	if err != nil {
		logger.Error("ghcapi.UpdateOrderHandler could not find move")
		moveID = uuid.Nil
	}

	// UpdateOrder event Trigger for the first updated move:
	_, err = event.TriggerEvent(event.Event{
		EndpointKey: event.GhcUpdateOrderEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.OrderUpdateEventKey, // Event that you want to trigger
		UpdatedObjectID: updatedOrder.ID,           // ID of the updated logical object (look at what the payload returns)
		MtoID:           moveID,                    // ID of the associated Move
		Request:         params.HTTPRequest,        // Pass on the http.Request
		DBConnection:    h.DB(),                    // Pass on the pop.Connection
		HandlerContext:  h,                         // Pass on the handlerContext
	})
	// If the event trigger fails, just log the error.
	if err != nil {
		logger.Error("ghcapi.UpdateOrderHandler could not generate the event")
	}

	orderPayload := payloads.Order(updatedOrder)

	return orderop.NewUpdateOrderOK().WithPayload(orderPayload)
}

// UpdateAllowanceHandler updates an order and entitlements via PATCH /orders/{orderId}/allowances
type UpdateAllowanceHandler struct {
	handlers.HandlerContext
	orderUpdater services.OrderUpdater
	orderFetcher services.OrderFetcher
}

// Handle ... updates an order from a request payload
func (h UpdateAllowanceHandler) Handle(params orderop.UpdateAllowanceParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	handleError := func(err error) middleware.Responder {
		logger.Error("error updating order allowance", zap.Error(err))
		switch err.(type) {
		case *services.BadDataError:
			return orderop.NewUpdateAllowanceBadRequest()
		case services.NotFoundError:
			return orderop.NewUpdateAllowanceNotFound()
		case services.InvalidInputError:
			payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceID(), validate.NewErrors())
			return orderop.NewUpdateAllowanceUnprocessableEntity().WithPayload(payload)
		case services.PreconditionFailedError:
			return orderop.NewUpdateAllowancePreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case services.ForbiddenError:
			return orderop.NewUpdateAllowanceForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return orderop.NewUpdateAllowanceInternalServerError()
		}
	}

	if !session.IsOfficeUser() || (!session.Roles.HasRole(roles.RoleTypePPMOfficeUsers) &&
		!session.Roles.HasRole(roles.RoleTypeTOO) && !session.Roles.HasRole(roles.RoleTypeTIO) &&
		!session.Roles.HasRole(roles.RoleTypeServicesCounselor)) {
		return handleError(services.NewForbiddenError("is not an user with roles ppm office, TOO, TIO or Service Counselor"))
	}

	// Parsing order id
	orderID, err := uuid.FromString(params.OrderID.String())
	if err != nil {
		return handleError(services.NewBadDataError("unable to parse order id param to uuid"))
	}

	// make sure order exists
	existingOrder, err := h.orderFetcher.FetchOrder(orderID)
	if err != nil {
		return handleError(services.NewNotFoundError(orderID, "while looking for order"))
	}

	// make sure eTag matches
	existingETag := etag.GenerateEtag(existingOrder.UpdatedAt)
	if existingETag != params.IfMatch {
		return handleError(services.NewPreconditionFailedError(orderID, query.StaleIdentifierError{StaleIdentifier: params.IfMatch}))
	}

	// good to update
	updatingOrder := OrderAllowance(*existingOrder, *params.Body, session)
	updatedOrder, err := h.orderUpdater.UpdateOrder(updatingOrder)
	if err != nil {
		return handleError(err)
	}

	// For capturing event information
	// Find the record where orderID matches order.ID
	var move models.Move
	query := h.DB().Where("orders_id = ?", updatedOrder.ID)
	err = query.First(&move)

	var moveID = move.ID
	if err != nil {
		logger.Error("ghcapi.UpdateAllowanceHandler could not find move")
		moveID = uuid.Nil
	}

	// UpdateOrder event Trigger for the first updated move:
	_, err = event.TriggerEvent(event.Event{
		EndpointKey: event.GhcUpdateAllowanceEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.OrderUpdateEventKey, // Event that you want to trigger
		UpdatedObjectID: updatedOrder.ID,           // ID of the updated logical object (look at what the payload returns)
		MtoID:           moveID,                    // ID of the associated Move
		Request:         params.HTTPRequest,        // Pass on the http.Request
		DBConnection:    h.DB(),                    // Pass on the pop.Connection
		HandlerContext:  h,                         // Pass on the handlerContext
	})
	// If the event trigger fails, just log the error.
	if err != nil {
		logger.Error("ghcapi.UpdateAllowanceHandler could not generate the event")
	}

	orderPayload := payloads.Order(updatedOrder)

	return orderop.NewUpdateAllowanceOK().WithPayload(orderPayload)
}

// Order transforms UpdateOrderPayload to Order model
func Order(existingOrder models.Order, payload ghcmessages.UpdateOrderPayload, session *auth.Session) models.Order {
	isServicesCounselor := session.Roles.HasRole(roles.RoleTypeServicesCounselor)
	order := existingOrder

	// update both order origin duty station and service member duty station
	if payload.OriginDutyStationID != nil {
		originDutyStationID := uuid.FromStringOrNil(payload.OriginDutyStationID.String())
		order.OriginDutyStationID = &originDutyStationID
		order.ServiceMember.DutyStationID = &originDutyStationID
	}

	if payload.NewDutyStationID != nil {
		newDutyStationID := uuid.FromStringOrNil(payload.NewDutyStationID.String())
		order.NewDutyStationID = newDutyStationID
	}

	if payload.DepartmentIndicator != nil && !isServicesCounselor {
		departmentIndicator := (*string)(payload.DepartmentIndicator)
		order.DepartmentIndicator = departmentIndicator
	}

	if payload.IssueDate != nil {
		order.IssueDate = time.Time(*payload.IssueDate)
	}

	if payload.OrdersNumber != nil && !isServicesCounselor {
		order.OrdersNumber = payload.OrdersNumber
	}

	if payload.OrdersTypeDetail != nil && !isServicesCounselor {
		orderTypeDetail := internalmessages.OrdersTypeDetail(*payload.OrdersTypeDetail)
		order.OrdersTypeDetail = &orderTypeDetail
	}

	if payload.ReportByDate != nil {
		order.ReportByDate = time.Time(*payload.ReportByDate)
	}

	if payload.Sac != nil && !isServicesCounselor {
		order.SAC = payload.Sac
	}

	if payload.Tac != nil && !isServicesCounselor {
		order.TAC = payload.Tac
	}

	order.OrdersType = internalmessages.OrdersType(payload.OrdersType)

	return order
}

// OrderAllowance transforms UpdateOrderPayload to Order model. Specifically for
// UpdateAllowance endpoint.
func OrderAllowance(existingOrder models.Order, payload ghcmessages.UpdateAllowancePayload, session *auth.Session) models.Order {
	isServicesCounselor := session.Roles.HasRole(roles.RoleTypeServicesCounselor)
	order := existingOrder

	if payload.ProGearWeight != nil {
		order.Entitlement.ProGearWeight = int(*payload.ProGearWeight)
	}

	if payload.ProGearWeightSpouse != nil {
		order.Entitlement.ProGearWeightSpouse = int(*payload.ProGearWeightSpouse)
	}

	if payload.RequiredMedicalEquipmentWeight != nil {
		order.Entitlement.RequiredMedicalEquipmentWeight = int(*payload.RequiredMedicalEquipmentWeight)
	}

	// branch for service member
	if payload.Agency != "" {
		serviceMemberAffiliation := models.ServiceMemberAffiliation(payload.Agency)
		order.ServiceMember.Affiliation = &serviceMemberAffiliation
	}

	// rank
	if payload.Grade != nil {
		grade := (*string)(payload.Grade)
		order.Grade = grade
	}

	if payload.OrganizationalClothingAndIndividualEquipment != nil {
		order.Entitlement.OrganizationalClothingAndIndividualEquipment = *payload.OrganizationalClothingAndIndividualEquipment
	}

	// only office users and TXO roles can edit authorized weight
	// omit value from update
	if payload.AuthorizedWeight != nil && !isServicesCounselor {
		dbAuthorizedWeight := swag.Int(int(*payload.AuthorizedWeight))
		order.Entitlement.DBAuthorizedWeight = dbAuthorizedWeight
	}

	if payload.DependentsAuthorized != nil {
		order.Entitlement.DependentsAuthorized = payload.DependentsAuthorized
	}

	return order
}
