package signedcertification

import (
	"database/sql"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// signedCertificationUpdater is the concrete struct implementing the services.SignedCertificationUpdater interface
type signedCertificationUpdater struct {
	checks []signedCertificationValidator
}

// NewSignedCertificationUpdater creates a new signedCertificationUpdater struct with the basic checks.
func NewSignedCertificationUpdater() services.SignedCertificationUpdater {
	return &signedCertificationUpdater{
		checks: basicSignedCertificationChecks(),
	}
}

// UpdateSignedCertification updates a signed certification.
func (s *signedCertificationUpdater) UpdateSignedCertification(appCtx appcontext.AppContext, signedCertification models.SignedCertification, eTag string) (*models.SignedCertification, error) {
	originalSignedCertification := &models.SignedCertification{}

	if err := appCtx.DB().Find(originalSignedCertification, signedCertification.ID); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(signedCertification.ID, "while looking for SignedCertification")
		default:
			return nil, apperror.NewQueryError("SignedCertification", err, "")
		}
	}

	if etag.GenerateEtag(originalSignedCertification.UpdatedAt) != eTag {
		return nil, apperror.NewPreconditionFailedError(signedCertification.ID, nil)
	}

	mergedSignedCertification := mergeSignedCertification(signedCertification, originalSignedCertification)

	if err := validateSignedCertification(appCtx, *mergedSignedCertification, originalSignedCertification, s.checks...); err != nil {
		return nil, err
	}

	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) (err error) {
		verrs, err := txnAppCtx.DB().ValidateAndUpdate(mergedSignedCertification)

		if verrs.HasAny() {
			return apperror.NewInvalidInputError(signedCertification.ID, nil, verrs, "Invalid input found while updating the signed certification.")
		} else if err != nil {
			return apperror.NewQueryError("SignedCertification", err, "Unable to update signed certification")
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return mergedSignedCertification, nil
}

// mergeSignedCertification merges the original signed certification with the new one.
func mergeSignedCertification(newSignedCertification models.SignedCertification, originalSignedCertification *models.SignedCertification) *models.SignedCertification {
	mergedSignedCertification := *originalSignedCertification

	if newSignedCertification.CertificationText != "" {
		mergedSignedCertification.CertificationText = newSignedCertification.CertificationText
	}

	if newSignedCertification.Signature != "" {
		mergedSignedCertification.Signature = newSignedCertification.Signature
	}

	if !newSignedCertification.Date.IsZero() {
		mergedSignedCertification.Date = newSignedCertification.Date
	}

	return &mergedSignedCertification
}
