package ghcapi

import (
	"fmt"

	"github.com/transcom/mymove/pkg/notifications"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models/roles"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
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
	"github.com/transcom/mymove/pkg/services/query"
)

// ListMTOShipmentsHandler returns a list of MTO Shipments
type ListMTOShipmentsHandler struct {
	handlers.HandlerContext
	services.ListFetcher
	services.Fetcher
}

// Handle listing mto shipments for the move task order
func (h ListMTOShipmentsHandler) Handle(params mtoshipmentops.ListMTOShipmentsParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)

	moveTaskOrderID, err := uuid.FromString(params.MoveTaskOrderID.String())
	// return any parsing error
	if err != nil {
		parsingError := fmt.Errorf("UUID Parsing for %s: %w", "MoveTaskOrderID", err).Error()
		logger.Error(parsingError)
		payload := payloadForValidationError("UUID(s) parsing error", parsingError, h.GetTraceID(), validate.NewErrors())

		return mtoshipmentops.NewListMTOShipmentsUnprocessableEntity().WithPayload(payload)
	}

	// check if move task order exists first
	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", moveTaskOrderID.String()),
	}

	moveTaskOrder := &models.Move{}
	err = h.Fetcher.FetchRecord(appCtx, moveTaskOrder, queryFilters)
	if err != nil {
		logger.Error("Error fetching move task order: ", zap.Error(fmt.Errorf("Move Task Order ID: %s", moveTaskOrder.ID)), zap.Error(err))

		return mtoshipmentops.NewListMTOShipmentsNotFound()
	}

	queryFilters = []services.QueryFilter{
		query.NewQueryFilter("move_id", "=", moveTaskOrderID.String()),
	}

	// TODO: These associations could be preloaded, but it will require Pop 5.3.4 to land first as it
	//   has a fix for using a "has_many" association that has a pointer-based foreign key (like the
	//   case with "MTOServiceItems.ReService"). There appear to be other changes that will need to be
	//   made for Pop 5.3.4 though (see https://ustcdp3.slack.com/archives/CP497TGAU/p1620421441217700).
	queryAssociations := query.NewQueryAssociations([]services.QueryAssociation{
		query.NewQueryAssociation("MTOServiceItems.ReService"),
		query.NewQueryAssociation("MTOAgents"),
		query.NewQueryAssociation("PickupAddress"),
		query.NewQueryAssociation("SecondaryPickupAddress"),
		query.NewQueryAssociation("DestinationAddress"),
		query.NewQueryAssociation("SecondaryPickupAddress"),
		query.NewQueryAssociation("SecondaryDeliveryAddress"),
		query.NewQueryAssociation("MTOServiceItems.Dimensions"),
		query.NewQueryAssociation("Reweigh"),
		query.NewQueryAssociation("SITExtensions"),
	})

	queryOrder := query.NewQueryOrder(swag.String("created_at"), swag.Bool(true))

	var shipments models.MTOShipments
	err = h.ListFetcher.FetchRecordList(appCtx, &shipments, queryFilters, queryAssociations, nil, queryOrder)
	// return any errors
	if err != nil {
		logger.Error("Error fetching mto shipments : ", zap.Error(err))

		return mtoshipmentops.NewListMTOShipmentsInternalServerError()
	}

	payload := payloads.MTOShipments(&shipments)
	return mtoshipmentops.NewListMTOShipmentsOK().WithPayload(*payload)
}

// CreateMTOShipmentHandler is the handler to create MTO shipments
type CreateMTOShipmentHandler struct {
	handlers.HandlerContext
	mtoShipmentCreator services.MTOShipmentCreator
}

