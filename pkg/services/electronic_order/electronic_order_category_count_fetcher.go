package electronicorder

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type electronicOrderCategoricalCountQueryBuilder interface {
	FetchCategoricalCountsFromOneModel(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter, andFilters *[]services.QueryFilter) (map[interface{}]int, error)
}

type electronicOrderCategoricalCountsFetcher struct {
	builder electronicOrderCategoricalCountQueryBuilder
}

// FetchElectronicOrderList uses the passed query builder to fetch a list of electronic_orders
func (o *electronicOrderCategoricalCountsFetcher) FetchElectronicOrderCategoricalCounts(appCtx appcontext.AppContext, filters []services.QueryFilter, andFilters *[]services.QueryFilter) (map[interface{}]int, error) {
	counts, err := o.builder.FetchCategoricalCountsFromOneModel(appCtx, models.ElectronicOrder{}, filters, andFilters)
	if err != nil {
		return nil, err
	}
	return counts, nil
}

// NewElectronicOrdersCategoricalCountsFetcher returns an implementation of OrdersListFetcher
func NewElectronicOrdersCategoricalCountsFetcher(builder electronicOrderCategoricalCountQueryBuilder) services.ElectronicOrderCategoryCountFetcher {
	return &electronicOrderCategoricalCountsFetcher{builder}
}
