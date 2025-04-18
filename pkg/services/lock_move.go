package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// MoveLocker is the exported interface for locking moves
//
//go:generate mockery --name MoveLocker
type MoveLocker interface {
	LockMove(appCtx appcontext.AppContext, move *models.Move, officeUserID uuid.UUID) (*models.Move, error)
	LockMoves(appCtx appcontext.AppContext, moveIds []uuid.UUID, officeUserID uuid.UUID) error
}

// MoveUnlocker is the exported interface for unlocking moves
//
//go:generate mockery --name MoveUnlocker
type MoveUnlocker interface {
	UnlockMove(appCtx appcontext.AppContext, move *models.Move, officeUserID uuid.UUID) (*models.Move, error)
	CheckForLockedMovesAndUnlock(appCtx appcontext.AppContext, officeUserID uuid.UUID) error
}