// Handle creates the mto shipment
func (h CreateMTOShipmentHandler) Handle(params mtoshipmentops.CreateMTOShipmentParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)
	payload := params.Body

	if payload == nil {
		logger.Error("Invalid mto shipment: params Body is nil")
		return mtoshipmentops.NewCreateMTOShipmentBadRequest()
	}

	mtoShipment := payloads.MTOShipmentModelFromCreate(payload)
	mtoShipment, err := h.mtoShipmentCreator.CreateMTOShipment(appCtx, mtoShipment, nil)

	if err != nil {
		logger.Error("ghcapi.CreateMTOShipmentHandler error", zap.Error(err))
		switch e := err.(type) {
		case services.NotFoundError:
			payload := ghcmessages.Error{
				Message: handlers.FmtString(err.Error()),
			}
			return mtoshipmentops.NewCreateMTOShipmentNotFound().WithPayload(&payload)
		case services.InvalidInputError:
			payload := payloadForValidationError("Validation errors", "CreateMTOShipment", h.GetTraceID(), e.ValidationErrors)
			return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payload)
		case services.QueryError:
			if e.Unwrap() != nil {
				// If you can unwrap, log the internal error (usually a pq error) for better debugging
				logger.Error("ghcapi.CreateMTOShipmentHandler query error", zap.Error(e.Unwrap()))
			}
			return mtoshipmentops.NewCreateMTOShipmentInternalServerError()
		default:
			return mtoshipmentops.NewCreateMTOShipmentInternalServerError()
		}
	}

	returnPayload := payloads.MTOShipment(mtoShipment)
	return mtoshipmentops.NewCreateMTOShipmentOK().WithPayload(returnPayload)
}

// UpdateShipmentHandler updates shipments
type UpdateShipmentHandler struct {
	handlers.HandlerContext
	services.Fetcher
	services.MTOShipmentUpdater
}

// Handle updates shipments
func (h UpdateShipmentHandler) Handle(params mtoshipmentops.UpdateMTOShipmentParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)

	payload := params.Body
	if payload == nil {
		logger.Error("Invalid mto shipment: params Body is nil")

		payload := payloadForValidationError(
			"Empty body error",
			"The MTO Shipment request body cannot be empty.",
			h.GetTraceID(),
			validate.NewErrors(),
		)

		return mtoshipmentops.NewUpdateMTOShipmentUnprocessableEntity().WithPayload(payload)
	}

	shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
	oldShipment, err := h.MTOShipmentUpdater.RetrieveMTOShipment(appCtx, shipmentID)

	if err != nil {
		logger.Error("ghcapi.UpdateShipmentHandler", zap.Error(err))
		switch err.(type) {
		case services.NotFoundError:
			return mtoshipmentops.NewUpdateMTOShipmentNotFound()
		default:
			msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceID())

			return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(
				&ghcmessages.Error{Message: &msg},
			)
		}
	}

	updateable, err := h.MTOShipmentUpdater.CheckIfMTOShipmentCanBeUpdated(appCtx, oldShipment, session)

	if err != nil {
		logger.Error("ghcapi.UpdateShipmentHandler", zap.Error(err))
		msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceID())
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
	mtoShipment.ID = shipmentID

	updatedMtoShipment, err := h.MTOShipmentUpdater.UpdateMTOShipmentOffice(appCtx, mtoShipment, params.IfMatch)

	if err != nil {
		logger.Error("ghcapi.UpdateShipmentHandler", zap.Error(err))

		switch e := err.(type) {
		case services.NotFoundError:
			return mtoshipmentops.NewUpdateMTOShipmentNotFound()
		case services.InvalidInputError:
			return mtoshipmentops.NewUpdateMTOShipmentUnprocessableEntity().WithPayload(
				payloadForValidationError(
					handlers.ValidationErrMessage,
					err.Error(),
					h.GetTraceID(),
					e.ValidationErrors,
				),
			)
		case services.PreconditionFailedError:
			msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceID())
			return mtoshipmentops.NewUpdateMTOShipmentPreconditionFailed().WithPayload(
				&ghcmessages.Error{Message: &msg},
			)
		case services.QueryError:
			if e.Unwrap() != nil {
				// If you can unwrap, log the internal error (usually a pq error) for better debugging
				logger.Error("ghcapi.UpdateShipmentHandler error", zap.Error(e.Unwrap()))
			}

			msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceID())

			return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(
				&ghcmessages.Error{Message: &msg},
			)
		default:
			msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceID())

			return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(
				&ghcmessages.Error{Message: &msg},
			)
		}
	}

	_, err = event.TriggerEvent(event.Event{
		EndpointKey: event.GhcUpdateMTOShipmentEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.MTOShipmentUpdateEventKey,    // Event that you want to trigger
		UpdatedObjectID: updatedMtoShipment.ID,              // ID of the updated logical object
		MtoID:           updatedMtoShipment.MoveTaskOrderID, // ID of the associated Move
		Request:         params.HTTPRequest,                 // Pass on the http.Request
		DBConnection:    h.DB(),                             // Pass on the pop.Connection
		HandlerContext:  h,                                  // Pass on the handlerContext
	})
	// If the event trigger fails, just log the error.
	if err != nil {
		logger.Error("ghcapi.UpdateMTOShipment could not generate the event")
	}

	returnPayload := payloads.MTOShipment(updatedMtoShipment)
	return mtoshipmentops.NewUpdateMTOShipmentOK().WithPayload(returnPayload)
}

