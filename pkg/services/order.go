package services

import "github.com/transcom/mymove/pkg/models"

// OrdersListFetcher is the exported interface for fetching multiple orders
//go:generate mockery -name OrderListFetcher
type OrderListFetcher interface {
	FetchOrderList(filters []QueryFilter) (models.Orders, error)
}
