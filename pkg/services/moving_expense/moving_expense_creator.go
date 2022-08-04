package movingexpense

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type movingExpenseCreator struct {
	checks []movingExpenseValidator
}

func NewMovingExpenseCreator() services.MovingExpenseCreator {
	return &movingExpenseCreator{
		checks: createChecks(),
	}
}

func (f *movingExpenseCreator) CreateMovingExpense(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.MovingExpense, error) {
	newMovingExpense := &models.MovingExpense{
		PPMShipmentID: ppmShipmentID,
		Document: models.Document{
			ServiceMemberID: appCtx.Session().ServiceMemberID,
		},
	}

	err := validateMovingExpense(appCtx, newMovingExpense, nil, f.checks...)

	if err != nil {
		return nil, err
	}

	txnErr := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		verrs, err := txnCtx.DB().Eager().ValidateAndCreate(newMovingExpense)

		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(uuid.Nil, err, verrs, "")
		} else if err != nil {
			return apperror.NewQueryError("Moving Expense", err, "")
		}

		return nil
	})

	if txnErr != nil {
		return nil, txnErr
	}

	return newMovingExpense, nil
}
