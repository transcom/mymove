package serviceitem

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type serviceItemListQueryBuilder interface {
	FetchMany(model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination) error
}

type serviceItemListFetcher struct {
	builder serviceItemListQueryBuilder
}

func (f *serviceItemListFetcher) FetchServiceItemList(filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination) (models.ServiceItems, error) {
	var serviceItems models.ServiceItems
	error := f.builder.FetchMany(&serviceItems, filters, associations, pagination)
	return serviceItems, error
}

func NewServiceItemListFetcher(builder serviceItemListQueryBuilder) services.ServiceItemListFetcher {
	return &serviceItemListFetcher{builder}
}
