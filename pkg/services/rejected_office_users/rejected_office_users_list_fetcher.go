package adminuser

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type rejectedOfficeUsersListQueryBuilder interface {
	FetchMany(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error
	Count(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) (int, error)
}

type rejectedOfficeUserListFetcher struct {
	builder rejectedOfficeUsersListQueryBuilder
}

// FetchAdminUserList uses the passed query builder to fetch a list of office users
func (o *rejectedOfficeUserListFetcher) FetchRejectedOfficeUsersList(appCtx appcontext.AppContext, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) (models.OfficeUsers, error) {
	var rejectedUsers models.OfficeUsers
	err := o.builder.FetchMany(appCtx, &rejectedUsers, filters, associations, pagination, ordering)
	return rejectedUsers, err
}

// FetchAdminUserList uses the passed query builder to fetch a list of office users
func (o *rejectedOfficeUserListFetcher) FetchRejectedOfficeUsersCount(appCtx appcontext.AppContext, filters []services.QueryFilter) (int, error) {
	var rejectedUsers models.OfficeUsers
	count, err := o.builder.Count(appCtx, &rejectedUsers, filters)
	return count, err
}

// NewAdminUserListFetcher returns an implementation of AdminUserListFetcher
func NewRejectedOfficeUsersListFetcher(builder rejectedOfficeUsersListQueryBuilder) services.RejectedOfficeUserListFetcher {
	return &rejectedOfficeUserListFetcher{builder}
}
