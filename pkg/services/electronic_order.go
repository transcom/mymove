package services

import "github.com/transcom/mymove/pkg/models"

// ElectronicOrderListFetcher is the exported interface for fetching multiple electronic orders
//go:generate mockery --name ElectronicOrderListFetcher --disable-version-string
type ElectronicOrderListFetcher interface {
	FetchElectronicOrderList(filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) (models.ElectronicOrders, error)
	FetchElectronicOrderCount(filters []QueryFilter) (int, error)
}

// ElectronicOrderCategoryCountFetcher is the exported interface for fetching counts of Electronic orders based on provided category QueryFilters.
//go:generate mockery --name ElectronicOrderCategoryCountFetcher --disable-version-string
type ElectronicOrderCategoryCountFetcher interface {
	FetchElectronicOrderCategoricalCounts(filters []QueryFilter, andFilters *[]QueryFilter) (map[interface{}]int, error)
}
