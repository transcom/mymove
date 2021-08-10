package services

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
)

// MTOServiceItemCreator is the exported interface for creating a mto service item
//go:generate mockery --name MTOServiceItemCreator --disable-version-string
type MTOServiceItemCreator interface {
	CreateMTOServiceItem(appCfg appconfig.AppConfig, serviceItem *models.MTOServiceItem) (*models.MTOServiceItems, *validate.Errors, error)
}

// MTOServiceItemUpdater is the exported interface for updating an mto service item
//go:generate mockery --name MTOServiceItemUpdater --disable-version-string
type MTOServiceItemUpdater interface {
	UpdateMTOServiceItemStatus(appCfg appconfig.AppConfig, mtoServiceItemID uuid.UUID, status models.MTOServiceItemStatus, rejectionReason *string, eTag string) (*models.MTOServiceItem, error)
	UpdateMTOServiceItem(appCfg appconfig.AppConfig, serviceItem *models.MTOServiceItem, eTag string, validator string) (*models.MTOServiceItem, error)
	UpdateMTOServiceItemBasic(appCfg appconfig.AppConfig, serviceItem *models.MTOServiceItem, eTag string) (*models.MTOServiceItem, error)
	UpdateMTOServiceItemPrime(appCfg appconfig.AppConfig, serviceItem *models.MTOServiceItem, eTag string) (*models.MTOServiceItem, error)
}
