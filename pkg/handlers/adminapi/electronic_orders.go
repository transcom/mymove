package adminapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	electronicorderop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/electronic_order"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func payloadForElectronicOrderModel(o models.ElectronicOrder) *adminmessages.ElectronicOrder {
	return &adminmessages.ElectronicOrder{
		ID:        handlers.FmtUUID(o.ID),
		Issuer:    adminmessages.Issuer(o.Issuer),
		CreatedAt: handlers.FmtDateTime(o.CreatedAt),
		UpdatedAt: handlers.FmtDateTime(o.UpdatedAt),
	}
}

type IndexElectronicOrdersHandler struct {
	handlers.HandlerContext
	services.ElectronicOrderListFetcher
	services.NewQueryFilter
	services.NewPagination
}

func (h IndexElectronicOrdersHandler) Handle(params electronicorderop.IndexElectronicOrdersParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	queryFilters := []services.QueryFilter{}

	var pagination services.Pagination
	if params.Page == nil {
		pagination = h.NewPagination(1, 25) // default number of records per page
	} else {
		page, perPage := *params.Page, *params.PerPage
		pagination = h.NewPagination(page, perPage)
	}

	electronicOrders, err := h.ElectronicOrderListFetcher.FetchElectronicOrderList(queryFilters, pagination)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	totalElectronicOrdersCount, err := h.DB().Count(&models.ElectronicOrder{})
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	queriedOfficeUsersCount := len(electronicOrders)

	payload := make(adminmessages.ElectronicOrders, queriedOfficeUsersCount)
	for i, s := range electronicOrders {
		payload[i] = payloadForElectronicOrderModel(s)
	}

	return electronicorderop.NewIndexElectronicOrdersOK().WithContentRange(fmt.Sprintf("electronic_orders %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedOfficeUsersCount, totalElectronicOrdersCount)).WithPayload(payload)
}
