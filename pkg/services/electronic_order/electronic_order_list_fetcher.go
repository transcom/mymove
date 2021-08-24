package electronicorder

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type electronicOrderListQueryBuilder interface {
	FetchMany(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error
	Count(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) (int, error)
}

type electronicOrderListFetcher struct {
	builder electronicOrderListQueryBuilder
}

// FetchElectronicOrderList uses the passed query builder to fetch a list of electronic_orders
func (o *electronicOrderListFetcher) FetchElectronicOrderList(appCtx appcontext.AppContext, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) (models.ElectronicOrders, error) {
	var electronicOrders models.ElectronicOrders
	error := o.builder.FetchMany(appCtx, &electronicOrders, filters, associations, pagination, ordering)
	return electronicOrders, error
}

// FetchElectronicOrderCount uses the passed query builder to count electronic_orders
func (o *electronicOrderListFetcher) FetchElectronicOrderCount(appCtx appcontext.AppContext, filters []services.QueryFilter) (int, error) {
	var electronicOrders models.ElectronicOrders
	count, error := o.builder.Count(appCtx, &electronicOrders, filters)
	return count, error
}

// NewElectronicOrderListFetcher returns an implementation of OrdersListFetcher
func NewElectronicOrderListFetcher(builder electronicOrderListQueryBuilder) services.ElectronicOrderListFetcher {
	return &electronicOrderListFetcher{builder}
}
