package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// PPMShipmentCreator creates a PPM shipment
//
//go:generate mockery --name PPMShipmentCreator --disable-version-string
type PPMShipmentCreator interface {
	CreatePPMShipmentWithDefaultCheck(appCtx appcontext.AppContext, ppmshipment *models.PPMShipment) (*models.PPMShipment, error)
}

// PPMShipmentUpdater updates a PPM shipment
//
//go:generate mockery --name PPMShipmentUpdater --disable-version-string
type PPMShipmentUpdater interface {
	UpdatePPMShipmentWithDefaultCheck(appCtx appcontext.AppContext, ppmshipment *models.PPMShipment, mtoShipmentID uuid.UUID) (*models.PPMShipment, error)
}

// PPMEstimator estimates the cost of a PPM shipment
//
//go:generate mockery --name PPMEstimator --disable-version-string
type PPMEstimator interface {
	EstimateIncentiveWithDefaultChecks(appCtx appcontext.AppContext, oldPPMShipment models.PPMShipment, newPPMShipment *models.PPMShipment) (*unit.Cents, *unit.Cents, error)
}

// PPMShipmentRouter routes a PPM shipment
//
//go:generate mockery --name PPMShipmentRouter --disable-version-string
type PPMShipmentRouter interface {
	SetToDraft(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) error
	Submit(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) error
	SendToCustomer(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) error
	SubmitCloseOutDocumentation(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) error
}

// PPMShipmentNewSubmitter handles a new submission for a PPM shipment
//
//go:generate mockery --name PPMShipmentNewSubmitter --disable-version-string
type PPMShipmentNewSubmitter interface {
	SubmitNewCustomerCloseOut(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID, signedCertification models.SignedCertification) (*models.PPMShipment, error)
}

// PPMShipmentUpdatedSubmitter handles a submission for a PPM shipment that has been submitted before
//
//go:generate mockery --name PPMShipmentUpdatedSubmitter --disable-version-string
type PPMShipmentUpdatedSubmitter interface {
	SubmitUpdatedCustomerCloseOut(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID, signedCertification models.SignedCertification, eTag string) (*models.PPMShipment, error)
}
