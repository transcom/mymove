package adminapi

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

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

// Handle retrieves a list of transportation service provider performance
func (h IndexTSPPsHandler) Handle(params tsppop.IndexTSPPsParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	// Here is where NewQueryFilter will be used to create Filters from the 'filter' query param
	queryFilters := h.generateQueryFilters(params.Filter, logger)

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

// generateQueryFilters is helper to convert filter params from a json string
// of the form `{"traffic_distribution_list_id": "8e4b3caf-98dc-462a-bbcc-1977d08a08eb" "transportation_service_provider_id": "8e4b3caf-98dc-462a-bbcc-1977d08a08eb"}`
// to an array of services.QueryFilter
func (h IndexTSPPsHandler) generateQueryFilters(filters *string, logger handlers.Logger) []services.QueryFilter {
	type Filter struct {
		TdlID string `json:"traffic_distribution_list_id"`
		TspID string `json:"transportation_service_provider_id"`
	}
	f := Filter{}
	var queryFilters []services.QueryFilter
	if filters == nil {
		return queryFilters
	}
	b := []byte(*filters)
	err := json.Unmarshal(b, &f)
	if err != nil {
		fs := fmt.Sprintf("%v", filters)
		logger.Warn("unable to decode param", zap.Error(err),
			zap.String("filters", fs))
	}
	if f.TdlID != "" {
		queryFilters = append(queryFilters, query.NewQueryFilter("traffic_distribution_list_id", "=", f.TdlID))
	}
	if f.TspID != "" {
		queryFilters = append(queryFilters, query.NewQueryFilter("transportation_service_provider_id", "=", f.TspID))
	}
	return queryFilters
}
