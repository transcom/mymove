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
func (o *accessCodeListFetcher) FetchAccessCodeList(filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) (models.AccessCodes, error) {
	var accessCodes models.AccessCodes

	err := o.builder.WithFilters(filters).WithAssociations(associations).WithPagination(pagination).Execute(&accessCodes)
	// err := o.builder.FetchMany(&accessCodes, filters, associations, pagination, ordering)
	if err != nil {
		return models.AccessCodes{}, err
	}

	return accessCodes, nil
}

// FetchAccessCodeCount uses the passed query builder to count access codes
func (o *accessCodeListFetcher) FetchAccessCodeCount(filters []services.QueryFilter) (int, error) {
	var accessCodes models.AccessCodes
	count, err := o.builder.WithFilters(filters).Count(&accessCodes)
	return count, err
}

// NewAccessCodeListFetcher returns an implementation of AccessCodeListFetcher
func NewAccessCodeListFetcher(builder accessCodeListQueryBuilder) services.AccessCodeListFetcher {
	return &accessCodeListFetcher{builder}
}
