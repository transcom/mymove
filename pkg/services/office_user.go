package services

import (
	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/models"
)

// OfficeUserFetcher is the exported interface for fetching a single office user
//go:generate mockery -name OfficeUserFetcher
type OfficeUserFetcher interface {
	FetchOfficeUser(filters []QueryFilter) (models.OfficeUser, error)
}

// OfficeUserListFetcher is the exported interface for fetching multiple office users
//go:generate mockery -name OfficeUserListFetcher
type OfficeUserListFetcher interface {
	FetchOfficeUserList(filters []QueryFilter) (models.OfficeUsers, error)
}

// OfficeUserCreator is the exported interface for creating an office user
//go:generate mockery -name OfficeUserCreator
type OfficeUserCreator interface {
	CreateOfficeUser(user *models.OfficeUser, transportationIDFilter []QueryFilter) (*models.OfficeUser, *validate.Errors, error)
}

// OfficeUserUpdater is the exported interface for creating an office user
//go:generate mockery -name OfficeUserUpdater
type OfficeUserUpdater interface {
	UpdateOfficeUser(user *models.OfficeUser) (*models.OfficeUser, *validate.Errors, error)
}
