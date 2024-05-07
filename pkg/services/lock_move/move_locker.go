package lock_move

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveLocker struct {
}

// NewMoveLocker creates a new moveLocker service
func NewMoveLocker() services.MoveLocker {
	return &moveLocker{}
}

// LockMove updates a move with relevant values of who has a move locked and the expiration of the lock pending it isn't unlocked before then
func (m moveLocker) LockMove(appCtx appcontext.AppContext, move *models.Move, officeUserID uuid.UUID) (*models.Move, error) {

	var err error
	if officeUserID == uuid.Nil {
		return &models.Move{}, apperror.NewQueryError("OfficeUserID", err, "No office user provided in request to lock move")
	}

	officeUser, err := models.FetchOfficeUserByID(appCtx.DB(), officeUserID)
	if err != nil {
		return nil, err
	}

	if move.LockedByOfficeUserID != &officeUserID {
		move.LockedByOfficeUserID = &officeUserID
	}

	if officeUser != nil {
		move.LockedByOfficeUser = officeUser
	}

	now := time.Now()
	expirationTime := now.Add(30 * time.Minute)
	move.LockExpiresAt = &expirationTime

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// save the move to the db
		verrs, saveErr := appCtx.DB().ValidateAndSave(move)
		if verrs != nil && verrs.HasAny() {
			invalidInputError := apperror.NewInvalidInputError(move.ID, nil, verrs, "Could not validate move while locking it.")

			return invalidInputError
		}
		if saveErr != nil {
			return err
		}

		return nil
	})

	if transactionError != nil {
		return &models.Move{}, transactionError
	}

	return move, nil
}