// DeleteShipmentHandler soft deletes a shipment
type DeleteShipmentHandler struct {
	handlers.HandlerContext
	services.ShipmentDeleter
}

// Handle soft deletes a shipment
func (h DeleteShipmentHandler) Handle(params shipmentops.DeleteShipmentParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)

	if !session.IsOfficeUser() || !session.Roles.HasRole(roles.RoleTypeServicesCounselor) {
		logger.Error("user is not authenticated with service counselor office role")
		return shipmentops.NewDeleteShipmentForbidden()
	}

	shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
	moveID, err := h.DeleteShipment(appCtx, shipmentID)
	if err != nil {
		logger.Error("ghcapi.DeleteShipmentHandler", zap.Error(err))

		switch err.(type) {
		case services.NotFoundError:
			return shipmentops.NewDeleteShipmentNotFound()
		case services.ForbiddenError:
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
	h.triggerShipmentDeletionEvent(shipmentID, moveID, params)

	return shipmentops.NewDeleteShipmentNoContent()
}

func (h DeleteShipmentHandler) triggerShipmentDeletionEvent(shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.DeleteShipmentParams) {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcDeleteShipmentEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ShipmentDeleteEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                   // ID of the updated logical object
		MtoID:           moveID,                       // ID of the associated Move
		Request:         params.HTTPRequest,           // Pass on the http.Request
		DBConnection:    h.DB(),                       // Pass on the pop.Connection
		HandlerContext:  h,                            // Pass on the handlerContext
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		logger.Error("ghcapi.DeleteShipmentHandler could not generate the event", zap.Error(err))
	}
}

// ApproveShipmentHandler approves a shipment
type ApproveShipmentHandler struct {
	handlers.HandlerContext
	services.ShipmentApprover
}

// Handle approves a shipment
func (h ApproveShipmentHandler) Handle(params shipmentops.ApproveShipmentParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)

	if !session.IsOfficeUser() || !session.Roles.HasRole(roles.RoleTypeTOO) {
		logger.Error("Only TOO role can approve shipments")
		return shipmentops.NewApproveShipmentForbidden()
	}

	shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
	eTag := params.IfMatch
	shipment, err := h.ApproveShipment(appCtx, shipmentID, eTag)

	if err != nil {
		logger.Error("ghcapi.ApproveShipmentHandler", zap.Error(err))

		switch e := err.(type) {
		case services.NotFoundError:
			return shipmentops.NewApproveShipmentNotFound()
		case services.InvalidInputError:
			payload := payloadForValidationError("Validation errors", "ApproveShipment", h.GetTraceID(), e.ValidationErrors)
			return shipmentops.NewApproveShipmentUnprocessableEntity().WithPayload(payload)
		case services.PreconditionFailedError:
			return shipmentops.NewApproveShipmentPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case services.ConflictError, mtoshipment.ConflictStatusError:
			return shipmentops.NewApproveShipmentConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return shipmentops.NewApproveShipmentInternalServerError()
		}
	}

	h.triggerShipmentApprovalEvent(shipmentID, shipment.MoveTaskOrderID, params)

	payload := payloads.MTOShipment(shipment)
	return shipmentops.NewApproveShipmentOK().WithPayload(payload)
}

