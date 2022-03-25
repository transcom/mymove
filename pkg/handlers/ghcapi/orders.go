package ghcapi

import (
	"database/sql"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	orderop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/order"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/event"
)

// GetOrdersHandler fetches the information of a specific order
type GetOrdersHandler struct {
	handlers.HandlerContext
	services.OrderFetcher
}

// Handle getting the information of a specific order
func (h GetOrdersHandler) Handle(params orderop.GetOrderParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			orderID, _ := uuid.FromString(params.OrderID.String())
			order, err := h.FetchOrder(appCtx, orderID)
			if err != nil {
				appCtx.Logger().Error("fetching order", zap.Error(err))
				switch err {
				case sql.ErrNoRows:
					return orderop.NewGetOrderNotFound(), err
				default:
					return orderop.NewGetOrderInternalServerError(), err
				}
			}
			orderPayload := payloads.Order(order)
			return orderop.NewGetOrderOK().WithPayload(orderPayload), nil
		})
}

// UpdateOrderHandler updates an order via PATCH /orders/{orderId}
type UpdateOrderHandler struct {
	handlers.HandlerContext
	orderUpdater services.OrderUpdater
	moveUpdater  services.MoveTaskOrderUpdater
}

// Handle ... updates an order from a request payload
func (h UpdateOrderHandler) Handle(params orderop.UpdateOrderParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {

			handleError := func(err error) middleware.Responder {
				appCtx.Logger().Error("error updating order", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return orderop.NewUpdateOrderNotFound()
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
					return orderop.NewUpdateOrderUnprocessableEntity().WithPayload(payload)
				case apperror.ConflictError:
					return orderop.NewUpdateOrderConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
				case apperror.PreconditionFailedError:
					return orderop.NewUpdateOrderPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
				case apperror.ForbiddenError:
					return orderop.NewUpdateOrderForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
				default:
					return orderop.NewUpdateOrderInternalServerError()
				}
			}

			if !appCtx.Session().IsOfficeUser() || (!appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) && !appCtx.Session().Roles.HasRole(roles.RoleTypeTIO)) {
				return handleError(apperror.NewForbiddenError("is not a TXO"))
			}

			orderID := uuid.FromStringOrNil(params.OrderID.String())
			updatedOrder, moveID, err := h.orderUpdater.UpdateOrderAsTOO(appCtx, orderID, *params.Body, params.IfMatch)
			if err != nil {
				return handleError(err)
			}

			h.triggerUpdateOrderEvent(appCtx, orderID, moveID, params)

			orderPayload := payloads.Order(updatedOrder)

			return orderop.NewUpdateOrderOK().WithPayload(orderPayload)
		})
}

// CounselingUpdateOrderHandler updates an order via PATCH /counseling/orders/{orderId}
type CounselingUpdateOrderHandler struct {
	handlers.HandlerContext
	orderUpdater services.OrderUpdater
}

// Handle ... updates an order as requested by a services counselor
func (h CounselingUpdateOrderHandler) Handle(
	params orderop.CounselingUpdateOrderParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// handleError is a reusable function to deal with multiple errors
			// when it comes to updating orders.
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("error updating order", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return orderop.NewCounselingUpdateOrderNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
					return orderop.NewCounselingUpdateOrderUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return orderop.NewCounselingUpdateOrderPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ForbiddenError:
					return orderop.NewCounselingUpdateOrderForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return orderop.NewCounselingUpdateOrderInternalServerError(), err
				}
			}

			if !appCtx.Session().IsOfficeUser() ||
				!appCtx.Session().Roles.HasRole(roles.RoleTypeServicesCounselor) {
				return handleError(apperror.NewForbiddenError("is not a Services Counselor"))
			}

			orderID := uuid.FromStringOrNil(params.OrderID.String())
			updatedOrder, moveID, err := h.orderUpdater.UpdateOrderAsCounselor(
				appCtx,
				orderID,
				*params.Body,
				params.IfMatch,
			)
			if err != nil {
				return handleError(err)
			}

			h.triggerCounselingUpdateOrderEvent(appCtx, orderID, moveID, params)

			orderPayload := payloads.Order(updatedOrder)

			return orderop.NewCounselingUpdateOrderOK().WithPayload(orderPayload), nil
		})
}

// UpdateAllowanceHandler updates an order and entitlements via PATCH /orders/{orderId}/allowances
type UpdateAllowanceHandler struct {
	handlers.HandlerContext
	orderUpdater services.OrderUpdater
}

