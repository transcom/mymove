package services

import (
	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
)

// OrganizationFetcher is the exported interface for fetching a single organization
type OrganizationFetcher interface {
	FetchOrganization(appCfg appconfig.AppConfig, filters []QueryFilter) (models.Organization, error)
}

// OrganizationListFetcher is the exported interface for fetching multiple organizations
//go:generate mockery --name OrganizationListFetcher --disable-version-string
type OrganizationListFetcher interface {
	FetchOrganizationList(appCfg appconfig.AppConfig, filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) (models.Organizations, error)
	FetchOrganizationCount(appCfg appconfig.AppConfig, filters []QueryFilter) (int, error)
}