func (h ApproveShipmentHandler) triggerShipmentApprovalEvent(shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.ApproveShipmentParams) {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcApproveShipmentEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ShipmentApproveEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                    // ID of the updated logical object
		MtoID:           moveID,                        // ID of the associated Move
		Request:         params.HTTPRequest,            // Pass on the http.Request
		DBConnection:    h.DB(),                        // Pass on the pop.Connection
		HandlerContext:  h,                             // Pass on the handlerContext
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		logger.Error("ghcapi.ApproveShipmentHandler could not generate the event", zap.Error(err))
	}
}

// RequestShipmentDiversionHandler Requests a shipment diversion
type RequestShipmentDiversionHandler struct {
	handlers.HandlerContext
	services.ShipmentDiversionRequester
}

// Handle Requests a shipment diversion
func (h RequestShipmentDiversionHandler) Handle(params shipmentops.RequestShipmentDiversionParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)

	if !session.IsOfficeUser() || !session.Roles.HasRole(roles.RoleTypeTOO) {
		logger.Error("Only TOO role can Request shipment diversions")
		return shipmentops.NewRequestShipmentDiversionForbidden()
	}

	shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
	eTag := params.IfMatch
	shipment, err := h.RequestShipmentDiversion(appCtx, shipmentID, eTag)

	if err != nil {
		logger.Error("ghcapi.RequestShipmentDiversionHandler", zap.Error(err))

		switch e := err.(type) {
		case services.NotFoundError:
			return shipmentops.NewRequestShipmentDiversionNotFound()
		case services.InvalidInputError:
			payload := payloadForValidationError("Validation errors", "RequestShipmentDiversion", h.GetTraceID(), e.ValidationErrors)
			return shipmentops.NewRequestShipmentDiversionUnprocessableEntity().WithPayload(payload)
		case services.PreconditionFailedError:
			return shipmentops.NewRequestShipmentDiversionPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case mtoshipment.ConflictStatusError:
			return shipmentops.NewRequestShipmentDiversionConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return shipmentops.NewRequestShipmentDiversionInternalServerError()
		}
	}

	h.triggerRequestShipmentDiversionEvent(shipmentID, shipment.MoveTaskOrderID, params)

	payload := payloads.MTOShipment(shipment)
	return shipmentops.NewRequestShipmentDiversionOK().WithPayload(payload)
}

func (h RequestShipmentDiversionHandler) triggerRequestShipmentDiversionEvent(shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.RequestShipmentDiversionParams) {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcRequestShipmentDiversionEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ShipmentRequestDiversionEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                             // ID of the updated logical object
		MtoID:           moveID,                                 // ID of the associated Move
		Request:         params.HTTPRequest,                     // Pass on the http.Request
		DBConnection:    h.DB(),                                 // Pass on the pop.Connection
		HandlerContext:  h,                                      // Pass on the handlerContext
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		logger.Error("ghcapi.RequestShipmentDiversionHandler could not generate the event", zap.Error(err))
	}
}

// ApproveShipmentDiversionHandler approves a shipment diversion
type ApproveShipmentDiversionHandler struct {
	handlers.HandlerContext
	services.ShipmentDiversionApprover
}

