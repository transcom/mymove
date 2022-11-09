package signedcertification

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// signedCertificationCreator is the concrete struct implementing the services.SignedCertificationCreator interface
type signedCertificationCreator struct {
	checks []signedCertificationValidator
}

// NewSignedCertificationCreator creates a new signedCertificationCreator struct with the basic checks.
func NewSignedCertificationCreator() services.SignedCertificationCreator {
	return &signedCertificationCreator{
		checks: basicSignedCertificationChecks(),
	}
}

// CreateSignedCertification creates a signed certification.
func (s *signedCertificationCreator) CreateSignedCertification(appCtx appcontext.AppContext, signedCertification models.SignedCertification) (*models.SignedCertification, error) {
	if err := validateSignedCertification(appCtx, signedCertification, nil, s.checks...); err != nil {
		return nil, err
	}

	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) (err error) {
		verrs, err := txnAppCtx.DB().ValidateAndCreate(&signedCertification)

		if verrs.HasAny() {
			return apperror.NewInvalidInputError(signedCertification.ID, nil, verrs, "Invalid input found while creating the signed certification.")
		} else if err != nil {
			return apperror.NewQueryError("SignedCertification", err, "Unable to create signed certification")
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return &signedCertification, nil
}
