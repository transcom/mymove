package move

import (
	"fmt"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

type validator interface {
	Validate(appCtx appcontext.AppContext, move models.Move, delta *models.Move) error
}

type validatorFunc func(appcontext.AppContext, models.Move, *models.Move) error

func (fn validatorFunc) Validate(appCtx appcontext.AppContext, move models.Move, delta *models.Move) error {
	return fn(appCtx, move, delta)
}

func validateMove(appCtx appcontext.AppContext, move models.Move, delta *models.Move, checks ...validator) (result error) {
	verrs := validate.NewErrors()
	for _, checker := range checks {
		if err := checker.Validate(appCtx, move, delta); err != nil {
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
		result = apperror.NewInvalidInputError(move.ID, nil, verrs, "Invalid input found while validating the move.")
	}
	return result
}

// basicChecks are the rules that should always run for move validation
func basicChecks() []validator {
	return []validator{
		checkMoveVisibility(),
	}
}

// primeChecks are the rules that should only run for validating Prime move actions
func primeChecks() []validator {
	return []validator{
		checkMoveVisibility(),
		checkPrimeAvailability(),
	}
}

// checkMoveVisibility verifies that the move in question is NOT deactivated or hidden to the user.
// The delta move in this case is checked to see if the move is being updated to be visible or not.
func checkMoveVisibility() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, move models.Move, delta *models.Move) error {
		isVisible := move.Show != nil && *move.Show

		if delta != nil && delta.Show != nil {
			isVisible = *delta.Show
		}

		if !isVisible {
			appCtx.Logger().Warn(fmt.Sprintf("Attempt to access deactivated move with ID: %s", move.ID.String()))
			return apperror.NewNotFoundError(move.ID, "for move")
		}
		return nil
	})
}

// checkPrimeAvailability verifies that the move in question is visible to the Prime.
// The delta move in this case checks to see if the Prime-availability date was being modified.
// However, there is no way to reset the AvailableToPrimeAt date at the moment.
func checkPrimeAvailability() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, move models.Move, delta *models.Move) error {
		isAvailable := move.AvailableToPrimeAt != nil && !move.AvailableToPrimeAt.IsZero()

		if delta != nil && delta.AvailableToPrimeAt != nil {
			isAvailable = !delta.AvailableToPrimeAt.IsZero()
		}

		if !isAvailable {
			appCtx.Logger().Warn(fmt.Sprintf("Attempt to access non-Prime move with ID: %s", move.ID.String()))
			return apperror.NewNotFoundError(move.ID, "for move")
		}
		return nil
	})
}

// remarksNeededForFinancialFlag checks to make sure that the FinancialReviewRemarks field is properly populated when the
// FinancialReviewFlag is set to true.
func remarksNeededForFinancialFlag() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, move models.Move, delta *models.Move) error {
		verrs := validate.NewErrors()
		// We don't want nil remarks if we're setting the flag to true.
		if delta != nil && delta.FinancialReviewFlag && delta.FinancialReviewRemarks == nil {
			message := fmt.Sprintf("Cannot flag a move for financial review without remarks for move id: %s", move.ID)
			verrs.Add("FinancialReviewRemarks", message)
		}
		// We also don't want an empty string for remarks if we're setting the flag to true.
		if delta != nil && delta.FinancialReviewFlag && delta.FinancialReviewRemarks != nil {
			if *delta.FinancialReviewRemarks == "" {
				message := fmt.Sprintf("Cannot flag a move for financial review with empty remarks for move id: %s", move.ID)
				verrs.Add("FinancialReviewRemarks", message)
			}
		}

		return verrs
	})
}

// checkFinancialFlagRemoval checks to make sure that the FinancialReviewFlagSetAt and FinancialReviewRemarks fields are set to
// nil when the FinancialReviewFlag is being set to false.
func checkFinancialFlagRemoval() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, move models.Move, delta *models.Move) error {
		verrs := validate.NewErrors()
		// If the flag is set to false, then the time stamp and remarks fields should be nil-ed out.
		if delta != nil && !delta.FinancialReviewFlag {
			if delta.FinancialReviewFlagSetAt != nil || delta.FinancialReviewRemarks != nil {
				message := fmt.Sprintf("When removing a financial review flag, the supporting fields "+
					"should be nulled out %s", move.ID)
				verrs.Add("FinancialReviewRemarks", message)
			}
		}

		return verrs

	})
}