// Handle approves a shipment diversion
func (h ApproveShipmentDiversionHandler) Handle(params shipmentops.ApproveShipmentDiversionParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)

	if !session.IsOfficeUser() || !session.Roles.HasRole(roles.RoleTypeTOO) {
		logger.Error("Only TOO role can approve shipment diversions")
		return shipmentops.NewApproveShipmentDiversionForbidden()
	}

	shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
	eTag := params.IfMatch
	shipment, err := h.ApproveShipmentDiversion(appCtx, shipmentID, eTag)

	if err != nil {
		logger.Error("ghcapi.ApproveShipmentDiversionHandler", zap.Error(err))

		switch e := err.(type) {
		case services.NotFoundError:
			return shipmentops.NewApproveShipmentDiversionNotFound()
		case services.InvalidInputError:
			payload := payloadForValidationError("Validation errors", "ApproveShipmentDiversion", h.GetTraceID(), e.ValidationErrors)
			return shipmentops.NewApproveShipmentDiversionUnprocessableEntity().WithPayload(payload)
		case services.PreconditionFailedError:
			return shipmentops.NewApproveShipmentDiversionPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case mtoshipment.ConflictStatusError:
			return shipmentops.NewApproveShipmentDiversionConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return shipmentops.NewApproveShipmentDiversionInternalServerError()
		}
	}

	h.triggerShipmentDiversionApprovalEvent(shipmentID, shipment.MoveTaskOrderID, params)

	payload := payloads.MTOShipment(shipment)
	return shipmentops.NewApproveShipmentDiversionOK().WithPayload(payload)
}

func (h ApproveShipmentDiversionHandler) triggerShipmentDiversionApprovalEvent(shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.ApproveShipmentDiversionParams) {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcApproveShipmentDiversionEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ShipmentApproveDiversionEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                             // ID of the updated logical object
		MtoID:           moveID,                                 // ID of the associated Move
		Request:         params.HTTPRequest,                     // Pass on the http.Request
		DBConnection:    h.DB(),                                 // Pass on the pop.Connection
		HandlerContext:  h,                                      // Pass on the handlerContext
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		logger.Error("ghcapi.ApproveShipmentDiversionHandler could not generate the event", zap.Error(err))
	}
}

// RejectShipmentHandler rejects a shipment
type RejectShipmentHandler struct {
	handlers.HandlerContext
	services.ShipmentRejecter
}

// Handle rejects a shipment
func (h RejectShipmentHandler) Handle(params shipmentops.RejectShipmentParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)

	if !session.IsOfficeUser() || !session.Roles.HasRole(roles.RoleTypeTOO) {
		logger.Error("Only TOO role can reject shipments")
		return shipmentops.NewRejectShipmentForbidden()
	}

	shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
	eTag := params.IfMatch
	rejectionReason := params.Body.RejectionReason
	shipment, err := h.RejectShipment(appCtx, shipmentID, eTag, rejectionReason)

	if err != nil {
		logger.Error("ghcapi.RejectShipmentHandler", zap.Error(err))

		switch e := err.(type) {
		case services.NotFoundError:
			return shipmentops.NewRejectShipmentNotFound()
		case services.InvalidInputError:
			payload := payloadForValidationError("Validation errors", "RejectShipment", h.GetTraceID(), e.ValidationErrors)
			return shipmentops.NewRejectShipmentUnprocessableEntity().WithPayload(payload)
		case services.PreconditionFailedError:
			return shipmentops.NewRejectShipmentPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case mtoshipment.ConflictStatusError:
			return shipmentops.NewRejectShipmentConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return shipmentops.NewRejectShipmentInternalServerError()
		}
	}

	h.triggerShipmentRejectionEvent(shipmentID, shipment.MoveTaskOrderID, params)

	payload := payloads.MTOShipment(shipment)
	return shipmentops.NewRejectShipmentOK().WithPayload(payload)
}

