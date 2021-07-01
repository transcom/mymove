package services

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
)

// OfficeUserFetcher is the exported interface for fetching a single office user
//go:generate mockery --name OfficeUserFetcher --disable-version-string
type OfficeUserFetcher interface {
	FetchOfficeUser(filters []QueryFilter) (models.OfficeUser, error)
}

// OfficeUserFetcherPop is the exported interface for fetching a single office user
//go:generate mockery --name OfficeUserFetcherPop --disable-version-string
type OfficeUserFetcherPop interface {
	FetchOfficeUserByID(id uuid.UUID) (models.OfficeUser, error)
}

// OfficeUserGblocFetcher is the exported interface for fetching the GBLOC of the
// currently signed in office user
//go:generate mockery --name OfficeUserGblocFetcher --disable-version-string
type OfficeUserGblocFetcher interface {
	FetchGblocForOfficeUser(id uuid.UUID) (string, error)
}

// OfficeUserCreator is the exported interface for creating an office user
//go:generate mockery --name OfficeUserCreator --disable-version-string
type OfficeUserCreator interface {
	CreateOfficeUser(user *models.OfficeUser, transportationIDFilter []QueryFilter) (*models.OfficeUser, *validate.Errors, error)
}

// OfficeUserUpdater is the exported interface for creating an office user
//go:generate mockery --name OfficeUserUpdater --disable-version-string
type OfficeUserUpdater interface {
	UpdateOfficeUser(id uuid.UUID, payload *adminmessages.OfficeUserUpdatePayload) (*models.OfficeUser, *validate.Errors, error)
}
