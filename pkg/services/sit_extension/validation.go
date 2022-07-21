package sitextension

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

// checks business requirements
// sitExtensionValidator describes a method for checking business requirements
type sitExtensionValidator interface {
	// Validate checks the sitExtension for adherence to business rules.
	// It is safe to return a *validate.Errors with zero errors as a success case.
	Validate(a appcontext.AppContext, sitExtension models.SITExtension, shipment *models.MTOShipment) error
}

// validateSITExtension checks a SIT extension against a set of business rule checks
func validateSITExtension(
	appCtx appcontext.AppContext,
	sitExtension models.SITExtension,
	shipment *models.MTOShipment,
	checks ...sitExtensionValidator,
) (result error) {
	verrs := validate.NewErrors()
	for _, checker := range checks {
		if err := checker.Validate(appCtx, sitExtension, shipment); err != nil {
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
		result = apperror.NewInvalidInputError(sitExtension.ID, nil, verrs, "Invalid input found while validating the sitExtension.")
	}
	return result
}

// sitExtensionValidatorFunc is an adapter type for converting a function into an implementation of sitExtensionValidator
type sitExtensionValidatorFunc func(appcontext.AppContext, models.SITExtension, *models.MTOShipment) error

// Validate fulfills the sitExtensionValidator interface
func (fn sitExtensionValidatorFunc) Validate(appCtx appcontext.AppContext, sit models.SITExtension, shipment *models.MTOShipment) error {
	return fn(appCtx, sit, shipment)
}
