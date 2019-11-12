package officeuser

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type officeUserListQueryBuilder interface {
	query.FetchMany
	Count(model interface{}, filters []services.QueryFilter) (int, error)
}

type officeUserListFetcher struct {
	builder officeUserListQueryBuilder
}

// FetchOfficeUserList uses the passed query builder to fetch a list of office users
func (o *officeUserListFetcher) FetchOfficeUserList(filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) (models.OfficeUsers, error) {
	var officeUsers models.OfficeUsers
	err := o.builder.WithModel(&officeUsers).WithFilters(filters).WithPagination(pagination).Execute()
	return officeUsers, err
}

// FetchOfficeUserCount uses the passed query builder to count office users
func (o *officeUserListFetcher) FetchOfficeUserCount(filters []services.QueryFilter) (int, error) {
	var officeUsers models.OfficeUsers
	count, error := o.builder.Count(&officeUsers, filters)
	return count, error
}

// NewOfficeUserListFetcher returns an implementation of OfficeUserListFetcher
func NewOfficeUserListFetcher(builder officeUserListQueryBuilder) services.OfficeUserListFetcher {
	return &officeUserListFetcher{builder}
}
