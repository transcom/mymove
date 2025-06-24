package services

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

// OfficeUserListFetcher is the exported interface for fetching multiple  office users
//
//go:generate mockery --name OfficeUserListFetcher
type OfficeUserListFetcher interface {
	FetchOfficeUsersList(appCtx appcontext.AppContext, filterFuncs []func(*pop.Query), pagination Pagination, ordering QueryOrder) (models.OfficeUsers, int, error)
	FetchOfficeUsersCount(appCtx appcontext.AppContext, filters []QueryFilter) (int, error)
}

// OfficeUserFetcher is the exported interface for fetching a single office user
//
//go:generate mockery --name OfficeUserFetcher
type OfficeUserFetcher interface {
	FetchOfficeUser(appCtx appcontext.AppContext, filters []QueryFilter) (models.OfficeUser, error)
}

// OfficeUserFetcherPop is the exported interface for fetching a single office user
//
//go:generate mockery --name OfficeUserFetcherPop
type OfficeUserFetcherPop interface {
	FetchOfficeUserByID(appCtx appcontext.AppContext, id uuid.UUID) (models.OfficeUser, error)
	FetchOfficeUserByIDWithTransportationOfficeAssignments(appCtx appcontext.AppContext, id uuid.UUID) (models.OfficeUser, error)
	FetchOfficeUsersByRoleAndOffice(appCtx appcontext.AppContext, role roles.RoleType) ([]models.OfficeUser, error)
	FetchSafetyMoveOfficeUsersByRoleAndOffice(appCtx appcontext.AppContext, role roles.RoleType) ([]models.OfficeUser, error)
	FetchOfficeUsersWithWorkloadByRoleAndOffice(appCtx appcontext.AppContext, role roles.RoleType, officeID uuid.UUID, queueType string) ([]models.OfficeUserWithWorkload, error)
}

// OfficeUserGblocFetcher is the exported interface for fetching the GBLOC of the
// currently signed in office user
//
//go:generate mockery --name OfficeUserGblocFetcher
type OfficeUserGblocFetcher interface {
	FetchGblocForOfficeUser(appCtx appcontext.AppContext, id uuid.UUID) (string, error)
}

// OfficeUserCreator is the exported interface for creating an office user
//
//go:generate mockery --name OfficeUserCreator
type OfficeUserCreator interface {
	CreateOfficeUser(appCtx appcontext.AppContext, user *models.OfficeUser, transportationIDFilter []QueryFilter) (*models.OfficeUser, *validate.Errors, error)
}

// OfficeUserUpdater is the exported interface for updating an office user
//
//go:generate mockery --name OfficeUserUpdater
type OfficeUserUpdater interface {
	UpdateOfficeUser(appCtx appcontext.AppContext, id uuid.UUID, payload *models.OfficeUser, primaryTransportationOfficeId uuid.UUID) (*models.OfficeUser, *validate.Errors, error)
}

// OfficeUserDeleter is the exported interface for hard deleting an office user and its associations (roles, privileges)
//
//go:generate mockery --name OfficeUserDeleter
type OfficeUserDeleter interface {
	DeleteOfficeUser(appCtx appcontext.AppContext, id uuid.UUID) error
}
