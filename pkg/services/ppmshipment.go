package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// PPMShipmentCreator creates a PPM shipment
//go:generate mockery --name PPMShipmentCreator --disable-version-string
type PPMShipmentCreator interface {
	CreatePPMShipmentWithDefaultCheck(appCtx appcontext.AppContext, ppmshipment *models.PPMShipment) (*models.PPMShipment, error)
}

// PPMShipmentUpdater updates a PPM shipment
//go:generate mockery --name PPMShipmentUpdater --disable-version-string
type PPMShipmentUpdater interface {
	UpdatePPMShipmentWithDefaultCheck(appCtx appcontext.AppContext, ppmshipment *models.PPMShipment, mtoShipmentID uuid.UUID) (*models.PPMShipment, error)
}

// PPMEstimator estimates the cost of a PPM shipment
//go:generate mockery --name PPMEstimator --disable-version-string
type PPMEstimator interface {
	EstimateIncentiveWithDefaultChecks(appCtx appcontext.AppContext, oldPPMShipment models.PPMShipment, newPPMShipment *models.PPMShipment) (*unit.Cents, error)
}
