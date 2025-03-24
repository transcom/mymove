package services

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// RejectedOfficeUserListFetcher is the exported interface for fetching multiple rejected office users
//
//go:generate mockery --name RejectedOfficeUserListFetcher
type RejectedOfficeUserListFetcher interface {
	FetchRejectedOfficeUsersList(appCtx appcontext.AppContext, filterFuncs []func(*pop.Query), pagination Pagination, ordering QueryOrder) (models.OfficeUsers, int, error)
	FetchRejectedOfficeUsersCount(appCtx appcontext.AppContext, filters []QueryFilter) (int, error)
}

// RejectedOfficeUserFetcher is the exported interface for fetching a single rejected office user
//
//go:generate mockery --name RejectedOfficeUserFetcher
type RejectedOfficeUserFetcher interface {
	FetchRejectedOfficeUser(appCtx appcontext.AppContext, filters []QueryFilter) (models.OfficeUser, error)
}
