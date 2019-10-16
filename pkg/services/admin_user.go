package services

import (
	"github.com/transcom/mymove/pkg/models"
)

// AdminUserListFetcher is the exported interface for fetching multiple admin users
//go:generate mockery -name AdminUserListFetcher
type AdminUserListFetcher interface {
	FetchAdminUserList(filters []QueryFilter, associations QueryAssociations, pagination Pagination) (models.AdminUsers, error)
}

// AdminUserFetcher is the exported interface for fetching a single admin user
//go:generate mockery -name AdminUserFetcher
type AdminUserFetcher interface {
	FetchAdminUser(filters []QueryFilter) (models.AdminUser, error)
}
