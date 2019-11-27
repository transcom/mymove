package fetch

import (
	"github.com/transcom/mymove/pkg/services"
)

type listQueryBuilder interface {
	FetchMany(model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error
	Count(model interface{}, filters []services.QueryFilter) (int, error)
}

type listFetcher struct {
	builder listQueryBuilder
}

// FetchRecordList uses the passed query builder to fetch a list of records
func (o *listFetcher) FetchRecordList(model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error {
	error := o.builder.FetchMany(model, filters, associations, pagination, ordering)
	return error
}

// FetchRecordCount uses the passed query builder to count records
func (o *listFetcher) FetchRecordCount(model interface{}, filters []services.QueryFilter) (int, error) {
	count, error := o.builder.Count(model, filters)
	return count, error
}

// NewListFetcher returns an implementation of ListFetcher
func NewListFetcher(builder listQueryBuilder) services.ListFetcher {
	return &listFetcher{builder}
}
