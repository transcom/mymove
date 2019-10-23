package services

import "github.com/transcom/mymove/pkg/models"

// ServiceItemListFetcher is the exported interface for fetching multiple transportation offices
//go:generate mockery -name ServiceItemListFetcher
type ServiceItemListFetcher interface {
	FetchServiceItemList(filters []QueryFilter, associations QueryAssociations, pagination Pagination) (models.ServiceItems, error)
}
