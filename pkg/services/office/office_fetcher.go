package office

import (
	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type officeQueryBuilder interface {
	FetchOne(appCfg appconfig.AppConfig, model interface{}, filters []services.QueryFilter) error
}

type officeFetcher struct {
	builder officeQueryBuilder
}

// FetchOffice fetches an office user for the given a slice of filters
func (o *officeFetcher) FetchOffice(appCfg appconfig.AppConfig, filters []services.QueryFilter) (models.TransportationOffice, error) {
	var office models.TransportationOffice
	error := o.builder.FetchOne(appCfg, &office, filters)
	return office, error
}

// NewOfficeFetcher return an implementaion of the OfficeFetcher interface
func NewOfficeFetcher(builder officeQueryBuilder) services.OfficeFetcher {
	return &officeFetcher{builder}
}
