package officeuser

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type officeUserListQueryBuilder interface {
	FetchMany(model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination) error
}

type officeUserListFetcher struct {
	builder officeUserListQueryBuilder
}

// FetchOfficeUserList uses the passed query builder to fetch a list of office users
func (o *officeUserListFetcher) FetchOfficeUserList(filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination) (models.OfficeUsers, error) {
	var officeUsers models.OfficeUsers
	error := o.builder.FetchMany(&officeUsers, filters, associations, pagination)
	return officeUsers, error
}

// NewOfficeUserListFetcher returns an implementation of OfficeUserListFetcher
func NewOfficeUserListFetcher(builder officeUserListQueryBuilder) services.OfficeUserListFetcher {
	return &officeUserListFetcher{builder}
}
