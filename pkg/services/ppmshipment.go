package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// PPMShipmentCreator creates a PPM shipment
type PPMShipmentCreator interface {
	CreatePPMShipmentWithDefaultCheck(appCtx appcontext.AppContext, ppmshipment *models.PPMShipment) (*models.PPMShipment, error)
}

// PPMShipmentUpdater updates a PPM shipment
type PPMShipmentUpdater interface {
	UpdatePPMShipmentCheck(appCtx appcontext.AppContext, ppmshipment *models.PPMShipment, eTag string) (*models.PPMShipment, error)
}
