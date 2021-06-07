package ghcapi

import (
	"fmt"

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
	err = h.Fetcher.FetchRecord(moveTaskOrder, queryFilters)
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
		query.NewQueryAssociation("DestinationAddress"),
		query.NewQueryAssociation("SecondaryPickupAddress"),
		query.NewQueryAssociation("SecondaryDeliveryAddress"),
	})

	var shipments models.MTOShipments
	err = h.ListFetcher.FetchRecordList(&shipments, queryFilters, queryAssociations, nil, nil)
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
	payload := params.Body

	if payload == nil {
		logger.Error("Invalid mto shipment: params Body is nil")
		return mtoshipmentops.NewCreateMTOShipmentBadRequest()
	}

	mtoShipment := payloads.MTOShipmentModelFromCreate(payload)
	mtoShipment, err := h.mtoShipmentCreator.CreateMTOShipment(mtoShipment, nil)

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
	oldShipment, err := h.MTOShipmentUpdater.RetrieveMTOShipment(shipmentID)

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

	updateable, err := h.MTOShipmentUpdater.CheckIfMTOShipmentCanBeUpdated(oldShipment, session)

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

	updatedMtoShipment, err := h.MTOShipmentUpdater.UpdateMTOShipment(mtoShipment, params.IfMatch)

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

// PatchShipmentHandler patches shipments
type PatchShipmentHandler struct {
	handlers.HandlerContext
	services.Fetcher
	services.MTOShipmentStatusUpdater
}

// Handle patches shipments
func (h PatchShipmentHandler) Handle(params mtoshipmentops.PatchMTOShipmentStatusParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
	bodyStatus := params.Body.Status
	status := models.MTOShipmentStatus(*bodyStatus)
	rejectionReason := params.Body.RejectionReason
	eTag := params.IfMatch
	shipment, err := h.UpdateMTOShipmentStatus(shipmentID, status, rejectionReason, eTag)
	if err != nil {
		logger.Error("UpdateMTOShipmentStatus error: ", zap.Error(err))

		switch e := err.(type) {
		case services.NotFoundError:
			return mtoshipmentops.NewPatchMTOShipmentStatusNotFound()
		case services.InvalidInputError:
			payload := payloadForValidationError("Validation errors", "UpdateShipmentMTOStatus", h.GetTraceID(), e.ValidationErrors)
			return mtoshipmentops.NewPatchMTOShipmentStatusUnprocessableEntity().WithPayload(payload)
		case services.PreconditionFailedError:
			return mtoshipmentops.NewPatchMTOShipmentStatusPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case mtoshipment.ConflictStatusError:
			return mtoshipmentops.NewPatchMTOShipmentStatusConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return mtoshipmentops.NewPatchMTOShipmentStatusInternalServerError()
		}
	}

	_, err = event.TriggerEvent(event.Event{
		EndpointKey: event.GhcPatchMTOShipmentStatusEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.MTOShipmentUpdateEventKey, // Event that you want to trigger
		UpdatedObjectID: shipment.ID,                     // ID of the updated logical object
		MtoID:           shipment.MoveTaskOrderID,        // ID of the associated Move
		Request:         params.HTTPRequest,              // Pass on the http.Request
		DBConnection:    h.DB(),                          // Pass on the pop.Connection
		HandlerContext:  h,                               // Pass on the handlerContext
	})
	// If the event trigger fails, just log the error.
	if err != nil {
		logger.Error("ghcapi.PatchShipmentHandler could not generate the event")
	}

	payload := payloads.MTOShipment(shipment)
	return mtoshipmentops.NewPatchMTOShipmentStatusOK().WithPayload(payload)
}

// DeleteShipmentHandler soft deletes a shipment
type DeleteShipmentHandler struct {
	handlers.HandlerContext
	services.ShipmentDeleter
}

// Handle soft deletes a shipment
func (h DeleteShipmentHandler) Handle(params shipmentops.DeleteShipmentParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	if !session.IsOfficeUser() || !session.Roles.HasRole(roles.RoleTypeServicesCounselor) {
		logger.Error("user is not authenticated with service counselor office role")
		return shipmentops.NewDeleteShipmentForbidden()
	}

	shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
	moveID, err := h.DeleteShipment(shipmentID)
	if err != nil {
		logger.Error("Error deleting shipment: ", zap.Error(err))

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
		logger.Error("ghcapi.DeleteMTOShipmentHandler could not generate the event")
	}
}
