package ghcapi

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/notifications"

	"github.com/transcom/mymove/pkg/models/roles"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"go.uber.org/zap"

	"github.com/gofrs/uuid"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_shipment"
	shipmentops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/shipment"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/event"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
)

// ListMTOShipmentsHandler returns a list of MTO Shipments
type ListMTOShipmentsHandler struct {
	handlers.HandlerContext
	services.MTOShipmentFetcher
	services.ShipmentSITStatus
}

// Handle listing mto shipments for the move task order
func (h ListMTOShipmentsHandler) Handle(params mtoshipmentops.ListMTOShipmentsParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {
			handleError := func(err error) middleware.Responder {
				appCtx.Logger().Error("ListMTOShipmentsHandler error", zap.Error(err))
				payload := &ghcmessages.Error{Message: handlers.FmtString(err.Error())}
				switch err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewListMTOShipmentsNotFound().WithPayload(payload)
				case apperror.ForbiddenError:
					return mtoshipmentops.NewListMTOShipmentsForbidden().WithPayload(payload)
				case apperror.QueryError:
					return mtoshipmentops.NewListMTOShipmentsInternalServerError()
				default:
					return mtoshipmentops.NewListMTOShipmentsInternalServerError()
				}
			}

			if !appCtx.Session().IsOfficeUser() || (!appCtx.Session().Roles.HasRole(roles.RoleTypeServicesCounselor) && !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) && !appCtx.Session().Roles.HasRole(roles.RoleTypeTIO)) {
				handleError(apperror.NewForbiddenError("user is not an office user or does not have a SC, TOO, or TIO role"))
			}

			moveID := uuid.FromStringOrNil(params.MoveTaskOrderID.String())

			shipments, err := h.ListMTOShipments(appCtx, moveID)
			if err != nil {
				return handleError(err)
			}
			mtoShipments := models.MTOShipments(shipments)

			shipmentSITStatuses := h.CalculateShipmentsSITStatuses(appCtx, shipments)

			sitStatusPayload := payloads.SITStatuses(shipmentSITStatuses)
			payload := payloads.MTOShipments(&mtoShipments, sitStatusPayload)
			return mtoshipmentops.NewListMTOShipmentsOK().WithPayload(*payload)
		})
}

// CreateMTOShipmentHandler is the handler to create MTO shipments
type CreateMTOShipmentHandler struct {
	handlers.HandlerContext
	mtoShipmentCreator services.MTOShipmentCreator
	shipmentStatus     services.ShipmentSITStatus
}

// Handle creates the mto shipment
func (h CreateMTOShipmentHandler) Handle(params mtoshipmentops.CreateMTOShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {
			payload := params.Body

			if payload == nil {
				appCtx.Logger().Error("Invalid mto shipment: params Body is nil")
				return mtoshipmentops.NewCreateMTOShipmentBadRequest()
			}

			handleError := func(err error) middleware.Responder {
				appCtx.Logger().Error("ghcapi.CreateMTOShipmentHandler error", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					payload := ghcmessages.Error{
						Message: handlers.FmtString(err.Error()),
					}
					return mtoshipmentops.NewCreateMTOShipmentNotFound().WithPayload(&payload)
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Validation errors", "CreateMTOShipment", h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
					return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payload)
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("ghcapi.CreateMTOShipmentHandler query error", zap.Error(e.Unwrap()))
					}
					return mtoshipmentops.NewCreateMTOShipmentInternalServerError()
				default:
					return mtoshipmentops.NewCreateMTOShipmentInternalServerError()
				}
			}

			mtoShipment := payloads.MTOShipmentModelFromCreate(payload)
			mtoShipment, err := h.mtoShipmentCreator.CreateMTOShipment(appCtx, mtoShipment, nil)

			if err != nil {
				return handleError(err)
			}

			if mtoShipment == nil {
				appCtx.Logger().Error("Unexpected nil shipment from CreateMTOShipment")
				return mtoshipmentops.NewCreateMTOShipmentInternalServerError()
			}

			sitAllowance, err := h.shipmentStatus.CalculateShipmentSITAllowance(appCtx, *mtoShipment)
			if err != nil {
				return handleError(err)
			}

			mtoShipment.SITDaysAllowance = &sitAllowance

			returnPayload := payloads.MTOShipment(mtoShipment, nil)
			return mtoshipmentops.NewCreateMTOShipmentOK().WithPayload(returnPayload)
		})
}

