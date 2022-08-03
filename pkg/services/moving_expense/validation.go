package movingexpense

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

type movingExpenseValidator interface {
	Validate(appCtx appcontext.AppContext, newMovingExpense *models.MovingExpense, originalMovingExpense *models.MovingExpense) error
}

type movingExpenseValidatorFunc func(appCtx appcontext.AppContext, newMovingExpense *models.MovingExpense, originalMovingExpense *models.MovingExpense) error

func (fn movingExpenseValidatorFunc) Validate(appCtx appcontext.AppContext, newMovingExpense *models.MovingExpense, originalMovingExpense *models.MovingExpense) error {
	return fn(appCtx, newMovingExpense, originalMovingExpense)
}

func validateMovingExpense(appCtx appcontext.AppContext, newMovingExpense *models.MovingExpense, originalMovingExpense *models.MovingExpense, checks ...movingExpenseValidator) error {
	verrs := validate.NewErrors()

	for _, check := range checks {
		if err := check.Validate(appCtx, newMovingExpense, originalMovingExpense); err != nil {
			switch e := err.(type) {
			case *validate.Errors:
				verrs.Append(e)
			default:
				return err
			}
		}
	}

	if verrs.HasAny() {
		var currentID uuid.UUID
		if newMovingExpense != nil {
			currentID = newMovingExpense.ID
		}
		return apperror.NewInvalidInputError(currentID, nil, verrs, "")
	}

	return nil
}
