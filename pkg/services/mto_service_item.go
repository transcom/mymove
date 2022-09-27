package services

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// MTOServiceItemCreator is the exported interface for creating a mto service item
//
//go:generate mockery --name MTOServiceItemCreator --disable-version-string
type MTOServiceItemCreator interface {
	CreateMTOServiceItem(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem) (*models.MTOServiceItems, *validate.Errors, error)
}

// MTOServiceItemUpdater is the exported interface for updating an mto service item
//
//go:generate mockery --name MTOServiceItemUpdater --disable-version-string
type MTOServiceItemUpdater interface {
	ApproveOrRejectServiceItem(appCtx appcontext.AppContext, mtoServiceItemID uuid.UUID, status models.MTOServiceItemStatus, rejectionReason *string, eTag string) (*models.MTOServiceItem, error)
	UpdateMTOServiceItem(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem, eTag string, validator string) (*models.MTOServiceItem, error)
	UpdateMTOServiceItemBasic(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem, eTag string) (*models.MTOServiceItem, error)
	UpdateMTOServiceItemPrime(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem, eTag string) (*models.MTOServiceItem, error)
}
