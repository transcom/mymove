package publicapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	tsppop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/transportation_service_provider_performance"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

// LogTransportationServiceProviderPerformanceHandler logs a TSPP record (for auditing purposes)
type LogTransportationServiceProviderPerformanceHandler struct {
	handlers.HandlerContext
	services.NewQueryFilter
	services.TransportationServiceProviderPerformanceFetcher
}

// Handle logging the TSPP record
func (h LogTransportationServiceProviderPerformanceHandler) Handle(params tsppop.LogTransportationServiceProviderPerformanceParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	tsppID, _ := uuid.FromString(params.TransportationServiceProviderPerformanceID.String())

	if !session.IsOfficeUser() {
		return tsppop.NewLogTransportationServiceProviderPerformanceForbidden()
	}

	queryFilters := []services.QueryFilter{
		h.NewQueryFilter("id", "=", tsppID.String()),
	}
	tspp, err := h.TransportationServiceProviderPerformanceFetcher.FetchTransportationServiceProviderPerformance(queryFilters)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			logger.Info("No record found for TSPP ID", zap.String("id", tsppID.String()))
		} else {
			logger.Error("DB Query", zap.Error(err))
			return tsppop.NewLogTransportationServiceProviderPerformanceInternalServerError()
		}
	} else {
		logger.Info("Record found for TSPP ID", zap.Object("TSPP", &tspp))
	}

	return tsppop.NewLogTransportationServiceProviderPerformanceNoContent()
}
