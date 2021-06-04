package adminapi

import (
	"fmt"

	"github.com/transcom/mymove/pkg/services/query"

	"github.com/go-openapi/runtime/middleware"

	tsppop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/transportation_service_provider_performances"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func payloadForTSPPModel(o models.TransportationServiceProviderPerformance) *adminmessages.TransportationServiceProviderPerformance {
	lhRate := o.LinehaulRate.Float64()
	sitRate := o.SITRate.Float64()

	return &adminmessages.TransportationServiceProviderPerformance{
		ID:                              handlers.FmtUUID(o.ID),
		TrafficDistributionListID:       handlers.FmtUUID(o.TrafficDistributionListID),
		TransportationServiceProviderID: handlers.FmtUUID(o.TransportationServiceProviderID),
		PerformancePeriodStart:          handlers.FmtDateTime(o.PerformancePeriodStart),
		PerformancePeriodEnd:            handlers.FmtDateTime(o.PerformancePeriodEnd),
		RateCycleStart:                  handlers.FmtDateTime(o.RateCycleStart),
		RateCycleEnd:                    handlers.FmtDateTime(o.RateCycleEnd),
		QualityBand:                     handlers.FmtIntPtrToInt64(o.QualityBand),
		OfferCount:                      handlers.FmtIntPtrToInt64(&o.OfferCount),
		BestValueScore:                  &o.BestValueScore,
		LinehaulRate:                    &lhRate,
		SitRate:                         &sitRate,
	}
}

// IndexTSPPsHandler returns a list of transportation service provider performance via GET /transportation_service_provider_performances
type IndexTSPPsHandler struct {
	handlers.HandlerContext
	services.TransportationServiceProviderPerformanceListFetcher
	services.NewQueryFilter
	services.NewPagination
}

var tsppFilterConverters = map[string]func(string) []services.QueryFilter{
	"traffic_distribution_list_id": func(content string) []services.QueryFilter {
		return []services.QueryFilter{query.NewQueryFilter("traffic_distribution_list_id", "=", content)}
	},
	"transportation_service_provider_id": func(content string) []services.QueryFilter {
		return []services.QueryFilter{query.NewQueryFilter("transportation_service_provider_id", "=", content)}
	},
}

// Handle retrieves a list of transportation service provider performance
func (h IndexTSPPsHandler) Handle(params tsppop.IndexTSPPsParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	// Here is where NewQueryFilter will be used to create Filters from the 'filter' query param
	queryFilters := generateQueryFilters(logger, params.Filter, tsppFilterConverters)

	pagination := h.NewPagination(params.Page, params.PerPage)
	ordering := query.NewQueryOrder(params.Sort, params.Order)

	tspps, err := h.TransportationServiceProviderPerformanceListFetcher.FetchTransportationServiceProviderPerformanceList(queryFilters, nil, pagination, ordering)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	totalTSPPsCount, err := h.TransportationServiceProviderPerformanceListFetcher.FetchTransportationServiceProviderPerformanceCount(queryFilters)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	queriedTSPPsCount := len(tspps)

	payload := make(adminmessages.TransportationServiceProviderPerformances, queriedTSPPsCount)
	for i, s := range tspps {
		payload[i] = payloadForTSPPModel(s)
	}

	return tsppop.NewIndexTSPPsOK().WithContentRange(fmt.Sprintf("tspps %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedTSPPsCount, totalTSPPsCount)).WithPayload(payload)
}

// GetTSPPHandler returns a transportation service provider performance via GET /transportation_service_provider_performances/{tspId}
type GetTSPPHandler struct {
	handlers.HandlerContext
	services.TransportationServiceProviderPerformanceFetcher
	services.NewQueryFilter
}

// Handle returns the payload for TSPP
func (h GetTSPPHandler) Handle(params tsppop.GetTSPPParams) middleware.Responder {
	_, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	tsppID := params.TsppID

	queryFilters := []services.QueryFilter{query.NewQueryFilter("id", "=", tsppID)}

	tspp, err := h.TransportationServiceProviderPerformanceFetcher.FetchTransportationServiceProviderPerformance(queryFilters)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	payload := payloadForTSPPModel(tspp)

	return tsppop.NewGetTSPPOK().WithPayload(payload)
}
