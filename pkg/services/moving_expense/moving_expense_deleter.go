package movingexpense

import (
	"database/sql"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type movingExpenseDeleter struct {
}

func NewMovingExpenseDeleter() services.MovingExpenseDeleter {
	return &movingExpenseDeleter{}
}

func (d *movingExpenseDeleter) DeleteMovingExpense(appCtx appcontext.AppContext, ppmID uuid.UUID, movingExpenseID uuid.UUID) error {
	var ppmShipment models.PPMShipment
	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"Shipment.MoveTaskOrder.Orders",
			"MovingExpenses",
		).
		Find(&ppmShipment, ppmID)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewNotFoundError(movingExpenseID, "while looking for MovingExpense")
		}
		return apperror.NewQueryError("MovingExpense fetch original", err, "")
	}

	if appCtx.Session().IsMilApp() {
		if ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID != appCtx.Session().ServiceMemberID {
			wrongServiceMemberIDErr := apperror.NewForbiddenError("Attempted delete by wrong service member")
			appCtx.Logger().Error("internalapi.DeleteMovingExpenseHandler", zap.Error(wrongServiceMemberIDErr))
			return wrongServiceMemberIDErr
		}
	}

	found := false
	for _, lineItem := range ppmShipment.MovingExpenses {
		if lineItem.ID == movingExpenseID {
			found = true
			break
		}
	}
	if !found {
		mismatchedPPMShipmentAndMovingExpenseIDErr := apperror.NewNotFoundError(movingExpenseID, "Moving expense does not exist on ppm shipment")
		appCtx.Logger().Error("internalapi.DeleteMovingExpenseHandler", zap.Error(mismatchedPPMShipmentAndMovingExpenseIDErr))
		return mismatchedPPMShipmentAndMovingExpenseIDErr
	}

	movingExpense, err := FetchMovingExpenseByID(appCtx, movingExpenseID)
	if err != nil {
		return err
	}

	transactionError := appCtx.NewTransaction(func(_ appcontext.AppContext) error {
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
