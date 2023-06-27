package adminuser

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type adminUserListQueryBuilder interface {
	FetchMany(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error
	Count(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) (int, error)
}

type adminUserListFetcher struct {
	builder adminUserListQueryBuilder
}

// FetchAdminUserList uses the passed query builder to fetch a list of office users
func (o *adminUserListFetcher) FetchAdminUserList(appCtx appcontext.AppContext, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) (models.AdminUsers, error) {
	var adminUsers models.AdminUsers
	err := o.builder.FetchMany(appCtx, &adminUsers, filters, associations, pagination, ordering)
	return adminUsers, err
}

// FetchAdminUserList uses the passed query builder to fetch a list of office users
func (o *adminUserListFetcher) FetchAdminUserCount(appCtx appcontext.AppContext, filters []services.QueryFilter) (int, error) {
	var adminUsers models.AdminUsers
	count, err := o.builder.Count(appCtx, &adminUsers, filters)
	return count, err
}

// NewAdminUserListFetcher returns an implementation of AdminUserListFetcher
func NewAdminUserListFetcher(builder adminUserListQueryBuilder) services.AdminUserListFetcher {
	return &adminUserListFetcher{builder}
}
