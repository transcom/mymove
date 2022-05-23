package adminapi

import (
	"fmt"
	"strings"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/go-openapi/runtime/middleware"

	electronicorderop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/electronic_order"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func payloadForElectronicOrderModel(o models.ElectronicOrder) *adminmessages.ElectronicOrder {
	return &adminmessages.ElectronicOrder{
		ID:           handlers.FmtUUID(o.ID),
		Issuer:       adminmessages.NewIssuer(adminmessages.Issuer(o.Issuer)),
		OrdersNumber: handlers.FmtString(o.OrdersNumber),
		CreatedAt:    handlers.FmtDateTime(o.CreatedAt),
		UpdatedAt:    handlers.FmtDateTime(o.UpdatedAt),
	}
}

// IndexElectronicOrdersHandler returns an index of electronic orders
type IndexElectronicOrdersHandler struct {
	handlers.HandlerConfig
	services.ElectronicOrderListFetcher
	services.NewQueryFilter
	services.NewPagination
}

// Handle returns an index of electronic orders
func (h IndexElectronicOrdersHandler) Handle(params electronicorderop.IndexElectronicOrdersParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			queryFilters := []services.QueryFilter{}

			pagination := h.NewPagination(params.Page, params.PerPage)
			ordering := query.NewQueryOrder(params.Sort, params.Order)

			electronicOrders, err := h.ElectronicOrderListFetcher.FetchElectronicOrderList(appCtx, queryFilters, nil, pagination, ordering)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			totalElectronicOrdersCount, err := h.ElectronicOrderListFetcher.FetchElectronicOrderCount(appCtx, queryFilters)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			queriedOfficeUsersCount := len(electronicOrders)

			payload := make(adminmessages.ElectronicOrders, queriedOfficeUsersCount)
			for i, s := range electronicOrders {
				payload[i] = payloadForElectronicOrderModel(s)
			}

			return electronicorderop.NewIndexElectronicOrdersOK().WithContentRange(fmt.Sprintf("electronic_orders %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedOfficeUsersCount, totalElectronicOrdersCount)).WithPayload(payload), nil
		})
}

// GetElectronicOrdersTotalsHandler returns totals of electronic orders
type GetElectronicOrdersTotalsHandler struct {
	handlers.HandlerConfig
	services.ElectronicOrderCategoryCountFetcher
	services.NewQueryFilter
}

func split(r rune) bool {
	return r == '.' || r == ':'
}

func translateComparator(s string) string {
	s = strings.ToLower(s)
	switch s {
	case "eq":
		return "="
	case "gt":
		return ">"
	case "lt":
		return "<"
	case "neq":
		return "!="
	case "lte":
		return "<="
	case "gte":
		return ">="
	}
	return s
}

// Handle returns electronic orders totals
func (h GetElectronicOrdersTotalsHandler) Handle(params electronicorderop.GetElectronicOrdersTotalsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			comparator := ""

			andQueryFilters := make([]services.QueryFilter, len(params.AndFilter))
			queryFilters := make([]services.QueryFilter, len(params.Filter))

			// Default behavior for this handler is going to be returning counts for each of the component services as categories
			if len(params.Filter) == 0 {
				queryFilters = []services.QueryFilter{
					h.NewQueryFilter("issuer", "=", models.IssuerAirForce),
					h.NewQueryFilter("issuer", "=", models.IssuerArmy),
					h.NewQueryFilter("issuer", "=", models.IssuerCoastGuard),
					h.NewQueryFilter("issuer", "=", models.IssuerNavy),
					h.NewQueryFilter("issuer", "=", models.IssuerMarineCorps),
				}
			} else {
				for i, filter := range params.Filter {
					queryFilterSplit := strings.FieldsFunc(filter, split)
					comparator = translateComparator(queryFilterSplit[1])
					queryFilters[i] = h.NewQueryFilter(queryFilterSplit[0], comparator, queryFilterSplit[2])
				}
			}

			if params.AndFilter != nil {
				for i, andFilter := range params.AndFilter {
					andFilterSplit := strings.FieldsFunc(andFilter, split)
					comparator = translateComparator(andFilterSplit[1])
					andQueryFilters[i] = h.NewQueryFilter(andFilterSplit[0], comparator, andFilterSplit[2])
				}
			}

			counts, err := h.ElectronicOrderCategoryCountFetcher.FetchElectronicOrderCategoricalCounts(appCtx, queryFilters, &andQueryFilters)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			payload := adminmessages.ElectronicOrdersTotals{}
			for key, count := range counts {
				count64 := int64(count)
				stringKey := fmt.Sprintf("%v", key)
				totalCount := adminmessages.ElectronicOrdersTotal{
					Category: stringKey,
					Count:    &count64,
				}
				payload = append(payload, &totalCount)
			}

			return electronicorderop.NewGetElectronicOrdersTotalsOK().WithPayload(payload), nil
		})
}
