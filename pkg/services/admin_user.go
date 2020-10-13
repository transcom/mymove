package services

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
)

// AdminUserListFetcher is the exported interface for fetching multiple admin users
//go:generate mockery -name AdminUserListFetcher
type AdminUserListFetcher interface {
	FetchAdminUserList(filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) (models.AdminUsers, error)
	FetchAdminUserCount(filters []QueryFilter) (int, error)
}

// AdminUserFetcher is the exported interface for fetching a single admin user
//go:generate mockery -name AdminUserFetcher
type AdminUserFetcher interface {
	FetchAdminUser(filters []QueryFilter) (models.AdminUser, error)
}

// AdminUserCreator is the exported interface for creating an admin user
//go:generate mockery -name AdminUserCreator
type AdminUserCreator interface {
	CreateAdminUser(user *models.AdminUser, organizationIDFilter []QueryFilter) (*models.AdminUser, *validate.Errors, error)
}

// AdminUserUpdater is the exported interface for creating an admin user
//go:generate mockery -name AdminUserUpdater
type AdminUserUpdater interface {
	UpdateAdminUser(id uuid.UUID, payload *adminmessages.AdminUserUpdatePayload) (*models.AdminUser, *validate.Errors, error)
}
