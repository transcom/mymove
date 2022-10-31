package ppmshipment

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ppmShipmentNewSubmitter is the concrete struct implementing the services.PPMShipmentNewSubmitter interface
type ppmShipmentNewSubmitter struct {
	services.SignedCertificationCreator
	services.PPMShipmentRouter
}

// NewPPMShipmentNewSubmitter creates a new ppmShipmentNewSubmitter
func NewPPMShipmentNewSubmitter(
	signedCertificationCreator services.SignedCertificationCreator,
	ppmShipmentRouter services.PPMShipmentRouter,
) services.PPMShipmentNewSubmitter {
	return &ppmShipmentNewSubmitter{
		SignedCertificationCreator: signedCertificationCreator,
		PPMShipmentRouter:          ppmShipmentRouter,
	}
}

// SubmitNewCustomerCloseOut saves a new customer signature for PPM documentation agreement and routes PPM shipment
func (p *ppmShipmentNewSubmitter) SubmitNewCustomerCloseOut(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID, signedCertification models.SignedCertification) (*models.PPMShipment, error) {
	if ppmShipmentID.IsNil() {
		return nil, apperror.NewBadDataError("PPM ID is required")
	}

	ppmShipment, err := models.GetPPMShipment(appCtx, ppmShipmentID)

	if err != nil {
		return nil, err
	}

	signedCertification.SubmittingUserID = appCtx.Session().UserID
	signedCertification.MoveID = ppmShipment.Shipment.MoveTaskOrderID
	signedCertification.PpmID = &ppmShipment.ID

	certType := models.SignedCertificationTypePPMPAYMENT
	signedCertification.CertificationType = &certType

	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		ppmShipment.SignedCertification, err = p.SignedCertificationCreator.CreateSignedCertification(txnAppCtx, signedCertification)

		if err != nil {
			return err
		}

		err = p.PPMShipmentRouter.SubmitCloseOutDocumentation(txnAppCtx, ppmShipment)

		if err != nil {
			return err
		}

		verrs, err := txnAppCtx.DB().ValidateAndUpdate(ppmShipment)

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

	return ppmShipment, nil
}
