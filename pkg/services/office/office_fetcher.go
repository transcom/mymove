package office

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type officeQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
}

type officeFetcher struct {
	builder officeQueryBuilder
}

// FetchOffice fetches an office user for the given a slice of filters
func (o *officeFetcher) FetchOffice(filters []services.QueryFilter) (models.TransportationOffice, error) {
	var office models.TransportationOffice
	error := o.builder.FetchOne(&office, filters)
	return office, error
}

// NewOfficeFetcher return an implementaion of the OfficeFetcher interface
func NewOfficeFetcher(builder officeQueryBuilder) services.OfficeFetcher {
	return &officeFetcher{builder}
}
