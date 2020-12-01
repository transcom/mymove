package ghcapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	mtoserviceitemop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/audit"
	"github.com/transcom/mymove/pkg/services/event"
	"github.com/transcom/mymove/pkg/services/query"
)

func payloadForClientError(title string, detail string, instance uuid.UUID) *ghcmessages.ClientError {
	return &ghcmessages.ClientError{
		Title:    handlers.FmtString(title),
		Detail:   handlers.FmtString(detail),
		Instance: handlers.FmtUUID(instance),
	}
}

func payloadForValidationError(title string, detail string, instance uuid.UUID, validationErrors *validate.Errors) *ghcmessages.ValidationError {
	return &ghcmessages.ValidationError{
		InvalidFields: handlers.NewValidationErrorsResponse(validationErrors).Errors,
		ClientError:   *payloadForClientError(title, detail, instance),
	}
}

// UpdateMTOServiceItemStatusHandler struct that describes updating service item status
type UpdateMTOServiceItemStatusHandler struct {
	handlers.HandlerContext
	services.MTOServiceItemUpdater
	services.Fetcher
}

// Handle handler that handles the handling for updating service item status
func (h UpdateMTOServiceItemStatusHandler) Handle(params mtoserviceitemop.UpdateMTOServiceItemStatusParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	existingMTOServiceItem := models.MTOServiceItem{}

	mtoServiceItemID, err := uuid.FromString(params.MtoServiceItemID)
	// return parsing errors
	if err != nil {
		parsingError := fmt.Errorf("UUID parsing failed for mtoServiceItem: %w", err).Error()
		logger.Error(parsingError)
		payload := payloadForValidationError("UUID(s) parsing error", parsingError, h.GetTraceID(), validate.NewErrors())

		return mtoserviceitemop.NewUpdateMTOServiceItemStatusUnprocessableEntity().WithPayload(payload)
	}

	// Fetch the existing service item
	filter := []services.QueryFilter{query.NewQueryFilter("id", "=", mtoServiceItemID)}
	err = h.Fetcher.FetchRecord(&existingMTOServiceItem, filter)

	if err != nil {
		logger.Error(fmt.Sprintf("Error finding MTOServiceItem for status update with ID: %s", mtoServiceItemID), zap.Error(err))
		return mtoserviceitemop.NewUpdateMTOServiceItemStatusNotFound()
	}

	// Capture update attempt in audit log
	_, err = audit.Capture(&existingMTOServiceItem, nil, logger, session, params.HTTPRequest)
	if err != nil {
		logger.Error("Auditing service error for service item update.", zap.Error(err))
		return mtoserviceitemop.NewUpdateMTOServiceItemStatusInternalServerError()
	}

	updatedMTOServiceItem, err := h.MTOServiceItemUpdater.UpdateMTOServiceItemStatus(mtoServiceItemID, models.MTOServiceItemStatus(params.Body.Status), params.Body.RejectionReason, params.IfMatch)

	if err != nil {
		switch err.(type) {
		case services.NotFoundError:
			payload := payloadForClientError("Unknown UUID(s)", "Unknown UUID(s) used to update a mto service item", h.GetTraceID())
			return mtoserviceitemop.NewUpdateMTOServiceItemStatusNotFound().WithPayload(payload)
		case services.PreconditionFailedError:
			return mtoserviceitemop.NewUpdateMTOServiceItemStatusPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		case services.InvalidInputError:
			payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceID(), validate.NewErrors())
			return mtoserviceitemop.NewUpdateMTOServiceItemStatusUnprocessableEntity().WithPayload(payload)
		default:
			logger.Error(fmt.Sprintf("Error saving payment request status for ID: %s: %s", mtoServiceItemID, err))
			return mtoserviceitemop.NewUpdateMTOServiceItemStatusInternalServerError()
		}
	}

	// trigger webhook event for Prime
	_, err = event.TriggerEvent(event.Event{
		EventKey:        event.MTOServiceItemUpdateEventKey,
		MtoID:           existingMTOServiceItem.MoveTaskOrder.ID,
		UpdatedObjectID: existingMTOServiceItem.ID,
		Request:         params.HTTPRequest,
		EndpointKey:     event.GhcUpdateMTOServiceItemStatusEndpointKey,
		DBConnection:    h.DB(),
		HandlerContext:  h,
	})

	if err != nil {
		logger.Error("ghcapi.UpdateMTOServiceItemStatusHandler could not generate the event")
	}

	payload := payloads.MTOServiceItemModel(updatedMTOServiceItem)
	return mtoserviceitemop.NewUpdateMTOServiceItemStatusOK().WithPayload(payload)
}

// ListMTOServiceItemsHandler struct that describes listing service items for the move task order
type ListMTOServiceItemsHandler struct {
	handlers.HandlerContext
	services.ListFetcher
	services.Fetcher
}

// Handle handler that lists mto service items for the move task order
func (h ListMTOServiceItemsHandler) Handle(params mtoserviceitemop.ListMTOServiceItemsParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	moveTaskOrderID, err := uuid.FromString(params.MoveTaskOrderID.String())
	// return any parsing error
	if err != nil {
		parsingError := fmt.Errorf("UUID Parsing for %s: %w", "MoveTaskOrderID", err).Error()
		logger.Error(parsingError)
		payload := payloadForValidationError("UUID(s) parsing error", parsingError, h.GetTraceID(), validate.NewErrors())

		return mtoserviceitemop.NewListMTOServiceItemsUnprocessableEntity().WithPayload(payload)
	}

	// check if move task order exists first
	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", moveTaskOrderID.String()),
	}

	moveTaskOrder := &models.Move{}
	err = h.Fetcher.FetchRecord(moveTaskOrder, queryFilters)
	if err != nil {
		logger.Error("Error fetching move task order: ", zap.Error(fmt.Errorf("Move Task Order ID: %s", moveTaskOrder.ID)), zap.Error(err))

		return mtoserviceitemop.NewListMTOServiceItemsNotFound()
	}

	queryFilters = []services.QueryFilter{
		query.NewQueryFilter("move_id", "=", moveTaskOrderID.String()),
	}
	queryAssociations := query.NewQueryAssociations([]services.QueryAssociation{
		query.NewQueryAssociation("ReService"),
		query.NewQueryAssociation("CustomerContacts"),
		query.NewQueryAssociation("Dimensions"),
	})

	var serviceItems models.MTOServiceItems
	err = h.ListFetcher.FetchRecordList(&serviceItems, queryFilters, queryAssociations, nil, nil)
	// return any errors
	if err != nil {
		logger.Error("Error fetching mto service items: ", zap.Error(err))

		return mtoserviceitemop.NewListMTOServiceItemsInternalServerError()
	}

	returnPayload := payloads.MTOServiceItemModels(serviceItems)
	return mtoserviceitemop.NewListMTOServiceItemsOK().WithPayload(returnPayload)
}
