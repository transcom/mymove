package ghcapi

import (
	"fmt"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/notifications"

	"github.com/transcom/mymove/pkg/models/roles"

	"github.com/go-openapi/runtime/middleware"
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
	handlers.HandlerConfig
	services.MTOShipmentFetcher
	services.ShipmentSITStatus
}

// Handle listing mto shipments for the move task order
func (h ListMTOShipmentsHandler) Handle(params mtoshipmentops.ListMTOShipmentsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ListMTOShipmentsHandler error", zap.Error(err))
				payload := &ghcmessages.Error{Message: handlers.FmtString(err.Error())}
				switch err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewListMTOShipmentsNotFound().WithPayload(payload), err
				case apperror.ForbiddenError:
					return mtoshipmentops.NewListMTOShipmentsForbidden().WithPayload(payload), err
				case apperror.QueryError:
					return mtoshipmentops.NewListMTOShipmentsInternalServerError(), err
				default:
					return mtoshipmentops.NewListMTOShipmentsInternalServerError(), err
				}
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
			return mtoshipmentops.NewListMTOShipmentsOK().WithPayload(*payload), nil
		})
}

// GetMTOShipmentHandler is the handler to fetch a single MTO shipment by ID
type GetMTOShipmentHandler struct {
	handlers.HandlerConfig
	mtoShipmentFetcher services.MTOShipmentFetcher
}

// Handle handles the handling of fetching a single MTO shipment by ID.
func (h GetMTOShipmentHandler) Handle(params mtoshipmentops.GetShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("GetShipment error", zap.Error(err))
				payload := &ghcmessages.Error{Message: handlers.FmtString(err.Error())}
				switch err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewGetShipmentNotFound().WithPayload(payload), err
				case apperror.ForbiddenError:
					return mtoshipmentops.NewGetShipmentForbidden().WithPayload(payload), err
				case apperror.QueryError:
					return mtoshipmentops.NewGetShipmentInternalServerError(), err
				default:
					return mtoshipmentops.NewGetShipmentInternalServerError(), err
				}
			}

			eagerAssociations := []string{"MoveTaskOrder",
				"PickupAddress",
				"DestinationAddress",
				"SecondaryPickupAddress",
				"SecondaryDeliveryAddress",
				"MTOAgents",
				"MTOServiceItems.CustomerContacts",
				"StorageFacility.Address"}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())

			mtoShipment, err := h.mtoShipmentFetcher.GetShipment(appCtx, shipmentID, eagerAssociations...)
			if err != nil {
				return handleError(err)
			}
			payload := payloads.MTOShipment(mtoShipment, nil)
			return mtoshipmentops.NewGetShipmentOK().WithPayload(payload), nil
		})
}

// CreateMTOShipmentHandler is the handler to create MTO shipments
type CreateMTOShipmentHandler struct {
	handlers.HandlerConfig
	shipmentCreator services.ShipmentCreator
	shipmentStatus  services.ShipmentSITStatus
}

// Handle creates the mto shipment
func (h CreateMTOShipmentHandler) Handle(params mtoshipmentops.CreateMTOShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			payload := params.Body

			if payload == nil {
				invalidShipmentError := apperror.NewBadDataError("Invalid mto shipment: params Body is nil")
				appCtx.Logger().Error(invalidShipmentError.Error())
				return mtoshipmentops.NewCreateMTOShipmentBadRequest(), invalidShipmentError
			}

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.CreateMTOShipmentHandler error", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					payload := ghcmessages.Error{
						Message: handlers.FmtString(err.Error()),
					}
					return mtoshipmentops.NewCreateMTOShipmentNotFound().WithPayload(&payload), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						"Validation errors",
						"CreateMTOShipment",
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors)
					return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payload), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("ghcapi.CreateMTOShipmentHandler query error", zap.Error(e.Unwrap()))
					}
					return mtoshipmentops.NewCreateMTOShipmentInternalServerError(), err
				default:
					return mtoshipmentops.NewCreateMTOShipmentInternalServerError(), err
				}
			}

			mtoShipment := payloads.MTOShipmentModelFromCreate(payload)

			var err error
			mtoShipment, err = h.shipmentCreator.CreateShipment(appCtx, mtoShipment)

			if err != nil {
				return handleError(err)
			}

			if mtoShipment == nil {
				shipmentNotCreatedError := apperror.NewInternalServerError("Unexpected nil shipment from CreateMTOShipment")
				appCtx.Logger().Error(shipmentNotCreatedError.Error())
				return mtoshipmentops.NewCreateMTOShipmentInternalServerError(), shipmentNotCreatedError
			}

			sitAllowance, err := h.shipmentStatus.CalculateShipmentSITAllowance(appCtx, *mtoShipment)
			if err != nil {
				return handleError(err)
			}

			mtoShipment.SITDaysAllowance = &sitAllowance

			returnPayload := payloads.MTOShipment(mtoShipment, nil)
			return mtoshipmentops.NewCreateMTOShipmentOK().WithPayload(returnPayload), nil
		})
}

