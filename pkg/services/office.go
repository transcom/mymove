package services

import (
	"github.com/transcom/mymove/pkg/models"
)

// OfficeFetcher is the exported interface for fetching a single transportation office
type OfficeFetcher interface {
	FetchOffice(filters []QueryFilter) (models.TransportationOffice, error)
}

// OfficeListFetcher is the exported interface for fetching multiple transportation offices
//go:generate mockery --name OfficeListFetcher --disable-version-string
type OfficeListFetcher interface {
	FetchOfficeList(filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) (models.TransportationOffices, error)
	FetchOfficeCount(filters []QueryFilter) (int, error)
}