// UpdateShipmentHandler updates shipments
type UpdateShipmentHandler struct {
	handlers.HandlerContext
	services.Fetcher
	services.MTOShipmentUpdater
	services.ShipmentSITStatus
}

// Handle updates shipments
func (h UpdateShipmentHandler) Handle(params mtoshipmentops.UpdateMTOShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {

			payload := params.Body
			if payload == nil {
				appCtx.Logger().Error("Invalid mto shipment: params Body is nil")

				payload := payloadForValidationError(
					"Empty body error",
					"The MTO Shipment request body cannot be empty.",
					h.GetTraceIDFromRequest(params.HTTPRequest),
					validate.NewErrors(),
				)

				return mtoshipmentops.NewUpdateMTOShipmentUnprocessableEntity().WithPayload(payload)
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			oldShipment, err := h.MTOShipmentUpdater.RetrieveMTOShipment(appCtx, shipmentID)

			if err != nil {
				appCtx.Logger().Error("ghcapi.UpdateShipmentHandler", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewUpdateMTOShipmentNotFound()
				default:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))

					return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(
						&ghcmessages.Error{Message: &msg},
					)
				}
			}

			updateable, err := h.MTOShipmentUpdater.CheckIfMTOShipmentCanBeUpdated(appCtx, oldShipment, appCtx.Session())

			if err != nil {
				appCtx.Logger().Error("ghcapi.UpdateShipmentHandler", zap.Error(err))
				msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
				return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(
					&ghcmessages.Error{Message: &msg},
				)
			}

			if !updateable {
				msg := fmt.Sprintf("%v is not updatable", shipmentID)
				return mtoshipmentops.NewUpdateMTOShipmentPreconditionFailed().WithPayload(
					&ghcmessages.Error{Message: &msg},
				)
			}

			mtoShipment := payloads.MTOShipmentModelFromUpdate(payload)

			//MTOShipmentModelFromUpdate defaults UsesExternalVendor to false if it's nil in the payload
			if payload.UsesExternalVendor == nil {
				mtoShipment.UsesExternalVendor = oldShipment.UsesExternalVendor
			}
			// booleans not passed will update to false
			mtoShipment.Diversion = oldShipment.Diversion

			mtoShipment.ID = shipmentID

			handleError := func(err error) middleware.Responder {
				appCtx.Logger().Error("ghcapi.UpdateShipmentHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewUpdateMTOShipmentNotFound()
				case apperror.InvalidInputError:
					return mtoshipmentops.NewUpdateMTOShipmentUnprocessableEntity().WithPayload(
						payloadForValidationError(
							handlers.ValidationErrMessage,
							err.Error(),
							h.GetTraceIDFromRequest(params.HTTPRequest),
							e.ValidationErrors,
						),
					)
				case apperror.PreconditionFailedError:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
					return mtoshipmentops.NewUpdateMTOShipmentPreconditionFailed().WithPayload(
						&ghcmessages.Error{Message: &msg},
					)
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("ghcapi.UpdateShipmentHandler error", zap.Error(e.Unwrap()))
					}

					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))

					return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(
						&ghcmessages.Error{Message: &msg},
					)
				default:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))

					return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(
						&ghcmessages.Error{Message: &msg},
					)
				}
			}
			updatedMtoShipment, err := h.MTOShipmentUpdater.UpdateMTOShipmentOffice(appCtx, mtoShipment, params.IfMatch)
			if err != nil {
				return handleError(err)
			}

			_, err = event.TriggerEvent(event.Event{
				EndpointKey: event.GhcUpdateMTOShipmentEndpointKey,
				// Endpoint that is being handled
				EventKey:        event.MTOShipmentUpdateEventKey,    // Event that you want to trigger
				UpdatedObjectID: updatedMtoShipment.ID,              // ID of the updated logical object
				MtoID:           updatedMtoShipment.MoveTaskOrderID, // ID of the associated Move
				AppContext:      appCtx,
				TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
			})
			// If the event trigger fails, just log the error.
			if err != nil {
				appCtx.Logger().Error("ghcapi.UpdateMTOShipment could not generate the event")
			}

			shipmentSITStatus, err := h.CalculateShipmentSITStatus(appCtx, *updatedMtoShipment)
			if err != nil {
				return handleError(err)
			}
			sitStatusPayload := payloads.SITStatus(shipmentSITStatus)

			returnPayload := payloads.MTOShipment(updatedMtoShipment, sitStatusPayload)
			return mtoshipmentops.NewUpdateMTOShipmentOK().WithPayload(returnPayload)
		})
}