// UpdateShipmentHandler updates shipments
type UpdateShipmentHandler struct {
	handlers.HandlerConfig
	services.ShipmentUpdater
	services.ShipmentSITStatus
}

// Handle updates shipments
func (h UpdateShipmentHandler) Handle(params mtoshipmentops.UpdateMTOShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			payload := params.Body
			if payload == nil {
				appCtx.Logger().Error("Invalid mto shipment: params Body is nil")
				emptyBodyError := apperror.NewBadDataError("The MTO Shipment request body cannot be empty.")
				payload := payloadForValidationError(
					"Empty body error",
					emptyBodyError.Error(),
					h.GetTraceIDFromRequest(params.HTTPRequest),
					validate.NewErrors(),
				)

				return mtoshipmentops.NewUpdateMTOShipmentUnprocessableEntity().WithPayload(payload), emptyBodyError
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			oldShipment, err := mtoshipment.FindShipment(appCtx, shipmentID)

			if err != nil {
				appCtx.Logger().Error("ghcapi.UpdateShipmentHandler", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewUpdateMTOShipmentNotFound(), err
				default:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))

					return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(
						&ghcmessages.Error{Message: &msg},
					), err
				}
			}

			mtoShipment := payloads.MTOShipmentModelFromUpdate(payload)
			mtoShipment.ID = shipmentID
			mtoShipment.ShipmentType = oldShipment.ShipmentType

			//MTOShipmentModelFromUpdate defaults UsesExternalVendor to false if it's nil in the payload
			if payload.UsesExternalVendor == nil {
				mtoShipment.UsesExternalVendor = oldShipment.UsesExternalVendor
			}
			// booleans not passed will update to false
			mtoShipment.Diversion = oldShipment.Diversion

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.UpdateShipmentHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewUpdateMTOShipmentNotFound(), err
				case apperror.ForbiddenError:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
					return mtoshipmentops.NewUpdateMTOShipmentForbidden().WithPayload(
						&ghcmessages.Error{Message: &msg},
					), err
				case apperror.InvalidInputError:
					return mtoshipmentops.NewUpdateMTOShipmentUnprocessableEntity().WithPayload(
						payloadForValidationError(
							handlers.ValidationErrMessage,
							err.Error(),
							h.GetTraceIDFromRequest(params.HTTPRequest),
							e.ValidationErrors,
						),
					), err
				case apperror.PreconditionFailedError:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
					return mtoshipmentops.NewUpdateMTOShipmentPreconditionFailed().WithPayload(
						&ghcmessages.Error{Message: &msg},
					), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("ghcapi.UpdateShipmentHandler error", zap.Error(e.Unwrap()))
					}

					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))

					return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(
						&ghcmessages.Error{Message: &msg},
					), err
				default:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))

					return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(
						&ghcmessages.Error{Message: &msg},
					), err
				}
			}
			updatedMtoShipment, err := h.ShipmentUpdater.UpdateShipment(appCtx, mtoShipment, params.IfMatch)
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
			return mtoshipmentops.NewUpdateMTOShipmentOK().WithPayload(returnPayload), nil
		})
}

