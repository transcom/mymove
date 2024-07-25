package move

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveCancellation struct{}

func NewMoveCancellation() services.MoveCancellation {
	return &moveCancellation{}
}

func (f *moveCancellation) CancelMove(appCtx appcontext.AppContext, moveID uuid.UUID) (*models.Move, error) {
	move := &models.Move{}
	err := appCtx.DB().Find(move, moveID)
	if err != nil {
		return nil, apperror.NewNotFoundError(moveID, "while looking for a move")
	}

	moveDelta := move
	moveDelta.Status = models.MoveStatusCANCELED

	if verrs := validateMove(appCtx, *move, moveDelta); verrs != nil {
		return move, verrs
	}

	txnErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		verrs, err := txnAppCtx.DB().ValidateAndUpdate(moveDelta)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(
				move.ID, err, verrs, "Validation errors found while setting move status")
		} else if err != nil {
			return apperror.NewQueryError("Move", err, "Failed to update status for move")
		}

		return nil
	})
	if txnErr != nil {
		return nil, txnErr
	}

	return move, nil
}
