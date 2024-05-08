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
}
