package ppmshipment

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/copier"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

// ppmShipmentNewSubmitter is the concrete struct implementing the services.PPMShipmentNewSubmitter interface
type ppmShipmentNewSubmitter struct {
	services.PPMShipmentFetcher
	services.SignedCertificationCreator
	services.PPMShipmentRouter
}

// NewPPMShipmentNewSubmitter creates a new ppmShipmentNewSubmitter
func NewPPMShipmentNewSubmitter(
	ppmShipmentFetcher services.PPMShipmentFetcher,
	signedCertificationCreator services.SignedCertificationCreator,
	ppmShipmentRouter services.PPMShipmentRouter,
) services.PPMShipmentNewSubmitter {
	return &ppmShipmentNewSubmitter{
		PPMShipmentFetcher:         ppmShipmentFetcher,
		SignedCertificationCreator: signedCertificationCreator,
		PPMShipmentRouter:          ppmShipmentRouter,
	}
}

// SubmitNewCustomerCloseOut saves a new customer signature for PPM documentation agreement and routes PPM shipment
func (p *ppmShipmentNewSubmitter) SubmitNewCustomerCloseOut(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID, signedCertification models.SignedCertification) (*models.PPMShipment, error) {
	if ppmShipmentID.IsNil() {
		return nil, apperror.NewBadDataError("PPM ID is required")
	}

	ppmShipment, err := p.GetPPMShipment(
		appCtx,
		ppmShipmentID,
		[]string{
			EagerPreloadAssociationShipment,
			EagerPreloadAssociationWeightTickets,
			EagerPreloadAssociationProgearWeightTickets,
			EagerPreloadAssociationMovingExpenses,
			EagerPreloadAssociationW2Address,
		}, []string{
			PostLoadAssociationSignedCertification,
			PostLoadAssociationWeightTicketUploads,
			PostLoadAssociationProgearWeightTicketUploads,
			PostLoadAssociationMovingExpenseUploads,
		},
	)

	if err != nil {
		return nil, err
	}

	signedCertification.SubmittingUserID = appCtx.Session().UserID
	signedCertification.MoveID = ppmShipment.Shipment.MoveTaskOrderID
	signedCertification.PpmID = &ppmShipment.ID

	certType := models.SignedCertificationTypePPMPAYMENT
	signedCertification.CertificationType = &certType

	var updatedPPMShipment models.PPMShipment

	err = copier.CopyWithOption(&updatedPPMShipment, ppmShipment, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		return nil, err
	}

	// initial allowable weight is equal to net weight of all shipments, E-05722
	var allowableWeight = unit.Pound(0)
	if len(updatedPPMShipment.WeightTickets) >= 1 {
		for _, weightTicket := range ppmShipment.WeightTickets {
			allowableWeight += *weightTicket.FullWeight - *weightTicket.EmptyWeight
		}
	}
	updatedPPMShipment.AllowableWeight = &allowableWeight

	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		updatedPPMShipment.SignedCertification, err = p.SignedCertificationCreator.CreateSignedCertification(txnAppCtx, signedCertification)

		if err != nil {
			return err
		}

		err = p.PPMShipmentRouter.SubmitCloseOutDocumentation(txnAppCtx, &updatedPPMShipment)

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
