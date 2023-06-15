package office

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type officeListQueryBuilder interface {
	FetchMany(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error
	Count(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) (int, error)
}

type officeListFetcher struct {
	builder officeListQueryBuilder
}

// FetchOfficeUserList uses the passed query builder to fetch a list of transportation offices
func (o *officeListFetcher) FetchOfficeList(appCtx appcontext.AppContext, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) (models.TransportationOffices, error) {
	var offices models.TransportationOffices
	err := o.builder.FetchMany(appCtx, &offices, filters, associations, pagination, ordering)
	return offices, err
}

// FetchOfficeUserCount uses the passed query builder to count the number of transportation offices
func (o *officeListFetcher) FetchOfficeCount(appCtx appcontext.AppContext, filters []services.QueryFilter) (int, error) {
	var offices models.TransportationOffices
	count, err := o.builder.Count(appCtx, &offices, filters)
	return count, err
}

// NewOfficeListFetcher returns an implementation of OfficeListFetcher
func NewOfficeListFetcher(builder officeListQueryBuilder) services.OfficeListFetcher {
	return &officeListFetcher{builder}
}
