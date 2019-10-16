package adminuser

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type adminUserQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
}

type adminUserFetcher struct {
	builder adminUserQueryBuilder
}

// FetchAdminUser fetches an admin user given a slice of filters
func (o *adminUserFetcher) FetchAdminUser(filters []services.QueryFilter) (models.AdminUser, error) {
	var adminUser models.AdminUser
	error := o.builder.FetchOne(&adminUser, filters)
	return adminUser, error
}

// NewAdminUserFetcher return an implementation of the AdminUserFetcher interface
func NewAdminUserFetcher(builder adminUserQueryBuilder) services.AdminUserFetcher {
	return &adminUserFetcher{builder}
}