// Handle ... updates an order from a request payload
func (h UpdateAllowanceHandler) Handle(params orderop.UpdateAllowanceParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// handleError is a reusable function to deal with multiple errors
			// when it comes to updating orders.
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("error updating order allowance", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return orderop.NewUpdateAllowanceNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
					return orderop.NewUpdateAllowanceUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return orderop.NewUpdateAllowancePreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ForbiddenError:
					return orderop.NewUpdateAllowanceForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return orderop.NewUpdateAllowanceInternalServerError(), err
				}
			}

			if !appCtx.Session().IsOfficeUser() ||
				!appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				return handleError(apperror.NewForbiddenError("is not a TOO"))
			}

			orderID := uuid.FromStringOrNil(params.OrderID.String())
			updatedOrder, moveID, err := h.orderUpdater.UpdateAllowanceAsTOO(
				appCtx,
				orderID,
				*params.Body,
				params.IfMatch,
			)
			if err != nil {
				return handleError(err)
			}

			h.triggerUpdatedAllowanceEvent(appCtx, orderID, moveID, params)

			orderPayload := payloads.Order(updatedOrder)

			return orderop.NewUpdateAllowanceOK().WithPayload(orderPayload), nil
		})
}

// CounselingUpdateAllowanceHandler updates an order and entitlements via PATCH /counseling/orders/{orderId}/allowances
type CounselingUpdateAllowanceHandler struct {
	handlers.HandlerContext
	orderUpdater services.OrderUpdater
}

// Handle ... updates an order from a request payload
func (h CounselingUpdateAllowanceHandler) Handle(
	params orderop.CounselingUpdateAllowanceParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// handleError is a reusable function to deal with multiple errors
			// when it comes to updating orders.
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("error updating order allowance", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return orderop.NewCounselingUpdateAllowanceNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
					return orderop.NewCounselingUpdateAllowanceUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return orderop.NewCounselingUpdateAllowancePreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ForbiddenError:
					return orderop.NewCounselingUpdateAllowanceForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return orderop.NewCounselingUpdateAllowanceInternalServerError(), err
				}
			}

			if !appCtx.Session().IsOfficeUser() ||
				!appCtx.Session().Roles.HasRole(roles.RoleTypeServicesCounselor) {
				return handleError(apperror.NewForbiddenError("is not a Services Counselor"))
			}

			orderID := uuid.FromStringOrNil(params.OrderID.String())
			updatedOrder, moveID, err := h.orderUpdater.UpdateAllowanceAsCounselor(
				appCtx,
				orderID,
				*params.Body,
				params.IfMatch,
			)
			if err != nil {
				return handleError(err)
			}

			h.triggerCounselingUpdateAllowanceEvent(appCtx, orderID, moveID, params)

			orderPayload := payloads.Order(updatedOrder)

			return orderop.NewCounselingUpdateAllowanceOK().WithPayload(orderPayload), nil
		})
}

// UpdateBillableWeightHandler updates the max billable weight on an order's entitlements via PATCH /orders/{orderId}/update-billable-weight
type UpdateBillableWeightHandler struct {
	handlers.HandlerContext
	excessWeightRiskManager services.ExcessWeightRiskManager
}

// Handle ... updates the authorized weight
func (h UpdateBillableWeightHandler) Handle(
	params orderop.UpdateBillableWeightParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// handleError is a reusable function to deal with multiple errors
			// when it comes to updating orders.
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("error updating max billable weight", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return orderop.NewUpdateBillableWeightNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(handlers.ValidationErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
					return orderop.NewUpdateBillableWeightUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return orderop.NewUpdateBillableWeightPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ForbiddenError:
					return orderop.NewUpdateBillableWeightForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return orderop.NewUpdateBillableWeightInternalServerError(), err
				}
			}

			if !appCtx.Session().IsOfficeUser() ||
				!appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				return handleError(apperror.NewForbiddenError("is not a TOO"))
			}

			orderID := uuid.FromStringOrNil(params.OrderID.String())
			dbAuthorizedWeight := swag.Int(int(*params.Body.AuthorizedWeight))
			updatedOrder, moveID, err := h.excessWeightRiskManager.UpdateBillableWeightAsTOO(
				appCtx,
				orderID,
				dbAuthorizedWeight,
				params.IfMatch,
			)
			if err != nil {
				return handleError(err)
			}

			h.triggerUpdatedBillableWeightEvent(appCtx, orderID, moveID, params)

			orderPayload := payloads.Order(updatedOrder)

			return orderop.NewUpdateBillableWeightOK().WithPayload(orderPayload), nil
		})
}

// UpdateMaxBillableWeightAsTIOHandler updates the max billable weight on an order's entitlements via PATCH /orders/{orderId}/update-billable-weight/tio
type UpdateMaxBillableWeightAsTIOHandler struct {
	handlers.HandlerContext
	excessWeightRiskManager services.ExcessWeightRiskManager
}

