package services

import "github.com/transcom/mymove/pkg/models"

// ElectronicOrderListFetcher is the exported interface for fetching multiple electronic orders
//go:generate mockery -name ElectronicOrderListFetcher
type ElectronicOrderListFetcher interface {
	FetchElectronicOrderList(filters []QueryFilter, pagination Pagination) (models.ElectronicOrders, error)
}
