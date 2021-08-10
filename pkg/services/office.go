package services

import (
	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
)

// OfficeFetcher is the exported interface for fetching a single transportation office
type OfficeFetcher interface {
	FetchOffice(appCfg appconfig.AppConfig, filters []QueryFilter) (models.TransportationOffice, error)
}

// OfficeListFetcher is the exported interface for fetching multiple transportation offices
//go:generate mockery --name OfficeListFetcher --disable-version-string
type OfficeListFetcher interface {
	FetchOfficeList(appCfg appconfig.AppConfig, filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) (models.TransportationOffices, error)
	FetchOfficeCount(appCfg appconfig.AppConfig, filters []QueryFilter) (int, error)
}
