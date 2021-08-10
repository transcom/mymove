package services

import (
	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
)

// TransportationServiceProviderPerformanceFetcher is the exported interface for fetching
// a single transportation service provider performance
//go:generate mockery --name TransportationServiceProviderPerformanceFetcher --disable-version-string
type TransportationServiceProviderPerformanceFetcher interface {
	FetchTransportationServiceProviderPerformance(appCfg appconfig.AppConfig, filters []QueryFilter) (models.TransportationServiceProviderPerformance, error)
}

// TransportationServiceProviderPerformanceListFetcher is the exported interface for fetching
// a list of transportation service provider performances
//go:generate mockery --name TransportationServiceProviderPerformanceListFetcher --disable-version-string
type TransportationServiceProviderPerformanceListFetcher interface {
	FetchTransportationServiceProviderPerformanceList(appCfg appconfig.AppConfig, filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) (models.TransportationServiceProviderPerformances, error)
	FetchTransportationServiceProviderPerformanceCount(appCfg appconfig.AppConfig, filters []QueryFilter) (int, error)
}
