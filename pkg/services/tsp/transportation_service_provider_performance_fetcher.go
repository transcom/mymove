package tsp

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type transportationServiceProviderPerformanceQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
}

type transportationServiceProviderPerformanceFetcher struct {
	builder transportationServiceProviderPerformanceQueryBuilder
}

// FetchTransportationServiceProviderPerformance fetches a transportation service provider performance given a slice of filters
func (o *transportationServiceProviderPerformanceFetcher) FetchTransportationServiceProviderPerformance(filters []services.QueryFilter) (models.TransportationServiceProviderPerformance, error) {
	var transportationServiceProviderPerformance models.TransportationServiceProviderPerformance
	error := o.builder.FetchOne(&transportationServiceProviderPerformance, filters)
	return transportationServiceProviderPerformance, error
}

// NewTransportationServiceProviderPerformanceFetcher return an implementation of the TransportationServiceProviderPerformanceFetcher interface
func NewTransportationServiceProviderPerformanceFetcher(builder transportationServiceProviderPerformanceQueryBuilder) services.TransportationServiceProviderPerformanceFetcher {
	return &transportationServiceProviderPerformanceFetcher{builder}
}
