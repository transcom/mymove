package office

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type officeListQueryBuilder interface {
	FetchMany(model interface{}, filters []services.QueryFilter, pagination services.Pagination) error
}

type officeListFetcher struct {
	builder officeListQueryBuilder
}

// FetchOfficeUserList uses the passed query builder to fetch a list of transportation offices
func (o *officeListFetcher) FetchOfficeList(filters []services.QueryFilter, pagination services.Pagination) (models.TransportationOffices, error) {
	var offices models.TransportationOffices
	error := o.builder.FetchMany(&offices, filters, pagination)
	return offices, error
}

// NewOfficeListFetcher returns an implementation of OfficeListFetcher
func NewOfficeListFetcher(builder officeListQueryBuilder) services.OfficeListFetcher {
	return &officeListFetcher{builder}
}
