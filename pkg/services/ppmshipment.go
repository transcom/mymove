package services

import (
	"io"

	"github.com/gofrs/uuid"
	"github.com/spf13/afero"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// PPMShipmentCreator creates a PPM shipment
//
//go:generate mockery --name PPMShipmentCreator
type PPMShipmentCreator interface {
	CreatePPMShipmentWithDefaultCheck(appCtx appcontext.AppContext, ppmshipment *models.PPMShipment) (*models.PPMShipment, error)
}

// PPMShipmentUpdater updates a PPM shipment
//
//go:generate mockery --name PPMShipmentUpdater
type PPMShipmentUpdater interface {
	UpdatePPMShipmentWithDefaultCheck(appCtx appcontext.AppContext, ppmshipment *models.PPMShipment, mtoShipmentID uuid.UUID) (*models.PPMShipment, error)
	UpdatePPMShipmentSITEstimatedCost(appCtx appcontext.AppContext, ppmshipment *models.PPMShipment) (*models.PPMShipment, error)
}

// PPMShipmentFetcher fetches a PPM shipment
//
//go:generate mockery --name PPMShipmentFetcher
type PPMShipmentFetcher interface {
	GetPPMShipment(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID, eagerPreloadAssociations []string, postloadAssociations []string) (*models.PPMShipment, error)
	PostloadAssociations(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment, postloadAssociations []string) error
}

// PPMDocumentFetcher fetches all documents associated with a PPM shipment
//
//go:generate mockery --name PPMDocumentFetcher
type PPMDocumentFetcher interface {
	GetPPMDocuments(appCtx appcontext.AppContext, mtoShipmentID uuid.UUID) (*models.PPMDocuments, error)
}

// PPMEstimator estimates the cost of a PPM shipment
//
//go:generate mockery --name PPMEstimator
type PPMEstimator interface {
	EstimateIncentiveWithDefaultChecks(appCtx appcontext.AppContext, oldPPMShipment models.PPMShipment, newPPMShipment *models.PPMShipment) (*unit.Cents, *unit.Cents, error)
	FinalIncentiveWithDefaultChecks(appCtx appcontext.AppContext, oldPPMShipment models.PPMShipment, newPPMShipment *models.PPMShipment) (*unit.Cents, error)
	PriceBreakdown(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) (unit.Cents, unit.Cents, unit.Cents, unit.Cents, unit.Cents, unit.Cents, unit.Cents, error)
	CalculatePPMSITEstimatedCost(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) (*unit.Cents, error)
	CalculatePPMSITEstimatedCostBreakdown(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) (*models.PPMSITEstimatedCostInfo, error)
}

// PPMShipmentRouter routes a PPM shipment
//
//go:generate mockery --name PPMShipmentRouter
type PPMShipmentRouter interface {
	SetToDraft(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) error
	Submit(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) error
	SendToCustomer(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) error
	SubmitCloseOutDocumentation(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) error
	SubmitReviewedDocuments(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) error
}

// PPMShipmentNewSubmitter handles a new submission for a PPM shipment
//
//go:generate mockery --name PPMShipmentNewSubmitter
type PPMShipmentNewSubmitter interface {
	SubmitNewCustomerCloseOut(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID, signedCertification models.SignedCertification) (*models.PPMShipment, error)
}

// PPMShipmentReviewDocuments handles a new submission for a PPM shipment
//
//go:generate mockery --name PPMShipmentReviewDocuments
type PPMShipmentReviewDocuments interface {
	SubmitReviewedDocuments(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.PPMShipment, error)
}

// PPMShipmentUpdatedSubmitter handles a submission for a PPM shipment that has been submitted before
//
//go:generate mockery --name PPMShipmentUpdatedSubmitter
type PPMShipmentUpdatedSubmitter interface {
	SubmitUpdatedCustomerCloseOut(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID, signedCertification models.SignedCertification, eTag string) (*models.PPMShipment, error)
}

// AOAPacketCreator creates an AOA packet for a PPM shipment
//
//go:generate mockery --name AOAPacketCreator
type AOAPacketCreator interface {
	VerifyAOAPacketInternal(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) error
	CreateAOAPacket(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID, isPaymentPacket bool) (afero.File, error)
}

// PaymentPacketCreator creates a payment packet for a PPM shipment
//
//go:generate mockery --name PaymentPacketCreator
type PaymentPacketCreator interface {
	Generate(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID, addBookmarks bool, addWaterMarks bool) (io.ReadCloser, error)
	GenerateDefault(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (io.ReadCloser, error)
}
