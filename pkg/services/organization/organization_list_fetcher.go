package organization

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type organizationListQueryBuilder interface {
	FetchMany(model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error
}

type organizationListFetcher struct {
	builder organizationListQueryBuilder
}

// FetchOrganizationUserList uses the passed query builder to fetch a list of transportation offices
func (o *organizationListFetcher) FetchOrganizationList(filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) (models.Organizations, error) {
	var organizations models.Organizations
	error := o.builder.FetchMany(&organizations, filters, associations, pagination, ordering)
	return organizations, error
}

// NewOrganizationListFetcher returns an implementation of OrganizationListFetcher
func NewOrganizationListFetcher(builder organizationListQueryBuilder) services.OrganizationListFetcher {
	return &organizationListFetcher{builder}
}
