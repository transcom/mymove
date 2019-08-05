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
