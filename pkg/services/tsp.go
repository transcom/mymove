package services

import (
	"github.com/transcom/mymove/pkg/models"
)

// TransportationServiceProviderPerformanceFetcher is the exported interface for fetching
// a single transportation service provider performance
//go:generate mockery -name TransportationServiceProviderPerformanceFetcher
type TransportationServiceProviderPerformanceFetcher interface {
	FetchTransportationServiceProviderPerformance(filters []QueryFilter) (models.TransportationServiceProviderPerformance, error)
}

// TransportationServiceProviderPerformanceListFetcher is the exported interface for fetching
// a list of transportation service provider performances
//go:generate mockery -name TransportationServiceProviderPerformanceListFetcher
type TransportationServiceProviderPerformanceListFetcher interface {
	FetchTransportationServiceProviderPerformanceList(filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) (models.TransportationServiceProviderPerformances, error)
	FetchTransportationServiceProviderPerformanceCount(filters []QueryFilter) (int, error)
}
