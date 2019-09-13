package adminapi

import (
	"fmt"
	"strings"

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

type GetElectronicOrdersTotalsHandler struct {
	handlers.HandlerContext
	services.ElectronicOrderCategoryCountFetcher
	services.NewQueryFilter
}

func split(r rune) bool {
	return r == '.' || r == ':'
}

func (h GetElectronicOrdersTotalsHandler) Handle(params electronicorderop.GetElectronicOrdersTotalsParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	queryFilters := []services.QueryFilter{}
	andQueryFilters := []services.QueryFilter{}

	for _, filter := range params.Filter {
		queryFilterSplit := strings.FieldsFunc(filter, split)
		queryFilters = append(queryFilters, h.NewQueryFilter(queryFilterSplit[0], queryFilterSplit[1], queryFilterSplit[2]))
	}

	if params.AndFilter != nil {
		for _, andFilter := range params.AndFilter {
			andFilterSplit := strings.FieldsFunc(andFilter, split)
			andQueryFilters = append(andQueryFilters, h.NewQueryFilter(andFilterSplit[0], andFilterSplit[1], andFilterSplit[2]))
		}
	}

	counts, err := h.ElectronicOrderCategoryCountFetcher.FetchElectronicOrderCategoricalCounts(queryFilters, &andQueryFilters)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	fmt.Println(counts)

	return electronicorderop.NewGetElectronicOrdersTotalsOK()
}