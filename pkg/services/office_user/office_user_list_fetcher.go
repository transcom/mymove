package officeuser

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type officeUserListQueryBuilder interface {
	query.FetchMany
}

type officeUserListFetcher struct {
	builder officeUserListQueryBuilder
}

// FetchOfficeUserList uses the passed query builder to fetch a list of office users
func (o *officeUserListFetcher) FetchOfficeUserList(filters []services.QueryFilter, pagination services.Pagination) (models.OfficeUsers, error) {
	var officeUsers models.OfficeUsers
	err := o.builder.WithModel(&officeUsers).WithFilters(filters).WithPagination(pagination).Execute()
	return officeUsers, err
}

// NewOfficeUserListFetcher returns an implementation of OfficeUserListFetcher
func NewOfficeUserListFetcher(builder officeUserListQueryBuilder) services.OfficeUserListFetcher {
	return &officeUserListFetcher{builder}
}
