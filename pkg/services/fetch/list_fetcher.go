package fetch

import (
	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/services"
)

type listQueryBuilder interface {
	FetchMany(appCfg appconfig.AppConfig, model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error
	Count(appCfg appconfig.AppConfig, model interface{}, filters []services.QueryFilter) (int, error)
}

type listFetcher struct {
	builder listQueryBuilder
}

// FetchRecordList uses the passed query builder to fetch a list of records
func (o *listFetcher) FetchRecordList(appCfg appconfig.AppConfig, model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error {
	error := o.builder.FetchMany(appCfg, model, filters, associations, pagination, ordering)
	return error
}

// FetchRecordCount uses the passed query builder to count records
func (o *listFetcher) FetchRecordCount(appCfg appconfig.AppConfig, model interface{}, filters []services.QueryFilter) (int, error) {
	count, error := o.builder.Count(appCfg, model, filters)
	return count, error
}

// NewListFetcher returns an implementation of ListFetcher
func NewListFetcher(builder listQueryBuilder) services.ListFetcher {
	return &listFetcher{builder}
}
