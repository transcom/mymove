package office

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type officeQueryBuilder interface {
	FetchOne(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error
}

type officeFetcher struct {
	builder officeQueryBuilder
}

// FetchOffice fetches an office user for the given a slice of filters
func (o *officeFetcher) FetchOffice(appCtx appcontext.AppContext, filters []services.QueryFilter) (models.TransportationOffice, error) {
	var office models.TransportationOffice
	err := o.builder.FetchOne(appCtx, &office, filters)
	return office, err
}

// NewOfficeFetcher return an implementation of the OfficeFetcher interface
func NewOfficeFetcher(builder officeQueryBuilder) services.OfficeFetcher {
	return &officeFetcher{builder}
}