// DeleteShipmentHandler soft deletes a shipment
type DeleteShipmentHandler struct {
	handlers.HandlerConfig
	services.ShipmentDeleter
}

// Handle soft deletes a shipment
func (h DeleteShipmentHandler) Handle(params shipmentops.DeleteShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeServicesCounselor) {
				forbiddenError := apperror.NewForbiddenError("user is not authenticated with service counselor office role")
				appCtx.Logger().Error(forbiddenError.Error())
				return shipmentops.NewDeleteShipmentForbidden(), forbiddenError
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			moveID, err := h.DeleteShipment(appCtx, shipmentID)
			if err != nil {
				appCtx.Logger().Error("ghcapi.DeleteShipmentHandler", zap.Error(err))

				switch err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewDeleteShipmentNotFound(), err
				case apperror.ConflictError:
					return shipmentops.NewDeleteShipmentConflict(), err
				case apperror.ForbiddenError:
					return shipmentops.NewDeleteShipmentForbidden(), err
				case apperror.UnprocessableEntityError:
					return shipmentops.NewDeleteShipmentUnprocessableEntity(), err
				default:
					return shipmentops.NewDeleteShipmentInternalServerError(), err
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

			return shipmentops.NewDeleteShipmentNoContent(), nil
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
	handlers.HandlerConfig
	services.ShipmentApprover
	services.ShipmentSITStatus
}

// Handle approves a shipment
func (h ApproveShipmentHandler) Handle(params shipmentops.ApproveShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				forbiddenError := apperror.NewForbiddenError("Only TOO role can approve shipments")
				appCtx.Logger().Error(forbiddenError.Error())
				return shipmentops.NewApproveShipmentForbidden(), forbiddenError
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			eTag := params.IfMatch

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.ApproveShipmentHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewApproveShipmentNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Validation errors", "ApproveShipment", h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
					return shipmentops.NewApproveShipmentUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return shipmentops.NewApproveShipmentPreconditionFailed().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ConflictError, mtoshipment.ConflictStatusError:
					return shipmentops.NewApproveShipmentConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return shipmentops.NewApproveShipmentInternalServerError(), err
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
			return shipmentops.NewApproveShipmentOK().WithPayload(payload), nil
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
	handlers.HandlerConfig
	services.ShipmentDiversionRequester
	services.ShipmentSITStatus
}

// Handle Requests a shipment diversion
func (h RequestShipmentDiversionHandler) Handle(params shipmentops.RequestShipmentDiversionParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				forbiddenError := apperror.NewForbiddenError("Only TOO role can Request shipment diversions")
				appCtx.Logger().Error(forbiddenError.Error())
				return shipmentops.NewRequestShipmentDiversionForbidden(), forbiddenError
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			eTag := params.IfMatch

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.RequestShipmentDiversionHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewRequestShipmentDiversionNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						"Validation errors",
						"RequestShipmentDiversion",
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors)
					return shipmentops.NewRequestShipmentDiversionUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return shipmentops.NewRequestShipmentDiversionPreconditionFailed().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case mtoshipment.ConflictStatusError:
					return shipmentops.NewRequestShipmentDiversionConflict().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return shipmentops.NewRequestShipmentDiversionInternalServerError(), err
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
			return shipmentops.NewRequestShipmentDiversionOK().WithPayload(payload), nil
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
	handlers.HandlerConfig
	services.ShipmentDiversionApprover
	services.ShipmentSITStatus
}

// Handle approves a shipment diversion
func (h ApproveShipmentDiversionHandler) Handle(params shipmentops.ApproveShipmentDiversionParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				forbiddenError := apperror.NewForbiddenError("Only TOO role can approve shipment diversions")
				appCtx.Logger().Error(forbiddenError.Error())
				return shipmentops.NewApproveShipmentDiversionForbidden(), forbiddenError
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			eTag := params.IfMatch

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.ApproveShipmentDiversionHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewApproveShipmentDiversionNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						"Validation errors",
						"ApproveShipmentDiversion",
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors)
					return shipmentops.NewApproveShipmentDiversionUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return shipmentops.NewApproveShipmentDiversionPreconditionFailed().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case mtoshipment.ConflictStatusError:
					return shipmentops.NewApproveShipmentDiversionConflict().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return shipmentops.NewApproveShipmentDiversionInternalServerError(), err
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
			return shipmentops.NewApproveShipmentDiversionOK().WithPayload(payload), nil
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
	handlers.HandlerConfig
	services.ShipmentRejecter
}

