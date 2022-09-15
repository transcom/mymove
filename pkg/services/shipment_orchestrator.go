package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// ShipmentCreator creates a shipment, taking into account different shipment types and their needs.
//go:generate mockery --name ShipmentCreator --disable-version-string
type ShipmentCreator interface {
	CreateShipment(appCtx appcontext.AppContext, shipment *models.MTOShipment) (*models.MTOShipment, error)
}

// ShipmentUpdater updates a shipment, taking into account different shipment types and their needs.
//go:generate mockery --name ShipmentUpdater --disable-version-string
type ShipmentUpdater interface {
	UpdateShipment(appCtx appcontext.AppContext, shipment *models.MTOShipment, eTag string) (*models.MTOShipment, error)
}
