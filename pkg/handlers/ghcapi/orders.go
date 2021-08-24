package ghcapi

import (
	"database/sql"
	"errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/move"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models/roles"
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
	handlers.HandlerConfig
	services.OrderFetcher
}

// Handle getting the information of a specific order
func (h GetOrdersHandler) Handle(params orderop.GetOrderParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	appCtx := h.AppContextFromRequest(params.HTTPRequest)
	orderID, _ := uuid.FromString(params.OrderID.String())
	order, err := h.FetchOrder(appCtx, orderID)
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

// UpdateOrderHandler updates an order via PATCH /orders/{orderId}
type UpdateOrderHandler struct {
	handlers.HandlerConfig
	orderUpdater services.OrderUpdater
	moveUpdater  services.MoveTaskOrderUpdater
}

func amendedOrdersRequiresApproval(params orderop.UpdateOrderParams, updatedOrder models.Order) bool {
	return params.Body.OrdersAcknowledgement != nil &&
		*params.Body.OrdersAcknowledgement &&
		updatedOrder.UploadedAmendedOrdersID != nil &&
		updatedOrder.AmendedOrdersAcknowledgedAt != nil
}

// Handle ... updates an order from a request payload
func (h UpdateOrderHandler) Handle(params orderop.UpdateOrderParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	appCtx := h.AppContextFromRequest(params.HTTPRequest)
	handleError := func(err error) middleware.Responder {
		logger.Error("error updating order", zap.Error(err))
		switch err.(type) {
		case services.NotFoundError:
			return orderop.NewUpdateOrderNotFound()
		case services.InvalidInputError:
			payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceID(), validate.NewErrors())
			return orderop.NewUpdateOrderUnprocessableEntity().WithPayload(payload)
		case services.ConflictError:
			return orderop.NewUpdateOrderConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case services.PreconditionFailedError:
			return orderop.NewUpdateOrderPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case services.ForbiddenError:
			return orderop.NewUpdateOrderForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return orderop.NewUpdateOrderInternalServerError()
		}
	}

	if !session.IsOfficeUser() || (!session.Roles.HasRole(roles.RoleTypeTOO) && !session.Roles.HasRole(roles.RoleTypeTIO)) {
		return handleError(services.NewForbiddenError("is not a TXO"))
	}

	orderID := uuid.FromStringOrNil(params.OrderID.String())
	updatedOrder, moveID, err := h.orderUpdater.UpdateOrderAsTOO(appCtx, orderID, *params.Body, params.IfMatch)
	if err != nil {
		return handleError(err)
	}

	h.triggerUpdateOrderEvent(appCtx, orderID, moveID, params)

	// the move status may need set back to approved if the amended orders upload caused it to be in approvals requested
	if amendedOrdersRequiresApproval(params, *updatedOrder) {
		moveRouter := move.NewMoveRouter()
		approvedMove, approveErr := moveRouter.ApproveAmendedOrders(appCtx, moveID, updatedOrder.ID)
		if approveErr != nil {
			if errors.Is(approveErr, models.ErrInvalidTransition) {
				return handleError(services.NewConflictError(moveID, approveErr.Error()))
			}
			return handleError(approveErr)
		}

		updateErr := h.moveUpdater.UpdateApprovedAmendedOrders(appCtx, approvedMove)
		if updateErr != nil {
			handleError(updateErr)
		}
	}

	orderPayload := payloads.Order(updatedOrder)

	return orderop.NewUpdateOrderOK().WithPayload(orderPayload)
}

// CounselingUpdateOrderHandler updates an order via PATCH /counseling/orders/{orderId}
type CounselingUpdateOrderHandler struct {
	handlers.HandlerConfig
	orderUpdater services.OrderUpdater
}

// Handle ... updates an order as requested by a services counselor
func (h CounselingUpdateOrderHandler) Handle(params orderop.CounselingUpdateOrderParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	handleError := func(err error) middleware.Responder {
		logger.Error("error updating order", zap.Error(err))
		switch err.(type) {
		case services.NotFoundError:
			return orderop.NewCounselingUpdateOrderNotFound()
		case services.InvalidInputError:
			payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceID(), validate.NewErrors())
			return orderop.NewCounselingUpdateOrderUnprocessableEntity().WithPayload(payload)
		case services.PreconditionFailedError:
			return orderop.NewCounselingUpdateOrderPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case services.ForbiddenError:
			return orderop.NewCounselingUpdateOrderForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return orderop.NewCounselingUpdateOrderInternalServerError()
		}
	}

	if !session.IsOfficeUser() || !session.Roles.HasRole(roles.RoleTypeServicesCounselor) {
		return handleError(services.NewForbiddenError("is not a Services Counselor"))
	}

	orderID := uuid.FromStringOrNil(params.OrderID.String())
	updatedOrder, moveID, err := h.orderUpdater.UpdateOrderAsCounselor(appCtx, orderID, *params.Body, params.IfMatch)
	if err != nil {
		return handleError(err)
	}

	h.triggerCounselingUpdateOrderEvent(appCtx, orderID, moveID, params)

	orderPayload := payloads.Order(updatedOrder)

	return orderop.NewCounselingUpdateOrderOK().WithPayload(orderPayload)
}

// UpdateAllowanceHandler updates an order and entitlements via PATCH /orders/{orderId}/allowances
type UpdateAllowanceHandler struct {
	handlers.HandlerConfig
	orderUpdater services.OrderUpdater
}

