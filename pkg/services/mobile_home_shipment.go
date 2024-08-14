package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// MobileHomeShipmentCreator creates a Mobile Home shipment
//
//go:generate mockery --name MobileHomeShipmentCreator
type MobileHomeShipmentCreator interface {
	CreateMobileHomeShipmentWithDefaultCheck(appCtx appcontext.AppContext, mobileHomeshipment *models.MobileHome) (*models.MobileHome, error)
}

// MobileHomeShipmentUpdater updates a Mobile Home shipment
//
//go:generate mockery --name MobileHomeShipmentUpdater
type MobileHomeShipmentUpdater interface {
	UpdateMobileHomeShipmentWithDefaultCheck(appCtx appcontext.AppContext, mobileHomeshipment *models.MobileHome, mtoShipmentID uuid.UUID) (*models.MobileHome, error)
}

// MobileHomeShipmentFetcher fetches a Mobile Home shipment
//
//go:generate mockery --name MobileHomeShipmentFetcher
type MobileHomeShipmentFetcher interface {
	GetMobileHomeShipment(appCtx appcontext.AppContext, mobileHomeShipmentID uuid.UUID, eagerPreloadAssociations []string, postloadAssociations []string) (*models.MobileHome, error)
	PostloadAssociations(appCtx appcontext.AppContext, mobileHomeShipment *models.MobileHome, postloadAssociations []string) error
}
