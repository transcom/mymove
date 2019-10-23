package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
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
	queryFilters := []services.QueryFilter{h.NewQueryFilter("move_task_order_id", "=", params.MoveTaskOrderID)}
	pagination := pagination.NewPagination(nil, nil)
	queryAssociations := []services.QueryAssociation{
		query.NewQueryAssociation(""),
	}
	associations := query.NewQueryAssociations(queryAssociations)

	serviceItems, err := h.ServiceItemListFetcher.FetchServiceItemList(queryFilters, associations, pagination)

	if err != nil {
		return serviceitemop.NewListServiceItemsInternalServerError()
	}

	payload := make(ghcmessages.ServiceItems, len(serviceItems))

	for i, s := range serviceItems {
		payload[i] = payloadForServiceItemModel(s)
	}

	return serviceitemop.NewListServiceItemsOK().WithPayload(payload)
}
