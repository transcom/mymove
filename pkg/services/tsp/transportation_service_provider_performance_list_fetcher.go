package tsp

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type transportationServiceProviderPerformanceListQueryBuilder interface {
	FetchMany(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error
	Count(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) (int, error)
}

type transportationServiceProviderPerformanceListFetcher struct {
	builder transportationServiceProviderPerformanceListQueryBuilder
}

// FetchTransportationServiceProviderPerformanceList fetches a transportation service provider performance given a slice of filters
func (o *transportationServiceProviderPerformanceListFetcher) FetchTransportationServiceProviderPerformanceList(appCtx appcontext.AppContext, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) (models.TransportationServiceProviderPerformances, error) {
	var tspps models.TransportationServiceProviderPerformances
	error := o.builder.FetchMany(appCtx, &tspps, filters, associations, pagination, ordering)
	return tspps, error
}

// FetchTransportationServiceProviderPerformanceCount counts the transportation service provider performance given a slice of filters
func (o *transportationServiceProviderPerformanceListFetcher) FetchTransportationServiceProviderPerformanceCount(appCtx appcontext.AppContext, filters []services.QueryFilter) (int, error) {
	var tspps models.TransportationServiceProviderPerformances
	count, error := o.builder.Count(appCtx, &tspps, filters)
	return count, error
}

// NewTransportationServiceProviderPerformanceListFetcher return an implementation of the TransportationServiceProviderPerformanceFetcher interface
func NewTransportationServiceProviderPerformanceListFetcher(builder transportationServiceProviderPerformanceListQueryBuilder) services.TransportationServiceProviderPerformanceListFetcher {
	return &transportationServiceProviderPerformanceListFetcher{builder}
}
