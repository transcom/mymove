package ghcapi

import (
	"fmt"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	mtoserviceitemop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/audit"
	"github.com/transcom/mymove/pkg/services/query"
)

func payloadForMTOServiceItemModel(s *models.MTOServiceItem) *ghcmessages.MTOServiceItem {
	if s == nil {
		return nil
	}

	return &ghcmessages.MTOServiceItem{
		ID:              handlers.FmtUUID(s.ID),
		MoveTaskOrderID: handlers.FmtUUID(s.MoveTaskOrderID),
		MtoShipmentID:   handlers.FmtUUIDPtr(s.MTOShipmentID),
		ReServiceID:     handlers.FmtUUID(s.ReServiceID),
		ReServiceCode:   handlers.FmtString(string(s.ReService.Code)),
		ReServiceName:   handlers.FmtStringPtr(&s.ReService.Name),
	}
}

func payloadForMTOServiceItemModels(s models.MTOServiceItems) ghcmessages.MTOServiceItems {
	serviceItems := ghcmessages.MTOServiceItems{}
	for _, item := range s {
		serviceItems = append(serviceItems, payloadForMTOServiceItemModel(&item))
	}

	return serviceItems
}

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

// CreateMTOServiceItemHandler struct that describes creating a mto service item handler
type CreateMTOServiceItemHandler struct {
	handlers.HandlerContext
	services.MTOServiceItemCreator
}

// Handle handler that creates a mto service item
func (h CreateMTOServiceItemHandler) Handle(params mtoserviceitemop.CreateMTOServiceItemParams) middleware.Responder {
	var errs []string
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	moveTaskOrderID, err := uuid.FromString(params.MoveTaskOrderID.String())
	if err != nil {
		errs = append(errs, fmt.Errorf("UUID Parsing for %s: %w", "MoveTaskOrderID", err).Error())
	}

	reServiceID, err := uuid.FromString(params.CreateMTOServiceItemBody.ReServiceID.String())
	if err != nil {
		errs = append(errs, fmt.Errorf("UUID Parsing for %s: %w", "ReServiceID", err).Error())
	}

	mtoShipmentID, err := uuid.FromString(params.CreateMTOServiceItemBody.MtoShipmentID.String())
	if err != nil {
		errs = append(errs, fmt.Errorf("UUID Parsing for %s: %w", "MtoShipmentID", err).Error())
	}

	// return any parsing errors for uuids
	if len(errs) > 0 {
		parsingError := strings.Join(errs, "\n")
		logger.Error(parsingError)
		payload := payloadForValidationError("UUID(s) parsing error", parsingError, h.GetTraceID(), validate.NewErrors())

		return mtoserviceitemop.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(payload)
	}

	serviceItem := models.MTOServiceItem{
		MoveTaskOrderID: moveTaskOrderID,
		ReServiceID:     reServiceID,
		MTOShipmentID:   &mtoShipmentID,
	}

	// Capture creation attempt in audit log
	_, err = audit.Capture(&serviceItem, nil, logger, session, params.HTTPRequest)
	if err != nil {
		logger.Error("Auditing service error for service item creation.", zap.Error(err))
		return mtoserviceitemop.NewCreateMTOServiceItemInternalServerError()
	}

	createdServiceItem, verrs, err := h.MTOServiceItemCreator.CreateMTOServiceItem(&serviceItem)
	if verrs != nil && verrs.HasAny() {
		logger.Error("Error validating mto service item: ", zap.Error(verrs))
		payload := payloadForValidationError(handlers.ValidationErrMessage, "The information you provided is invalid.", h.GetTraceID(), verrs)

		return mtoserviceitemop.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(payload)
	}

	// return any errors
	if err != nil {
		logger.Error("Error creating mto service item: ", zap.Error(err))

		if strings.Contains(errors.Cause(err).Error(), models.ViolatesForeignKeyConstraint) {
			payload := payloadForClientError("Unknown UUID(s)", "Unknown UUID(s) used to create a mto service item.", h.GetTraceID())

			return mtoserviceitemop.NewCreateMTOServiceItemNotFound().WithPayload(payload)
		}

		return mtoserviceitemop.NewCreateMTOServiceItemInternalServerError()
	}

	returnPayload := payloadForMTOServiceItemModel(createdServiceItem)
	return mtoserviceitemop.NewCreateMTOServiceItemCreated().WithPayload(returnPayload)
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

	moveTaskOrder := &models.MoveTaskOrder{}
	err = h.Fetcher.FetchRecord(moveTaskOrder, queryFilters)
	if err != nil {
		logger.Error("Error fetching move task order: ", zap.Error(fmt.Errorf("Move Task Order ID: %s", moveTaskOrder.ID)), zap.Error(err))

		return mtoserviceitemop.NewListMTOServiceItemsNotFound()
	}

	queryFilters = []services.QueryFilter{
		query.NewQueryFilter("move_task_order_id", "=", moveTaskOrderID.String()),
	}
	queryAssociations := query.NewQueryAssociations([]services.QueryAssociation{
		query.NewQueryAssociation("ReService"),
	})

	var serviceItems models.MTOServiceItems
	err = h.ListFetcher.FetchRecordList(&serviceItems, queryFilters, queryAssociations, nil, nil)
	// return any errors
	if err != nil {
		logger.Error("Error fetching mto service items: ", zap.Error(err))

		return mtoserviceitemop.NewListMTOServiceItemsInternalServerError()
	}

	returnPayload := payloadForMTOServiceItemModels(serviceItems)
	return mtoserviceitemop.NewListMTOServiceItemsOK().WithPayload(returnPayload)
}
