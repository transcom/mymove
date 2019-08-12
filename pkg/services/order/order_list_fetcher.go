package order

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type orderListQueryBuilder interface {
	FetchMany(model interface{}, filters []services.QueryFilter) error
}

type orderListFetcher struct {
	builder orderListQueryBuilder
}

// FetchOrderList uses the passed query builder to fetch a list of orders
func (o *orderListFetcher) FetchOrderList(filters []services.QueryFilter) (models.Orders, error) {
	var orders models.Orders
	error := o.builder.FetchMany(&orders, filters)
	return orders, error
}

// NewOrderListFetcher returns an implementation of OrdersListFetcher
func NewOrderListFetcher(builder orderListQueryBuilder) services.OrderListFetcher {
	return &orderListFetcher{builder}
}
