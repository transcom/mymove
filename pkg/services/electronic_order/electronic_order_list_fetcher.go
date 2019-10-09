package electronicorder

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type electronicOrderListQueryBuilder interface {
	FetchMany(model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination) error
}

type electronicOrderListFetcher struct {
	builder electronicOrderListQueryBuilder
}

// FetchElectronicOrderList uses the passed query builder to fetch a list of electronic_orders
func (o *electronicOrderListFetcher) FetchElectronicOrderList(filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination) (models.ElectronicOrders, error) {
	var electronicOrders models.ElectronicOrders
	error := o.builder.FetchMany(&electronicOrders, filters, associations, pagination)
	return electronicOrders, error
}

// NewElectronicOrderListFetcher returns an implementation of OrdersListFetcher
func NewElectronicOrderListFetcher(builder electronicOrderListQueryBuilder) services.ElectronicOrderListFetcher {
	return &electronicOrderListFetcher{builder}
}
