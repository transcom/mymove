package lockmove

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveUnlocker struct {
}

// NewMoveLocker creates a new moveLocker service
func NewMoveUnlocker() services.MoveUnlocker {
	return &moveUnlocker{}
}

// UnlockMove updates a move by checking if there are values in the lock_expires_at and locked_by columns and nils them out
// this service object is called when loading queues
func (m moveUnlocker) UnlockMove(appCtx appcontext.AppContext, move *models.Move, officeUserID uuid.UUID) (*models.Move, error) {

	if move == nil {
		return nil, apperror.NewQueryError("Move", nil, "No move provided in request to unlock move")
	}

	if officeUserID == uuid.Nil {
		return nil, apperror.NewQueryError("OfficeUserID", nil, "No office user provided in request to unlock move")
	}

	// nil out all of the columns since the office user is no longer in the move
	if move.LockExpiresAt != nil {
		move.LockExpiresAt = nil
	}

	if move.LockedByOfficeUserID != nil {
		move.LockedByOfficeUserID = nil
	}

	if move.LockedByOfficeUser != nil {
		move.LockedByOfficeUser = nil
	}

	// Store move before update
	var moveBeforeUpdate = *move

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if err := appCtx.DB().RawQuery("UPDATE moves SET locked_by=?, lock_expires_at=?, updated_at=? WHERE id=?", nil, nil, moveBeforeUpdate.UpdatedAt, move.ID).Exec(); err != nil {
			return err
		}

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return move, nil
}

// CheckForUnlockedMovesAndUnlock finds moves with the officeUserID in the locked_by column for the move
// this service object is called when a user logs out
func (m moveUnlocker) CheckForLockedMovesAndUnlock(appCtx appcontext.AppContext, officeUserID uuid.UUID) error {

	if officeUserID == uuid.Nil {
		return apperror.NewQueryError("OfficeUserID", nil, "No office user provided in request to unlock move")
	}

	// get all moves where locked_by matches officeUserID
	var moves []models.Move
	query := appCtx.DB().Where("locked_by = ?", officeUserID)
	err := query.Eager(
		"LockedByOfficeUser",
	).
		All(&moves)
	if err != nil {
		return err
	}

	// iterate through each move and clear the values by using our existing service object above
	if appCtx.Session().IsOfficeUser() && len(moves) > 0 {
		for _, move := range moves {
			lockedOfficeUserID := move.LockedByOfficeUserID
			if lockedOfficeUserID != nil && *lockedOfficeUserID == officeUserID {
				copyOfMove := move
				_, err := m.UnlockMove(appCtx, &copyOfMove, officeUserID)
				if err != nil {
					return err
				}
			}
		}
	}

	return err
}
