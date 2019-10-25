package serviceitem

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
)

type serviceItemListQueryBuilder interface {
	FetchMany(model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination) error
}

type serviceItemListFetcher struct {
	builder serviceItemListQueryBuilder
}

func (f *serviceItemListFetcher) FetchServiceItemList(moveTaskOrderID uuid.UUID) (models.ServiceItems, error) {

	pagination := pagination.NewPagination(nil, nil)
	associations := query.NewQueryAssociations([]services.QueryAssociation{})
	filters := []services.QueryFilter{query.NewQueryFilter("move_task_order_id", "=", moveTaskOrderID)}
	var serviceItems models.ServiceItems
	error := f.builder.FetchMany(&serviceItems, filters, associations, pagination)
	return serviceItems, error
}

func NewServiceItemListFetcher(builder serviceItemListQueryBuilder) services.ServiceItemListFetcher {
	return &serviceItemListFetcher{builder}
}
