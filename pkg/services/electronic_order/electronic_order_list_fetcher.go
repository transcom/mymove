package electronicorder

import (
	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type electronicOrderListQueryBuilder interface {
	FetchMany(appCfg appconfig.AppConfig, model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error
	Count(appCfg appconfig.AppConfig, model interface{}, filters []services.QueryFilter) (int, error)
}

type electronicOrderListFetcher struct {
	builder electronicOrderListQueryBuilder
}

// FetchElectronicOrderList uses the passed query builder to fetch a list of electronic_orders
func (o *electronicOrderListFetcher) FetchElectronicOrderList(appCfg appconfig.AppConfig, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) (models.ElectronicOrders, error) {
	var electronicOrders models.ElectronicOrders
	error := o.builder.FetchMany(appCfg, &electronicOrders, filters, associations, pagination, ordering)
	return electronicOrders, error
}

// FetchElectronicOrderCount uses the passed query builder to count electronic_orders
func (o *electronicOrderListFetcher) FetchElectronicOrderCount(appCfg appconfig.AppConfig, filters []services.QueryFilter) (int, error) {
	var electronicOrders models.ElectronicOrders
	count, error := o.builder.Count(appCfg, &electronicOrders, filters)
	return count, error
}

// NewElectronicOrderListFetcher returns an implementation of OrdersListFetcher
func NewElectronicOrderListFetcher(builder electronicOrderListQueryBuilder) services.ElectronicOrderListFetcher {
	return &electronicOrderListFetcher{builder}
}
