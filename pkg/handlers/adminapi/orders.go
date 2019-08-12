package adminapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	orderop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/order"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func payloadForOrderModel(o models.Order) *adminmessages.Order {
	return &adminmessages.Order{
		ID:                  handlers.FmtUUID(o.ID),
		CreatedAt:           handlers.FmtDateTime(o.CreatedAt),
		UpdatedAt:           handlers.FmtDateTime(o.UpdatedAt),
		IssueDate:           handlers.FmtDate(o.IssueDate),
		ReportByDate:        handlers.FmtDate(o.ReportByDate),
		DepartmentIndicator: (*adminmessages.DeptIndicator)(o.DepartmentIndicator),
	}
}

type IndexOrdersHandler struct {
	handlers.HandlerContext
	services.OrderListFetcher
	services.NewQueryFilter
}

func (h IndexOrdersHandler) Handle(params orderop.IndexOrdersParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	queryFilters := []services.QueryFilter{}

	orders, err := h.OrderListFetcher.FetchOrderList(queryFilters)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	ordersCount := len(orders)
	payload := make(adminmessages.Orders, ordersCount)
	for i, s := range orders {
		payload[i] = payloadForOrderModel(s)
	}

	return orderop.NewIndexOrdersOK().WithContentRange(fmt.Sprintf("orders 0-%d/%d", ordersCount, ordersCount)).WithPayload(payload)
}
