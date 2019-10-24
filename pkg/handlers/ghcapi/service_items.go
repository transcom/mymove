package ghcapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	serviceitemop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/service_item"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
)

func payloadForServiceItemModel(s models.ServiceItem) *ghcmessages.ServiceItem {
	return &ghcmessages.ServiceItem{
		ID: handlers.FmtUUID(s.ID),
	}
}

type ListServiceItemsHandler struct {
	handlers.HandlerContext
	services.ServiceItemListFetcher
	services.NewQueryFilter
}

func (h ListServiceItemsHandler) Handle(params serviceitemop.ListServiceItemsParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	id, err := uuid.FromString(params.MoveTaskOrderID)

	if err != nil {
		logger.Error(fmt.Sprintf("UUID Parsing for %s", params.MoveTaskOrderID), zap.Error(err))
	}

	queryFilters := []services.QueryFilter{h.NewQueryFilter("move_task_order_id", "=", id)}
	pagination := pagination.NewPagination(nil, nil)
	associations := query.NewQueryAssociations([]services.QueryAssociation{})

	serviceItems, err := h.ServiceItemListFetcher.FetchServiceItemList(queryFilters, associations, pagination)

	if err != nil {
		logger.Error("Unable to fetch records:", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	payload := make(ghcmessages.ServiceItems, len(serviceItems))

	for i, s := range serviceItems {
		payload[i] = payloadForServiceItemModel(s)
	}

	return serviceitemop.NewListServiceItemsOK().WithPayload(payload)
}

type CreateServiceItemHandler struct {
	handlers.HandlerContext
	services.ServiceItemCreator
	services.NewQueryFilter
}

func (h CreateServiceItemHandler) Handle(params serviceitemop.CreateServiceItemParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	moveTaskOrderID, err := uuid.FromString(params.MoveTaskOrderID)
	if err != nil {
		logger.Error(fmt.Sprintf("UUID Parsing for %s", params.MoveTaskOrderID), zap.Error(err))
	}

	serviceItem := models.ServiceItem{
		MoveTaskOrderID: moveTaskOrderID,
	}

	transportationIDFilter := []services.QueryFilter{
		h.NewQueryFilter("id", "=", moveTaskOrderID),
	}

	createdServiceItem, verrs, err := h.ServiceItemCreator.CreateServiceItem(&serviceItem, transportationIDFilter)
	if err != nil || verrs != nil {
		logger.Error("Error saving service item", zap.Error(verrs))
		return serviceitemop.NewCreateServiceItemInternalServerError()
	}

	returnPayload := payloadForServiceItemModel(*createdServiceItem)
	return serviceitemop.NewCreateServiceItemCreated().WithPayload(returnPayload)
}
