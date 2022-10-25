package signedcertification

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

// signedCertificationValidator is the interface for validating a signed certification against business rules
type signedCertificationValidator interface {
	// Validate checks the newSignedCertification for adherence to business rules.
	// It is safe to return a *validate.Errors with zero added errors as a success case.
	Validate(appCtx appcontext.AppContext, newSignedCertification models.SignedCertification) error
}

// validateSignedCertification checks a signed certification against a passed-in set of business rule checks
func validateSignedCertification(
	appCtx appcontext.AppContext,
	newSignedCertification models.SignedCertification,
	checks ...signedCertificationValidator,
) (result error) {
	verrs := validate.NewErrors()
	for _, checker := range checks {
		if err := checker.Validate(appCtx, newSignedCertification); err != nil {
			switch e := err.(type) {
			case *validate.Errors:
				// accumulate validation errors
				verrs.Append(e)
			default:
				// non-validation errors have priority,
				// and short-circuit doing any further checks
				return err
			}
		}
	}
	if verrs.HasAny() {
		result = apperror.NewInvalidInputError(newSignedCertification.ID, nil, verrs, "Invalid input found while validating the signed certification.")
	}

	return result
}

// signedCertificationValidatorFunc is an adapter type for converting a function into an implementation of signedCertificationValidator
type signedCertificationValidatorFunc func(appcontext.AppContext, models.SignedCertification) error

// Validate fulfills the signedCertificationValidator interface
func (fn signedCertificationValidatorFunc) Validate(appCtx appcontext.AppContext, newSignedCertification models.SignedCertification) error {
	return fn(appCtx, newSignedCertification)
}