func (h RejectShipmentHandler) triggerShipmentRejectionEvent(shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.RejectShipmentParams) {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcRejectShipmentEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ShipmentRejectEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                   // ID of the updated logical object
		MtoID:           moveID,                       // ID of the associated Move
		Request:         params.HTTPRequest,           // Pass on the http.Request
		DBConnection:    h.DB(),                       // Pass on the pop.Connection
		HandlerContext:  h,                            // Pass on the handlerContext
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		logger.Error("ghcapi.RejectShipmentHandler could not generate the event", zap.Error(err))
	}
}

// RequestShipmentCancellationHandler Requests a shipment diversion
type RequestShipmentCancellationHandler struct {
	handlers.HandlerContext
	services.ShipmentCancellationRequester
}

// Handle Requests a shipment diversion
func (h RequestShipmentCancellationHandler) Handle(params shipmentops.RequestShipmentCancellationParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)

	if !session.IsOfficeUser() || !session.Roles.HasRole(roles.RoleTypeTOO) {
		logger.Error("Only TOO role can Request shipment diversions")
		return shipmentops.NewRequestShipmentCancellationForbidden()
	}

	shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
	eTag := params.IfMatch
	shipment, err := h.RequestShipmentCancellation(appCtx, shipmentID, eTag)

	if err != nil {
		logger.Error("ghcapi.RequestShipmentCancellationHandler", zap.Error(err))

		switch e := err.(type) {
		case services.NotFoundError:
			return shipmentops.NewRequestShipmentCancellationNotFound()
		case services.InvalidInputError:
			payload := payloadForValidationError("Validation errors", "RequestShipmentCancellation", h.GetTraceID(), e.ValidationErrors)
			return shipmentops.NewRequestShipmentCancellationUnprocessableEntity().WithPayload(payload)
		case services.PreconditionFailedError:
			return shipmentops.NewRequestShipmentCancellationPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case mtoshipment.ConflictStatusError:
			return shipmentops.NewRequestShipmentCancellationConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return shipmentops.NewRequestShipmentCancellationInternalServerError()
		}
	}

	h.triggerRequestShipmentCancellationEvent(shipmentID, shipment.MoveTaskOrderID, params)

	payload := payloads.MTOShipment(shipment)
	return shipmentops.NewRequestShipmentCancellationOK().WithPayload(payload)
}

func (h RequestShipmentCancellationHandler) triggerRequestShipmentCancellationEvent(shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.RequestShipmentCancellationParams) {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcRequestShipmentCancellationEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ShipmentRequestCancellationEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                                // ID of the updated logical object
		MtoID:           moveID,                                    // ID of the associated Move
		Request:         params.HTTPRequest,                        // Pass on the http.Request
		DBConnection:    h.DB(),                                    // Pass on the pop.Connection
		HandlerContext:  h,                                         // Pass on the handlerContext
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		logger.Error("ghcapi.RequestShipmentCancellationHandler could not generate the event", zap.Error(err))
	}
}

// RequestShipmentReweighHandler Requests a shipment reweigh
type RequestShipmentReweighHandler struct {
	handlers.HandlerContext
	services.ShipmentReweighRequester
}

// Handle Requests a shipment reweigh
func (h RequestShipmentReweighHandler) Handle(params shipmentops.RequestShipmentReweighParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)

	if !session.IsOfficeUser() || !session.Roles.HasRole(roles.RoleTypeTOO) {
		logger.Error("Only TOO role can Request a shipment reweigh")
		return shipmentops.NewRequestShipmentReweighForbidden()
	}

	shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
	reweigh, err := h.RequestShipmentReweigh(appCtx, shipmentID, models.ReweighRequesterTOO)

	if err != nil {
		logger.Error("ghcapi.RequestShipmentReweighHandler", zap.Error(err))

		switch e := err.(type) {
		case services.NotFoundError:
			return shipmentops.NewRequestShipmentReweighNotFound()
		case services.InvalidInputError:
			payload := payloadForValidationError("Validation errors", "RequestShipmentReweigh", h.GetTraceID(), e.ValidationErrors)
			return shipmentops.NewRequestShipmentReweighUnprocessableEntity().WithPayload(payload)
		case services.ConflictError:
			return shipmentops.NewRequestShipmentReweighConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return shipmentops.NewRequestShipmentReweighInternalServerError()
		}
	}

	moveID := reweigh.Shipment.MoveTaskOrderID
	h.triggerRequestShipmentReweighEvent(shipmentID, moveID, params)

	err = h.NotificationSender().SendNotification(
		notifications.NewReweighRequested(h.DB(), logger, session, moveID),
	)
	if err != nil {
		logger.Error("problem sending email to user", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}
	payload := payloads.Reweigh(reweigh)
	return shipmentops.NewRequestShipmentReweighOK().WithPayload(payload)
}

