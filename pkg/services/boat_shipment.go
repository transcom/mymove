package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// BoatShipmentCreator creates a Boat shipment
//
//go:generate mockery --name BoatShipmentCreator
type BoatShipmentCreator interface {
	CreateBoatShipmentWithDefaultCheck(appCtx appcontext.AppContext, boatshipment *models.BoatShipment) (*models.BoatShipment, error)
}

// BoatShipmentUpdater updates a Boat shipment
//
//go:generate mockery --name BoatShipmentUpdater
type BoatShipmentUpdater interface {
	UpdateBoatShipmentWithDefaultCheck(appCtx appcontext.AppContext, boatshipment *models.BoatShipment, mtoShipmentID uuid.UUID) (*models.BoatShipment, error)
	// UpdateBoatShipmentSITEstimatedCost(appCtx appcontext.AppContext, boatshipment *models.BoatShipment) (*models.BoatShipment, error)
}

// BoatShipmentFetcher fetches a Boat shipment
//
//go:generate mockery --name BoatShipmentFetcher
type BoatShipmentFetcher interface {
	GetBoatShipment(appCtx appcontext.AppContext, boatShipmentID uuid.UUID, eagerPreloadAssociations []string, postloadAssociations []string) (*models.BoatShipment, error)
	PostloadAssociations(appCtx appcontext.AppContext, boatShipment *models.BoatShipment, postloadAssociations []string) error
}
