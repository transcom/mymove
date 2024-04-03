package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// RequestedOfficeUserListFetcher is the exported interface for fetching multiple requested office users
//
//go:generate mockery --name RequestedOfficeUserListFetcher
type RequestedOfficeUserListFetcher interface {
	FetchRequestedOfficeUsersList(appCtx appcontext.AppContext, filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) (models.OfficeUsers, error)
	FetchRequestedOfficeUsersCount(appCtx appcontext.AppContext, filters []QueryFilter) (int, error)
}

// RequestedOfficeUserFetcher is the exported interface for fetching a single requested office user
//
//go:generate mockery --name RequestedOfficeUserFetcher
type RequestedOfficeUserFetcher interface {
	FetchRequestedOfficeUser(appCtx appcontext.AppContext, filters []QueryFilter) (models.OfficeUser, error)
}