// Handle ... updates an order from a request payload
func (h UpdateAllowanceHandler) Handle(params orderop.UpdateAllowanceParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	appCtx := h.AppContextFromRequest(params.HTTPRequest)
	handleError := func(err error) middleware.Responder {
		logger.Error("error updating order allowance", zap.Error(err))
		switch err.(type) {
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

	if !session.IsOfficeUser() || !session.Roles.HasRole(roles.RoleTypeTOO) {
		return handleError(services.NewForbiddenError("is not a TOO"))
	}

	orderID := uuid.FromStringOrNil(params.OrderID.String())
	updatedOrder, moveID, err := h.orderUpdater.UpdateAllowanceAsTOO(appCtx, orderID, *params.Body, params.IfMatch)
	if err != nil {
		return handleError(err)
	}

	h.triggerUpdatedAllowanceEvent(appCtx, orderID, moveID, params)

	orderPayload := payloads.Order(updatedOrder)

	return orderop.NewUpdateAllowanceOK().WithPayload(orderPayload)
}

// CounselingUpdateAllowanceHandler updates an order and entitlements via PATCH /counseling/orders/{orderId}/allowances
type CounselingUpdateAllowanceHandler struct {
	handlers.HandlerConfig
	orderUpdater services.OrderUpdater
}

// Handle ... updates an order from a request payload
func (h CounselingUpdateAllowanceHandler) Handle(params orderop.CounselingUpdateAllowanceParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	appCtx := h.AppContextFromRequest(params.HTTPRequest)
	handleError := func(err error) middleware.Responder {
		logger.Error("error updating order allowance", zap.Error(err))
		switch err.(type) {
		case services.NotFoundError:
			return orderop.NewCounselingUpdateAllowanceNotFound()
		case services.InvalidInputError:
			payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceID(), validate.NewErrors())
			return orderop.NewCounselingUpdateAllowanceUnprocessableEntity().WithPayload(payload)
		case services.PreconditionFailedError:
			return orderop.NewCounselingUpdateAllowancePreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case services.ForbiddenError:
			return orderop.NewCounselingUpdateAllowanceForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return orderop.NewCounselingUpdateAllowanceInternalServerError()
		}
	}

	if !session.IsOfficeUser() || !session.Roles.HasRole(roles.RoleTypeServicesCounselor) {
		return handleError(services.NewForbiddenError("is not a Services Counselor"))
	}

	orderID := uuid.FromStringOrNil(params.OrderID.String())
	updatedOrder, moveID, err := h.orderUpdater.UpdateAllowanceAsCounselor(appCtx, orderID, *params.Body, params.IfMatch)
	if err != nil {
		return handleError(err)
	}

	h.triggerCounselingUpdateAllowanceEvent(appCtx, orderID, moveID, params)

	orderPayload := payloads.Order(updatedOrder)

	return orderop.NewCounselingUpdateAllowanceOK().WithPayload(orderPayload)
}

func (h UpdateOrderHandler) triggerUpdateOrderEvent(appCtx appcontext.AppContext, orderID uuid.UUID, moveID uuid.UUID, params orderop.UpdateOrderParams) {
	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcUpdateOrderEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.OrderUpdateEventKey, // Event that you want to trigger
		UpdatedObjectID: orderID,                   // ID of the updated logical object
		MtoID:           moveID,                    // ID of the associated Move
		Request:         params.HTTPRequest,        // Pass on the http.Request
		DBConnection:    appCtx.DB(),               // Pass on the pop.Connection
		HandlerContext:  h,                         // Pass on the handlerContext
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.UpdateOrderHandler could not generate the event")
	}
}

func (h CounselingUpdateOrderHandler) triggerCounselingUpdateOrderEvent(appCtx appcontext.AppContext, orderID uuid.UUID, moveID uuid.UUID, params orderop.CounselingUpdateOrderParams) {
	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcCounselingUpdateOrderEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.OrderUpdateEventKey, // Event that you want to trigger
		UpdatedObjectID: orderID,                   // ID of the updated logical object
		MtoID:           moveID,                    // ID of the associated Move
		Request:         params.HTTPRequest,        // Pass on the http.Request
		DBConnection:    appCtx.DB(),               // Pass on the pop.Connection
		HandlerContext:  h,                         // Pass on the handlerContext
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.UpdateAllowanceHandler could not generate the event")
	}
}

func (h UpdateAllowanceHandler) triggerUpdatedAllowanceEvent(appCtx appcontext.AppContext, orderID uuid.UUID, moveID uuid.UUID, params orderop.UpdateAllowanceParams) {
	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcUpdateAllowanceEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.OrderUpdateEventKey, // Event that you want to trigger
		UpdatedObjectID: orderID,                   // ID of the updated logical object
		MtoID:           moveID,                    // ID of the associated Move
		Request:         params.HTTPRequest,        // Pass on the http.Request
		DBConnection:    appCtx.DB(),               // Pass on the pop.Connection
		HandlerContext:  h,                         // Pass on the handlerContext
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.UpdateAllowanceHandler could not generate the event")
	}
}

func (h CounselingUpdateAllowanceHandler) triggerCounselingUpdateAllowanceEvent(appCtx appcontext.AppContext, orderID uuid.UUID, moveID uuid.UUID, params orderop.CounselingUpdateAllowanceParams) {
	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcCounselingUpdateAllowanceEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.OrderUpdateEventKey, // Event that you want to trigger
		UpdatedObjectID: orderID,                   // ID of the updated logical object
		MtoID:           moveID,                    // ID of the associated Move
		Request:         params.HTTPRequest,        // Pass on the http.Request
		DBConnection:    appCtx.DB(),               // Pass on the pop.Connection
		HandlerContext:  h,                         // Pass on the handlerContext
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.CounselingUpdateAllowanceHandler could not generate the event")
	}
}
