package services

import (
	"github.com/transcom/mymove/pkg/models"
)

// OrganizationFetcher is the exported interface for fetching a single organization
type OrganizationFetcher interface {
	FetchOrganization(filters []QueryFilter) (models.Organization, error)
}

// OrganizationListFetcher is the exported interface for fetching multiple organization
//go:generate mockery -name OrganizationFetcher
type OrganizationListFetcher interface {
	FetchOrganizationList(filters []QueryFilter, associations QueryAssociations, pagination Pagination) (models.Organizations, error)
}
