package services

import (
	"github.com/transcom/mymove/pkg/models"
)

// AdminUserListFetcher is the exported interface for fetching multiple office users
//go:generate mockery -name AdminUserListFetcher
type AdminUserListFetcher interface {
	FetchAdminUserList(filters []QueryFilter, pagination Pagination) (models.AdminUsers, error)
}