// Handle rejects a shipment
func (h RejectShipmentHandler) Handle(params shipmentops.RejectShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				forbiddenError := apperror.NewForbiddenError("Only TOO role can reject shipments")
				appCtx.Logger().Error(forbiddenError.Error())
				return shipmentops.NewRejectShipmentForbidden(), forbiddenError
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			eTag := params.IfMatch
			rejectionReason := params.Body.RejectionReason
			shipment, err := h.RejectShipment(appCtx, shipmentID, eTag, rejectionReason)

			if err != nil {
				appCtx.Logger().Error("ghcapi.RejectShipmentHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewRejectShipmentNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						"Validation errors",
						"RejectShipment",
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors)
					return shipmentops.NewRejectShipmentUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return shipmentops.NewRejectShipmentPreconditionFailed().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case mtoshipment.ConflictStatusError:
					return shipmentops.NewRejectShipmentConflict().
							WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}),
						err
				default:
					return shipmentops.NewRejectShipmentInternalServerError(), err
				}
			}

			h.triggerShipmentRejectionEvent(appCtx, shipmentID, shipment.MoveTaskOrderID, params)

			payload := payloads.MTOShipment(shipment, nil)
			return shipmentops.NewRejectShipmentOK().WithPayload(payload), nil
		})
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
	handlers.HandlerConfig
	services.ShipmentCancellationRequester
	services.ShipmentSITStatus
}

// Handle Requests a shipment diversion
func (h RequestShipmentCancellationHandler) Handle(params shipmentops.RequestShipmentCancellationParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				forbiddenError := apperror.NewForbiddenError("Only TOO role can Request shipment diversions")
				appCtx.Logger().Error(forbiddenError.Error())
				return shipmentops.NewRequestShipmentCancellationForbidden(), forbiddenError
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			eTag := params.IfMatch

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.RequestShipmentCancellationHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewRequestShipmentCancellationNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						"Validation errors",
						"RequestShipmentCancellation",
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors)
					return shipmentops.NewRequestShipmentCancellationUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return shipmentops.NewRequestShipmentCancellationPreconditionFailed().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case mtoshipment.ConflictStatusError:
					return shipmentops.NewRequestShipmentCancellationConflict().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return shipmentops.NewRequestShipmentCancellationInternalServerError(), err
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
			return shipmentops.NewRequestShipmentCancellationOK().WithPayload(payload), nil
		})
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
	handlers.HandlerConfig
	services.ShipmentReweighRequester
	services.ShipmentSITStatus
	services.MTOShipmentUpdater
}