// DeleteShipmentHandler soft deletes a shipment
type DeleteShipmentHandler struct {
	handlers.HandlerContext
	services.ShipmentDeleter
}

// Handle soft deletes a shipment
func (h DeleteShipmentHandler) Handle(params shipmentops.DeleteShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeServicesCounselor) {
				appCtx.Logger().Error("user is not authenticated with service counselor office role")
				return shipmentops.NewDeleteShipmentForbidden()
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			moveID, err := h.DeleteShipment(appCtx, shipmentID)
			if err != nil {
				appCtx.Logger().Error("ghcapi.DeleteShipmentHandler", zap.Error(err))

				switch err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewDeleteShipmentNotFound()
				case apperror.ForbiddenError:
					return shipmentops.NewDeleteShipmentForbidden()
				default:
					return shipmentops.NewDeleteShipmentInternalServerError()
				}
			}

			// Note that this is currently not sending any notifications because
			// the move isn't available to the Prime yet. See the objectEventHandler
			// function in pkg/services/event/notification.go.
			// We added this now because eventually, we will want to save events in
			// the DB for auditing purposes. When that happens, this code in the handler
			// will not change. However, we should make sure to add a test in
			// mto_shipment_test.go that verifies the audit got saved.
			h.triggerShipmentDeletionEvent(appCtx, shipmentID, moveID, params)

			return shipmentops.NewDeleteShipmentNoContent()
		})
}

func (h DeleteShipmentHandler) triggerShipmentDeletionEvent(appCtx appcontext.AppContext, shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.DeleteShipmentParams) {

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcDeleteShipmentEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ShipmentDeleteEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                   // ID of the updated logical object
		MtoID:           moveID,                       // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.DeleteShipmentHandler could not generate the event", zap.Error(err))
	}
}

// ApproveShipmentHandler approves a shipment
type ApproveShipmentHandler struct {
	handlers.HandlerContext
	services.ShipmentApprover
	services.ShipmentSITStatus
}

// Handle approves a shipment
func (h ApproveShipmentHandler) Handle(params shipmentops.ApproveShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				appCtx.Logger().Error("Only TOO role can approve shipments")
				return shipmentops.NewApproveShipmentForbidden()
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			eTag := params.IfMatch

			handleError := func(err error) middleware.Responder {
				appCtx.Logger().Error("ghcapi.ApproveShipmentHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewApproveShipmentNotFound()
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Validation errors", "ApproveShipment", h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
					return shipmentops.NewApproveShipmentUnprocessableEntity().WithPayload(payload)
				case apperror.PreconditionFailedError:
					return shipmentops.NewApproveShipmentPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
				case apperror.ConflictError, mtoshipment.ConflictStatusError:
					return shipmentops.NewApproveShipmentConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
				default:
					return shipmentops.NewApproveShipmentInternalServerError()
				}
			}

			shipment, err := h.ApproveShipment(appCtx, shipmentID, eTag)
			if err != nil {
				return handleError(err)
			}

			h.triggerShipmentApprovalEvent(appCtx, shipmentID, shipment.MoveTaskOrderID, params)

			shipmentSITStatus, err := h.CalculateShipmentSITStatus(appCtx, *shipment)
			if err != nil {
				return handleError(err)
			}
			sitStatusPayload := payloads.SITStatus(shipmentSITStatus)

			payload := payloads.MTOShipment(shipment, sitStatusPayload)
			return shipmentops.NewApproveShipmentOK().WithPayload(payload)
		})
}

