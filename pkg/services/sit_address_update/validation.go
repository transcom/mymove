package sitaddressupdate

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

// sitAddressUpdateValidator defines the interface for checking business rules for a SIT address update
type sitAddressUpdateValidator interface {
	// Validate checks the new SITAddressUpdate for adherence to business rules.
	// It is safe to return a *validate.Errors with zero errors as a success case.
	Validate(appCtx appcontext.AppContext, sitAddressUpdate *models.SITAddressUpdate) error
}

// validateSignedCertification checks a signed certification against a passed-in set of business rule checks
func validateSITAddressUpdate(
	appCtx appcontext.AppContext,
	sitAddressUpdate *models.SITAddressUpdate,
	checks ...sitAddressUpdateValidator,
) (result error) {
	verrs := validate.NewErrors()
	for _, check := range checks {
		if err := check.Validate(appCtx, sitAddressUpdate); err != nil {
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
		result = apperror.NewInvalidInputError(sitAddressUpdate.ID, nil, verrs, "Invalid input found while validating the SIT address update.")
	}

	return result
}

// sitAddressUpdateValidatorFunc is an adapter that will convert a function into an implementation of sitAddressUpdateValidator
type sitAddressUpdateValidatorFunc func(appcontext.AppContext, *models.SITAddressUpdate) error

// Validate fulfills the sitAddressUpdateValidator interface
func (fn sitAddressUpdateValidatorFunc) Validate(appCtx appcontext.AppContext, sitAddressUpdate *models.SITAddressUpdate) error {
	return fn(appCtx, sitAddressUpdate)
}