// Handle Requests a shipment reweigh
func (h RequestShipmentReweighHandler) Handle(params shipmentops.RequestShipmentReweighParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				forbiddenError := apperror.NewForbiddenError("Only TOO role can Request a shipment reweigh")
				appCtx.Logger().Error(forbiddenError.Error())
				return shipmentops.NewRequestShipmentReweighForbidden(), forbiddenError
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			reweigh, err := h.RequestShipmentReweigh(appCtx, shipmentID, models.ReweighRequesterTOO)

			if err != nil {
				appCtx.Logger().Error("ghcapi.RequestShipmentReweighHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewRequestShipmentReweighNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						"Validation errors", "RequestShipmentReweigh",
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors)
					return shipmentops.NewRequestShipmentReweighUnprocessableEntity().WithPayload(payload), err
				case apperror.ConflictError:
					return shipmentops.NewRequestShipmentReweighConflict().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return shipmentops.NewRequestShipmentReweighInternalServerError(), err
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

			shipment, err := mtoshipment.FindShipment(appCtx, shipmentID)
			if err != nil {
				return handleError(err), err
			}

			moveID := shipment.MoveTaskOrderID
			h.triggerRequestShipmentReweighEvent(appCtx, shipmentID, moveID, params)

			err = h.NotificationSender().SendNotification(appCtx,
				notifications.NewReweighRequested(moveID, *shipment),
			)
			if err != nil {
				appCtx.Logger().Error("problem sending email to user", zap.Error(err))
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			shipmentSITStatus, err := h.CalculateShipmentSITStatus(appCtx, reweigh.Shipment)
			if err != nil {
				return handleError(err), err
			}
			sitStatusPayload := payloads.SITStatus(shipmentSITStatus)

			payload := payloads.Reweigh(reweigh, sitStatusPayload)
			return shipmentops.NewRequestShipmentReweighOK().WithPayload(payload), nil
		})
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
	handlers.HandlerConfig
	services.SITExtensionApprover
	services.ShipmentSITStatus
}

// Handle ... approves the SIT extension
func (h ApproveSITExtensionHandler) Handle(params shipmentops.ApproveSITExtensionParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("error approving SIT extension", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewApproveSITExtensionNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						handlers.ValidationErrMessage,
						err.Error(),
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors)
					return shipmentops.NewApproveSITExtensionUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return shipmentops.NewApproveSITExtensionPreconditionFailed().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ForbiddenError:
					return shipmentops.NewApproveSITExtensionForbidden().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return shipmentops.NewApproveSITExtensionInternalServerError(), err
				}
			}

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				forbiddenError := apperror.NewForbiddenError("is not a TOO")
				return handleError(forbiddenError)
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
			return shipmentops.NewApproveSITExtensionOK().WithPayload(shipmentPayload), nil
		})
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
	handlers.HandlerConfig
	services.SITExtensionDenier
	services.ShipmentSITStatus
}

// Handle ... denies the SIT extension
func (h DenySITExtensionHandler) Handle(params shipmentops.DenySITExtensionParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("error denying SIT extension", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return shipmentops.NewDenySITExtensionNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						handlers.ValidationErrMessage,
						err.Error(),
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors)
					return shipmentops.NewDenySITExtensionUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return shipmentops.NewDenySITExtensionPreconditionFailed().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ForbiddenError:
					return shipmentops.NewDenySITExtensionForbidden().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return shipmentops.NewDenySITExtensionInternalServerError(), err
				}
			}

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				forbiddenError := apperror.NewForbiddenError("is not a TOO")
				return handleError(forbiddenError)
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

			return shipmentops.NewDenySITExtensionOK().WithPayload(shipmentPayload), nil
		})
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
	handlers.HandlerConfig
	services.SITExtensionCreatorAsTOO
	services.ShipmentSITStatus
}

// Handle creates the approved SIT extension
func (h CreateSITExtensionAsTOOHandler) Handle(params shipmentops.CreateSITExtensionAsTOOParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			payload := params.Body
			shipmentID := params.ShipmentID

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.CreateApprovedSITExtension error", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					payload := ghcmessages.Error{
						Message: handlers.FmtString(err.Error()),
					}
					return shipmentops.NewCreateSITExtensionAsTOONotFound().WithPayload(&payload), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						"Validation errors",
						"CreateApprovedSITExtension",
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors)
					return shipmentops.NewCreateSITExtensionAsTOOUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return shipmentops.NewDenySITExtensionPreconditionFailed().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("ghcapi.CreateApprovedSITExtension query error", zap.Error(e.Unwrap()))
					}
					return shipmentops.NewCreateSITExtensionAsTOOInternalServerError(), err
				case apperror.ForbiddenError:
					return shipmentops.NewCreateSITExtensionAsTOOForbidden().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return shipmentops.NewCreateSITExtensionAsTOOInternalServerError(), err
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
			return shipmentops.NewCreateSITExtensionAsTOOOK().WithPayload(returnPayload), nil
		})
}
