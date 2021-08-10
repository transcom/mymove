package office

import (
	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type officeListQueryBuilder interface {
	FetchMany(appCfg appconfig.AppConfig, model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error
	Count(appCfg appconfig.AppConfig, model interface{}, filters []services.QueryFilter) (int, error)
}

type officeListFetcher struct {
	builder officeListQueryBuilder
}

// FetchOfficeUserList uses the passed query builder to fetch a list of transportation offices
func (o *officeListFetcher) FetchOfficeList(appCfg appconfig.AppConfig, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) (models.TransportationOffices, error) {
	var offices models.TransportationOffices
	error := o.builder.FetchMany(appCfg, &offices, filters, associations, pagination, ordering)
	return offices, error
}

// FetchOfficeUserCount uses the passed query builder to count the number of transportation offices
func (o *officeListFetcher) FetchOfficeCount(appCfg appconfig.AppConfig, filters []services.QueryFilter) (int, error) {
	var offices models.TransportationOffices
	count, error := o.builder.Count(appCfg, &offices, filters)
	return count, error
}

// NewOfficeListFetcher returns an implementation of OfficeListFetcher
func NewOfficeListFetcher(builder officeListQueryBuilder) services.OfficeListFetcher {
	return &officeListFetcher{builder}
}
