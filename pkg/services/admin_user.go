package services

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
)

// AdminUserListFetcher is the exported interface for fetching multiple admin users
//go:generate mockery --name AdminUserListFetcher --disable-version-string
type AdminUserListFetcher interface {
	FetchAdminUserList(appCfg appconfig.AppConfig, filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) (models.AdminUsers, error)
	FetchAdminUserCount(appCfg appconfig.AppConfig, filters []QueryFilter) (int, error)
}

// AdminUserFetcher is the exported interface for fetching a single admin user
//go:generate mockery --name AdminUserFetcher --disable-version-string
type AdminUserFetcher interface {
	FetchAdminUser(appCfg appconfig.AppConfig, filters []QueryFilter) (models.AdminUser, error)
}

// AdminUserCreator is the exported interface for creating an admin user
//go:generate mockery --name AdminUserCreator --disable-version-string
type AdminUserCreator interface {
	CreateAdminUser(appCfg appconfig.AppConfig, user *models.AdminUser, organizationIDFilter []QueryFilter) (*models.AdminUser, *validate.Errors, error)
}

// AdminUserUpdater is the exported interface for creating an admin user
//go:generate mockery --name AdminUserUpdater --disable-version-string
type AdminUserUpdater interface {
	UpdateAdminUser(appCfg appconfig.AppConfig, id uuid.UUID, payload *adminmessages.AdminUserUpdatePayload) (*models.AdminUser, *validate.Errors, error)
}
