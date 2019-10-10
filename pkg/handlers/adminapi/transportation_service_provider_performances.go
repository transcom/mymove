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
	return &adminmessages.TransportationServiceProviderPerformance{
		ID:                              handlers.FmtUUID(o.ID),
		TrafficDistributionListID:       *handlers.FmtUUID(o.TrafficDistributionListID),
		TransportationServiceProviderID: *handlers.FmtUUID(o.TransportationServiceProviderID),
		PerformancePeriodStart:          *handlers.FmtDateTime(o.PerformancePeriodStart),
		PerformancePeriodEnd:            *handlers.FmtDateTime(o.PerformancePeriodEnd),
		RateCycleStart:                  *handlers.FmtDateTime(o.RateCycleStart),
		RateCycleEnd:                    *handlers.FmtDateTime(o.RateCycleEnd),
		QualityBand:                     handlers.FmtIntPtrToInt64(o.QualityBand),
		OfferCount:                      int64(o.OfferCount),
		BestValueScore:                  o.BestValueScore,
		LinehaulRate:                    o.LinehaulRate.Float64(),
		SITRate:                         o.SITRate.Float64(),
	}
}

// IndexTSPPsHandler returns a list of office users via GET /transportation_service_provider_performances
type IndexTSPPsHandler struct {
	handlers.HandlerContext
	services.TransportationServiceProviderPerformanceListFetcher
	services.NewQueryFilter
	services.NewPagination
}

// Handle retrieves a list of office users
func (h IndexTSPPsHandler) Handle(params tsppop.IndexTSPPsParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	// Here is where NewQueryFilter will be used to create Filters from the 'filter' query param
	queryFilters := []services.QueryFilter{}

	pagination := h.NewPagination(params.Page, params.PerPage)
	associations := query.NewQueryAssociations([]services.QueryAssociation{})

	tspps, err := h.TransportationServiceProviderPerformanceListFetcher.FetchTransportationServiceProviderPerformanceList(queryFilters, associations, pagination)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	totalTSPPsCount, err := h.DB().Count(&models.TransportationServiceProviderPerformance{})
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