func (h RequestShipmentReweighHandler) triggerRequestShipmentReweighEvent(shipmentID uuid.UUID, moveID uuid.UUID, params shipmentops.RequestShipmentReweighParams) {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcRequestShipmentReweighEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.ShipmentRequestReweighEventKey, // Event that you want to trigger
		UpdatedObjectID: shipmentID,                           // ID of the updated logical object
		MtoID:           moveID,                               // ID of the associated Move
		Request:         params.HTTPRequest,                   // Pass on the http.Request
		DBConnection:    h.DB(),                               // Pass on the pop.Connection
		HandlerContext:  h,                                    // Pass on the handlerContext
	})

	// If the event trigger fails, just log the error.
	if err != nil {
		logger.Error("ghcapi.RequestShipmentReweighHandler could not generate the event", zap.Error(err))
	}
}

// ApproveSITExtensionHandler approves a SIT extension
type ApproveSITExtensionHandler struct {
	handlers.HandlerContext
	services.SITExtensionApprover
}

// Handle ... approves the SIT extension
func (h ApproveSITExtensionHandler) Handle(params shipmentops.ApproveSitExtensionParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)
	handleError := func(err error) middleware.Responder {
		logger.Error("error approving SIT extension", zap.Error(err))
		switch e := err.(type) {
		case services.NotFoundError:
			return shipmentops.NewApproveSitExtensionNotFound()
		case services.InvalidInputError:
			payload := payloadForValidationError(handlers.ValidationErrMessage, err.Error(), h.GetTraceID(), e.ValidationErrors)
			return shipmentops.NewApproveSitExtensionUnprocessableEntity().WithPayload(payload)
		case services.PreconditionFailedError:
			return shipmentops.NewApproveSitExtensionPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case services.ForbiddenError:
			return shipmentops.NewApproveSitExtensionForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return shipmentops.NewApproveSitExtensionInternalServerError()
		}
	}

	if !session.IsOfficeUser() || !session.Roles.HasRole(roles.RoleTypeTOO) {
		return handleError(services.NewForbiddenError("is not a TOO"))
	}

	shipmentID := uuid.FromStringOrNil(string(params.ShipmentID))
	sitExtensionID := uuid.FromStringOrNil(string(params.SitExtensionID))
	approvedDays := int(*params.Body.ApprovedDays)
	officeRemarks := params.Body.OfficeRemarks
	updatedShipment, err := h.SITExtensionApprover.ApproveSITExtension(appCtx, shipmentID, sitExtensionID, approvedDays, officeRemarks, params.IfMatch)
	if err != nil {
		return handleError(err)
	}

	shipmentPayload := payloads.MTOShipment(updatedShipment)

	return shipmentops.NewApproveSitExtensionOK().WithPayload(shipmentPayload)
}

// DenySITExtensionHandler denies a SIT extension
type DenySITExtensionHandler struct {
	handlers.HandlerContext
	services.SITExtensionDenier
}

