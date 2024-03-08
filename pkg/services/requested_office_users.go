package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// RequestedOfficeUserListFetcher is the exported interface for fetching multiple admin users
//
//go:generate mockery --name RequestedOfficeUserListFetcher
type RequestedOfficeUserListFetcher interface {
	FetchRequestedOfficeUsersList(appCtx appcontext.AppContext, filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) (models.OfficeUsers, error)
	FetchRequestedOfficeUsersCount(appCtx appcontext.AppContext, filters []QueryFilter) (int, error)
}
