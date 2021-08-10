package accesscode

import (
	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type accessCodeListQueryBuilder interface {
	FetchMany(appCfg appconfig.AppConfig, model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error
	Count(appCfg appconfig.AppConfig, model interface{}, filters []services.QueryFilter) (int, error)
}

type accessCodeListFetcher struct {
	builder accessCodeListQueryBuilder
}

// FetchAccessCodeList uses the passed query builder to fetch a list of access codes
func (o *accessCodeListFetcher) FetchAccessCodeList(appCfg appconfig.AppConfig, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) (models.AccessCodes, error) {
	var accessCodes models.AccessCodes

	err := o.builder.FetchMany(appCfg, &accessCodes, filters, associations, pagination, ordering)
	if err != nil {
		return models.AccessCodes{}, err
	}

	return accessCodes, nil
}

// FetchAccessCodeCount uses the passed query builder to count access codes
func (o *accessCodeListFetcher) FetchAccessCodeCount(appCfg appconfig.AppConfig, filters []services.QueryFilter) (int, error) {
	var accessCodes models.AccessCodes
	count, err := o.builder.Count(appCfg, &accessCodes, filters)
	return count, err
}

// NewAccessCodeListFetcher returns an implementation of AccessCodeListFetcher
func NewAccessCodeListFetcher(builder accessCodeListQueryBuilder) services.AccessCodeListFetcher {
	return &accessCodeListFetcher{builder}
}
