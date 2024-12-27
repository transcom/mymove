package services

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
)

// RequestedOfficeUserListFetcher is the exported interface for fetching multiple requested office users
//
//go:generate mockery --name RequestedOfficeUserListFetcher
type RequestedOfficeUserListFetcher interface {
	FetchRequestedOfficeUsersList(appCtx appcontext.AppContext, filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) (models.OfficeUsers, error)
	FetchRequestedOfficeUsersCount(appCtx appcontext.AppContext, filters []QueryFilter) (int, error)
}

// RequestedOfficeUserFetcher is the exported interface for fetching a single requested office user
//
//go:generate mockery --name RequestedOfficeUserFetcher
type RequestedOfficeUserFetcher interface {
	FetchRequestedOfficeUser(appCtx appcontext.AppContext, filters []QueryFilter) (models.OfficeUser, error)
}

// RequestedOfficeUserFetcherPop is the exported interface for fetching a single office user
//
//go:generate mockery --name RequestedOfficeUserFetcherPop
type RequestedOfficeUserFetcherPop interface {
	FetchRequestedOfficeUserByID(appCtx appcontext.AppContext, id uuid.UUID) (models.OfficeUser, error)
}

// RequestedOfficeUserFetcher is the exported interface for updating a requested office user
//
//go:generate mockery --name RequestedOfficeUserUpdater
type RequestedOfficeUserUpdater interface {
	UpdateRequestedOfficeUser(appCtx appcontext.AppContext, id uuid.UUID, payload *adminmessages.RequestedOfficeUserUpdate) (*models.OfficeUser, *validate.Errors, error)
}
