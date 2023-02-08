package movingexpense

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/services"
)

type movingExpenseDeleter struct {
}

func NewMovingExpenseDeleter() services.MovingExpenseDeleter {
	return &movingExpenseDeleter{}
}

func (d *movingExpenseDeleter) DeleteMovingExpense(appCtx appcontext.AppContext, movingExpenseID uuid.UUID) error {
	movingExpense, err := FetchMovingExpenseByID(appCtx, movingExpenseID)
	if err != nil {
		return err
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// movingExpense.Document is a belongs_to relation, so will not be automatically
		// deleted when we call SoftDestroy on the moving expense
		err = utilities.SoftDestroy(appCtx.DB(), &movingExpense.Document)
		if err != nil {
			return err
		}
		err = utilities.SoftDestroy(appCtx.DB(), movingExpense)
		if err != nil {
			return err
		}

		return nil
	})

	if transactionError != nil {
		return transactionError
	}
	return nil
}
