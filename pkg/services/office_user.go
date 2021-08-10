package services

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
)

// OfficeUserFetcher is the exported interface for fetching a single office user
//go:generate mockery --name OfficeUserFetcher --disable-version-string
type OfficeUserFetcher interface {
	FetchOfficeUser(appCfg appconfig.AppConfig, filters []QueryFilter) (models.OfficeUser, error)
}

// OfficeUserFetcherPop is the exported interface for fetching a single office user
//go:generate mockery --name OfficeUserFetcherPop --disable-version-string
type OfficeUserFetcherPop interface {
	FetchOfficeUserByID(appCfg appconfig.AppConfig, id uuid.UUID) (models.OfficeUser, error)
}

// OfficeUserGblocFetcher is the exported interface for fetching the GBLOC of the
// currently signed in office user
//go:generate mockery --name OfficeUserGblocFetcher --disable-version-string
type OfficeUserGblocFetcher interface {
	FetchGblocForOfficeUser(appCfg appconfig.AppConfig, id uuid.UUID) (string, error)
}

// OfficeUserCreator is the exported interface for creating an office user
//go:generate mockery --name OfficeUserCreator --disable-version-string
type OfficeUserCreator interface {
	CreateOfficeUser(appCfg appconfig.AppConfig, user *models.OfficeUser, transportationIDFilter []QueryFilter) (*models.OfficeUser, *validate.Errors, error)
}

// OfficeUserUpdater is the exported interface for creating an office user
//go:generate mockery --name OfficeUserUpdater --disable-version-string
type OfficeUserUpdater interface {
	UpdateOfficeUser(appCfg appconfig.AppConfig, id uuid.UUID, payload *adminmessages.OfficeUserUpdatePayload) (*models.OfficeUser, *validate.Errors, error)
}