func (h ApproveShipmentHandler) triggerShipmentApprovalEvent(appCtx appcontext.AppContext, shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.ApproveShipmentParams) {

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcApproveShipmentEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ShipmentApproveEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                    // ID of the updated logical object
		MtoID:           moveID,                        // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.ApproveShipmentHandler could not generate the event", zap.Error(err))
	}
}

// RequestShipmentDiversionHandler Requests a shipment diversion
type RequestShipmentDiversionHandler struct {
	handlers.HandlerContext
	services.ShipmentDiversionRequester
	services.ShipmentSITStatus
}

// Handle Requests a shipment diversion
func (h RequestShipmentDiversionHandler) Handle(params shipmentops.RequestShipmentDiversionParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				appCtx.Logger().Error("Only TOO role can Request shipment diversions")
				return shipmentops.NewRequestShipmentDiversionForbidden()
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			eTag := params.IfMatch

			handleError := func(err error) middleware.Responder {
				appCtx.Logger().Error("ghcapi.RequestShipmentDiversionHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewRequestShipmentDiversionNotFound()
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Validation errors", "RequestShipmentDiversion", h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
					return shipmentops.NewRequestShipmentDiversionUnprocessableEntity().WithPayload(payload)
				case apperror.PreconditionFailedError:
					return shipmentops.NewRequestShipmentDiversionPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
				case mtoshipment.ConflictStatusError:
					return shipmentops.NewRequestShipmentDiversionConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
				default:
					return shipmentops.NewRequestShipmentDiversionInternalServerError()
				}
			}

			shipment, err := h.RequestShipmentDiversion(appCtx, shipmentID, eTag)
			if err != nil {
				return handleError(err)
			}

			h.triggerRequestShipmentDiversionEvent(appCtx, shipmentID, shipment.MoveTaskOrderID, params)

			shipmentSITStatus, err := h.CalculateShipmentSITStatus(appCtx, *shipment)
			if err != nil {
				return handleError(err)
			}
			sitStatusPayload := payloads.SITStatus(shipmentSITStatus)

			payload := payloads.MTOShipment(shipment, sitStatusPayload)
			return shipmentops.NewRequestShipmentDiversionOK().WithPayload(payload)
		})
}

func (h RequestShipmentDiversionHandler) triggerRequestShipmentDiversionEvent(appCtx appcontext.AppContext, shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.RequestShipmentDiversionParams) {

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcRequestShipmentDiversionEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ShipmentRequestDiversionEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                             // ID of the updated logical object
		MtoID:           moveID,                                 // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.RequestShipmentDiversionHandler could not generate the event", zap.Error(err))
	}
}

// ApproveShipmentDiversionHandler approves a shipment diversion
type ApproveShipmentDiversionHandler struct {
	handlers.HandlerContext
	services.ShipmentDiversionApprover
	services.ShipmentSITStatus
}

// Handle approves a shipment diversion
func (h ApproveShipmentDiversionHandler) Handle(params shipmentops.ApproveShipmentDiversionParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				appCtx.Logger().Error("Only TOO role can approve shipment diversions")
				return shipmentops.NewApproveShipmentDiversionForbidden()
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			eTag := params.IfMatch

			handleError := func(err error) middleware.Responder {
				appCtx.Logger().Error("ghcapi.ApproveShipmentDiversionHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewApproveShipmentDiversionNotFound()
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Validation errors", "ApproveShipmentDiversion", h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
					return shipmentops.NewApproveShipmentDiversionUnprocessableEntity().WithPayload(payload)
				case apperror.PreconditionFailedError:
					return shipmentops.NewApproveShipmentDiversionPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
				case mtoshipment.ConflictStatusError:
					return shipmentops.NewApproveShipmentDiversionConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
				default:
					return shipmentops.NewApproveShipmentDiversionInternalServerError()
				}
			}

			shipment, err := h.ApproveShipmentDiversion(appCtx, shipmentID, eTag)
			if err != nil {
				return handleError(err)
			}

			h.triggerShipmentDiversionApprovalEvent(appCtx, shipmentID, shipment.MoveTaskOrderID, params)

			shipmentSITStatus, err := h.CalculateShipmentSITStatus(appCtx, *shipment)
			if err != nil {
				return handleError(err)
			}
			sitStatusPayload := payloads.SITStatus(shipmentSITStatus)

			payload := payloads.MTOShipment(shipment, sitStatusPayload)
			return shipmentops.NewApproveShipmentDiversionOK().WithPayload(payload)
		})
}

