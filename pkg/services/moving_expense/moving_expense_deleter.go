package movingexpense

import (
	"errors"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/services"
)

type movingExpenseDeleter struct {
}

func NewMovingExpenseDeleter() services.MovingExpenseDeleter {
	return &movingExpenseDeleter{}
}

func handleSoftDestroyError(err error) error {
	if err == nil {
		return nil
	}
	switch err.Error() {
	case "error updating model":
		return apperror.NewUnprocessableEntityError("while updating model")
	case "this model does not have deleted_at field":
		return apperror.NewPreconditionFailedError(uuid.Nil, errors.New("model or sub table missing deleted_at field"))
	default:
		return apperror.NewInternalServerError("failed attempt to soft delete model")
	}
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
			return handleSoftDestroyError(err)
		}
		err = utilities.SoftDestroy(appCtx.DB(), movingExpense)
		if err != nil {
			return handleSoftDestroyError(err)
		}

		return nil
	})

	if transactionError != nil {
		return transactionError
	}
	return nil
}
