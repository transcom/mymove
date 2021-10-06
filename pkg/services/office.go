package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// OfficeFetcher is the exported interface for fetching a single transportation office
type OfficeFetcher interface {
	FetchOffice(appCtx appcontext.AppContext, filters []QueryFilter) (models.TransportationOffice, error)
}

// OfficeListFetcher is the exported interface for fetching multiple transportation offices
//go:generate mockery --name OfficeListFetcher --disable-version-string
type OfficeListFetcher interface {
	FetchOfficeList(appCtx appcontext.AppContext, filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) (models.TransportationOffices, error)
	FetchOfficeCount(appCtx appcontext.AppContext, filters []QueryFilter) (int, error)
}