func (h ApproveShipmentDiversionHandler) triggerShipmentDiversionApprovalEvent(appCtx appcontext.AppContext, shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.ApproveShipmentDiversionParams) {

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcApproveShipmentDiversionEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ShipmentApproveDiversionEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                             // ID of the updated logical object
		MtoID:           moveID,                                 // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.ApproveShipmentDiversionHandler could not generate the event", zap.Error(err))
	}
}

// RejectShipmentHandler rejects a shipment
type RejectShipmentHandler struct {
	handlers.HandlerContext
	services.ShipmentRejecter
}

// Handle rejects a shipment
func (h RejectShipmentHandler) Handle(params shipmentops.RejectShipmentParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
		appCtx.Logger().Error("Only TOO role can reject shipments")
		return shipmentops.NewRejectShipmentForbidden()
	}

	shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
	eTag := params.IfMatch
	rejectionReason := params.Body.RejectionReason
	shipment, err := h.RejectShipment(appCtx, shipmentID, eTag, rejectionReason)

	if err != nil {
		appCtx.Logger().Error("ghcapi.RejectShipmentHandler", zap.Error(err))

		switch e := err.(type) {
		case apperror.NotFoundError:
			return shipmentops.NewRejectShipmentNotFound()
		case apperror.InvalidInputError:
			payload := payloadForValidationError("Validation errors", "RejectShipment", h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
			return shipmentops.NewRejectShipmentUnprocessableEntity().WithPayload(payload)
		case apperror.PreconditionFailedError:
			return shipmentops.NewRejectShipmentPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case mtoshipment.ConflictStatusError:
			return shipmentops.NewRejectShipmentConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return shipmentops.NewRejectShipmentInternalServerError()
		}
	}

	h.triggerShipmentRejectionEvent(appCtx, shipmentID, shipment.MoveTaskOrderID, params)

	payload := payloads.MTOShipment(shipment, nil)
	return shipmentops.NewRejectShipmentOK().WithPayload(payload)
}

func (h RejectShipmentHandler) triggerShipmentRejectionEvent(appCtx appcontext.AppContext, shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.RejectShipmentParams) {

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcRejectShipmentEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ShipmentRejectEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                   // ID of the updated logical object
		MtoID:           moveID,                       // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.RejectShipmentHandler could not generate the event", zap.Error(err))
	}
}

// RequestShipmentCancellationHandler Requests a shipment diversion
type RequestShipmentCancellationHandler struct {
	handlers.HandlerContext
	services.ShipmentCancellationRequester
	services.ShipmentSITStatus
}

