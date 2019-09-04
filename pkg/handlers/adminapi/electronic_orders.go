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
}

func (h IndexElectronicOrdersHandler) Handle(params electronicorderop.IndexElectronicOrdersParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	// queryFilters := []services.QueryFilter{}

	// electronicOrders, err := h.ElectronicOrderListFetcher.FetchElectronicOrderList(queryFilters)
	// TODO: Remove when we ship pagination in the query builder
	query := `SELECT id, issuer, created_at, updated_at from electronic_orders LIMIT 100`
	electronicOrders := models.ElectronicOrders{}
	err := h.DB().RawQuery(query).All(&electronicOrders)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	electronicOrdersCount := len(electronicOrders)
	payload := make(adminmessages.ElectronicOrders, electronicOrdersCount)
	for i, s := range electronicOrders {
		payload[i] = payloadForElectronicOrderModel(s)
	}

	return electronicorderop.NewIndexElectronicOrdersOK().WithContentRange(fmt.Sprintf("electronic_orders 0-%d/%d", electronicOrdersCount, electronicOrdersCount)).WithPayload(payload)
}
