package accesscode

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type accessCodeListQueryBuilder interface {
	query.FetchMany
}

type accessCodeListFetcher struct {
	builder accessCodeListQueryBuilder
}

// FetchAccessCodeList uses the passed query builder to fetch a list of access codes
func (o *accessCodeListFetcher) FetchAccessCodeList(filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination) (models.AccessCodes, error) {
	var accessCodes models.AccessCodes

	err := o.builder.WithModel(&accessCodes).WithFilters(filters).WithAssociations(associations).WithPagination(pagination).Execute()
	if err != nil {
		return models.AccessCodes{}, err
	}

	return accessCodes, nil
}

// NewAccessCodeListFetcher returns an implementation of AccessCodeListFetcher
func NewAccessCodeListFetcher(builder accessCodeListQueryBuilder) services.AccessCodeListFetcher {
	return &accessCodeListFetcher{builder}
}
