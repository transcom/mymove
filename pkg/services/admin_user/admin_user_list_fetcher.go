package adminuser

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type adminUserListQueryBuilder interface {
	FetchMany(model interface{}, filters []services.QueryFilter, pagination services.Pagination) error
}

type adminUserListFetcher struct {
	builder adminUserListQueryBuilder
}

// FetchAdminUserList uses the passed query builder to fetch a list of office users
func (o *adminUserListFetcher) FetchAdminUserList(filters []services.QueryFilter, pagination services.Pagination) (models.AdminUsers, error) {
	var adminUsers models.AdminUsers
	error := o.builder.FetchMany(&adminUsers, filters, pagination)
	return adminUsers, error
}

// NewAdminUserListFetcher returns an implementation of AdminUserListFetcher
func NewAdminUserListFetcher(builder adminUserListQueryBuilder) services.AdminUserListFetcher {
	return &adminUserListFetcher{builder}
}
