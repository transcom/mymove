package services

import (
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MTOServiceItemCreator is the exported interface for creating a mto service item
//go:generate mockery --name MTOServiceItemCreator --disable-version-string
type MTOServiceItemCreator interface {
	CreateMTOServiceItem(serviceItem *models.MTOServiceItem) (*models.MTOServiceItems, *validate.Errors, error)
	SetConnection(db *pop.Connection)
}

// MTOServiceItemUpdater is the exported interface for updating an mto service item
//go:generate mockery --name MTOServiceItemUpdater --disable-version-string
type MTOServiceItemUpdater interface {
	UpdateMTOServiceItemStatus(mtoServiceItemID uuid.UUID, status models.MTOServiceItemStatus, rejectionReason *string, eTag string) (*models.MTOServiceItem, error)
	UpdateMTOServiceItem(db *pop.Connection, serviceItem *models.MTOServiceItem, eTag string, validator string) (*models.MTOServiceItem, error)
	UpdateMTOServiceItemBasic(db *pop.Connection, serviceItem *models.MTOServiceItem, eTag string) (*models.MTOServiceItem, error)
	UpdateMTOServiceItemPrime(db *pop.Connection, serviceItem *models.MTOServiceItem, eTag string) (*models.MTOServiceItem, error)
}
