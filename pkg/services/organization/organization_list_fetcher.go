package organization

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type organizationListQueryBuilder interface {
	FetchMany(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error
	Count(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) (int, error)
}

type organizationListFetcher struct {
	builder organizationListQueryBuilder
}

// FetchOrganizationUserList uses the passed query builder to fetch a list of transportation offices
func (o *organizationListFetcher) FetchOrganizationList(appCtx appcontext.AppContext, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) (models.Organizations, error) {
	var organizations models.Organizations
	error := o.builder.FetchMany(appCtx, &organizations, filters, associations, pagination, ordering)
	return organizations, error
}

// FetchOrganizationUserList uses the passed query builder to fetch a list of transportation offices
func (o *organizationListFetcher) FetchOrganizationCount(appCtx appcontext.AppContext, filters []services.QueryFilter) (int, error) {
	var organizations models.Organizations
	count, error := o.builder.Count(appCtx, &organizations, filters)
	return count, error
}

// NewOrganizationListFetcher returns an implementation of OrganizationListFetcher
func NewOrganizationListFetcher(builder organizationListQueryBuilder) services.OrganizationListFetcher {
	return &organizationListFetcher{builder}
}
