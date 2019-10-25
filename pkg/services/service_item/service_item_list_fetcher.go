package serviceitem

import (
	"github.com/gofrs/uuid"

	serviceitemop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/service_item"
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

func (f *serviceItemListFetcher) FetchServiceItemList(params interface{}) (models.ServiceItems, error) {
	parameters := params.(serviceitemop.ListServiceItemsParams)
	id, err := uuid.FromString(parameters.MoveTaskOrderID)

	if err != nil {
		return nil, err
	}

	pagination := pagination.NewPagination(nil, nil)
	associations := query.NewQueryAssociations([]services.QueryAssociation{})
	filters := []services.QueryFilter{query.NewQueryFilter("move_task_order_id", "=", id)}
	var serviceItems models.ServiceItems
	error := f.builder.FetchMany(&serviceItems, filters, associations, pagination)
	return serviceItems, error
}

func NewServiceItemListFetcher(builder serviceItemListQueryBuilder) services.ServiceItemListFetcher {
	return &serviceItemListFetcher{builder}
}
