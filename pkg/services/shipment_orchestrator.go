package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
)

// ShipmentCreator creates a shipment, taking into account different shipment types and their needs.
//
//go:generate mockery --name ShipmentCreator
type ShipmentCreator interface {
	CreateShipment(appCtx appcontext.AppContext, shipment *models.MTOShipment) (*models.MTOShipment, error)
}

// ShipmentUpdater updates a shipment, taking into account different shipment types and their needs.
//
//go:generate mockery --name ShipmentUpdater
type ShipmentUpdater interface {
	UpdateShipment(appCtx appcontext.AppContext, shipment *models.MTOShipment, eTag string, api string, planner route.Planner) (*models.MTOShipment, error)
}
