package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// RejectedOfficeUserListFetcher is the exported interface for fetching multiple rejected rejected office users
//
//go:generate mockery --name RejectedOfficeUserListFetcher
type RejectedOfficeUserListFetcher interface {
	FetchRejectedOfficeUsersList(appCtx appcontext.AppContext, filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) (models.OfficeUsers, error)
	FetchRejectedOfficeUsersCount(appCtx appcontext.AppContext, filters []QueryFilter) (int, error)
}

// RejectedOfficeUserFetcher is the exported interface for fetching a single rejected rejected office user
//
//go:generate mockery --name RejectedOfficeUserFetcher
type RejectedOfficeUserFetcher interface {
	FetchRejectedOfficeUser(appCtx appcontext.AppContext, filters []QueryFilter) (models.OfficeUser, error)
}
