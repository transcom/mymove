package ppmshipment

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/copier"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ppmShipmentReviewDocuments implements the services.PPMShipmentReviewDocuments interface
type ppmShipmentReviewDocuments struct {
	services.PPMShipmentRouter
	services.PPMDocumentFetcher
}

// NewPPMShipmentReviewDocuments creates a new ppmShipmentReviewDocuments
func NewPPMShipmentReviewDocuments(
	ppmShipmentRouter services.PPMShipmentRouter,
	ppmDocumentFetcher services.PPMDocumentFetcher,
) services.PPMShipmentReviewDocuments {
	return &ppmShipmentReviewDocuments{
		PPMShipmentRouter:  ppmShipmentRouter,
		PPMDocumentFetcher: ppmDocumentFetcher,
	}
}

// SubmitReviewedDocuments saves a new customer signature for PPM documentation agreement and routes PPM shipment
func (p *ppmShipmentReviewDocuments) SubmitReviewedDocuments(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.PPMShipment, error) {
	if ppmShipmentID.IsNil() {
		return nil, apperror.NewBadDataError("PPM ID is required")
	}

	ppmShipment, err := FindPPMShipment(appCtx, ppmShipmentID)

	if err != nil {
		return nil, err
	}

	ppmDocuments, err := p.PPMDocumentFetcher.GetPPMDocuments(appCtx, ppmShipment.ShipmentID)
	if err != nil {
		return nil, err
	}

	var updatedPPMShipment models.PPMShipment

	err = copier.CopyWithOption(&updatedPPMShipment, ppmShipment, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		return nil, err
	}

	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		err = p.PPMShipmentRouter.SubmitReviewedDocuments(txnAppCtx, &updatedPPMShipment, *ppmDocuments)

		if err != nil {
			return err
		}

		err = validatePPMShipment(appCtx, updatedPPMShipment, ppmShipment, &ppmShipment.Shipment, PPMShipmentUpdaterChecks...)

		if err != nil {
			return err
		}

		verrs, err := txnAppCtx.DB().ValidateAndUpdate(&updatedPPMShipment)

		if verrs.HasAny() {
			return apperror.NewInvalidInputError(ppmShipment.ID, err, verrs, "unable to validate PPMShipment")
		} else if err != nil {
			return apperror.NewQueryError("PPMShipment", err, "unable to update PPMShipment")
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return &updatedPPMShipment, nil
}