// Handle Requests a shipment diversion
func (h RequestShipmentCancellationHandler) Handle(params shipmentops.RequestShipmentCancellationParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
		appCtx.Logger().Error("Only TOO role can Request shipment diversions")
		return shipmentops.NewRequestShipmentCancellationForbidden()
	}

	shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
	eTag := params.IfMatch

	handleError := func(err error) middleware.Responder {
		appCtx.Logger().Error("ghcapi.RequestShipmentCancellationHandler", zap.Error(err))

		switch e := err.(type) {
		case apperror.NotFoundError:
			return shipmentops.NewRequestShipmentCancellationNotFound()
		case apperror.InvalidInputError:
			payload := payloadForValidationError("Validation errors", "RequestShipmentCancellation", h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
			return shipmentops.NewRequestShipmentCancellationUnprocessableEntity().WithPayload(payload)
		case apperror.PreconditionFailedError:
			return shipmentops.NewRequestShipmentCancellationPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case mtoshipment.ConflictStatusError:
			return shipmentops.NewRequestShipmentCancellationConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return shipmentops.NewRequestShipmentCancellationInternalServerError()
		}
	}

	shipment, err := h.RequestShipmentCancellation(appCtx, shipmentID, eTag)
	if err != nil {
		return handleError(err)
	}

	h.triggerRequestShipmentCancellationEvent(appCtx, shipmentID, shipment.MoveTaskOrderID, params)

	shipmentSITStatus, err := h.CalculateShipmentSITStatus(appCtx, *shipment)
	if err != nil {
		return handleError(err)
	}
	sitStatusPayload := payloads.SITStatus(shipmentSITStatus)

	payload := payloads.MTOShipment(shipment, sitStatusPayload)
	return shipmentops.NewRequestShipmentCancellationOK().WithPayload(payload)
}

func (h RequestShipmentCancellationHandler) triggerRequestShipmentCancellationEvent(appCtx appcontext.AppContext, shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.RequestShipmentCancellationParams) {

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcRequestShipmentCancellationEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ShipmentRequestCancellationEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                                // ID of the updated logical object
		MtoID:           moveID,                                    // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.RequestShipmentCancellationHandler could not generate the event", zap.Error(err))
	}
}

// RequestShipmentReweighHandler Requests a shipment reweigh
type RequestShipmentReweighHandler struct {
	handlers.HandlerContext
	services.ShipmentReweighRequester
	services.ShipmentSITStatus
	services.MTOShipmentUpdater
}

// Handle Requests a shipment reweigh
func (h RequestShipmentReweighHandler) Handle(params shipmentops.RequestShipmentReweighParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
		appCtx.Logger().Error("Only TOO role can Request a shipment reweigh")
		return shipmentops.NewRequestShipmentReweighForbidden()
	}

	shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
	reweigh, err := h.RequestShipmentReweigh(appCtx, shipmentID, models.ReweighRequesterTOO)

	if err != nil {
		appCtx.Logger().Error("ghcapi.RequestShipmentReweighHandler", zap.Error(err))

		switch e := err.(type) {
		case apperror.NotFoundError:
			return shipmentops.NewRequestShipmentReweighNotFound()
		case apperror.InvalidInputError:
			payload := payloadForValidationError("Validation errors", "RequestShipmentReweigh", h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
			return shipmentops.NewRequestShipmentReweighUnprocessableEntity().WithPayload(payload)
		case apperror.ConflictError:
			return shipmentops.NewRequestShipmentReweighConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return shipmentops.NewRequestShipmentReweighInternalServerError()
		}
	}

	handleError := func(err error) middleware.Responder {
		appCtx.Logger().Error("ghcapi.UpdateShipmentHandler", zap.Error(err))

		switch err.(type) {
		case apperror.NotFoundError:
			return mtoshipmentops.NewUpdateMTOShipmentNotFound()
		default:
			msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))

			return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(
				&ghcmessages.Error{Message: &msg},
			)
		}
	}

	shipment, err := h.MTOShipmentUpdater.RetrieveMTOShipment(appCtx, shipmentID)
	if err != nil {
		return handleError(err)
	}

	moveID := shipment.MoveTaskOrderID
	h.triggerRequestShipmentReweighEvent(appCtx, shipmentID, moveID, params)

	err = h.NotificationSender().SendNotification(appCtx,
		notifications.NewReweighRequested(moveID, *shipment),
	)
	if err != nil {
		appCtx.Logger().Error("problem sending email to user", zap.Error(err))
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	shipmentSITStatus, err := h.CalculateShipmentSITStatus(appCtx, reweigh.Shipment)
	if err != nil {
		return handleError(err)
	}
	sitStatusPayload := payloads.SITStatus(shipmentSITStatus)

	payload := payloads.Reweigh(reweigh, sitStatusPayload)
	return shipmentops.NewRequestShipmentReweighOK().WithPayload(payload)
}

