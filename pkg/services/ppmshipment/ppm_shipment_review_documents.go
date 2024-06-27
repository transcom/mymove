package ppmshipment

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/copier"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ppmShipmentReviewDocuments implements the services.PPMShipmentReviewDocuments interface
type ppmShipmentReviewDocuments struct {
	services.PPMShipmentRouter
	services.SignedCertificationCreator
	services.SignedCertificationUpdater
}

// NewPPMShipmentReviewDocuments creates a new ppmShipmentReviewDocuments
func NewPPMShipmentReviewDocuments(
	ppmShipmentRouter services.PPMShipmentRouter,
	signedCertificationCreator services.SignedCertificationCreator,
	signedCertificationUpdater services.SignedCertificationUpdater,
) services.PPMShipmentReviewDocuments {
	return &ppmShipmentReviewDocuments{
		PPMShipmentRouter:          ppmShipmentRouter,
		SignedCertificationCreator: signedCertificationCreator,
		SignedCertificationUpdater: signedCertificationUpdater,
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

	var updatedPPMShipment models.PPMShipment

	err = copier.CopyWithOption(&updatedPPMShipment, ppmShipment, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		return nil, err
	}

	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		err = p.PPMShipmentRouter.SubmitReviewedDocuments(txnAppCtx, &updatedPPMShipment)

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

		err = p.signCertificationPPMCloseout(appCtx, updatedPPMShipment.Shipment.MoveTaskOrderID, updatedPPMShipment.ID)

		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return &updatedPPMShipment, nil
}

func (p *ppmShipmentReviewDocuments) signCertificationPPMCloseout(appCtx appcontext.AppContext, moveID uuid.UUID, ppmShipmentID uuid.UUID) error {
	// Retrieve if PPM has certificate
	signedCertifications, err := models.FetchSignedCertificationPPMByType(appCtx.DB(), appCtx.Session(), moveID, ppmShipmentID, models.SignedCertificationTypeCloseoutReviewedPPMPAYMENT)
	if err != nil {
		return err
	}

	signatureText := fmt.Sprintf("%s %s", appCtx.Session().FirstName, appCtx.Session().LastName)

	if len(signedCertifications) == 0 {
		// Add new certificate
		certType := models.SignedCertificationTypeCloseoutReviewedPPMPAYMENT
		now := time.Now()
		signedCertification := models.SignedCertification{
			SubmittingUserID:  appCtx.Session().UserID,
			MoveID:            moveID,
			PpmID:             models.UUIDPointer(ppmShipmentID),
			CertificationType: &certType,
			CertificationText: "Confirmed: Reviewed Closeout PPM PAYMENT ",
			Signature:         signatureText,
			Date:              now,
		}
		_, err := p.SignedCertificationCreator.CreateSignedCertification(appCtx, signedCertification)
		if err != nil {
			return err
		}
	} else {
		// Update existing certificate. Note, reviews can occur N times.
		eTag := etag.GenerateEtag(signedCertifications[0].UpdatedAt)
		// Update with current counselor information
		signedCertifications[0].SubmittingUserID = appCtx.Session().UserID
		signedCertifications[0].Signature = signatureText
		_, err := p.SignedCertificationUpdater.UpdateSignedCertification(appCtx, *signedCertifications[0], eTag)
		if err != nil {
			return err
		}
	}

	return nil
}
