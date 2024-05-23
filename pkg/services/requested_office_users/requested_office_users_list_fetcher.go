package adminuser

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type requestedOfficeUsersListQueryBuilder interface {
	FetchMany(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error
	Count(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) (int, error)
}

type requestedOfficeUserListFetcher struct {
	builder requestedOfficeUsersListQueryBuilder
}

// FetchAdminUserList uses the passed query builder to fetch a list of office users
func (o *requestedOfficeUserListFetcher) FetchRequestedOfficeUsersList(appCtx appcontext.AppContext, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) (models.OfficeUsers, error) {
	var requestedUsers models.OfficeUsers
	err := o.builder.FetchMany(appCtx, &requestedUsers, filters, associations, pagination, ordering)
	return requestedUsers, err
}

// FetchAdminUserList uses the passed query builder to fetch a list of office users
func (o *requestedOfficeUserListFetcher) FetchRequestedOfficeUsersCount(appCtx appcontext.AppContext, filters []services.QueryFilter) (int, error) {
	var requestedUsers models.OfficeUsers
	count, err := o.builder.Count(appCtx, &requestedUsers, filters)
	return count, err
}

// NewAdminUserListFetcher returns an implementation of AdminUserListFetcher
func NewRequestedOfficeUsersListFetcher(builder requestedOfficeUsersListQueryBuilder) services.RequestedOfficeUserListFetcher {
	return &requestedOfficeUserListFetcher{builder}
}