func (h RequestShipmentReweighHandler) triggerRequestShipmentReweighEvent(appCtx appcontext.AppContext, shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.RequestShipmentReweighParams) {

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcRequestShipmentReweighEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ShipmentRequestReweighEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                           // ID of the updated logical object
		MtoID:           moveID,                               // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.RequestShipmentReweighHandler could not generate the event", zap.Error(err))
	}
}

// ApproveSITExtensionHandler approves a SIT extension
type ApproveSITExtensionHandler struct {
	handlers.HandlerContext
	services.SITExtensionApprover
	services.ShipmentSITStatus
}

// Handle ... approves the SIT extension
func (h ApproveSITExtensionHandler) Handle(params shipmentops.ApproveSITExtensionParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)
	handleError := func(err error) middleware.Responder {
		appCtx.Logger().Error("error approving SIT extension", zap.Error(err))
		switch e := err.(type) {
		case apperror.NotFoundError:
			return shipmentops.NewApproveSITExtensionNotFound()
		case apperror.InvalidInputError:
			payload := payloadForValidationError(handlers.ValidationErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
			return shipmentops.NewApproveSITExtensionUnprocessableEntity().WithPayload(payload)
		case apperror.PreconditionFailedError:
			return shipmentops.NewApproveSITExtensionPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case apperror.ForbiddenError:
			return shipmentops.NewApproveSITExtensionForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return shipmentops.NewApproveSITExtensionInternalServerError()
		}
	}

	if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
		return handleError(apperror.NewForbiddenError("is not a TOO"))
	}

	shipmentID := uuid.FromStringOrNil(string(params.ShipmentID))
	sitExtensionID := uuid.FromStringOrNil(string(params.SitExtensionID))
	approvedDays := int(*params.Body.ApprovedDays)
	officeRemarks := params.Body.OfficeRemarks
	updatedShipment, err := h.SITExtensionApprover.ApproveSITExtension(appCtx, shipmentID, sitExtensionID, approvedDays, officeRemarks, params.IfMatch)
	if err != nil {
		return handleError(err)
	}

	shipmentSITStatus, err := h.CalculateShipmentSITStatus(appCtx, *updatedShipment)
	if err != nil {
		return handleError(err)
	}
	sitStatusPayload := payloads.SITStatus(shipmentSITStatus)

	shipmentPayload := payloads.MTOShipment(updatedShipment, sitStatusPayload)

	h.triggerApproveSITExtensionEvent(appCtx, shipmentID, updatedShipment.MoveTaskOrderID, params)
	return shipmentops.NewApproveSITExtensionOK().WithPayload(shipmentPayload)
}

func (h ApproveSITExtensionHandler) triggerApproveSITExtensionEvent(appCtx appcontext.AppContext, shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.ApproveSITExtensionParams) {

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcApproveSITExtensionEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ApproveSITExtensionEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                        // ID of the updated logical object
		MtoID:           moveID,                            // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.ApproveSITExtensionHandler could not generate the event", zap.Error(err))
	}
}

// DenySITExtensionHandler denies a SIT extension
type DenySITExtensionHandler struct {
	handlers.HandlerContext
	services.SITExtensionDenier
	services.ShipmentSITStatus
}

