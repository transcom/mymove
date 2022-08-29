package reweigh

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// checks business requirements
// reweighValidator describes a method for checking business requirements
type reweighValidator interface {
	// Validate checks the newReweigh for adherence to business rules. The
	// oldReweigh parameter is expected to be nil in creation use cases.
	// It is safe to return a *validate.Errors with zero added errors as
	// a success case.
	Validate(a appcontext.AppContext, newReweigh models.Reweigh, oldReweigh *models.Reweigh, shipment *models.MTOShipment) error
}

// validateReweigh checks a reweigh against a passed-in set of business rule checks
func validateReweigh(
	appCtx appcontext.AppContext,
	newReweigh models.Reweigh,
	oldReweigh *models.Reweigh,
	shipment *models.MTOShipment,
	checks ...reweighValidator,
) (result error) {
	verrs := validate.NewErrors()
	for _, checker := range checks {
		if err := checker.Validate(appCtx, newReweigh, oldReweigh, shipment); err != nil {
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
		result = apperror.NewInvalidInputError(newReweigh.ID, nil, verrs, "Invalid input found while validating the reweigh.")
	}
	return result
}

// reweighValidatorFunc is an adapter type for converting a function into an implementation of reweighValidator
type reweighValidatorFunc func(appcontext.AppContext, models.Reweigh, *models.Reweigh, *models.MTOShipment) error

// Validate fulfills the reweighValidator interface
func (fn reweighValidatorFunc) Validate(appCtx appcontext.AppContext, newer models.Reweigh, older *models.Reweigh, ship *models.MTOShipment) error {
	return fn(appCtx, newer, older, ship)
}

// mergeReweigh compares NewReweigh and OldReweigh and updates a new MTOReweigh instance with all data
// (changed and unchanged) filled in. Does not return an error, data must be checked for validation before this step.
func mergeReweigh(newReweigh models.Reweigh, oldReweigh *models.Reweigh) *models.Reweigh {
	if oldReweigh == nil {
		return &newReweigh
	}

	reweigh := *oldReweigh

	reweigh.Weight = services.SetOptionalPoundField(newReweigh.Weight, reweigh.Weight)
	reweigh.VerificationReason = services.SetOptionalStringField(newReweigh.VerificationReason, reweigh.VerificationReason)
	reweigh.VerificationProvidedAt = services.SetOptionalDateTimeField(newReweigh.VerificationProvidedAt, reweigh.VerificationProvidedAt)

	return &reweigh
}

// reweighChanged returns true if the reweigh weight has changed, otherwise, returns false
func reweighChanged(newReweigh models.Reweigh, oldReweigh models.Reweigh) bool {
	changed := false

	// Compare updated weight to previous weight to see if anything changed
	if newReweigh.Weight != nil && oldReweigh.Weight != nil {
		if *newReweigh.Weight != *oldReweigh.Weight {
			// change in value
			changed = true
		}
	} else if newReweigh.Weight != nil || oldReweigh.Weight != nil {
		// changing from nil to some value (not nil) or
		// changing from some value (not nil) to nil
		changed = true
	}

	return changed
}
