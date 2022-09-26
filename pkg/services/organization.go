package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// OrganizationFetcher is the exported interface for fetching a single organization
type OrganizationFetcher interface {
	FetchOrganization(appCtx appcontext.AppContext, filters []QueryFilter) (models.Organization, error)
}

// OrganizationListFetcher is the exported interface for fetching multiple organizations
//
//go:generate mockery --name OrganizationListFetcher --disable-version-string
type OrganizationListFetcher interface {
	FetchOrganizationList(appCtx appcontext.AppContext, filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) (models.Organizations, error)
	FetchOrganizationCount(appCtx appcontext.AppContext, filters []QueryFilter) (int, error)
}
