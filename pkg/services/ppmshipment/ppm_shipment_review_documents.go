package ppmshipment

import (
	"fmt"
	"strconv"
	"strings"
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
		updatedPPMShipment.Shipment.MoveTaskOrder.SCCloseoutAssignedID = nil
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

		verrs, err = appCtx.DB().ValidateAndSave(&updatedPPMShipment.Shipment.MoveTaskOrder)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(updatedPPMShipment.ID, nil, verrs, "")
		}
		if err != nil {
			return err
		}

		err = p.signCertificationPPMCloseout(txnAppCtx, updatedPPMShipment.Shipment.MoveTaskOrderID, updatedPPMShipment.ID)
		if err != nil {
			return err
		}

		err = models.CalculatePPMCloseoutSummary(txnAppCtx.DB(), updatedPPMShipment.ID, false)
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

func getPriceParts(rawPrice string, expectedDecimalPlaces int) (int, int, error) {
	// Get rid of a dollar sign if there is one.
	basePrice := strings.Replace(rawPrice, "$", "", -1)
	basePrice = strings.Replace(basePrice, ",", "", -1)

	// Split the string on the decimal point.
	priceParts := strings.Split(basePrice, ".")
	if len(priceParts) != 2 {
		return 0, 0, fmt.Errorf("expected 2 price parts but found %d for price [%s]", len(priceParts), rawPrice)
	}

	integerPart, err := strconv.Atoi(priceParts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("could not convert integer part of price [%s]", rawPrice)
	}

	if len(priceParts[1]) != expectedDecimalPlaces {
		return 0, 0, fmt.Errorf("expected %d decimal places but found %d for price [%s]", expectedDecimalPlaces,
			len(priceParts[1]), rawPrice)
	}

	fractionalPart, err := strconv.Atoi(priceParts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("could not convert fractional part of price [%s]", rawPrice)
	}

	return integerPart, fractionalPart, nil
}

func priceToCents(rawPrice string) (int, error) {
	s := strings.TrimSpace(rawPrice)
	if !strings.Contains(s, "$") {
		return 0, nil
	}
	integerPart, fractionalPart, err := getPriceParts(rawPrice, 2)
	if err != nil {
		return 0, fmt.Errorf("could not parse price [%s]: %w", rawPrice, err)
	}

	cents := (integerPart * 100) + fractionalPart
	return cents, nil
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
		cert, err := p.SignedCertificationCreator.CreateSignedCertification(appCtx, signedCertification)
		fmt.Println(cert)
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