// Handle ... denies the SIT extension
func (h DenySITExtensionHandler) Handle(params shipmentops.DenySitExtensionParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)
	handleError := func(err error) middleware.Responder {
		logger.Error("error denying SIT extension", zap.Error(err))
		switch e := err.(type) {
		case services.NotFoundError:
			return shipmentops.NewDenySitExtensionNotFound()
		case services.InvalidInputError:
			payload := payloadForValidationError(handlers.ValidationErrMessage, err.Error(), h.GetTraceID(), e.ValidationErrors)
			return shipmentops.NewDenySitExtensionUnprocessableEntity().WithPayload(payload)
		case services.PreconditionFailedError:
			return shipmentops.NewDenySitExtensionPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case services.ForbiddenError:
			return shipmentops.NewDenySitExtensionForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return shipmentops.NewDenySitExtensionInternalServerError()
		}
	}

	if !session.IsOfficeUser() || !session.Roles.HasRole(roles.RoleTypeTOO) {
		return handleError(services.NewForbiddenError("is not a TOO"))
	}

	shipmentID := uuid.FromStringOrNil(string(params.ShipmentID))
	sitExtensionID := uuid.FromStringOrNil(string(params.SitExtensionID))
	officeRemarks := params.Body.OfficeRemarks
	updatedShipment, err := h.SITExtensionDenier.DenySITExtension(appCtx, shipmentID, sitExtensionID, officeRemarks, params.IfMatch)
	if err != nil {
		return handleError(err)
	}

	shipmentPayload := payloads.MTOShipment(updatedShipment)

	return shipmentops.NewDenySitExtensionOK().WithPayload(shipmentPayload)
}

// CreateSITExtensionAsTOOHandler creates a SIT extension in the approved state
type CreateSITExtensionAsTOOHandler struct {
	handlers.HandlerContext
	services.SITExtensionCreatorAsTOO
}

// Handle creates the approved SIT extension
func (h CreateSITExtensionAsTOOHandler) Handle(params shipmentops.CreateSitExtensionAsTOOParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)
	payload := params.Body
	shipmentID := params.ShipmentID

	handleError := func(err error) middleware.Responder {
		logger.Error("ghcapi.CreateApprovedSITExtension error", zap.Error(err))
		switch e := err.(type) {
		case services.NotFoundError:
			payload := ghcmessages.Error{
				Message: handlers.FmtString(err.Error()),
			}
			return shipmentops.NewCreateSitExtensionAsTOONotFound().WithPayload(&payload)
		case services.InvalidInputError:
			payload := payloadForValidationError("Validation errors", "CreateApprovedSITExtension", h.GetTraceID(), e.ValidationErrors)
			return shipmentops.NewCreateSitExtensionAsTOOUnprocessableEntity().WithPayload(payload)
		case services.PreconditionFailedError:
			return shipmentops.NewDenySitExtensionPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case services.QueryError:
			if e.Unwrap() != nil {
				// If you can unwrap, log the internal error (usually a pq error) for better debugging
				logger.Error("ghcapi.CreateApprovedSITExtension query error", zap.Error(e.Unwrap()))
			}
			return shipmentops.NewCreateSitExtensionAsTOOInternalServerError()
		case services.ForbiddenError:
			return shipmentops.NewCreateSitExtensionAsTOOForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return shipmentops.NewCreateSitExtensionAsTOOInternalServerError()
		}
	}

	sitExtension := payloads.ApprovedSITExtensionFromCreate(payload, shipmentID)
	shipment, err := h.SITExtensionCreatorAsTOO.CreateSITExtensionAsTOO(appCtx, sitExtension, sitExtension.MTOShipmentID, params.IfMatch)
	if err != nil {
		return handleError(err)
	}

	if !session.IsOfficeUser() || !session.Roles.HasRole(roles.RoleTypeTOO) {
		return handleError(services.NewForbiddenError("is not a TOO"))
	}

	returnPayload := payloads.MTOShipment(shipment)
	return shipmentops.NewCreateSitExtensionAsTOOOK().WithPayload(returnPayload)
}