// Handle ... denies the SIT extension
func (h DenySITExtensionHandler) Handle(params shipmentops.DenySITExtensionParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)
	handleError := func(err error) middleware.Responder {
		appCtx.Logger().Error("error denying SIT extension", zap.Error(err))
		switch e := err.(type) {
		case apperror.NotFoundError:
			return shipmentops.NewDenySITExtensionNotFound()
		case apperror.InvalidInputError:
			payload := payloadForValidationError(handlers.ValidationErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
			return shipmentops.NewDenySITExtensionUnprocessableEntity().WithPayload(payload)
		case apperror.PreconditionFailedError:
			return shipmentops.NewDenySITExtensionPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case apperror.ForbiddenError:
			return shipmentops.NewDenySITExtensionForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return shipmentops.NewDenySITExtensionInternalServerError()
		}
	}

	if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
		return handleError(apperror.NewForbiddenError("is not a TOO"))
	}

	shipmentID := uuid.FromStringOrNil(string(params.ShipmentID))
	sitExtensionID := uuid.FromStringOrNil(string(params.SitExtensionID))
	officeRemarks := params.Body.OfficeRemarks
	updatedShipment, err := h.SITExtensionDenier.DenySITExtension(appCtx, shipmentID, sitExtensionID, officeRemarks, params.IfMatch)
	if err != nil {
		return handleError(err)
	}

	shipmentSITStatus, err := h.CalculateShipmentSITStatus(appCtx, *updatedShipment)
	if err != nil {
		return handleError(err)
	}

	sitStatusPayload := payloads.SITStatus(shipmentSITStatus)
	shipmentPayload := payloads.MTOShipment(updatedShipment, sitStatusPayload)

	h.triggerDenySITExtensionEvent(appCtx, shipmentID, updatedShipment.MoveTaskOrderID, params)

	return shipmentops.NewDenySITExtensionOK().WithPayload(shipmentPayload)
}

func (h DenySITExtensionHandler) triggerDenySITExtensionEvent(appCtx appcontext.AppContext, shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.DenySITExtensionParams) {

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcDenySITExtensionEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.DenySITExtensionEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                     // ID of the updated logical object
		MtoID:           moveID,                         // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.DenySITExtensionHandler could not generate the event", zap.Error(err))
	}
}

// CreateSITExtensionAsTOOHandler creates a SIT extension in the approved state
type CreateSITExtensionAsTOOHandler struct {
	handlers.HandlerContext
	services.SITExtensionCreatorAsTOO
	services.ShipmentSITStatus
}

// Handle creates the approved SIT extension
func (h CreateSITExtensionAsTOOHandler) Handle(params shipmentops.CreateSITExtensionAsTOOParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)
	payload := params.Body
	shipmentID := params.ShipmentID

	handleError := func(err error) middleware.Responder {
		appCtx.Logger().Error("ghcapi.CreateApprovedSITExtension error", zap.Error(err))
		switch e := err.(type) {
		case apperror.NotFoundError:
			payload := ghcmessages.Error{
				Message: handlers.FmtString(err.Error()),
			}
			return shipmentops.NewCreateSITExtensionAsTOONotFound().WithPayload(&payload)
		case apperror.InvalidInputError:
			payload := payloadForValidationError("Validation errors", "CreateApprovedSITExtension", h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
			return shipmentops.NewCreateSITExtensionAsTOOUnprocessableEntity().WithPayload(payload)
		case apperror.PreconditionFailedError:
			return shipmentops.NewDenySITExtensionPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case apperror.QueryError:
			if e.Unwrap() != nil {
				// If you can unwrap, log the internal error (usually a pq error) for better debugging
				appCtx.Logger().Error("ghcapi.CreateApprovedSITExtension query error", zap.Error(e.Unwrap()))
			}
			return shipmentops.NewCreateSITExtensionAsTOOInternalServerError()
		case apperror.ForbiddenError:
			return shipmentops.NewCreateSITExtensionAsTOOForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return shipmentops.NewCreateSITExtensionAsTOOInternalServerError()
		}
	}

	sitExtension := payloads.ApprovedSITExtensionFromCreate(payload, shipmentID)
	shipment, err := h.SITExtensionCreatorAsTOO.CreateSITExtensionAsTOO(appCtx, sitExtension, sitExtension.MTOShipmentID, params.IfMatch)
	if err != nil {
		return handleError(err)
	}

	if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
		return handleError(apperror.NewForbiddenError("is not a TOO"))
	}

	shipmentSITStatus, err := h.CalculateShipmentSITStatus(appCtx, *shipment)
	if err != nil {
		return handleError(err)
	}

	sitStatusPayload := payloads.SITStatus(shipmentSITStatus)
	returnPayload := payloads.MTOShipment(shipment, sitStatusPayload)
	return shipmentops.NewCreateSITExtensionAsTOOOK().WithPayload(returnPayload)
}