// Handle ... updates the authorized weight
func (h UpdateMaxBillableWeightAsTIOHandler) Handle(
	params orderop.UpdateMaxBillableWeightAsTIOParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// handleError is a reusable function to deal with multiple errors
			// when it comes to updating orders.
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("error updating max billable weight", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return orderop.NewUpdateMaxBillableWeightAsTIONotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(handlers.ValidationErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
					return orderop.NewUpdateMaxBillableWeightAsTIOUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return orderop.NewUpdateMaxBillableWeightAsTIOPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ForbiddenError:
					return orderop.NewUpdateMaxBillableWeightAsTIOForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return orderop.NewUpdateMaxBillableWeightAsTIOInternalServerError(), err
				}
			}

			if !appCtx.Session().IsOfficeUser() ||
				!appCtx.Session().Roles.HasRole(roles.RoleTypeTIO) {
				return handleError(apperror.NewForbiddenError("is not a TIO"))
			}

			orderID := uuid.FromStringOrNil(params.OrderID.String())
			dbAuthorizedWeight := swag.Int(int(*params.Body.AuthorizedWeight))
			remarks := params.Body.TioRemarks
			updatedOrder, moveID, err := h.excessWeightRiskManager.UpdateMaxBillableWeightAsTIO(
				appCtx,
				orderID,
				dbAuthorizedWeight,
				remarks,
				params.IfMatch,
			)
			if err != nil {
				return handleError(err)
			}

			h.triggerUpdatedMaxBillableWeightAsTIOEvent(appCtx, orderID, moveID, params)

			orderPayload := payloads.Order(updatedOrder)

			return orderop.NewUpdateMaxBillableWeightAsTIOOK().WithPayload(orderPayload), nil
		})
}

// AcknowledgeExcessWeightRiskHandler is called when a TOO dismissed the alert to acknowledge the excess weight risk via POST /orders/{orderId}/acknowledge-excess-weight-risk
type AcknowledgeExcessWeightRiskHandler struct {
	handlers.HandlerContext
	excessWeightRiskManager services.ExcessWeightRiskManager
}

// Handle ... updates the authorized weight
func (h AcknowledgeExcessWeightRiskHandler) Handle(
	params orderop.AcknowledgeExcessWeightRiskParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// handleError is a reusable function to deal with multiple errors
			// when it comes to updating orders.
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("error acknowledging excess weight risk", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return orderop.NewAcknowledgeExcessWeightRiskNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(handlers.ValidationErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
					return orderop.NewAcknowledgeExcessWeightRiskUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return orderop.NewAcknowledgeExcessWeightRiskPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ForbiddenError:
					return orderop.NewAcknowledgeExcessWeightRiskForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return orderop.NewAcknowledgeExcessWeightRiskInternalServerError(), err
				}
			}

			if !appCtx.Session().IsOfficeUser() ||
				!appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				return handleError(apperror.NewForbiddenError("is not a TOO"))
			}

			orderID := uuid.FromStringOrNil(params.OrderID.String())
			updatedMove, err := h.excessWeightRiskManager.AcknowledgeExcessWeightRisk(
				appCtx,
				orderID,
				params.IfMatch,
			)
			if err != nil {
				return handleError(err)
			}

			h.triggerAcknowledgeExcessWeightRiskEvent(appCtx, updatedMove.ID, params)

			movePayload := payloads.Move(updatedMove)

			return orderop.NewAcknowledgeExcessWeightRiskOK().WithPayload(movePayload), nil
		})
}

func (h UpdateOrderHandler) triggerUpdateOrderEvent(appCtx appcontext.AppContext, orderID uuid.UUID, moveID uuid.UUID, params orderop.UpdateOrderParams) {
	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcUpdateOrderEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.OrderUpdateEventKey, // Event that you want to trigger
		UpdatedObjectID: orderID,                   // ID of the updated logical object
		MtoID:           moveID,                    // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
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
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
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
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
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
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.CounselingUpdateAllowanceHandler could not generate the event")
	}
}

func (h UpdateBillableWeightHandler) triggerUpdatedBillableWeightEvent(appCtx appcontext.AppContext, orderID uuid.UUID, moveID uuid.UUID, params orderop.UpdateBillableWeightParams) {
	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcUpdateBillableWeightEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.OrderUpdateEventKey, // Event that you want to trigger
		UpdatedObjectID: orderID,                   // ID of the updated logical object
		MtoID:           moveID,                    // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.UpdateBillableWeightHandler could not generate the event")
	}
}

func (h UpdateMaxBillableWeightAsTIOHandler) triggerUpdatedMaxBillableWeightAsTIOEvent(appCtx appcontext.AppContext, orderID uuid.UUID, moveID uuid.UUID, params orderop.UpdateMaxBillableWeightAsTIOParams) {
	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcUpdateMaxBillableWeightAsTIOEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.OrderUpdateEventKey, // Event that you want to trigger
		UpdatedObjectID: orderID,                   // ID of the updated logical object
		MtoID:           moveID,                    // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.UpdateMaxBillableWeightAsTIOHandler could not generate the event")
	}
}

func (h AcknowledgeExcessWeightRiskHandler) triggerAcknowledgeExcessWeightRiskEvent(appCtx appcontext.AppContext, moveID uuid.UUID, params orderop.AcknowledgeExcessWeightRiskParams) {
	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcAcknowledgeExcessWeightRiskEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.MoveTaskOrderUpdateEventKey, // Event that you want to trigger
		UpdatedObjectID: moveID,                            // ID of the updated logical object
		MtoID:           moveID,                            // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.UpdateBillableWeightHandler could not generate the event")
	}
}
